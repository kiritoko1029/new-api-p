[Root](../CLAUDE.md) > **relay**

# relay -- AI API Relay/Proxy Engine

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

The relay module is the core of new-api's proxy functionality. It receives AI API requests from clients, selects the best upstream channel, converts request/response formats between different providers, handles streaming, and manages billing.

## Entry & Startup

- Relay is initialized in `main.go` via `middleware.Distribute()` which triggers channel selection
- The `controller/relay.go` handler is the main entry point, called from `router/relay-router.go`
- Task relay handlers: `controller.RelayMidjourney`, `controller.RelayTask`

## Public Interfaces

### Relay Formats (`types/relay_format.go`)
- `RelayFormatOpenAI` -- Standard OpenAI chat completions
- `RelayFormatOpenAIResponses` -- OpenAI Responses API
- `RelayFormatOpenAIResponsesCompaction` -- Responses API compaction
- `RelayFormatOpenAIImage` -- Image generation/editing
- `RelayFormatOpenAIAudio` -- Audio TTS/STT
- `RelayFormatOpenAIRealtime` -- WebSocket realtime
- `RelayFormatClaude` -- Claude/Anthropic messages API
- `RelayFormatGemini` -- Google Gemini native format
- `RelayFormatEmbedding` -- Text embeddings
- `RelayFormatRerank` -- Rerank API

### Channel Adaptor Interface (`relay/channel/adapter.go`)
```go
type Adaptor interface {
    Init(info *relaycommon.RelayInfo)
    GetRequestURL(info *relaycommon.RelayInfo) (string, error)
    SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
    ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
    ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error)
    ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)
    ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error)
    ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)
    ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error)
    DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
    DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
    GetModelList() []string
    GetChannelName() string
    ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error)
    ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error)
}
```

### TaskAdaptor Interface (for async providers like Midjourney, Suno, Kling)
```go
type TaskAdaptor interface {
    Init(info *relaycommon.RelayInfo)
    ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError
    EstimateBilling(c *gin.Context, info *relaycommon.RelayInfo) map[string]float64
    AdjustBillingOnSubmit(info *relaycommon.RelayInfo, taskData []byte) map[string]float64
    AdjustBillingOnComplete(info *relaycommon.RelayInfo, taskData []byte) float64
    // ... plus task lifecycle methods
}
```

## Supported Providers (42+)

| Provider | Directory | Channel Type | Formats |
|----------|-----------|-------------|---------|
| OpenAI | `openai/` | 1 | Chat, Responses, Audio, Image, Embedding, Realtime |
| Anthropic Claude | `claude/` | 14 | Messages (native) + OpenAI-compatible |
| Google Gemini | `gemini/` | 24 | Native + OpenAI-compatible |
| Azure OpenAI | (via openai) | 3 | Chat, Responses |
| AWS Bedrock | `aws/` | 33 | Multiple providers via Bedrock |
| Vertex AI | `vertex/` | 41 | Gemini via Vertex |
| DeepSeek | `deepseek/` | 43 | Chat |
| Mistral | `mistral/` | 42 | Chat |
| Cohere | `cohere/` | 34 | Chat |
| Perplexity | `perplexity/` | 27 | Chat |
| Ollama | `ollama/` | 4 | Chat, TTS |
| Cloudflare | `cloudflare/` | 39 | Chat |
| Ali (Qwen) | `ali/` | 17 | Chat, Image, Rerank |
| Baidu | `baidu/`, `baidu_v2/` | 15, 46 | Chat |
| Zhipu | `zhipu/`, `zhipu_4v/` | 16, 26 | Chat |
| Tencent | `tencent/` | 23 | Chat |
| Xunfei | `xunfei/` | 18 | Chat |
| Moonshot | `moonshot/` | 25 | Chat |
| MiniMax | `minimax/` | 35 | Chat, TTS |
| xAI (Grok) | `xai/` | 48 | Chat |
| Coze | `coze/` | 49 | Chat |
| Dify | `dify/` | 37 | Chat |
| Jina | `jina/` | 38 | Rerank |
| SiliconFlow | `siliconflow/` | 40 | Chat |
| Volcengine | `volcengine/` | 45 | Chat |
| Replicate | `replicate/` | 56 | Image |
| MokaAI | `mokaai/` | 44 | Chat |
| Codex | `codex/` | 57 | OAuth-based key management |

**Task providers** (async image/video/audio generation):
| Provider | Directory | Channel Type |
|----------|-----------|-------------|
| Midjourney | (in controller/) | 2 |
| Suno | `task/suno/` | 36 |
| Kling | `task/kling/` | 50 |
| Jimeng | `task/jimeng/` | 51 |
| Vidu | `task/vidu/` | 52 |
| Hailuo | `task/hailuo/` | -- |
| Doubao Video | `task/doubao/` | 54 |
| Sora | `task/sora/` | 55 |
| Ali (tasks) | `task/ali/` | 17 |
| Gemini (tasks) | `task/gemini/` | 24 |
| Vertex (tasks) | `task/vertex/` | 41 |

## Key Dependencies

- `dto/` -- Request/response structs for all API formats
- `model/` -- Channel selection, ability cache
- `service/` -- Billing, token counting
- `middleware/` -- Distribution, rate limiting
- `types/` -- Relay format constants, error types

## Data Models

- `relay/common/relay_info.go` -- `RelayInfo` struct: holds all context for a relay request (channel, user, token, model, format)
- `relay/common/billing.go` -- Usage billing helpers
- `relay/common/request_conversion.go` -- Request format conversion utilities

## Testing

- `relay/channel/api_request_test.go` -- API request building
- `relay/channel/aws/relay_aws_test.go` -- AWS Bedrock relay
- `relay/channel/claude/message_delta_usage_patch_test.go` -- Claude streaming usage patch
- `relay/channel/claude/relay_claude_test.go` -- Claude relay tests
- `relay/channel/gemini/relay_gemini_usage_test.go` -- Gemini usage parsing
- `relay/common/override_test.go` -- Parameter override tests
- `relay/common/relay_info_test.go` -- RelayInfo construction tests
- `relay/helper/stream_scanner_test.go` -- SSE stream scanner tests

## FAQ

**Q: How do I add a new provider?**
A: Create `relay/channel/<provider>/adaptor.go` implementing the `Adaptor` interface, add channel type to `constant/channel.go`, register in the channel type switch.

**Q: How does streaming work?**
A: Each adaptor's `DoResponse` method handles streaming by writing SSE events to `gin.Context.Writer`. The `relay/helper/stream_scanner.go` provides a reusable scanner for OpenAI-format SSE streams.

**Q: How are channels selected?**
A: The `middleware/distributor.go` middleware calls `service/channel_select.go` which ranks channels by priority, weight, model availability, and channel affinity.

## Related Files

- `relay/audio_handler.go` -- Audio relay handler
- `relay/channel/adapter.go` -- Adaptor interface definition
- `relay/channel/api_request.go` -- Common API request building
- `relay/common/relay_info.go` -- RelayInfo context struct
- `relay/common/billing.go` -- Billing utilities
- `relay/helper/common.go` -- Common relay helpers
- `relay/helper/price.go` -- Price calculation
- `relay/helper/model_mapped.go` -- Model mapping
