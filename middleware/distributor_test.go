package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSelectChannelWithSessionCheckUsingSelector_ReturnsSessionLimitErrorWithoutFallback(t *testing.T) {
	require.NoError(t, i18n.Init())
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	callCount := 0
	channel, group, err := selectChannelWithSessionCheckUsingSelector(
		ctx,
		"claude-3-7-sonnet",
		"default",
		"session-1",
		func(retry int) (*model.Channel, string, error) {
			callCount++
			return &model.Channel{
				Id:            11,
				Name:          "Claude A",
				Type:          constant.ChannelTypeAnthropic,
				OtherSettings: `{"claude_max_sessions":2,"claude_session_ttl_minutes":30}`,
			}, "default", nil
		},
		func(ch *model.Channel, sessionID string) bool {
			return true
		},
	)

	require.Nil(t, channel)
	require.Equal(t, "default", group)
	require.Equal(t, 1, callCount)

	var apiErr *types.NewAPIError
	require.Error(t, err)
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, types.ErrorCodeChannelSessionLimitReached, apiErr.GetErrorCode())
	require.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	require.Contains(t, apiErr.Error(), "Claude A")
}

func TestSelectChannelWithSessionCheckUsingSelector_ReturnsFirstAvailableCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	callCount := 0
	channel, group, err := selectChannelWithSessionCheckUsingSelector(
		ctx,
		"claude-3-7-sonnet",
		"default",
		"session-1",
		func(retry int) (*model.Channel, string, error) {
			callCount++
			return &model.Channel{
				Id:   22,
				Name: "Claude B",
				Type: constant.ChannelTypeAnthropic,
			}, "default", nil
		},
		func(ch *model.Channel, sessionID string) bool {
			return false
		},
	)

	require.NoError(t, err)
	require.NotNil(t, channel)
	require.Equal(t, 22, channel.Id)
	require.Equal(t, "default", group)
	require.Equal(t, 1, callCount)
}
