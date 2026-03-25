[Root](../CLAUDE.md) > **electron**

# electron -- Desktop App Wrapper

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Provides an Electron-based desktop application that wraps the Go backend binary and the web frontend. Enables running new-api as a native desktop application on macOS, Windows, and Linux.

## Key Files

| File | Description |
|------|-------------|
| `main.js` | Electron main process -- creates window, loads web UI |
| `preload.js` | Preload script for secure IPC |
| `create-tray-icon.js` | System tray icon generation |
| `package.json` | Electron dependencies and build configuration |

## Build Commands

```bash
cd electron
npm install
npm run dev-app       # Development mode
npm run build:mac     # Build macOS DMG + ZIP
npm run build:win     # Build Windows NSIS + portable
npm run build:linux   # Build Linux AppImage + deb
```

## Architecture

The Electron app:
1. Starts the Go binary (`new-api` or `new-api.exe`) as a child process
2. Opens a BrowserWindow pointing to `http://localhost:3000`
3. Provides system tray integration for background running
4. Packages the Go binary + web dist as extra resources

## Dependencies

- `electron` v35.7.5
- `electron-builder` v26.7.0
- `cross-env` -- Cross-platform environment variables

## Related Files

- `electron/main.js` -- Main process
- `electron/package.json` -- Build configuration
