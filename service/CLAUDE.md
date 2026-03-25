[Root](../CLAUDE.md) > **service**

# service -- Business Logic

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Contains core business logic that sits between controllers and models. Handles billing calculations, token counting, channel selection, request conversion, payment processing, notification, file handling, and various background tasks.

## Key Services

| File | Description |
|------|-------------|
| `billing.go` | Core billing logic -- calculates quota based on model, tokens, and ratios |
| `billing_session.go` | Billing session management |
| `channel.go` | Channel management operations |
| `channel_select.go` | Channel selection algorithm (priority, weight, affinity) |
| `channel_affinity.go` | Channel affinity -- sticky routing to preferred channels |
| `quota.go` | Quota management and consumption |
| `text_quota.go` | Text-based quota calculation |
| `token_counter.go` | Token counting using tiktoken |
| `token_estimator.go` | Token estimation for pre-billing |
| `tokenizer.go` | Tokenizer initialization and management |
| `convert.go` | Request/response format conversion |
| `http.go`, `http_client.go` | HTTP client management (connection pooling, timeouts) |
| `error.go` | Error handling and classification |
| `sensitive.go` | Content sensitivity filtering |
| `epay.go` | Epay payment integration |
| `funding_source.go` | Funding source management |
| `task.go` | Async task management |
| `task_billing.go` | Task-based billing (pre-charge + settlement) |
| `task_polling.go` | Task polling adaptor interface and registry |
| `image.go` | Image processing utilities |
| `audio.go` | Audio processing utilities |
| `download.go` | File download handling |
| `file_service.go` | File storage service |
| `file_decoder.go` | Audio file format detection and decoding |
| `notify-limit.go` | Notification rate limiting |
| `user_notify.go` | User notification (email) |
| `violation_fee.go` | Violation fee processing |
| `webhook.go` | Webhook management |
| `log_info_generate.go` | Log info generation |
| `subscription_reset_task.go` | Subscription quota reset background task |
| `codex_oauth.go` | Codex OAuth flow |
| `codex_credential_refresh.go` | Codex credential auto-refresh |
| `codex_credential_refresh_task.go` | Codex credential refresh background task |
| `codex_wham_usage.go` | Codex usage tracking |
| `openaicompat/` | OpenAI Responses API <-> Chat Completions conversion |
| `passkey/` | WebAuthn passkey service |

## Testing

- `service/error_test.go` -- Error handling tests
- `service/task_billing_test.go` -- Task billing calculation tests
- `service/text_quota_test.go` -- Text quota calculation tests
- `service/channel_affinity_template_test.go` -- Channel affinity template tests
- `service/channel_affinity_usage_cache_test.go` -- Channel affinity cache tests

## FAQ

**Q: How does billing work?**
A: 1) Pre-charge estimated quota before sending request to upstream. 2) After response, calculate actual usage from token counts. 3) Settle the difference (refund excess or charge deficit). For streaming, billing happens incrementally.

**Q: How does channel selection work?**
A: `service/channel_select.go` ranks channels by: priority (descending), then weight (random within same priority), filtered by model availability, group access, and channel status. Channel affinity can override this to stick users to specific channels.

## Related Files

- `service/billing.go` -- Core billing
- `service/channel_select.go` -- Channel selection
- `service/token_counter.go` -- Token counting
