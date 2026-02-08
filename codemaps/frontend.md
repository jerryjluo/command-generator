# Frontend Structure

> Last updated: 2026-02-01

## Technology Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| React | 19.2.0 | UI framework |
| React Router | 7.13.0 | Client-side routing |
| TypeScript | 5.9.3 | Type safety |
| Vite | 7.2.4 | Build tool & dev server |
| Tailwind CSS | 4.1.18 | Styling |
| ESLint | 9.39.1 | Code quality |

## Directory Structure

```
web/
├── package.json              # Dependencies & scripts
├── tsconfig.json             # TypeScript config
├── vite.config.ts            # Vite bundler config
├── tailwind.config.js        # Tailwind CSS config
├── eslintrc.mjs              # ESLint config
├── index.html                # HTML template
├── public/                   # Static assets
├── dist/                     # Built output (embedded in Go)
└── src/
    ├── main.tsx              # React entry point
    ├── App.tsx               # Root component with routing
    ├── index.css             # Global styles (Tailwind)
    ├── api/
    │   └── logs.ts           # API client functions
    ├── types/
    │   └── log.ts            # TypeScript interfaces
    ├── hooks/
    │   ├── useLogs.ts        # Fetch logs list
    │   ├── useLogDetail.ts   # Fetch individual log
    │   └── useFilters.ts     # Filter state management
    ├── components/
    │   ├── Layout.tsx        # Main layout wrapper
    │   ├── LogTable/         # Log list view
    │   ├── LogDetail/        # Detail view
    │   └── common/           # Shared components
    └── assets/               # Static assets
```

---

## Component Tree

```
BrowserRouter
└── App
    └── Layout
        └── Routes
            ├── "/" → LogListPage
            │   └── LogTable
            │       ├── Filters
            │       │   ├── StatusFilter
            │       │   ├── ModelFilter
            │       │   ├── SearchInput
            │       │   └── DateRangeFilter
            │       └── SortHeader
            │
            └── "/logs/:id" → LogDetailPage
                └── LogDetail
                    └── TabPanel
                        ├── ResponseTab
                        ├── SystemPromptTab
                        ├── UserPromptTab
                        ├── UserQueryTab
                        ├── BuildToolsTab
                        ├── TmuxContextTab
                        ├── DocumentationContextTab
                        └── PreferencesTab
```

---

## Route Configuration

| Path | Component | Description |
|------|-----------|-------------|
| `/` | `LogListPage` | Paginated log table with filters |
| `/logs/:id` | `LogDetailPage` | Full session log details |

---

## Components

### `components/Layout.tsx`

Main layout wrapper with header/navigation.

### `components/LogTable/`

| Component | Purpose |
|-----------|---------|
| `LogTable.tsx` | Table display with rows |
| `Filters.tsx` | Filter controls container |
| `SortHeader.tsx` | Sortable column headers |
| `index.ts` | Barrel export |

### `components/LogDetail/`

| Component | Purpose |
|-----------|---------|
| `LogDetail.tsx` | Detail page wrapper |
| `TabPanel.tsx` | Tab navigation UI |
| `index.ts` | Barrel export |

### `components/LogDetail/tabs/`

| Tab | Content |
|-----|---------|
| `ResponseTab` | Generated command + explanation |
| `SystemPromptTab` | Full system prompt sent to Claude |
| `UserPromptTab` | User prompt with context |
| `UserQueryTab` | Original natural language query |
| `BuildToolsTab` | Detected build commands |
| `TmuxContextTab` | Captured terminal scrollback |
| `DocumentationContextTab` | Project documentation |
| `PreferencesTab` | User's claude.md content |

### `components/common/`

| Component | Purpose |
|-----------|---------|
| `Badge.tsx` | Status badge (accepted/rejected/quit) |
| `TimeAgo.tsx` | Relative timestamp display |
| `CodeBlock.tsx` | Syntax-highlighted code |

---

## Custom Hooks

### `useLogs.ts`

```typescript
function useLogs(params: FilterParams): {
    logs: LogSummary[];
    total: number;
    loading: boolean;
    error: Error | null;
    refetch: () => void;
}
```

Fetches paginated log list based on filter parameters.

### `useLogDetail.ts`

```typescript
function useLogDetail(id: string): {
    log: SessionLog | null;
    loading: boolean;
    error: Error | null;
}
```

Fetches full session log by ID.

### `useFilters.ts`

```typescript
function useFilters(): {
    filters: FilterParams;
    setStatus: (status: string) => void;
    setModel: (model: string) => void;
    setSearch: (search: string) => void;
    setDateRange: (from: string, to: string) => void;
    setSort: (field: string, order: 'asc' | 'desc') => void;
    setPagination: (limit: number, offset: number) => void;
    resetFilters: () => void;
}
```

Manages filter state and syncs with URL search params.

---

## API Client (`api/logs.ts`)

```typescript
const BASE_URL = '/api/v1';

// Fetch paginated logs with filters
async function fetchLogs(params: FilterParams): Promise<LogListResponse> {
    const query = new URLSearchParams(params);
    const response = await fetch(`${BASE_URL}/logs?${query}`);
    return response.json();
}

// Fetch single log by ID
async function fetchLogById(id: string): Promise<SessionLog> {
    const response = await fetch(`${BASE_URL}/logs/${id}`);
    return response.json();
}
```

---

## Type Definitions (`types/log.ts`)

```typescript
// Core interfaces
interface SessionLog { id, user_query, context_sources, iterations, metadata }
interface LogSummary { id, user_query, final_status, model, timestamp, ... }
interface FilterParams { status?, model?, search?, from?, to?, sort?, ... }
interface LogListResponse { logs, total, limit, offset }

// Nested types
interface Iteration { feedback, model_input, model_output, timestamp }
interface ContextSources { claude_md_content, terminal_context, ... }
interface Metadata { timestamp, model, final_status, iteration_count, ... }
interface TmuxInfo { in_tmux, session?, window?, pane? }
```

---

## Build Configuration

### `vite.config.ts`

```typescript
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: 'dist',
        // Output embedded in Go binary
    }
});
```

### `tsconfig.json`

```json
{
    "compilerOptions": {
        "target": "ES2020",
        "module": "ESNext",
        "strict": true,
        "jsx": "react-jsx"
    }
}
```

### `tailwind.config.js`

```javascript
module.exports = {
    content: ['./src/**/*.{ts,tsx}'],
    theme: { extend: {} },
    plugins: []
};
```

---

## npm Scripts

| Script | Command | Purpose |
|--------|---------|---------|
| `dev` | `vite` | Hot reload dev server |
| `build` | `tsc -b && vite build` | Production build |
| `lint` | `eslint .` | Code linting |
| `preview` | `vite preview` | Preview built output |

---

## Dependencies

### Production

```json
{
    "react": "19.2.0",
    "react-dom": "19.2.0",
    "react-router-dom": "7.13.0"
}
```

### Development

```json
{
    "@types/react": "...",
    "@types/react-dom": "...",
    "@vitejs/plugin-react": "...",
    "autoprefixer": "10.4.24",
    "eslint": "9.39.1",
    "postcss": "...",
    "tailwindcss": "4.1.18",
    "typescript": "5.9.3",
    "vite": "7.2.4"
}
```

---

## Deployment

The frontend is embedded into the Go binary:

1. `npm run build` compiles to `web/dist/`
2. Go build embeds `web/dist/` via `//go:embed`
3. Single binary serves both API and SPA
4. SPA fallback routes all paths to `index.html`
