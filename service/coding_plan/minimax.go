package coding_plan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// MiniMaxCodingPlanProvider MiniMax编程套餐提供者
type MiniMaxCodingPlanProvider struct {
	baseURL      string
	channelType  int
	providerName string
}

// NewMiniMaxCodingPlanProvider 创建MiniMax编程套餐提供者
func NewMiniMaxCodingPlanProvider() *MiniMaxCodingPlanProvider {
	return &MiniMaxCodingPlanProvider{
		baseURL:      "https://www.minimaxi.com",
		channelType:  60, // ChannelTypeMiniMaxCodingPlan
		providerName: "MiniMax编程套餐",
	}
}

func (p *MiniMaxCodingPlanProvider) GetUsage(ctx context.Context, apiKey string) (*CodingPlanUsage, error) {
	url := fmt.Sprintf("%s/v1/api/openplatform/coding_plan/remains", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// MiniMax使用Bearer Token认证
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB limit
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var minimaxResp minimaxCodingPlanResponse
	if err := common.Unmarshal(body, &minimaxResp); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	// 检查API返回状态
	if minimaxResp.BaseResp.StatusCode != 0 {
		return nil, fmt.Errorf("API error: status_code=%d, status_msg=%s", minimaxResp.BaseResp.StatusCode, minimaxResp.BaseResp.StatusMsg)
	}

	// 转换为通用格式
	usage := &CodingPlanUsage{
		Level:     "pro", // MiniMax不返回level字段，默认pro
		QueryTime: time.Now().UnixMilli(),
	}

	// 使用最大的重置时间作为全局重置时间
	for _, model := range minimaxResp.ModelRemains {
		// 当前区间重置时间
		if model.EndTime > usage.NextResetTime {
			usage.NextResetTime = model.EndTime
		}
		// 周重置时间
		if model.WeeklyEndTime > usage.NextResetTime {
			usage.NextResetTime = model.WeeklyEndTime
		}

		// 将每个模型的用量转换为CodingPlanLimit
		// 当前区间用量
		if model.CurrentIntervalTotalCount > 0 || model.CurrentIntervalUsageCount > 0 {
			intervalLimit := CodingPlanLimit{
				Type:          "INTERVAL_" + model.ModelName,
				Unit:          4, // 4=日/区间
				Number:        model.CurrentIntervalTotalCount,
				Usage:         model.CurrentIntervalTotalCount,
				CurrentValue:  model.CurrentIntervalUsageCount,
				Remaining:     model.CurrentIntervalTotalCount - model.CurrentIntervalUsageCount,
				Percentage:    0,
				NextResetTime: model.EndTime,
			}
			if model.CurrentIntervalTotalCount > 0 {
				intervalLimit.Percentage = (model.CurrentIntervalUsageCount * 100) / model.CurrentIntervalTotalCount
			}
			usage.Limits = append(usage.Limits, intervalLimit)
		}

		// 周用量
		if model.CurrentWeeklyTotalCount > 0 || model.CurrentWeeklyUsageCount > 0 {
			weeklyLimit := CodingPlanLimit{
				Type:          "WEEKLY_" + model.ModelName,
				Unit:          3, // 3=周
				Number:        model.CurrentWeeklyTotalCount,
				Usage:         model.CurrentWeeklyTotalCount,
				CurrentValue:  model.CurrentWeeklyUsageCount,
				Remaining:     model.CurrentWeeklyTotalCount - model.CurrentWeeklyUsageCount,
				Percentage:    0,
				NextResetTime: model.WeeklyEndTime,
			}
			if model.CurrentWeeklyTotalCount > 0 {
				weeklyLimit.Percentage = (model.CurrentWeeklyUsageCount * 100) / model.CurrentWeeklyTotalCount
			}
			usage.Limits = append(usage.Limits, weeklyLimit)
		}
	}

	return usage, nil
}

func (p *MiniMaxCodingPlanProvider) GetProviderName() string {
	return p.providerName
}

func (p *MiniMaxCodingPlanProvider) GetChannelType() int {
	return p.channelType
}

// minimaxCodingPlanResponse MiniMax编程套餐响应结构
type minimaxCodingPlanResponse struct {
	ModelRemains []struct {
		StartTime                 int64  `json:"start_time"`
		EndTime                   int64  `json:"end_time"`
		RemainsTime              int64  `json:"remains_time"`
		CurrentIntervalTotalCount int    `json:"current_interval_total_count"`
		CurrentIntervalUsageCount int    `json:"current_interval_usage_count"`
		ModelName                 string `json:"model_name"`
		CurrentWeeklyTotalCount   int    `json:"current_weekly_total_count"`
		CurrentWeeklyUsageCount   int    `json:"current_weekly_usage_count"`
		WeeklyStartTime          int64  `json:"weekly_start_time"`
		WeeklyEndTime            int64  `json:"weekly_end_time"`
		WeeklyRemainsTime        int64  `json:"weekly_remains_time"`
	} `json:"model_remains"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}
