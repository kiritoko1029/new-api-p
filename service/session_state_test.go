package service

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/require"
)

func TestBuildActiveSessionEntriesFromScores_FiltersExpiredAndSorts(t *testing.T) {
	now := time.Unix(1_710_000_000, 0)
	entries := buildActiveSessionEntriesFromScores(now, 15, map[string]float64{
		"session-newer":  float64(now.Add(-2 * time.Minute).Unix()),
		"session-expire": float64(now.Add(-16 * time.Minute).Unix()),
		"session-older":  float64(now.Add(-10 * time.Minute).Unix()),
	})

	require.Len(t, entries, 2)
	require.Equal(t, "se...er", entries[0].SessionIDMasked)
	require.Equal(t, now.Add(-2*time.Minute).Unix(), entries[0].LastActiveAt)
	require.Equal(t, now.Add(13*time.Minute).Unix(), entries[0].ExpiresAt)
	require.EqualValues(t, 13*60, entries[0].RemainingSeconds)
	require.Equal(t, "se...er", entries[1].SessionIDMasked)
	require.Equal(t, now.Add(-10*time.Minute).Unix(), entries[1].LastActiveAt)
}

func TestBuildClaudeChannelSessionState_UsesConfiguredLimitsAndDetails(t *testing.T) {
	now := time.Unix(1_710_000_000, 0)
	channel := &model.Channel{
		Id:            7,
		Name:          "Claude Session Channel",
		Type:          constant.ChannelTypeAnthropic,
		OtherSettings: `{"claude_max_sessions":2,"claude_session_ttl_minutes":30}`,
	}

	state := buildClaudeChannelSessionState(
		channel,
		now,
		map[string]float64{
			"session-abc123": float64(now.Add(-5 * time.Minute).Unix()),
		},
		true,
	)

	require.Equal(t, 7, state.ChannelID)
	require.Equal(t, "Claude Session Channel", state.ChannelName)
	require.Equal(t, 1, state.ActiveSessions)
	require.Equal(t, 2, state.MaxSessions)
	require.Equal(t, 30, state.TTLMinutes)
	require.Len(t, state.Sessions, 1)
	require.Equal(t, "se...23", state.Sessions[0].SessionIDMasked)
	require.Equal(t, now.Unix(), state.UpdatedAt)
}
