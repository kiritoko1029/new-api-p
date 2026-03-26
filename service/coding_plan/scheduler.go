package coding_plan

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

const refreshInterval = 5 * time.Minute

// StartUsageRefreshScheduler 启动编程套餐用量定时刷新
func StartUsageRefreshScheduler() {
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		// 首次启动延迟 30 秒
		time.Sleep(30 * time.Second)
		RefreshAllUsage()
		for range ticker.C {
			RefreshAllUsage()
		}
	}()
	common.SysLog("编程套餐用量定时刷新已启动，间隔: " + refreshInterval.String())
}

// RefreshAllUsage 刷新所有编程套餐渠道的用量
func RefreshAllUsage() {
	channels, err := model.GetCodingPlanChannels()
	if err != nil {
		common.SysError("获取编程套餐渠道失败: " + err.Error())
		return
	}

	for _, channel := range channels {
		provider := GetProvider(channel.Type)
		if provider == nil {
			continue
		}

		keys := channel.GetKeys()
		if len(keys) == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		usage, err := provider.GetUsage(ctx, keys[0])
		cancel()

		if err != nil {
			common.SysError(fmt.Sprintf("渠道[%d]%s用量查询失败: %s", channel.Id, channel.Name, err.Error()))
			// 记录错误到数据库
			errorUsage := &CodingPlanUsage{
				QueryTime: time.Now().UnixMilli(),
				Error:     err.Error(),
			}
			if jsonData, err := common.Marshal(errorUsage); err == nil {
				_ = model.UpdateChannelCodingPlanUsage(channel.Id, string(jsonData))
			}
			continue
		}

		if jsonData, err := common.Marshal(usage); err == nil {
			_ = model.UpdateChannelCodingPlanUsage(channel.Id, string(jsonData))
		}
	}
}

// RefreshChannelUsage 刷新单个渠道的编程套餐用量
func RefreshChannelUsage(channelId int) (*CodingPlanUsage, error) {
	channel, err := model.GetCodingPlanChannelById(channelId)
	if err != nil {
		return nil, fmt.Errorf("渠道不存在: %w", err)
	}

	provider := GetProvider(channel.Type)
	if provider == nil {
		return nil, fmt.Errorf("不支持的渠道类型: %d", channel.Type)
	}

	keys := channel.GetKeys()
	if len(keys) == 0 {
		return nil, fmt.Errorf("渠道没有配置 API Key")
	}

	// 解码 base64 key（渠道 key 存储时可能被编码）
	apiKey := keys[0]
	if decoded, err := base64.StdEncoding.DecodeString(apiKey); err == nil {
		apiKey = string(decoded)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	usage, err := provider.GetUsage(ctx, apiKey)
	cancel()

	if err != nil {
		return nil, err
	}

	if jsonData, err := common.Marshal(usage); err == nil {
		_ = model.UpdateChannelCodingPlanUsage(channelId, string(jsonData))
	}

	return usage, nil
}
