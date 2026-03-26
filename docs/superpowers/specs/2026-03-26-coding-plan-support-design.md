# 编程套餐支持功能设计文档

## 概述

在渠道管理中添加对编程套餐的支持，允许用户配置智谱、MiniMax 等厂商的编程套餐渠道，并实时显示套餐用量信息。

## 需求背景

编程套餐是部分 AI 厂商（如智谱、MiniMax）提供的特殊计费方案，与传统按 Token 计费不同，编程套餐通常包含：
- 固定时长的调用次数限额（如 5 小时、周限额）
- MCP 工具调用次数限制
- 套餐等级信息

用户需要在渠道列表中直接查看这些用量信息，以便及时了解套餐使用情况。

## 设计目标

1. **可扩展性** - 架构支持轻松添加新的编程套餐类型
2. **统一显示** - 在渠道列表中用统一的格式展示不同套餐的用量
3. **自动刷新** - 支持定时自动刷新 + 手动立即刷新
4. **配置灵活** - 刷新频率可在系统设置中配置

---

## 架构设计

### 1. 编程套餐类型常量

在 `constant/channel.go` 中新增编程套餐渠道类型：

```go
const (
    // ... 现有类型 ...
    ChannelTypeZhipuCodingPlan            = 58 // 智谱编程套餐（国内）
    ChannelTypeZhipuCodingPlanInternational = 59 // 智谱编程套餐（国际版）
    ChannelTypeMiniMaxCodingPlan          = 60 // MiniMax编程套餐（预留）
    // 预留 61-69 用于未来其他编程套餐
)
```

### 2. 编程套餐接口抽象

创建 `service/coding_plan/` 目录，定义统一的编程套餐接口：

```go
// service/coding_plan/interface.go

// CodingPlanUsage 编程套餐用量信息
type CodingPlanUsage struct {
    // 通用字段
    Level         string `json:"level"`           // 套餐等级：pro, plus, free 等
    NextResetTime int64  `json:"next_reset_time"` // 下次重置时间（Unix毫秒时间戳）

    // 限额列表（不同套餐可能有多个限额维度）
    Limits []CodingPlanLimit `json:"limits"`

    // 原始响应（用于调试）
    RawResponse string `json:"raw_response,omitempty"`

    // 查询时间
    QueryTime int64 `json:"query_time"`

    // 查询错误
    Error string `json:"error,omitempty"`
}

// CodingPlanLimit 单个限额维度
type CodingPlanLimit struct {
    Type          string                 `json:"type"`           // 限额类型：TIME_LIMIT, TOKENS_LIMIT
    Unit          int                    `json:"unit"`           // 时间单位：3=周, 5=5小时
    Number        int                    `json:"number"`         // 限额数量
    Usage         int                    `json:"usage"`          // 总限额（如 1000 次）
    CurrentValue  int                    `json:"current_value"`  // 当前已使用值
    Remaining     int                    `json:"remaining"`      // 剩余
    Percentage    int                    `json:"percentage"`     // 已用百分比
    NextResetTime int64                  `json:"next_reset_time"` // 该限额的重置时间
    UsageDetails  []CodingPlanUsageDetail `json:"usage_details"`  // 各模型/工具使用详情
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
```

### 3. 智谱编程套餐实现

#### API 端点

智谱编程套餐提供三个用量查询端点：

| 端点 | 描述 | 用途 |
|------|------|------|
| `/api/monitor/usage/quota/limit` | 配额限制 | 主要端点，获取限额和使用情况 |
| `/api/monitor/usage/model-usage` | 模型使用 | 各模型调用统计 |
| `/api/monitor/usage/tool-usage` | 工具使用 | MCP 工具调用统计 |

**Base URL**:
- 国内版: `https://open.bigmodel.cn` 或 `https://dev.bigmodel.cn`
- 国际版: `https://api.z.ai`

**认证方式**: `Authorization: {apiKey}` (直接使用 token，非 Bearer)

#### 响应结构

`/api/monitor/usage/quota/limit` 返回：
```json
{
    "limits": [
        {
            "type": "TIME_LIMIT",
            "unit": 5,
            "number": 1,
            "usage": 1000,
            "currentValue": 29,
            "remaining": 971,
            "percentage": 2,
            "nextResetTime": 1774921618998,
            "usageDetails": [
                {"modelCode": "search-prime", "usage": 8},
                {"modelCode": "web-reader", "usage": 21}
            ]
        },
        {
            "type": "TOKENS_LIMIT",
            "unit": 3,
            "number": 5,
            "percentage": 5,
            "nextResetTime": 1774528750587
        }
    ],
    "level": "pro"
}
```

**字段说明**:
- `type`: `TIME_LIMIT` = MCP 工具调用次数限制, `TOKENS_LIMIT` = Token 用量限制
- `unit`: 时间单位 - `3` = 周, `5` = 5小时
- `level`: 套餐等级 - `pro`, `plus`, `free` 等

#### Go 实现

```go
// service/coding_plan/zhipu.go

type ZhipuCodingPlanProvider struct {
    baseURL string
}

func NewZhipuCodingPlanProvider(isInternational bool) *ZhipuCodingPlanProvider {
    baseURL := "https://open.bigmodel.cn"
    if isInternational {
        baseURL = "https://api.z.ai"
    }
    return &ZhipuCodingPlanProvider{baseURL: baseURL}
}

func (p *ZhipuCodingPlanProvider) GetUsage(ctx context.Context, apiKey string) (*CodingPlanUsage, error) {
    // GET /api/monitor/usage/quota/limit
    // Header: Authorization: {apiKey} (非 Bearer)
    url := fmt.Sprintf("%s/api/monitor/usage/quota/limit", p.baseURL)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    // 注意：智谱使用直接 token 认证，非 Bearer
    req.Header.Set("Authorization", apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
    }

    // 解析响应
    var zhipuResp ZhipuQuotaLimitResponse
    if err := json.Unmarshal(body, &zhipuResp); err != nil {
        return nil, err
    }

    // 转换为通用格式
    usage := &CodingPlanUsage{
        Level:     zhipuResp.Level,
        QueryTime: time.Now().UnixMilli(),
    }

    for _, limit := range zhipuResp.Limits {
        usage.Limits = append(usage.Limits, CodingPlanLimit{
            Type:          limit.Type,
            Unit:          limit.Unit,
            Number:        limit.Number,
            Usage:         limit.Usage,
            CurrentValue:  limit.CurrentValue,
            Remaining:     limit.Remaining,
            Percentage:    limit.Percentage,
            NextResetTime: limit.NextResetTime,
            UsageDetails:  limit.UsageDetails,
        })
        if limit.NextResetTime > usage.NextResetTime {
            usage.NextResetTime = limit.NextResetTime
        }
    }

    return usage, nil
}

// 智谱响应结构
type ZhipuQuotaLimitResponse struct {
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
}
```

### 4. 编程套餐注册表

```go
// service/coding_plan/registry.go

var providers = make(map[int]CodingPlanProvider)

func RegisterProvider(channelType int, provider CodingPlanProvider) {
    providers[channelType] = provider
}

func GetProvider(channelType int) CodingPlanProvider {
    return providers[channelType]
}

func IsCodingPlanChannel(channelType int) bool {
    _, ok := providers[channelType]
    return ok
}

func init() {
    // 注册智谱编程套餐
    RegisterProvider(constant.ChannelTypeZhipuCodingPlan,
        NewZhipuCodingPlanProvider(false))
    RegisterProvider(constant.ChannelTypeZhipuCodingPlanInternational,
        NewZhipuCodingPlanProvider(true))
}
```

### 5. 数据模型扩展

在 `model/channel.go` 中扩展 Channel 模型，添加编程套餐用量缓存字段：

```go
type Channel struct {
    // ... 现有字段 ...

    // 编程套餐用量缓存（JSON存储）
    CodingPlanUsage string `json:"coding_plan_usage" gorm:"type:text"`
    // 用量查询时间
    CodingPlanUsageUpdatedTime int64 `json:"coding_plan_usage_updated_time" gorm:"bigint"`
}
```

用量缓存结构：
```go
type ChannelCodingPlanUsage struct {
    Level         string                `json:"level"`
    NextResetTime int64                 `json:"next_reset_time"`
    Limits        []CodingPlanLimit     `json:"limits"`
    QueryTime     int64                 `json:"query_time"`
    Error         string                `json:"error,omitempty"`
}
```

### 6. API 端点

新增以下 API 端点：

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/channel/:id/coding-plan-usage` | GET | 手动刷新单个渠道的编程套餐用量 |
| `/api/channel/coding-plan-usage/batch` | POST | 批量刷新多个渠道的用量 |

### 7. 定时刷新任务

在 `service/coding_plan/scheduler.go` 中实现定时刷新：

```go
type CodingPlanUsageScheduler struct {
    interval time.Duration
    stopCh   chan struct{}
}

func (s *CodingPlanUsageScheduler) Start() {
    ticker := time.NewTicker(s.interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                s.refreshAllCodingPlanUsage()
            case <-s.stopCh:
                ticker.Stop()
                return
            }
        }
    }()
}

func (s *CodingPlanUsageScheduler) refreshAllCodingPlanUsage() {
    // 获取所有编程套餐渠道
    channels, err := model.GetCodingPlanChannels()
    if err != nil {
        return
    }

    // 并发刷新（限制并发数）
    sem := make(chan struct{}, 5)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(channel *model.Channel) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            usage, err := refreshChannelCodingPlanUsage(channel)
            if err != nil {
                log.Printf("刷新渠道 %d 用量失败: %v", channel.Id, err)
            }
            _ = usage // 更新数据库
        }(ch)
    }
    wg.Wait()
}
```

### 8. 系统设置

在 `setting/performance_setting/` 中添加编程套餐刷新间隔配置：

```go
var CodingPlanUsageRefreshInterval = 60 * time.Minute // 默认1小时
```

---

## 前端设计

### 1. 渠道类型选项

在 `web/src/constants/index.js` 中添加编程套餐渠道类型：

```javascript
{
    value: 58,
    label: 'ZhipuCodingPlan',
    color: 'amber',
    icon: 'zhipu'
},
{
    value: 59,
    label: 'ZhipuCodingPlanInternational',
    color: 'amber',
    icon: 'zhipu'
},
// 预留 MiniMax
```

### 2. 用量列渲染

修改 `ChannelsColumnDefs.jsx` 中的余额列，针对编程套餐渠道显示用量信息：

```jsx
const renderCodingPlanUsage = (record, t) => {
    const usage = parseCodingPlanUsage(record.coding_plan_usage);
    if (!usage) {
        return (
            <Tag color='grey' shape='circle'>
                {t('未查询')}
            </Tag>
        );
    }

    if (usage.error) {
        return (
            <Tooltip content={usage.error}>
                <Tag color='red' shape='circle'>
                    {t('查询失败')}
                </Tag>
            </Tooltip>
        );
    }

    // 渲染限额信息
    return (
        <div className="flex flex-col gap-1">
            {usage.limits.map((limit, idx) => (
                <Tooltip key={idx} content={formatLimitDetail(limit, t)}>
                    <Tag
                        color={getLimitColor(limit.percentage)}
                        shape='circle'
                        className='cursor-pointer'
                        onClick={() => refreshUsage(record.id)}
                    >
                        {formatLimitBrief(limit, t)}
                    </Tag>
                </Tooltip>
            ))}
            <span className="text-xs text-gray-500">
                {usage.level} | {formatTimeToReset(usage.next_reset_time)}
            </span>
        </div>
    );
};

// 限额简要显示
const formatLimitBrief = (limit, t) => {
    switch (limit.type) {
    case 'TIME_LIMIT':
        return `${limit.current_value}/${limit.usage} ${t('次')} (${limit.remaining}${t('剩余')})`;
    case 'TOKENS_LIMIT':
        return `Token: ${limit.percentage}%`;
    default:
        return `${limit.percentage}%`;
    }
};
```

### 3. 刷新交互

- 点击用量标签：立即刷新该渠道用量
- 刷新后显示 loading 状态
- 支持批量刷新（在渠道操作栏添加"刷新用量"按钮）

---

## 实现计划

### Phase 1: 后端基础架构
1. 新增渠道类型常量
2. 创建编程套餐接口和服务层
3. 实现智谱编程套餐用量查询
4. 扩展 Channel 模型（添加用量缓存字段）
5. 添加 API 端点

### Phase 2: 定时刷新
1. 实现定时刷新调度器
2. 添加系统设置配置项
3. 启动时自动加载调度器

### Phase 3: 前端实现
1. 添加编程套餐渠道类型选项
2. 修改渠道列表用量列渲染
3. 实现手动刷新交互
4. 添加刷新状态反馈

### Phase 4: 测试与优化
1. 单元测试
2. 集成测试
3. 性能优化

---

## 扩展指南

### 添加新的编程套餐类型

1. 在 `constant/channel.go` 添加新的渠道类型常量
2. 在 `service/coding_plan/` 创建新的 Provider 实现
3. 在 `init()` 中注册 Provider
4. 在前端 `constants/index.js` 添加渠道类型选项
5. 数据库迁移自动处理

### MiniMax 编程套餐接入示例

```go
// service/coding_plan/minimax.go
type MiniMaxCodingPlanProvider struct{}

func (p *MiniMaxCodingPlanProvider) GetUsage(ctx context.Context, apiKey string) (*CodingPlanUsage, error) {
    // 实现 MiniMax 用量查询
}

func init() {
    RegisterProvider(constant.ChannelTypeMiniMaxCodingPlan,
        &MiniMaxCodingPlanProvider{})
}
```

---

## 风险与缓解

| 风险 | 缓解措施 |
|------|----------|
| 上游 API 不可用 | 缓存上次成功结果，显示错误状态 |
| 刷新频率过高被封 | 可配置刷新间隔，限制并发数 |
| 用量格式变化 | 保留原始响应用于调试 |

---

## 测试要点

1. **单元测试**
   - 智谱用量响应解析
   - Provider 注册与获取
   - 用量格式化函数

2. **集成测试**
   - API 端点调用
   - 定时刷新触发
   - 数据库读写

3. **手动测试**
   - 创建编程套餐渠道
   - 用量刷新交互
   - 错误状态显示
