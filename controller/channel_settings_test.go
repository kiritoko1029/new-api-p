package controller

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupChannelControllerTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gin.SetMode(gin.TestMode)
	common.UsingSQLite = true
	common.UsingMySQL = false
	common.UsingPostgreSQL = false
	common.RedisEnabled = false

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	model.DB = db
	model.LOG_DB = db

	if err := db.AutoMigrate(&model.Channel{}, &model.Ability{}); err != nil {
		t.Fatalf("failed to migrate channel tables: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

func seedChannel(t *testing.T, db *gorm.DB, channel *model.Channel) *model.Channel {
	t.Helper()

	if err := db.Create(channel).Error; err != nil {
		t.Fatalf("failed to create channel: %v", err)
	}
	return channel
}

func TestUpdateChannelPersistsSettingsAndGetChannelReturnsThem(t *testing.T) {
	db := setupChannelControllerTestDB(t)
	channel := seedChannel(t, db, &model.Channel{
		Type:          14,
		Key:           "claude-key",
		Status:        common.ChannelStatusEnabled,
		Name:          "claude-session-test",
		Models:        "claude-3-7-sonnet",
		Group:         "default",
		OtherSettings: `{"claude_max_sessions":2,"claude_session_ttl_minutes":30}`,
	})

	if !db.Migrator().HasColumn(&model.Channel{}, "settings") {
		t.Fatalf("expected channels table to include settings column")
	}

	updateBody := map[string]any{
		"id":       channel.Id,
		"type":     14,
		"settings": `{"claude_max_sessions":5,"claude_session_ttl_minutes":45}`,
	}

	updateCtx, updateRecorder := newAuthenticatedContext(t, http.MethodPut, "/api/channel/", updateBody, 1)
	UpdateChannel(updateCtx)

	updateResponse := decodeAPIResponse(t, updateRecorder)
	if !updateResponse.Success {
		t.Fatalf("expected update success, got message: %s", updateResponse.Message)
	}

	reloaded, err := model.GetChannelById(channel.Id, true)
	if err != nil {
		t.Fatalf("failed to reload updated channel: %v", err)
	}

	updatedSettings := reloaded.GetOtherSettings()
	if updatedSettings.GetClaudeMaxSessions() != 5 {
		t.Fatalf("expected claude_max_sessions=5 after update, got %d (raw=%s)", updatedSettings.GetClaudeMaxSessions(), reloaded.OtherSettings)
	}
	if updatedSettings.GetClaudeSessionTTLMinutes() != 45 {
		t.Fatalf("expected claude_session_ttl_minutes=45 after update, got %d (raw=%s)", updatedSettings.GetClaudeSessionTTLMinutes(), reloaded.OtherSettings)
	}

	getCtx, getRecorder := newAuthenticatedContext(t, http.MethodGet, "/api/channel/"+strconv.Itoa(channel.Id), nil, 1)
	getCtx.Params = gin.Params{{Key: "id", Value: strconv.Itoa(channel.Id)}}
	GetChannel(getCtx)

	getResponse := decodeAPIResponse(t, getRecorder)
	if !getResponse.Success {
		t.Fatalf("expected get success, got message: %s", getResponse.Message)
	}

	returnedChannel := &model.Channel{}
	if err := common.Unmarshal(getResponse.Data, returnedChannel); err != nil {
		t.Fatalf("failed to decode returned channel: %v", err)
	}

	returnedSettings := returnedChannel.GetOtherSettings()
	if returnedSettings.GetClaudeMaxSessions() != 5 {
		t.Fatalf("expected GET response claude_max_sessions=5, got %d (raw=%s)", returnedSettings.GetClaudeMaxSessions(), returnedChannel.OtherSettings)
	}
	if returnedSettings.GetClaudeSessionTTLMinutes() != 45 {
		t.Fatalf("expected GET response claude_session_ttl_minutes=45, got %d (raw=%s)", returnedSettings.GetClaudeSessionTTLMinutes(), returnedChannel.OtherSettings)
	}
}

func TestChannelOtherSettingsJSONRoundTrip(t *testing.T) {
	channel := &model.Channel{}
	maxSessions := 7
	ttlMinutes := 60

	channel.SetOtherSettings(dto.ChannelOtherSettings{
		ClaudeMaxSessions:       &maxSessions,
		ClaudeSessionTTLMinutes: &ttlMinutes,
	})

	settings := channel.GetOtherSettings()
	if settings.GetClaudeMaxSessions() != 7 {
		t.Fatalf("expected round-trip claude_max_sessions=7, got %d", settings.GetClaudeMaxSessions())
	}
	if settings.GetClaudeSessionTTLMinutes() != 60 {
		t.Fatalf("expected round-trip claude_session_ttl_minutes=60, got %d", settings.GetClaudeSessionTTLMinutes())
	}
}
