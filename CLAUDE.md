# CLAUDE.md

> Last updated: 2026-02-08

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **command_generator** (`cmd`), a CLI tool that generates shell commands from natural language using Claude AI. The primary interface is a fish shell key binding (**Ctrl+G**) that places the accepted command directly on the user's prompt line, ready to edit or execute. It also works as a standalone CLI that copies to clipboard. Users describe what they want in plain English, Claude generates the command with an explanation, and users can iteratively refine it with feedback until accepted.

## Build Commands

Uses **mise** as the task runner:

```bash
mise run build          # Build Go binary to ./cmd
mise run install-cli    # Install binary to ~/.local/bin/cmd
mise run web-build      # Build React frontend
mise run web-dev        # Dev server for frontend (hot reload)
mise run build-all      # Build both backend and frontend
mise run install        # Full installation (CLI + web assets + fish integration)
mise run uninstall      # Remove binary and fish integration
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

1. Parse flags (`--model`, `--context-lines`, `--output`, `--logs`, `--help`)
2. Set up SIGINT handler for clean Ctrl+C exit (important for shell key binding integration)
3. Get query from args or interactive prompt ("What do you need?")
4. Capture tmux terminal context (scrollback)
5. Detect build tools in current directory
6. Detect documentation files for additional context
7. Call Claude API via `claude` CLI with JSON schema for structured output
8. Display command + explanation
9. Interactive loop: Accept (A), Reject with feedback (R), Quit (Q)
10. On accept: write command to `--output` file, or copy to clipboard
11. Log all interactions to JSON files

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

### Shell Integration (`shell/`)

Fish shell integration (`shell/cmd.fish`):
- `cmd-generate` function bound to **Ctrl+G**
- Uses `stty sane` to reset terminal from fish's raw mode before running `cmd`
- Reads from `/dev/tty` for proper terminal I/O
- Writes accepted command to a temp file via `--output`, then places it on the fish prompt via `commandline -r`
- Installed to `~/.config/fish/conf.d/cmd.fish` by `mise run install`

### Key Data Paths

| Path | Purpose |
|------|---------|
| `~/.config/cmd/claude.md` | User preferences (prepended to system prompt) |
| `~/.local/share/cmd/logs/` | JSON session logs |
| `~/.local/bin/cmd` | Installed binary location |
| `~/.config/fish/conf.d/cmd.fish` | Fish shell integration (Ctrl+G binding) |

## Code Patterns

- **Structured Claude output**: Uses `--json-schema` flag with `GenerateResponse` struct (command + explanation fields)
- **JSON extraction**: `extractJSON()` in claude.go handles markdown code blocks in responses
- **Atomic file writes**: Logging uses temp file + rename pattern for thread safety
- **Parser interface**: Build tool detection is extensible via common interface
