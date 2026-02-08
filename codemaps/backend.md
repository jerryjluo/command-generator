# Backend Structure

> Last updated: 2026-02-01

## Directory Structure

```
/
├── main.go                     # CLI entry point (287 lines)
├── go.mod                      # Module: github.com/jerryluo/cmd
├── go.sum                      # Dependency lock
├── mise.toml                   # Task runner config
├── cmd                         # Compiled binary
└── internal/
    ├── buildtools/             # Build tool detection
    │   ├── buildtools.go       # Detection orchestration
    │   ├── parser.go           # Parser interface
    │   ├── makefile.go         # Makefile parser
    │   ├── package_json.go     # npm scripts parser
    │   ├── mise.go             # mise tasks parser
    │   ├── justfile.go         # Justfile parser
    │   ├── taskfile.go         # Taskfile parser
    │   ├── cargo.go            # Cargo.toml parser
    │   ├── pyproject.go        # pyproject.toml parser
    │   └── docker_compose.go   # docker-compose parser
    ├── claude/
    │   └── claude.go           # Claude API integration
    ├── clipboard/
    │   └── clipboard.go        # Cross-platform clipboard
    ├── config/
    │   └── config.go           # User configuration
    ├── docs/
    │   ├── docs.go             # Documentation detection
    │   ├── parser.go           # Doc file parsing
    │   └── docs_test.go        # Tests
    ├── logging/
    │   └── logging.go          # Session logging
    ├── server/
    │   ├── server.go           # HTTP server
    │   ├── handlers.go         # API handlers
    │   └── middleware.go       # CORS, logging
    └── terminal/
        └── context.go          # tmux context capture
```

## Entry Point (`main.go`)

### Responsibilities

1. **Flag Parsing**: `--model`, `--context-lines`, `--logs`, `--help`
2. **Mode Selection**: Log viewer vs command generation
3. **Context Gathering**: Combines all context sources
4. **Interactive Loop**: Accept/Reject/Quit handling
5. **Clipboard Integration**: Copy accepted commands

### Key Functions

| Function | Purpose |
|----------|---------|
| `main()` | Entry point, orchestrates entire flow |
| `runLogViewer()` | Starts embedded web server |
| `printUsage()` | Displays help message |
| `readSingleKey()` | Raw terminal input for A/R/Q |

### Embedded Assets

```go
//go:embed web/dist
var webAssets embed.FS
```

Embeds compiled React app for standalone binary deployment.

---

## Package Details

### `internal/buildtools/`

Detects build tools and extracts available commands.

**Interface:**
```go
type Parser interface {
    FileName() string
    Parse(content []byte) (*Tool, error)
}
```

**Registered Parsers:**

| File | Parser | Commands Extracted |
|------|--------|-------------------|
| `Makefile` | MakefileParser | Make targets |
| `package.json` | PackageJsonParser | npm scripts |
| `mise.toml` | MiseParser | mise tasks |
| `Justfile` | JustfileParser | just recipes |
| `Taskfile.yml` | TaskfileParser | task commands |
| `Cargo.toml` | CargoParser | cargo commands |
| `pyproject.toml` | PyprojectParser | Python scripts |
| `docker-compose.yml` | DockerComposeParser | Services |

**Key Function:**
```go
func Detect(dir string) *DetectionResult
```

### `internal/claude/`

Communicates with Claude CLI using JSON schema for structured output.

**Key Function:**
```go
func GenerateCommand(
    model, claudeMdContent, terminalContext,
    buildToolsContext, docsContext, userQuery, feedback string,
) (*GenerateResult, error)
```

**Claude CLI Invocation:**
```bash
claude -p --model <model> --output-format json \
    --append-system-prompt <prompt> \
    --json-schema <schema>
```

**JSON Schema Used:**
```json
{
    "type": "object",
    "properties": {
        "command": {"type": "string"},
        "explanation": {"type": "string"}
    },
    "required": ["command", "explanation"]
}
```

### `internal/clipboard/`

Platform-specific clipboard operations.

| Platform | Tool Used |
|----------|-----------|
| macOS | `pbcopy` |
| Linux | `xclip` → `xsel` fallback |

### `internal/config/`

Manages user preferences.

**Paths:**
- Config dir: `~/.config/cmd/`
- Preferences: `~/.config/cmd/claude.md`

**Default claude.md:**
```markdown
# Command Generation Preferences

- Generate commands for macOS/zsh unless context suggests otherwise
- Prefer modern CLI tools (ripgrep over grep, fd over find, etc.)
- Use safe defaults (e.g., prefer interactive flags like -i)
```

### `internal/docs/`

Detects documentation files in project.

**Supported Files:**
- `README.md`
- `CONTRIBUTING.md`
- `.github/CONTRIBUTING.md`
- `docs/` directory contents
- `INSTALL.md`
- `.github/INSTALL.md`

### `internal/logging/`

JSON session logging with atomic writes.

**Log Location:** `~/.local/share/cmd/logs/`

**Key Features:**
- Atomic writes (temp file + rename)
- Thread-safe with mutex
- Supports search and filtering

**Key Functions:**
```go
func NewLogger(...) *Logger
func (l *Logger) AddIteration(...)
func (l *Logger) Finalize(status, feedback) error
func ListLogs() ([]LogSummary, error)
func SearchLogs(query string) ([]LogSummary, error)
```

### `internal/server/`

HTTP server for log viewer.

**Port:** 8765 (configurable)

**Routes:**
| Route | Handler |
|-------|---------|
| `GET /api/v1/logs` | List logs with filtering |
| `GET /api/v1/logs/{id}` | Get log details |
| `GET /*` | Serve React SPA |

**Middleware:**
- CORS headers
- Request logging

### `internal/terminal/`

Captures tmux context.

**Constants:**
```go
const ScrollbackLines = 100
```

**Key Functions:**
```go
func CaptureContext(lines int) (string, string, error)
func GetTmuxInfo() TmuxInfo
func InTmux() bool
```

**tmux Commands Used:**
- `tmux capture-pane -p -S -N` (scrollback)
- `tmux display-message -p` (session info)

---

## Dependencies

```
github.com/jerryluo/cmd
go 1.25.3

require:
    github.com/BurntSushi/toml v1.6.0  # TOML parsing (mise, cargo, pyproject)
    golang.org/x/sys v0.40.0           # System calls
    golang.org/x/term v0.39.0          # Terminal handling
    gopkg.in/yaml.v3 v3.0.1            # YAML parsing (taskfile, docker-compose)
```

---

## Build Tasks (`mise.toml`)

| Task | Description |
|------|-------------|
| `build` | Build binary with embedded web assets |
| `install` | Install to `~/.local/bin/cmd` |
| `web-install` | Install npm dependencies |
| `web-build` | Build React frontend |
| `web-dev` | Start Vite dev server |
| `build-all` | Build backend + frontend |
