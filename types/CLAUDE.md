[Root](../CLAUDE.md) > **types**

# types -- Type Definitions

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Defines shared type definitions used across the codebase, particularly relay format types, error types, and data structures that do not belong in `dto/` or `model/`.

## Key Types

| File | Description |
|------|-------------|
| `relay_format.go` | `RelayFormat` constants -- identifies the API format for relay (OpenAI, Claude, Gemini, etc.) |
| `error.go` | `NewAPIError` -- unified error type with status code, message, and type |
| `channel_error.go` | Channel-specific error types |
| `file_source.go` | File source type constants |
| `file_data.go` | File data structures |
| `price_data.go` | Price data structures |
| `request_meta.go` | Request metadata types |
| `rw_map.go` | Thread-safe read-write map |
| `set.go` | Set data structure |

## Related Files

- `types/relay_format.go` -- Relay format constants
- `types/error.go` -- Error types
