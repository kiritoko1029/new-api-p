package coding_plan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// ZhipuCodingPlanProvider 智谱编程套餐提供者
type ZhipuCodingPlanProvider struct {
	baseURL      string
	channelType  int
	providerName string
}

// NewZhipuCodingPlanProvider 创建智谱编程套餐提供者
// isInternational: true=国际版, false=国内版
func NewZhipuCodingPlanProvider(isInternational bool) *ZhipuCodingPlanProvider {
	baseURL := "https://open.bigmodel.cn"
	channelType := 58 // ChannelTypeZhipuCodingPlan
	providerName := "智谱编程套餐"

	if isInternational {
		baseURL = "https://api.z.ai"
		channelType = 59 // ChannelTypeZhipuCodingPlanInternational
		providerName = "智谱编程套餐(国际版)"
	}

	return &ZhipuCodingPlanProvider{
		baseURL:      baseURL,
		channelType:  channelType,
		providerName: providerName,
	}
}

func (p *ZhipuCodingPlanProvider) GetUsage(ctx context.Context, apiKey string) (*CodingPlanUsage, error) {
	url := fmt.Sprintf("%s/api/monitor/usage/quota/limit", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 智谱使用直接 token 认证，非 Bearer
	req.Header.Set("Authorization", apiKey)
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
	var zhipuResp zhipuQuotaLimitResponse
	if err := common.Unmarshal(body, &zhipuResp); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	// 检查 API 返回状态
	if !zhipuResp.Success || zhipuResp.Code != 200 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", zhipuResp.Code, zhipuResp.Msg)
	}

	// 转换为通用格式
	usage := &CodingPlanUsage{
		Level:     zhipuResp.Data.Level,
		QueryTime: time.Now().UnixMilli(),
	}

	for _, limit := range zhipuResp.Data.Limits {
		planLimit := CodingPlanLimit{
			Type:          limit.Type,
			Unit:          limit.Unit,
			Number:        limit.Number,
			Usage:         limit.Usage,
			CurrentValue:  limit.CurrentValue,
			Remaining:     limit.Remaining,
			Percentage:    limit.Percentage,
			NextResetTime: limit.NextResetTime,
		}

		for _, detail := range limit.UsageDetails {
			planLimit.UsageDetails = append(planLimit.UsageDetails, CodingPlanUsageDetail{
				ModelCode: detail.ModelCode,
				Usage:     detail.Usage,
			})
		}

		usage.Limits = append(usage.Limits, planLimit)

		// 使用最大的 NextResetTime 作为全局重置时间
		if limit.NextResetTime > usage.NextResetTime {
			usage.NextResetTime = limit.NextResetTime
		}
	}

	return usage, nil
}

func (p *ZhipuCodingPlanProvider) GetProviderName() string {
	return p.providerName
}

func (p *ZhipuCodingPlanProvider) GetChannelType() int {
	return p.channelType
}

// zhipuQuotaLimitResponse 智谱配额限制响应结构
type zhipuQuotaLimitResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		Limits []struct {
			Type          string `json:"type"`
			Unit          int    `json:"unit"`
			Number        int    `json:"number"`
			Usage         int    `json:"usage"`
			CurrentValue  int    `json:"currentValue"`
			Remaining     int    `json:"remaining"`
			Percentage    int    `json:"percentage"`
			NextResetTime int64  `json:"nextResetTime"`
			UsageDetails  []struct {
				ModelCode string `json:"modelCode"`
				Usage     int    `json:"usage"`
			} `json:"usageDetails"`
		} `json:"limits"`
		Level string `json:"level"`
	} `json:"data"`
}
