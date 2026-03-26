package coding_plan

import (
	"context"
)

// CodingPlanUsage 编程套餐用量信息
type CodingPlanUsage struct {
	// 通用字段
	Level         string             `json:"level"`           // 套餐等级：pro, plus, free 等
	NextResetTime int64              `json:"next_reset_time"` // 下次重置时间（Unix毫秒时间戳）

	// 限额列表（不同套餐可能有多个限额维度）
	Limits []CodingPlanLimit `json:"limits"`

	// 查询时间
	QueryTime int64 `json:"query_time"`

	// 查询错误
	Error string `json:"error,omitempty"`
}

// CodingPlanLimit 单个限额维度
type CodingPlanLimit struct {
	Type          string                  `json:"type"`            // 限额类型：TIME_LIMIT, TOKENS_LIMIT
	Unit          int                     `json:"unit"`            // 时间单位：3=周, 5=5小时
	Number        int                     `json:"number"`          // 限额数量
	Usage         int                     `json:"usage"`           // 总限额（如 1000 次）
	CurrentValue  int                     `json:"current_value"`   // 当前已使用值
	Remaining     int                     `json:"remaining"`       // 剩余
	Percentage    int                     `json:"percentage"`      // 已用百分比
	NextResetTime int64                   `json:"next_reset_time"` // 该限额的重置时间
	UsageDetails  []CodingPlanUsageDetail `json:"usage_details"`   // 各模型/工具使用详情
}

// CodingPlanUsageDetail 模型/工具使用详情
type CodingPlanUsageDetail struct {
	ModelCode string `json:"model_code"` // 模型/工具代码
	Usage     int    `json:"usage"`      // 使用量
}

// CodingPlanProvider 编程套餐提供者接口
type CodingPlanProvider interface {
	// GetUsage 查询用量信息
	GetUsage(ctx context.Context, apiKey string) (*CodingPlanUsage, error)

	// GetProviderName 获取提供者名称
	GetProviderName() string

	// GetChannelType 获取关联的渠道类型
	GetChannelType() int
}

// IsCodingPlanChannel 判断是否为编程套餐渠道类型
func IsCodingPlanChannel(channelType int) bool {
	switch channelType {
	case 58, 59: // ChannelTypeZhipuCodingPlan, ChannelTypeZhipuCodingPlanInternational
		return true
	default:
		return false
	}
}
