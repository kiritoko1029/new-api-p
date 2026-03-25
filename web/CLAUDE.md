[Root](../CLAUDE.md) > **web**

# web -- React Frontend (Admin Dashboard)

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Provides the complete admin dashboard and user-facing UI for new-api. Includes user authentication, channel management, token management, billing/pricing, playground for testing AI models, subscription management, log viewing, settings configuration, and a responsive layout with dark/light theme support.

## Entry & Startup

- `web/src/App.jsx` -- Root component; defines all routes
- `web/index.html` -- HTML entry point
- `web/vite.config.js` -- Vite configuration with proxy to backend
- `web/src/components/layout/PageLayout.jsx` -- Main layout wrapper (sidebar, headerbar, content)

Dev: `bun run dev` (starts Vite dev server with proxy to localhost:3000)
Build: `bun run build` (outputs to `web/dist/`)

The built `web/dist/` is embedded into the Go binary via `//go:embed web/dist` in `main.go`.

## Public Interfaces (Routes/Pages)

| Page | Path | Description |
|------|------|-------------|
| Home/Dashboard | `/` | User dashboard with stats, charts, announcements |
| Channel | `/channel` | Channel management (CRUD, test, tag) |
| Token | `/token` | API token management |
| User | `/user` | User management (admin) |
| Log | `/log` | Request log viewer |
| Setting | `/setting` | System settings (multi-tab: Operation, Ratio, Model, Drawing, Chat, Dashboard, Payment, Rate Limit, Performance, System) |
| Pricing | `/pricing` | Public model pricing display |
| TopUp | `/topup` | Quota top-up page |
| Redemption | `/redemption` | Redemption code management (admin) |
| Playground | `/playground` | AI model testing playground |
| Midjourney | `/midjourney` | Midjourney task logs |
| Task | `/task` | Async task logs |
| Model | `/model` | Model metadata management |
| ModelDeployment | `/model-deployment` | io.net model deployment management |
| Subscription | `/subscription` | Subscription plan management |
| About | `/about` | About page |
| Privacy | `/privacy-policy` | Privacy policy |
| UserAgreement | `/user-agreement` | User agreement |
| Setup | `/setup` | Initial setup wizard |
| Chat | `/chat` | AI chat interface |
| Chat2Link | `/chat2link` | Chat link generator |
| PricingPage | `/pricing` | Public pricing display |

## Key Dependencies

- `@douyinfe/semi-ui` -- Primary UI component library
- `react-router-dom` -- Client-side routing
- `axios` -- HTTP client (via `helpers/api.js`)
- `i18next` + `react-i18next` -- Internationalization
- `@visactor/react-vchart` + `@visactor/vchart` -- Charts
- `marked` + `react-markdown` -- Markdown rendering
- `mermaid` -- Mermaid diagram rendering
- `qrcode.react` -- QR code generation
- `sse.js` -- Server-Sent Events client (for playground streaming)
- `lucide-react` + `react-icons` -- Icon libraries

## Component Architecture

```
web/src/
  App.jsx                    -- Route definitions
  components/
    auth/                    -- Login, Register, OAuth callback, 2FA, Password reset
    common/                  -- Shared UI components (CardPro, Loading, JSONEditor, modals)
    dashboard/               -- Dashboard panels (Stats, Charts, Announcements, FAQ)
    layout/                  -- Page layout (Sidebar, HeaderBar, Footer, NoticeModal)
    playground/              -- AI playground (Chat, Code viewer, SSE viewer, Settings)
    settings/                -- Settings panels (System, Operation, Ratio, Model, Payment, etc.)
    setup/                   -- Setup wizard (multi-step)
    table/                   -- Data tables with filters, actions, modals:
      channels/              -- Channel management table
      models/                -- Model metadata table
      redemptions/           -- Redemption code table
      subscriptions/         -- Subscription table
      task-logs/             -- Task log table
      mj-logs/               -- Midjourney log table
      model-deployments/     -- Deployment table
      model-pricing/         -- Pricing display (card + table views)
      users/                 -- User table
      tokens/                -- Token table
      logs/                  -- Log table
  pages/                     -- Page components (route targets)
    Setting/                 -- Settings sub-pages (multi-level)
  helpers/                   -- API client, utilities (api.js, quota.js, render.js, etc.)
  hooks/                     -- Custom React hooks
  constants/                 -- Frontend constants
  context/                   -- React context providers (User, Status)
  i18n/                      -- i18n configuration and locale files
  services/                  -- Frontend service layer (secureVerification.js)
```

## Testing

- No dedicated test framework; ESLint and Prettier for code quality
- `bun run eslint` -- Run ESLint checks
- `bun run lint` -- Run Prettier checks

## FAQ

**Q: How do I add a new page?**
A: Create component in `web/src/pages/<PageName>/index.jsx`, add route in `App.jsx`, add sidebar entry in `components/layout/SiderBar.jsx`.

**Q: How do I add a new settings tab?**
A: Create component in `web/src/pages/Setting/<Category>/`, add tab entry in `pages/Setting/index.jsx`.

**Q: How does i18n work in the frontend?**
A: Use `useTranslation()` hook from `react-i18next`. Translation keys are Chinese strings. Locale files are in `web/src/i18n/locales/`. Run `bun run i18n:extract` to find missing keys.

**Q: How does the playground streaming work?**
A: The playground uses `sse.js` to connect to `/v1/chat/completions` with `stream: true`. The `ChatArea.jsx` component processes SSE events and renders them incrementally.

## Related Files

- `web/src/App.jsx` -- Route definitions
- `web/src/helpers/api.js` -- Axios API client
- `web/src/components/layout/SiderBar.jsx` -- Sidebar navigation
- `web/src/i18n/i18n.js` -- i18n configuration
- `web/vite.config.js` -- Vite config
- `web/tailwind.config.js` -- Tailwind CSS config
- `web/package.json` -- Dependencies and scripts
