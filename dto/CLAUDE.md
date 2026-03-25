[Root](../CLAUDE.md) > **dto**

# dto -- Data Transfer Objects

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Defines request and response structs for all API formats. These DTOs are used by the relay layer to convert between different provider formats and by controllers for request/response serialization.

## Key DTOs

| File | Description |
|------|-------------|
| `openai_request.go` | `GeneralOpenAIRequest` -- Unified OpenAI-format request (chat, completions) |
| `openai_response.go` | `OpenAIResponse`, `OpenAITokenUsage` -- OpenAI-format response |
| `openai_image.go` | `ImageRequest`, `ImageResponse` -- Image generation |
| `openai_video.go` | Video generation request/response |
| `audio.go` | `AudioRequest` -- Audio TTS/STT request |
| `embedding.go` | `EmbeddingRequest`, `EmbeddingResponse` -- Text embeddings |
| `rerank.go` | `RerankRequest`, `RerankResponse` -- Rerank API |
| `claude.go` | `ClaudeRequest` -- Anthropic Claude native format |
| `gemini.go` | `GeminiChatRequest` -- Google Gemini native format |
| `realtime.go` | Realtime WebSocket DTOs |
| `midjourney.go` | Midjourney task request/response |
| `suno.go` | Suno task request/response |
| `task.go` | `TaskError`, generic task DTOs |
| `error.go` | Error response DTOs |
| `notify.go` | Notification DTOs |
| `pricing.go` | Pricing display DTOs |
| `sensitive.go` | Content sensitivity DTOs |
| `channel_settings.go` | Channel settings (JSON struct) |
| `user_settings.go` | User preference settings |
| `playground.go` | Playground request DTOs |
| `request_common.go` | Shared request fields |
| `openai_compaction.go` | Responses API compaction DTOs |
| `openai_responses_compaction_request.go` | Responses compaction request |
| `values.go` | Shared value types |
| `ratio_sync.go` | Ratio sync DTOs |

## Testing

- `dto/gemini_generation_config_test.go` -- Gemini config parsing
- `dto/openai_request_zero_value_test.go` -- Zero value preservation tests (Rule 6)

## FAQ

**Q: Why do some fields use pointer types?**
A: See Rule 6 in root CLAUDE.md. Optional scalar fields use `*int`, `*bool`, etc. with `omitempty` to distinguish "absent" from "zero value". This prevents silently dropping explicit zero values when re-marshaling to upstream providers.

## Related Files

- `dto/openai_request.go` -- Most-used request struct
- `dto/openai_response.go` -- Standard response format
