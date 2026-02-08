# Frontend Structure

> Last updated: 2026-02-08

## Technology Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| React | 19.2.0 | UI framework |
| React Router | 7.13.0 | Client-side routing |
| TypeScript | 5.9.3 | Type safety |
| Vite | 7.2.4 | Build tool & dev server |
| Tailwind CSS | 4.1.18 | Utility-first styling |
| ESLint | 9.39.1 | Code quality |

## Directory Structure

```
web/
├── package.json              # Dependencies & scripts
├── tsconfig.json             # TypeScript config
├── vite.config.ts            # Vite bundler config
├── index.html                # HTML template
├── dist/                     # Built output (embedded in Go binary)
└── src/
    ├── main.tsx              # React entry point (BrowserRouter)
    ├── App.tsx               # Root component with routing
    ├── index.css             # Global styles (Tailwind)
    ├── api/
    │   └── logs.ts           # API client functions
    ├── types/
    │   └── log.ts            # TypeScript interfaces
    ├── hooks/
    │   ├── useLogs.ts        # Fetch logs list
    │   ├── useLogDetail.ts   # Fetch individual log
    │   └── useFilters.ts     # Filter state via URL search params
    └── components/
        ├── Layout.tsx        # Main layout wrapper (header + content)
        ├── common/           # Shared components
        │   ├── Badge.tsx     # Badge, StatusBadge, ModelBadge
        │   ├── TimeAgo.tsx   # Relative/absolute timestamp toggle
        │   └── CodeBlock.tsx # Code display with copy button
        ├── LogTable/         # Log list view
        │   ├── LogTable.tsx  # Table with sortable columns
        │   ├── Filters.tsx   # Search, status, model, date range filters
        │   ├── SortHeader.tsx # Clickable column header with sort indicator
        │   └── index.ts      # Barrel export
        └── LogDetail/        # Detail view
            ├── LogDetail.tsx # Detail page with header + tabs
            ├── TabPanel.tsx  # Generic tab navigation component
            ├── index.ts      # Barrel export
            └── tabs/
                ├── ResponseTab.tsx            # Command + explanation
                ├── SystemPromptTab.tsx        # Full system prompt
                ├── UserPromptTab.tsx          # User prompt with context
                ├── UserQueryTab.tsx           # Original query
                ├── TmuxContextTab.tsx         # Terminal scrollback + tmux info
                ├── DocumentationContextTab.tsx # Project docs context
                ├── BuildToolsTab.tsx          # Detected build commands
                └── PreferencesTab.tsx         # claude.md content
```

---

## Component Tree

```
BrowserRouter (main.tsx)
└── App
    └── Layout
        └── Routes
            ├── "/" → LogListPage
            │   └── LogTable
            │       ├── Filters
            │       │   ├── Search input (debounced, 300ms)
            │       │   ├── Status select (accepted/rejected/quit)
            │       │   ├── Model select (opus/sonnet/haiku)
            │       │   ├── Date range (from/to)
            │       │   └── Clear button
            │       ├── SortHeader (per column)
            │       └── Table rows → StatusBadge, ModelBadge, TimeAgo
            │
            └── "/logs/:id" → LogDetailPage
                └── LogDetail
                    ├── Header (query, time, iterations, StatusBadge, ModelBadge)
                    ├── Iteration History (if multiple iterations)
                    └── TabPanel
                        ├── ResponseTab
                        ├── SystemPromptTab
                        ├── UserPromptTab
                        ├── UserQueryTab
                        ├── TmuxContextTab
                        ├── DocumentationContextTab
                        ├── BuildToolsTab
                        └── PreferencesTab
```

---

## Route Configuration

| Path | Component | Description |
|------|-----------|-------------|
| `/` | `LogListPage` | Paginated log table with filters and sorting |
| `/logs/:id` | `LogDetailPage` | Full session log details with tabbed content |

---

## Components

### `components/Layout.tsx`

Main layout wrapper with header ("cmd Log Viewer") and content area. Max width 7xl.

### `components/LogTable/`

| Component | Exports | Purpose |
|-----------|---------|---------|
| `LogTable.tsx` | `LogTable` | Table display with clickable rows, loading/empty/error states |
| `Filters.tsx` | `Filters` | Search (debounced 300ms), status/model dropdowns, date range, clear |
| `SortHeader.tsx` | `SortHeader` | Sortable column header with ascending/descending indicator |
| `index.ts` | Barrel | Re-exports `LogTable`, `Filters`, `SortHeader` |

### `components/LogDetail/`

| Component | Exports | Purpose |
|-----------|---------|---------|
| `LogDetail.tsx` | `LogDetail` | Detail page: header with metadata, iteration history, tab panel |
| `TabPanel.tsx` | `TabPanel` | Generic tab navigation (accepts array of `{id, label, content}`) |
| `index.ts` | Barrel | Re-exports `LogDetail`, `TabPanel` |

### `components/LogDetail/tabs/`

| Tab | Props | Content |
|-----|-------|---------|
| `ResponseTab` | `command`, `explanation`, `rawResponse` | Generated command + explanation |
| `SystemPromptTab` | `content` | Full system prompt sent to Claude |
| `UserPromptTab` | `content` | User prompt with all context |
| `UserQueryTab` | `content` | Original natural language query |
| `TmuxContextTab` | `terminalContext`, `tmuxInfo` | Captured scrollback + session info |
| `DocumentationContextTab` | `documentationContext` | Project documentation sections |
| `BuildToolsTab` | `userPrompt` | Build tools context (extracted from user prompt) |
| `PreferencesTab` | `claudeMdContent` | User's claude.md preferences |

### `components/common/`

| Component | Exports | Purpose |
|-----------|---------|---------|
| `Badge.tsx` | `Badge`, `StatusBadge`, `ModelBadge` | Colored badges (success/error/warning/info variants) |
| `TimeAgo.tsx` | `TimeAgo` | Click-toggleable relative/absolute timestamp |
| `CodeBlock.tsx` | `CodeBlock` | Dark-themed code display with hover copy button |

---

## Custom Hooks

### `useLogs.ts`

```typescript
function useLogs(filters: FilterParams): {
    logs: LogSummary[];
    total: number;
    loading: boolean;
    error: string | null;
    refetch: () => void;
}
```

Fetches paginated log list. Re-fetches when `filters` change (compared via `JSON.stringify`).

### `useLogDetail.ts`

```typescript
function useLogDetail(id: string | null): {
    log: SessionLog | null;
    loading: boolean;
    error: string | null;
    refetch: () => void;
}
```

Fetches full session log by ID. No-ops when `id` is null.

### `useFilters.ts`

```typescript
function useFilters(): {
    filters: FilterParams;                                              // Current filter state from URL
    setFilter: (key: keyof FilterParams, value: string | number | undefined) => void;  // Set single filter
    clearFilters: () => void;                                           // Reset all filters
    setSort: (field: string) => void;                                   // Toggle sort field/direction
}
```

Manages filter state via URL search params (`useSearchParams`). Resets offset when non-pagination filters change. Sort toggles direction when clicking the same field.

---

## API Client (`api/logs.ts`)

```typescript
const API_BASE = '/api/v1';

// Fetch paginated logs with filters
async function fetchLogs(params: FilterParams): Promise<LogListResponse>

// Fetch single log by ID
async function fetchLogById(id: string): Promise<SessionLog>
```

Both functions build URL with search params, call `fetch()`, check `response.ok`, and throw on error with server error message.

---

## Type Definitions (`types/log.ts`)

All types mirror the Go backend JSON responses:

```typescript
// Core types
interface SessionLog { id, user_query, context_sources, iterations, metadata }
interface LogSummary { id, user_query, final_status, model, timestamp, iteration_count, command_preview, tmux_session? }
interface LogListResponse { logs, total, limit, offset }
interface FilterParams { status?, model?, search?, from?, to?, sort?, order?, limit?, offset? }

// Nested types
interface Iteration { feedback, model_input, model_output, timestamp }
interface ContextSources { claude_md_content, terminal_context, documentation_context }
interface Metadata { timestamp, model, final_status, final_feedback?, iteration_count, tmux_info }
interface ModelInput { system_prompt, user_prompt }
interface ModelOutput { raw_response, command, explanation }
interface TmuxInfo { in_tmux, session?, window?, pane? }
```

---

## npm Scripts

| Script | Command | Purpose |
|--------|---------|---------|
| `dev` | `vite` | Hot reload dev server |
| `build` | `tsc -b && vite build` | Type check + production build |
| `lint` | `eslint .` | Code linting |
| `preview` | `vite preview` | Preview built output |

---

## Dependencies

### Production

| Package | Version |
|---------|---------|
| `react` | ^19.2.0 |
| `react-dom` | ^19.2.0 |
| `react-router-dom` | ^7.13.0 |

### Development

| Package | Version | Purpose |
|---------|---------|---------|
| `typescript` | ~5.9.3 | Type checking |
| `vite` | ^7.2.4 | Build tool |
| `@vitejs/plugin-react` | ^5.1.1 | React support for Vite |
| `tailwindcss` | ^4.1.18 | CSS framework |
| `@tailwindcss/postcss` | ^4.1.18 | PostCSS integration |
| `postcss` | ^8.5.6 | CSS processing |
| `autoprefixer` | ^10.4.24 | CSS vendor prefixes |
| `eslint` | ^9.39.1 | Linting |
| `eslint-plugin-react-hooks` | ^7.0.1 | React hooks rules |
| `eslint-plugin-react-refresh` | ^0.4.24 | Fast refresh lint rules |
| `typescript-eslint` | ^8.46.4 | TypeScript ESLint support |
| `globals` | ^16.5.0 | Global variable definitions |
| `@types/react` | ^19.2.5 | React type definitions |
| `@types/react-dom` | ^19.2.3 | React DOM type definitions |
| `@types/node` | ^24.10.1 | Node.js type definitions |

---

## Deployment

The frontend is embedded into the Go binary:

1. `npm run build` compiles to `web/dist/`
2. Go build embeds `web/dist/` via `//go:embed`
3. Single binary serves both API and SPA
4. SPA fallback routes all non-file paths to `index.html`
