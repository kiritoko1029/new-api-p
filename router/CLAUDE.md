[Root](../CLAUDE.md) > **router**

# router -- HTTP Route Registration

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Registers all HTTP routes on the Gin engine. Organized into four route groups: API, Relay, Dashboard, and Web (static file serving).

## Entry & Startup

- `router/main.go` -- `SetRouter()` called from `main.go`; delegates to sub-routers
- Frontend is served from embedded `web/dist` (via `//go:embed`) or redirected to `FRONTEND_BASE_URL`

## Route Groups

| File | Function | Prefix | Auth | Description |
|------|----------|--------|------|-------------|
| `api-router.go` | `SetApiRouter()` | `/api` | Varies | Management API: users, channels, tokens, logs, settings, payments |
| `relay-router.go` | `SetRelayRouter()` | `/v1`, `/mj`, `/suno`, `/v1beta` | Token | AI API relay: chat, image, audio, embedding, rerank, tasks |
| `dashboard.go` | `SetDashboardRouter()` | `/dashboard`, `/v1/dashboard` | Token | Billing dashboard API (OpenAI-compatible) |
| `video-router.go` | `SetVideoRouter()` | `/v1/video` | Token | Video generation relay |
| `web-router.go` | `SetWebRouter()` | `/*` | None | Static file serving for SPA |

## Key Routes

### Relay Routes (AI API proxy)
- `POST /v1/chat/completions` -- Chat completions (OpenAI format)
- `POST /v1/completions` -- Text completions
- `POST /v1/responses` -- OpenAI Responses API
- `POST /v1/messages` -- Claude/Anthropic messages
- `POST /v1/models/*path` -- Gemini native format
- `POST /v1/images/generations` -- Image generation
- `POST /v1/embeddings` -- Text embeddings
- `POST /v1/audio/speech`, `/v1/audio/transcriptions`, `/v1/audio/translations` -- Audio
- `POST /v1/rerank` -- Rerank
- `GET /v1/realtime` -- WebSocket realtime
- `POST /mj/submit/*` -- Midjourney task submission
- `POST /suno/submit/*` -- Suno task submission

### API Routes (Management)
- `POST /api/user/register`, `/api/user/login` -- Auth
- `GET /api/channel/`, `POST /api/channel/` -- Channel CRUD (admin)
- `GET /api/token/`, `POST /api/token/` -- Token CRUD (user)
- `GET /api/log/` -- Log viewing (admin/user)
- `GET/PUT /api/option/` -- System settings (root)
- `GET /api/pricing` -- Public pricing

## Related Files

- `router/api-router.go` -- Management API routes
- `router/relay-router.go` -- AI API relay routes
- `router/dashboard.go` -- Dashboard billing routes
- `router/web-router.go` -- Static file serving
