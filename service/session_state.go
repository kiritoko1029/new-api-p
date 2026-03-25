package service

import (
	"sort"
	"time"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
)

type ChannelSessionEntry struct {
	SessionIDMasked  string `json:"session_id_masked"`
	LastActiveAt     int64  `json:"last_active_at"`
	ExpiresAt        int64  `json:"expires_at"`
	RemainingSeconds int64  `json:"remaining_seconds"`
}

type ChannelSessionState struct {
	ChannelID      int                   `json:"channel_id"`
	ChannelName    string                `json:"channel_name"`
	ChannelType    int                   `json:"channel_type"`
	ActiveSessions int                   `json:"active_sessions"`
	MaxSessions    int                   `json:"max_sessions"`
	TTLMinutes     int                   `json:"ttl_minutes"`
	IsFull         bool                  `json:"is_full"`
	Sessions       []ChannelSessionEntry `json:"sessions,omitempty"`
	UpdatedAt      int64                 `json:"updated_at"`
}

type ChannelSessionStreamPayload struct {
	Channels  []ChannelSessionState `json:"channels"`
	UpdatedAt int64                 `json:"updated_at"`
}

func maskSessionID(sessionID string) string {
	if len(sessionID) <= 4 {
		return sessionID
	}
	return sessionID[:2] + "..." + sessionID[len(sessionID)-2:]
}

func buildActiveSessionEntriesFromScores(now time.Time, ttlMinutes int, scoreMap map[string]float64) []ChannelSessionEntry {
	minActiveAt := now.Add(-time.Duration(ttlMinutes) * time.Minute).Unix()
	nowUnix := now.Unix()
	entries := make([]ChannelSessionEntry, 0, len(scoreMap))

	for sessionID, score := range scoreMap {
		lastActiveAt := int64(score)
		if lastActiveAt < minActiveAt {
			continue
		}
		expiresAt := lastActiveAt + int64(ttlMinutes*60)
		remainingSeconds := expiresAt - nowUnix
		if remainingSeconds <= 0 {
			continue
		}
		entries = append(entries, ChannelSessionEntry{
			SessionIDMasked:  maskSessionID(sessionID),
			LastActiveAt:     lastActiveAt,
			ExpiresAt:        expiresAt,
			RemainingSeconds: remainingSeconds,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].LastActiveAt == entries[j].LastActiveAt {
			if entries[i].ExpiresAt == entries[j].ExpiresAt {
				return entries[i].SessionIDMasked < entries[j].SessionIDMasked
			}
			return entries[i].ExpiresAt > entries[j].ExpiresAt
		}
		return entries[i].LastActiveAt > entries[j].LastActiveAt
	})

	return entries
}

func buildClaudeChannelSessionState(channel *model.Channel, now time.Time, scoreMap map[string]float64, includeDetails bool) ChannelSessionState {
	otherSettings := channel.GetOtherSettings()
	ttlMinutes := otherSettings.GetClaudeSessionTTLMinutes()
	maxSessions := otherSettings.GetClaudeMaxSessions()
	entries := buildActiveSessionEntriesFromScores(now, ttlMinutes, scoreMap)

	state := ChannelSessionState{
		ChannelID:      channel.Id,
		ChannelName:    channel.Name,
		ChannelType:    channel.Type,
		ActiveSessions: len(entries),
		MaxSessions:    maxSessions,
		TTLMinutes:     ttlMinutes,
		IsFull:         maxSessions > 0 && len(entries) >= maxSessions,
		UpdatedAt:      now.Unix(),
	}
	if includeDetails {
		state.Sessions = entries
	}
	return state
}

func BuildClaudeChannelSessionStreamPayload(channelIDs []int, includeDetails bool) (ChannelSessionStreamPayload, error) {
	now := time.Now()
	channels, err := getClaudeChannelsForSessionState(channelIDs)
	if err != nil {
		return ChannelSessionStreamPayload{}, err
	}

	states := make([]ChannelSessionState, 0, len(channels))
	for _, channel := range channels {
		if channel == nil {
			continue
		}
		otherSettings := channel.GetOtherSettings()
		ttlMinutes := otherSettings.GetClaudeSessionTTLMinutes()
		scoreMap := GetChannelSessionScoreMap(channel.Id, ttlMinutes)
		states = append(states, buildClaudeChannelSessionState(channel, now, scoreMap, includeDetails))
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].ChannelID < states[j].ChannelID
	})

	return ChannelSessionStreamPayload{
		Channels:  states,
		UpdatedAt: now.Unix(),
	}, nil
}

func getClaudeChannelsForSessionState(channelIDs []int) ([]*model.Channel, error) {
	if len(channelIDs) == 0 {
		channels := make([]*model.Channel, 0)
		err := model.DB.Where("type = ?", constant.ChannelTypeAnthropic).Find(&channels).Error
		return channels, err
	}

	channels, err := model.GetChannelsByIds(channelIDs)
	if err != nil {
		return nil, err
	}

	filtered := make([]*model.Channel, 0, len(channels))
	for _, channel := range channels {
		if channel == nil || channel.Type != constant.ChannelTypeAnthropic {
			continue
		}
		filtered = append(filtered, channel)
	}
	return filtered, nil
}
