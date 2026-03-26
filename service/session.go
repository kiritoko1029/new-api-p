package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/go-redis/redis/v8"
)

const (
	channelSessionsPrefix  = "channel_sessions:"
	channelSessionUsersPrefix = "channel_session_users:"
	// maxSessionContentLen limits the first user message content used for session ID hashing.
	// Truncating avoids unnecessarily large hash inputs while maintaining uniqueness.
	maxSessionContentLen = 500
	// sessionKeyTTL is the Redis key-level TTL for eventual garbage collection.
	// Individual session expiry is handled via score-based (timestamp) filtering.
	sessionKeyTTL = 24 * time.Hour
)

// sessionRequest is a minimal struct for extracting messages from request bodies.
// Supports both OpenAI and Claude native formats since both use the same messages structure.
type sessionRequest struct {
	Messages []sessionMessage `json:"messages"`
}

type sessionMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

// ComputeSessionID derives a stable session identifier from the first user message and token ID.
// Same first user message + same token = same session (conversation continuation).
func ComputeSessionID(firstUserMessage string, tokenID int) string {
	content := firstUserMessage
	if len(content) > maxSessionContentLen {
		content = content[:maxSessionContentLen]
	}
	h := sha256.Sum256([]byte(content + ":" + strconv.Itoa(tokenID)))
	return fmt.Sprintf("%x", h)
}

// ExtractSessionInfo parses the request body to extract the session ID.
// Returns the session ID and true if a valid session was found, empty string and false otherwise.
func ExtractSessionInfo(body []byte, tokenID int) (string, bool) {
	if len(body) == 0 {
		return "", false
	}
	var req sessionRequest
	if err := common.Unmarshal(body, &req); err != nil || len(req.Messages) == 0 {
		return "", false
	}
	firstMsg := extractFirstUserMessageText(req.Messages)
	if firstMsg == "" {
		return "", false
	}
	return ComputeSessionID(firstMsg, tokenID), true
}

// extractFirstUserMessageText finds the first "user" role message and extracts its text content.
func extractFirstUserMessageText(messages []sessionMessage) string {
	for _, msg := range messages {
		if strings.EqualFold(msg.Role, "user") {
			return extractTextContent(msg.Content)
		}
	}
	return ""
}

// extractTextContent handles both string and array-of-objects content formats.
func extractTextContent(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		var result strings.Builder
		for _, item := range v {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if m["type"] == "text" {
				if text, ok := m["text"].(string); ok {
					result.WriteString(text)
				}
			}
		}
		return result.String()
	case json.RawMessage:
		// Handle raw JSON content (from intermediate parsing)
		var str string
		if common.Unmarshal(v, &str) == nil {
			return str
		}
	}
	return ""
}

// sessionRedisKey returns the Redis sorted set key for a channel's sessions.
func sessionRedisKey(channelID int) string {
	return channelSessionsPrefix + strconv.Itoa(channelID)
}

// sessionUsersRedisKey returns the Redis hash key for storing session -> masked username mapping.
func sessionUsersRedisKey(channelID int) string {
	return channelSessionUsersPrefix + strconv.Itoa(channelID)
}

// MaskUsername masks a username for display.
// For 2-char names: first char is replaced with *, e.g., "张三" -> "*三"
// For names with 3+ chars: middle characters are replaced with *, e.g., "张三丰" -> "张*丰"
func MaskUsername(username string) string {
	if username == "" {
		return ""
	}
	runes := []rune(username)
	length := len(runes)

	if length == 1 {
		return "*"
	}
	if length == 2 {
		// 2-char: replace first char with *
		return "*" + string(runes[1])
	}
	// 3+ chars: replace middle characters with *
	if length == 3 {
		return string(runes[0]) + "*" + string(runes[2])
	}
	// For longer names, mask the middle portion
	return string(runes[0]) + strings.Repeat("*", length-2) + string(runes[length-1])
}

// GetChannelActiveSessionCount returns the number of active (non-expired) sessions for a channel.
// Returns 0 if Redis is not enabled or an error occurs.
func GetChannelActiveSessionCount(channelID int, ttlMinutes int) int {
	if !common.RedisEnabled || common.RDB == nil {
		return 0
	}
	ctx := context.Background()
	key := sessionRedisKey(channelID)
	minScore := float64(time.Now().Add(-time.Duration(ttlMinutes) * time.Minute).Unix())

	// Clean up expired sessions first (lazy cleanup)
	common.RDB.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", minScore))

	count, err := common.RDB.ZCount(ctx, key, fmt.Sprintf("%f", minScore), "+inf").Result()
	if err != nil {
		common.SysError("failed to count active sessions: " + err.Error())
		return 0
	}
	return int(count)
}

// GetChannelSessionScoreMap returns the active session -> last-active timestamp map for a channel.
// Expired sessions are lazily cleaned before reading.
func GetChannelSessionScoreMap(channelID int, ttlMinutes int) map[string]float64 {
	if !common.RedisEnabled || common.RDB == nil {
		return map[string]float64{}
	}
	ctx := context.Background()
	key := sessionRedisKey(channelID)
	minScore := float64(time.Now().Add(-time.Duration(ttlMinutes) * time.Minute).Unix())
	minScoreStr := fmt.Sprintf("%f", minScore)

	common.RDB.ZRemRangeByScore(ctx, key, "-inf", minScoreStr)

	results, err := common.RDB.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: minScoreStr,
		Max: "+inf",
	}).Result()
	if err != nil {
		common.SysError("failed to list active sessions: " + err.Error())
		return map[string]float64{}
	}

	scoreMap := make(map[string]float64, len(results))
	for _, result := range results {
		member := fmt.Sprintf("%v", result.Member)
		if member == "" {
			continue
		}
		scoreMap[member] = result.Score
	}
	return scoreMap
}

// GetChannelSessionUserMap returns the active session -> masked username map for a channel.
func GetChannelSessionUserMap(channelID int) map[string]string {
	if !common.RedisEnabled || common.RDB == nil {
		return map[string]string{}
	}
	ctx := context.Background()
	usersKey := sessionUsersRedisKey(channelID)

	result, err := common.RDB.HGetAll(ctx, usersKey).Result()
	if err != nil {
		common.SysError("failed to get session user map: " + err.Error())
		return map[string]string{}
	}
	return result
}

// IsSessionActive checks whether a specific session exists on a channel and is still active (non-expired).
func IsSessionActive(channelID int, sessionID string, ttlMinutes int) bool {
	if !common.RedisEnabled || common.RDB == nil || sessionID == "" {
		return false
	}
	ctx := context.Background()
	key := sessionRedisKey(channelID)
	minScore := float64(time.Now().Add(-time.Duration(ttlMinutes) * time.Minute).Unix())

	score, err := common.RDB.ZScore(ctx, key, sessionID).Result()
	if err != nil {
		return false
	}
	return score >= minScore
}

// RegisterOrUpdateSession adds or refreshes a session for a channel in Redis.
// The session's score (timestamp) is updated to the current time.
// It also stores the masked username in a separate hash for display purposes.
func RegisterOrUpdateSession(channelID int, sessionID string, maskedUsername string) error {
	if !common.RedisEnabled || common.RDB == nil || sessionID == "" {
		return nil
	}
	ctx := context.Background()
	key := sessionRedisKey(channelID)
	usersKey := sessionUsersRedisKey(channelID)
	now := float64(time.Now().Unix())

	_, err := common.RDB.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZAdd(ctx, key, &redis.Z{Score: now, Member: sessionID})
		pipe.Expire(ctx, key, sessionKeyTTL)
		if maskedUsername != "" {
			pipe.HSet(ctx, usersKey, sessionID, maskedUsername)
			pipe.Expire(ctx, usersKey, sessionKeyTTL)
		}
		return nil
	})
	return err
}

// IsChannelSessionFull checks whether a channel has reached its max session limit.
// Returns false (not full) if maxSessions is 0 or negative (unlimited).
func IsChannelSessionFull(channelID int, maxSessions int, ttlMinutes int) bool {
	if maxSessions <= 0 {
		return false
	}
	return GetChannelActiveSessionCount(channelID, ttlMinutes) >= maxSessions
}
