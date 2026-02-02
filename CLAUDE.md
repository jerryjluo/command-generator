# CLAUDE.md

> Last updated: 2026-02-01

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **command_generator** (`cmd`), a CLI tool that generates shell commands from natural language using Claude AI. Users describe what they want in plain English, Claude generates the command with an explanation, and users can iteratively refine it with feedback until accepted.

## Build Commands

Uses **mise** as the task runner:

```bash
mise run build          # Build Go binary to ./cmd
mise run install-cli    # Install binary to ~/.local/bin/cmd
mise run web-build      # Build React frontend
mise run web-dev        # Dev server for frontend (hot reload)
mise run build-all      # Build both backend and frontend
mise run install        # Full installation (CLI + web assets)
```

Frontend (in `web/` directory):
```bash
npm install             # Install dependencies
npm run dev             # Vite dev server
npm run build           # Production build
npm run lint            # ESLint
```

## Architecture

### CLI Flow (main.go)

1. Parse flags (`--model`, `--context-lines`, `--logs`, `--help`)
2. Capture tmux terminal context (scrollback)
3. Detect build tools in current directory
4. Detect documentation files for additional context
5. Call Claude API via `claude` CLI with JSON schema for structured output
6. Display command + explanation
7. Interactive loop: Accept (A), Reject with feedback (R), Quit (Q)
8. Copy accepted command to clipboard
9. Log all interactions to JSON files

### Internal Packages

| Package | Purpose |
|---------|---------|
| `internal/claude/` | Claude API integration via `claude` CLI tool with JSON schema output |
| `internal/buildtools/` | Detects 8 build systems (Makefile, package.json, mise, just, task, cargo, pyproject, docker-compose) and extracts available commands |
| `internal/config/` | Loads `~/.config/cmd/claude.md` user preferences |
| `internal/terminal/` | Captures tmux scrollback via `tmux capture-pane` |
| `internal/logging/` | JSON session logging to `~/.local/share/cmd/logs/` |
| `internal/clipboard/` | Cross-platform clipboard (pbcopy/xclip) |
| `internal/server/` | HTTP server for web log viewer |
| `internal/docs/` | Detects documentation files (README.md, CONTRIBUTING.md, etc.) for context |

### Build Tools Detection

Each parser in `internal/buildtools/` implements the `Parser` interface:
- `FileName() string` - returns the config file to look for
- `Parse(content []byte) ([]Command, error)` - extracts available commands

To add a new build system, create a new file implementing this interface and register it in `buildtools.go`.

### Web Frontend (`web/`)

React + TypeScript + Tailwind CSS SPA for browsing generation logs:
- `/` - Log list with filtering, sorting, pagination
- `/logs/{id}` - Detailed log view with tabs (response, prompts, context)

API served by Go backend on port 8765:
- `GET /api/v1/logs` - List logs with query params for filtering
- `GET /api/v1/logs/{id}` - Full log details

### Key Data Paths

| Path | Purpose |
|------|---------|
| `~/.config/cmd/claude.md` | User preferences (prepended to system prompt) |
| `~/.local/share/cmd/logs/` | JSON session logs |
| `~/.local/bin/cmd` | Installed binary location |

## Code Patterns

- **Structured Claude output**: Uses `--json-schema` flag with `GenerateResponse` struct (command + explanation fields)
- **JSON extraction**: `extractJSON()` in claude.go handles markdown code blocks in responses
- **Atomic file writes**: Logging uses temp file + rename pattern for thread safety
- **Parser interface**: Build tool detection is extensible via common interface
