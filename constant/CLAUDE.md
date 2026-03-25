[Root](../CLAUDE.md) > **constant**

# constant -- Constants

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Defines all constant values used across the backend, including channel type identifiers, API type mappings, context keys, environment variable names, and task platform constants.

## Key Constants

| File | Description |
|------|-------------|
| `channel.go` | `ChannelType*` constants (1-58) mapping to each AI provider; `ChannelBaseURLs` default URLs |
| `api_type.go` | API type identifiers |
| `azure.go` | Azure-specific constants |
| `cache_key.go` | Redis/memory cache key templates |
| `context_key.go` | Gin context key strings |
| `endpoint_type.go` | Endpoint type constants |
| `env.go` | Environment variable name constants |
| `finish_reason.go` | LLM finish reason constants |
| `midjourney.go` | Midjourney action/type constants |
| `multi_key_mode.go` | Multi-key rotation mode constants |
| `setup.go` | Setup wizard constants |
| `task.go` | `TaskPlatform` constants for async task providers |
| `waffo_pay_method.go` | Waffo payment method constants |

## FAQ

**Q: How many channel types are there?**
A: Currently 58 (0-57), defined in `constant/channel.go`. The last entry `ChannelTypeDummy` is a sentinel for counting.

## Related Files

- `constant/channel.go` -- Channel type definitions (most referenced)
- `constant/task.go` -- Task platform types
