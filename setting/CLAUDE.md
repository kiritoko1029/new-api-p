[Root](../CLAUDE.md) > **setting**

# setting -- Configuration Management

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Manages all system configuration settings. Settings are stored in the database (`model.Option`) and cached in memory. Organized into logical sub-packages for different setting categories.

## Sub-Packages

| Sub-package | Files | Description |
|-------------|-------|-------------|
| `config/` | `config.go` | Base configuration types and helpers |
| `console_setting/` | `config.go`, `validation.go` | Console (frontend) settings with validation |
| `model_setting/` | `claude.go`, `gemini.go`, `global.go`, `grok.go`, `qwen.go` | Per-provider model settings (max tokens, reasoning effort, etc.) |
| `operation_setting/` | `operation_setting.go`, `general_setting.go`, `checkin_setting.go`, `monitor_setting.go`, `payment_setting.go`, `quota_setting.go`, `token_setting.go`, `channel_affinity_setting.go`, `status_code_ranges.go`, `tools.go` | Operational settings (checkin, monitoring, payment, quotas, etc.) |
| `ratio_setting/` | `model_ratio.go`, `group_ratio.go`, `cache_ratio.go`, `expose_ratio.go`, `compact_suffix.go` | Model and group pricing ratios |
| `system_setting/` | `system_setting_old.go`, `fetch_setting.go`, `discord.go`, `legal.go`, `oidc.go`, `passkey.go` | System-wide settings (OAuth, legal, etc.) |
| `performance_setting/` | `config.go` | Performance tuning configuration |
| `reasoning/` | `suffix.go` | Reasoning suffix configuration |
| Root files | `auto_group.go`, `chat.go`, `midjourney.go`, `sensitive.go`, `rate_limit.go`, `payment_stripe.go`, `payment_creem.go`, `payment_waffo.go`, `user_usable_group.go` | Top-level settings |

## Testing

- `setting/operation_setting/status_code_ranges_test.go` -- Status code range parsing tests

## FAQ

**Q: How are settings loaded?**
A: `model.InitOptionMap()` loads all options from DB into memory at startup. `model.SyncOptions()` periodically refreshes the cache. Individual settings are accessed via `model.GetOptionByKey()`.

## Related Files

- `setting/ratio_setting/model_ratio.go` -- Model pricing ratios
- `setting/operation_setting/operation_setting.go` -- Operation settings aggregation
- `setting/model_setting/global.go` -- Global model settings
