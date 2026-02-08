# Architecture Overview

> Last updated: 2026-02-08

## System Overview

**command_generator** (`cmd`) is a CLI tool that generates shell commands from natural language using Claude AI. It combines a Go backend for CLI operations with a React frontend for log visualization.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User Interface                              │
├──────────────────────────────┬──────────────────────────────────────┤
│       CLI (main.go)          │      Web Viewer (React SPA)          │
│  - Natural language input    │  - Log browsing & filtering          │
│  - Interactive A/R/Q loop    │  - Session detail viewing            │
│  - Clipboard / file output   │  - Served on localhost:8765          │
├──────────────────────────────┘──────────────────────────────────────┤
│       Shell Integration (fish)                                      │
│  - Ctrl+G keybinding                                                │
│  - Writes to temp file via --output, places on fish prompt          │
└──────────────┬──────────────────────────────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────────────────────────────┐
│                          Go Backend                                   │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│   claude/   │ buildtools/ │  terminal/  │   logging/  │   server/   │
│  API client │  8 parsers  │ tmux context│  JSON logs  │  HTTP API   │
├─────────────┼─────────────┤             ├─────────────┼─────────────┤
│   config/   │   docs/     │             │  clipboard/ │             │
│  user prefs │ doc parsing │             │  copy cmd   │             │
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

### Shell Integration Flow (Fish)

```
Ctrl+G → stty sane → cmd --output /tmp/cmd-output.XXXX → Read temp file → commandline -r
```

### Log Viewer Flow

```
cmd --logs → Start Server → Embed React → Serve API → Open Browser → Display Logs
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
└── internal/terminal     # TmuxInfo type
```

## Technology Stack

| Layer | Technology |
|-------|------------|
| CLI | Go 1.25.3 |
| Build Runner | mise |
| Frontend | React 19 + TypeScript + Vite 7 |
| Styling | Tailwind CSS 4 |
| Routing | React Router 7 |
| API | REST (Go net/http) |
| Storage | JSON files on disk |
| External | `claude` CLI tool |
| Shell | Fish shell integration (Ctrl+G) |

## Key Configuration Paths

| Path | Purpose |
|------|---------|
| `~/.config/cmd/claude.md` | User preferences (prepended to system prompt) |
| `~/.local/share/cmd/logs/` | Session logs (JSON files) |
| `~/.local/bin/cmd` | Installed binary |
| `~/.config/fish/conf.d/cmd.fish` | Fish shell integration (Ctrl+G binding) |

## Build Process

```
mise run build
       │
       ├── 1. web-build (npm run build in web/)
       │      └── Compiles React → web/dist/
       │
       └── 2. go build -o cmd .
              └── Embeds web/dist/ via //go:embed
              └── Outputs single binary: ./cmd
```

## External Integrations

- **Claude API**: Via `claude` CLI subprocess with `--json-schema` for structured output
- **Clipboard**: Platform-specific (`pbcopy` on macOS, `xclip`/`xsel` on Linux)
- **Browser**: Platform-specific (`open` on macOS, `xdg-open` on Linux, `cmd /c start` on Windows)
- **Terminal**: `tmux capture-pane` for scrollback, `tmux display-message` for session info
- **Shell**: Fish shell function `cmd-generate` bound to `\cg` (Ctrl+G)
