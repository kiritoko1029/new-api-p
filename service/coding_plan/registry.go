package coding_plan

import (
	"sync"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[int]CodingPlanProvider)
)

// RegisterProvider 注册编程套餐提供者
func RegisterProvider(channelType int, provider CodingPlanProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[channelType] = provider
}

// GetProvider 获取编程套餐提供者
func GetProvider(channelType int) CodingPlanProvider {
	providersMu.RLock()
	defer providersMu.RUnlock()
	return providers[channelType]
}

// GetAllProviders 获取所有注册的提供者
func GetAllProviders() map[int]CodingPlanProvider {
	providersMu.RLock()
	defer providersMu.RUnlock()
	result := make(map[int]CodingPlanProvider, len(providers))
	for k, v := range providers {
		result[k] = v
	}
	return result
}

func init() {
	// 注册智谱编程套餐（国内版）
	RegisterProvider(58, NewZhipuCodingPlanProvider(false))
	// 注册智谱编程套餐（国际版）
	RegisterProvider(59, NewZhipuCodingPlanProvider(true))
	// 注册MiniMax编程套餐
	RegisterProvider(60, NewMiniMaxCodingPlanProvider())
}
