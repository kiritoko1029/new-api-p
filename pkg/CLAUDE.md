[Root](../CLAUDE.md) > **pkg**

# pkg -- Internal Packages

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Contains self-contained internal packages that provide specific functionality to the rest of the codebase.

## Sub-Packages

### `pkg/cachex/` -- Hybrid Cache

A hybrid caching system that combines in-memory and Redis caching.

| File | Description |
|------|-------------|
| `hybrid_cache.go` | `HybridCache` -- two-level cache (memory + Redis) |
| `namespace.go` | Cache namespace isolation |
| `codec.go` | Serialization codec for cache values |

### `pkg/ionet/` -- io.net Integration

Client library for the io.net GPU cloud platform, used for model deployment management.

| File | Description |
|------|-------------|
| `client.go` | HTTP client for io.net API |
| `types.go` | io.net API types |
| `container.go` | Container management |
| `deployment.go` | Deployment operations |
| `hardware.go` | Hardware type queries |
| `jsonutil.go` | JSON utilities for io.net responses |

## Related Files

- `pkg/cachex/hybrid_cache.go` -- Hybrid cache implementation
- `pkg/ionet/client.go` -- io.net API client
