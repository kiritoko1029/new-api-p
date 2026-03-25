[Root](../CLAUDE.md) > **oauth**

# oauth -- OAuth Provider Implementations

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Implements OAuth2 authentication flows for various identity providers. Providers are registered via `init()` functions and discovered through a central registry.

## Key Files

| File | Description |
|------|-------------|
| `provider.go` | `OAuthProvider` interface definition |
| `types.go` | Shared OAuth types (user info, tokens) |
| `registry.go` | Provider registry -- registers and looks up providers by name |
| `github.go` | GitHub OAuth |
| `discord.go` | Discord OAuth |
| `oidc.go` | Generic OIDC OAuth |
| `linuxdo.go` | LinuxDO OAuth |
| `generic.go` | Generic/custom OAuth (for dynamically configured providers) |

## How It Works

1. Each provider file has an `init()` function that calls `registry.Register()` with the provider name and implementation
2. `oauth.LoadCustomProviders()` loads additional providers from the database (stored in `model.CustomOAuthProvider`)
3. The router imports `_ "github.com/QuantumNous/new-api/oauth"` to trigger init()
4. `controller/oauth.go` handles the OAuth callback flow

## Related Files

- `oauth/registry.go` -- Provider registry
- `oauth/provider.go` -- Provider interface
