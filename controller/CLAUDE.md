[Root](../CLAUDE.md) > **controller**

# controller -- Request Handlers

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

HTTP request handlers for the Gin router. Controllers receive requests, validate input, call service/model layer, and return JSON responses. This layer should be thin -- business logic belongs in `service/`.

## Entry & Startup

Controllers are registered in the router packages:
- `router/api-router.go` -- API routes (users, channels, tokens, logs, settings, etc.)
- `router/relay-router.go` -- Relay routes (AI API proxying)
- `router/dashboard.go` -- Dashboard billing routes
- `router/video-router.go` -- Video relay routes

## Key Controllers

| File | Description |
|------|-------------|
| `relay.go` | Main relay handler; delegates to `relay/` engine |
| `channel.go` | Channel CRUD, test, balance update |
| `channel-test.go` | Channel connectivity testing |
| `channel_upstream_update.go` | Upstream model detection/updates |
| `user.go` (implicit in `router/api-router.go`) | User CRUD (via userRoute group) |
| `token.go` | Token CRUD and key retrieval |
| `log.go` | Log querying and statistics |
| `billing.go` | Subscription and usage billing |
| `option.go` | System option get/update |
| `oauth.go` | OAuth authentication flow |
| `passkey.go` | WebAuthn passkey operations |
| `midjourney.go` | Midjourney task relay |
| `task.go` | Generic async task relay |
| `model.go` | Model listing |
| `model_meta.go` | Model metadata CRUD |
| `model_sync.go` | Upstream model sync |
| `pricing.go` | Public pricing display |
| `subscription.go` | Subscription plan management |
| `subscription_payment_stripe.go` | Stripe payment integration |
| `subscription_payment_epay.go` | Epay payment integration |
| `subscription_payment_creem.go` | Creem payment integration |
| `playground.go` | AI model playground |
| `deployment.go` | io.net model deployment management |
| `telegram.go` | Telegram OAuth |
| `setup.go` | Initial setup wizard |
| `checkin.go` | Daily check-in |
| `redemption.go` | Redemption code management |
| `performance.go` | Performance monitoring |
| `ratio_config.go` | Ratio configuration |
| `ratio_sync.go` | Upstream ratio sync |

## Testing

- `controller/channel_upstream_update_test.go` -- Channel upstream model update tests
- `controller/token_test.go` -- Token operations tests

## FAQ

**Q: How do I add a new API endpoint?**
A: 1) Create handler function in a controller file. 2) Register route in `router/api-router.go` (or appropriate router file). 3) Add middleware (auth, rate limit) as needed.

**Q: How does authentication work in controllers?**
A: Middleware (`middleware/auth.go`) sets user/token info in the Gin context. Controllers access it via `c.GetInt("id")` (user ID), `c.GetString("token_name")`, etc.

## Related Files

- `controller/relay.go` -- Main relay handler
- `controller/channel.go` -- Channel management
- `controller/option.go` -- System settings
- `controller/subscription.go` -- Subscription management
