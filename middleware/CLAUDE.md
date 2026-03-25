[Root](../CLAUDE.md) > **middleware**

# middleware -- Request Pipeline

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Gin middleware that forms the request processing pipeline. Handles authentication, authorization, rate limiting, channel distribution, CORS, compression, logging, and request tracking.

## Key Middleware

| File | Middleware | Description |
|------|-----------|-------------|
| `auth.go` | `UserAuth()`, `AdminAuth()`, `RootAuth()`, `TokenAuth()`, `TokenAuthReadOnly()`, `TryUserAuth()` | JWT/session authentication at different authorization levels |
| `distributor.go` | `Distribute()` | **Core middleware** -- selects the best channel for relay requests based on priority, weight, affinity |
| `rate-limit.go` | `GlobalAPIRateLimit()`, `CriticalRateLimit()`, `SearchRateLimit()` | IP-based and user-based rate limiting |
| `model-rate-limit.go` | `ModelRequestRateLimit()` | Per-model rate limiting |
| `cors.go` | `CORS()` | Cross-origin resource sharing |
| `gzip.go` | (via gin-contrib/gzip) | Response compression |
| `logger.go` | `SetUpLogger()` | Request logging |
| `recover.go` | (via gin.CustomRecovery) | Panic recovery |
| `request-id.go` | `RequestId()` | Adds unique request ID |
| `i18n.go` | `I18n()` | Sets language from Accept-Language header or user preference |
| `cache.go` | `DisableCache()` | Disables caching for sensitive endpoints |
| `turnstile-check.go` | `TurnstileCheck()` | Cloudflare Turnstile CAPTCHA verification |
| `stats.go` | `StatsMiddleware()` | Request statistics collection |
| `performance.go` | `SystemPerformanceCheck()` | Blocks requests when system is overloaded |
| `body_cleanup.go` | `BodyStorageCleanup()` | Cleans up stored request bodies |
| `email-verification-rate-limit.go` | `EmailVerificationRateLimit()` | Rate limiting for email verification |
| `secure_verification.go` | `SecureVerificationRequired()` | Requires secure verification for sensitive operations |
| `jimeng_adapter.go` | Jimeng-specific middleware | Request adaptation for Jimeng API |
| `kling_adapter.go` | Kling-specific middleware | Request adaptation for Kling API |
| `utils.go` | `RouteTag()` | Tags routes for logging/metrics |

## Execution Order (Relay Requests)

1. `RequestId()` -- Assign request ID
2. `PoweredBy()` -- Set X-Powered-By header
3. `I18n()` -- Set language
4. Logger setup
5. Session store
6. `CORS()` -- Handle CORS
7. `DecompressRequestMiddleware()` -- Decompress request body
8. `BodyStorageCleanup()` -- Schedule body cleanup
9. `StatsMiddleware()` -- Collect stats
10. `TokenAuth()` -- Authenticate via API key
11. `ModelRequestRateLimit()` -- Check model-specific rate limit
12. `Distribute()` -- **Select channel and set context**

## FAQ

**Q: How does the Distribute middleware work?**
A: It parses the request to determine the model, then calls `service/channel_select.go` to find the best channel. It sets channel info in the Gin context for the relay handler to use.

**Q: How are auth levels different?**
A: `UserAuth()` requires any logged-in user. `AdminAuth()` requires admin role. `RootAuth()` requires root (super admin) role. `TokenAuth()` authenticates via Bearer token (API key). `TryUserAuth()` sets user context if token is valid but does not reject unauthenticated requests.

## Related Files

- `middleware/distributor.go` -- Channel distribution (core)
- `middleware/auth.go` -- Authentication
- `middleware/rate-limit.go` -- Rate limiting
