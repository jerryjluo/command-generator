# Architecture Overview

> Last updated: 2026-02-01

## System Overview

**command_generator** (`cmd`) is a CLI tool that generates shell commands from natural language using Claude AI. It combines a Go backend for CLI operations with a React frontend for log visualization.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User Interface                               │
├──────────────────────────────┬──────────────────────────────────────┤
│       CLI (main.go)          │      Web Viewer (React SPA)          │
│  - Natural language input    │  - Log browsing & filtering          │
│  - Interactive A/R/Q loop    │  - Session detail viewing            │
│  - Clipboard integration     │  - Served on localhost:8765          │
└──────────────┬───────────────┴──────────────────┬───────────────────┘
               │                                   │
               ▼                                   ▼
┌──────────────────────────────────────────────────────────────────────┐
│                          Go Backend                                   │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│   claude/   │ buildtools/ │  terminal/  │   logging/  │   server/   │
│  API client │  8 parsers  │ tmux context│  JSON logs  │  HTTP API   │
└──────┬──────┴──────┬──────┴──────┬──────┴──────┬──────┴──────┬──────┘
       │             │             │             │             │
       ▼             ▼             ▼             ▼             ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐
│ claude CLI  │ │ Build files │ │    tmux     │ │ ~/.local/share/cmd/ │
│  (external) │ │ in project  │ │   session   │ │       logs/         │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────────────┘
```

## Core Flows

### Command Generation Flow

```
User Input → Load Config → Capture Context → Call Claude → Display → Log
     │            │              │                │           │       │
     ▼            ▼              ▼                ▼           ▼       ▼
 "list files" ~/.config/cmd  tmux scrollback  claude CLI   A/R/Q   JSON
              /claude.md     + build tools    --json-schema  loop   file
                             + docs
```

### Log Viewer Flow

```
cmd --logs → Start Server → Embed React → Serve API → Display Logs
                  │              │            │
                  ▼              ▼            ▼
            port 8765     web/dist/*    /api/v1/logs
```

## Package Dependencies

```
main.go
├── internal/buildtools   # Detect build systems
├── internal/claude       # Claude API integration
├── internal/clipboard    # Cross-platform clipboard
├── internal/config       # User preferences
├── internal/docs         # Documentation detection
├── internal/logging      # Session logging
├── internal/server       # HTTP server + API
└── internal/terminal     # tmux context capture

server/handlers.go
└── internal/logging      # Read/search logs

logging/logging.go
└── internal/terminal     # Store tmux info
```

## Technology Stack

| Layer | Technology |
|-------|------------|
| CLI | Go 1.25.3 |
| Build Runner | mise |
| Frontend | React 19 + TypeScript + Vite |
| Styling | Tailwind CSS |
| API | REST (Go net/http) |
| Storage | JSON files |
| External | `claude` CLI tool |

## Key Configuration Paths

| Path | Purpose |
|------|---------|
| `~/.config/cmd/claude.md` | User preferences |
| `~/.local/share/cmd/logs/` | Session logs |
| `~/.local/bin/cmd` | Installed binary |

## Build Process

```
mise run build
       │
       ├── 1. web-build (npm run build)
       │      └── Compiles React → web/dist/
       │
       └── 2. go build
              └── Embeds web/dist/ via //go:embed
              └── Outputs single binary: ./cmd
```

## External Integrations

- **Claude API**: Via `claude` CLI subprocess with `--json-schema` for structured output
- **Clipboard**: Platform-specific (`pbcopy` on macOS, `xclip`/`xsel` on Linux)
- **Browser**: Platform-specific (`open` on macOS, `xdg-open` on Linux)
- **Terminal**: `tmux capture-pane` for context capture
