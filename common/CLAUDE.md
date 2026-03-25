[Root](../CLAUDE.md) > **common**

# common -- Shared Utilities

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Provides shared utility functions used across the entire backend. This is the "kitchen sink" package for cross-cutting concerns that do not belong in any specific domain module.

## Key Files

| File | Description |
|------|-------------|
| `json.go` | **Critical** -- JSON marshal/unmarshal wrappers (Rule 1). All business code must use these instead of `encoding/json` directly. |
| `constants.go` | Shared constants (user roles, status codes, etc.) |
| `env.go` | Environment variable loading and defaults |
| `init.go` | Package initialization |
| `database.go` | Database type detection (`UsingPostgreSQL`, `UsingSQLite`, `UsingMySQL`) |
| `redis.go` | Redis client initialization and helpers |
| `rate-limit.go` | Rate limiter implementation |
| `crypto.go` | Password hashing, random string generation |
| `email.go` | Email sending (SMTP) |
| `email-outlook-auth.go` | Outlook OAuth email authentication |
| `ip.go` | IP address utilities |
| `logger.go` | (moved to `logger/` package) |
| `sys_log.go` | System logging helpers |
| `system_monitor.go` | System resource monitoring (CPU, memory, goroutines) |
| `quota.go` | Quota calculation helpers |
| `utils.go` | General utility functions |
| `str.go` | String manipulation utilities |
| `hash.go` | Hashing utilities |
| `validate.go` | Input validation |
| `verification.go` | Email/code verification |
| `totp.go` | TOTP 2FA implementation |
| `audio.go` | Audio format detection and processing |
| `body_storage.go` | Request body storage for re-reading |
| `copy.go` | Deep copy utilities |
| `gin.go` | Gin-specific helpers |
| `go-channel.go` | Go channel utilities |
| `gopool.go` | Goroutine pool wrapper |
| `page_info.go` | Pagination helpers |
| `model.go` | Model-related shared functions |
| `api_type.go` | API type constants and mappings |
| `endpoint_type.go` | Endpoint type constants |
| `endpoint_defaults.go` | Default endpoint configurations |
| `custom-event.go` | Custom event types |
| `embed-file-system.go` | Embedded file system helpers |
| `disk_cache.go`, `disk_cache_config.go` | Disk-based caching |
| `performance_config.go` | Performance tuning configuration |
| `pprof.go` | Pprof profiling helpers |
| `pyro.go` | Pyroscope profiling integration |
| `ssrf_protection.go` | SSRF protection for outbound requests |
| `url_validator.go` | URL validation |
| `system_monitor_unix.go`, `system_monitor_windows.go` | Platform-specific monitoring |
| `topup-ratio.go` | Top-up ratio configuration |
| `limiter/limiter.go` | Rate limiter implementation |

## Testing

- `common/url_validator_test.go` -- URL validation tests

## FAQ

**Q: Why must I use common.Marshal instead of json.Marshal?**
A: See Rule 1 in root CLAUDE.md. The wrappers allow future swapping of the JSON library (e.g., to sonic or jsoniter) without changing business code.

**Q: How does database detection work?**
A: `common.InitEnv()` reads the `SQL_DSN` environment variable and sets `UsingPostgreSQL`, `UsingSQLite`, `UsingMySQL` flags accordingly. These flags are used throughout the codebase for database-specific logic.

## Related Files

- `common/json.go` -- JSON wrappers (most referenced)
- `common/env.go` -- Environment variable loading
- `common/database.go` -- Database type detection
- `common/redis.go` -- Redis client
