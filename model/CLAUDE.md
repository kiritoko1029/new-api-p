[Root](../CLAUDE.md) > **model**

# model -- Data Models and Database Access

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Defines all GORM data models, handles database initialization (SQLite/MySQL/PostgreSQL), runs auto-migrations, manages in-memory channel/token caches with periodic sync, and provides CRUD operations for all entities.

## Entry & Startup

- `model/main.go` -- `InitDB()` initializes the database connection, runs auto-migrations, creates root account
- `model/main.go` -- `InitLogDB()` initializes a separate log database (optional, via `LOG_SQL_DSN`)
- Called from `main.go` -> `InitResources()` -> `model.InitDB()`

## Public Interfaces

### Core Models

| Model | File | Description |
|-------|------|-------------|
| `User` | `user.go` | User accounts with role, quota, status, 2FA, passkey support |
| `Channel` | `channel.go` | AI provider channels with key, models, group, priority, weight |
| `Token` | `token.go` | API tokens for authentication (user-scoped) |
| `Log` | `log.go` | Request logs with usage, billing, model info |
| `Pricing` | `pricing.go` | Model pricing (input/output token rates) |
| `ModelMeta` | `model_meta.go` | Extended model metadata |
| `Redemption` | `redemption.go` | Redemption codes for quota top-up |
| `Subscription` | `subscription.go` | Subscription plans and user subscriptions |
| `TopUp` | `topup.go` | Top-up transaction records |
| `Checkin` | `checkin.go` | Daily check-in records |
| `Task` | `task.go` | Async task records (Midjourney, Suno, etc.) |
| `Midjourney` | `midjourney.go` | Midjourney-specific task records |
| `Option` | `option.go` | System options (key-value settings) |
| `VendorMeta` | `vendor_meta.go` | Vendor metadata |
| `Passkey` | `passkey.go` | WebAuthn passkey credentials |
| `CustomOAuthProvider` | `custom_oauth_provider.go` | Custom OAuth provider definitions |
| `UserOAuthBinding` | `user_oauth_binding.go` | User-OAuth provider bindings |
| `PrefillGroup` | `prefill_group.go` | Prefill group definitions for channels |
| `MissingModels` | `missing_models.go` | Tracked missing models from upstream |

### Cache Layer

- `channel_cache.go` -- In-memory channel cache with `InitChannelCache()`, `SyncChannelCache()`
- `token_cache.go` -- In-memory token cache
- `user_cache.go` -- In-memory user cache
- Sync frequency controlled by `SYNC_FREQUENCY` env var (default: 60s)

### DB Compatibility Helpers

- `main.go` exports: `commonGroupCol`, `commonKeyCol`, `commonTrueVal`, `commonFalseVal`, `logKeyCol`, `logGroupCol`
- `common/database.go` exports: `UsingPostgreSQL`, `UsingSQLite`, `UsingMySQL`, `DatabaseTypePostgreSQL`, etc.

## Key Dependencies

- `gorm.io/gorm` -- ORM
- `gorm.io/driver/mysql`, `gorm.io/driver/postgres`, `github.com/glebarez/sqlite` -- DB drivers
- `common/` -- Database type detection, JSON helpers, constants
- `dto/` -- Channel settings, pricing DTOs

## Data Model Notes

- The `Channel` model stores multiple API keys as comma-separated in the `Key` field; parsed into `Keys []string` at cache time
- The `Models` field is a comma-separated string of model names; the `Group` field is also comma-separated
- JSON fields (`ChannelInfo`, `Setting`, etc.) are stored as TEXT for cross-DB compatibility
- The `Log` table can be stored in a separate database (configured via `LOG_SQL_DSN`)

## Testing

- `model/task_cas_test.go` -- Task CAS (compare-and-swap) operations

## FAQ

**Q: How do I add a new model/table?**
A: Add a struct in a new file under `model/`, then add `DB.AutoMigrate(&NewModel{})` in `model/main.go`'s `InitDB()` function.

**Q: How does the cache sync work?**
A: `model.SyncChannelCache()` runs in a goroutine, periodically fetching all channels from DB and rebuilding the in-memory cache. Similarly for tokens and users. The sync frequency is controlled by `SYNC_FREQUENCY` (default 60s).

**Q: Why are JSON fields stored as TEXT?**
A: For cross-database compatibility. PostgreSQL has native JSONB, but SQLite and MySQL handle JSON differently. TEXT is universally supported.

## Related Files

- `model/main.go` -- DB init, migrations, column name helpers
- `model/channel.go` -- Channel model
- `model/user.go` -- User model
- `model/token.go` -- Token model
- `model/log.go` -- Log model
- `model/pricing.go` -- Pricing model
- `model/ability.go` -- Channel ability/cache management
- `model/option.go` -- System options
- `model/utils.go` -- DB utility functions
- `model/usedata.go` -- Usage data aggregation
