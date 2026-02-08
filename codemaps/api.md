# API Reference

> Last updated: 2026-02-01

## HTTP API (Log Viewer)

Base URL: `http://localhost:8765/api/v1`

### Endpoints

#### List Logs

```
GET /api/v1/logs
```

Returns paginated list of session logs with filtering and sorting.

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Filter by final status: `accepted`, `rejected`, `quit` |
| `model` | string | Filter by Claude model used |
| `search` | string | Full-text search on query and command |
| `from` | string | Start date (RFC3339 format) |
| `to` | string | End date (RFC3339 format) |
| `sort` | string | Sort field: `timestamp`, `model` |
| `order` | string | Sort order: `asc`, `desc` |
| `limit` | int | Results per page (default: 50) |
| `offset` | int | Skip N results for pagination |

**Response:**

```json
{
    "logs": [
        {
            "id": "2026-02-01T10-30-00-abc123",
            "user_query": "list all go files",
            "final_status": "accepted",
            "model": "claude-sonnet-4-20250514",
            "timestamp": "2026-02-01T10:30:00Z",
            "iteration_count": 1,
            "command_preview": "fd -e go"
        }
    ],
    "total": 150,
    "limit": 50,
    "offset": 0
}
```

#### Get Log Detail

```
GET /api/v1/logs/{id}
```

Returns full session log including all iterations and context.

**Response:**

```json
{
    "id": "2026-02-01T10-30-00-abc123",
    "user_query": "list all go files",
    "context_sources": {
        "claude_md_content": "...",
        "terminal_context": "...",
        "documentation_context": "..."
    },
    "iterations": [...],
    "metadata": {...}
}
```

---

## Internal Go APIs

### Claude Package (`internal/claude/`)

```go
// GenerateCommand calls Claude API to generate a shell command
func GenerateCommand(
    model string,              // Claude model to use
    claudeMdContent string,    // User preferences from ~/.config/cmd/claude.md
    terminalContext string,    // Captured tmux scrollback
    buildToolsContext string,  // Detected build commands
    docsContext string,        // Documentation file contents
    userQuery string,          // Natural language request
    feedback string,           // Optional refinement feedback
) (*GenerateResult, error)
```

### Build Tools Package (`internal/buildtools/`)

```go
// Detect scans a directory for build tool configurations
func Detect(dir string) *DetectionResult

// FormatForPrompt converts detection result to prompt-ready text
func (r *DetectionResult) FormatForPrompt() string

// Parser interface for adding new build tool support
type Parser interface {
    FileName() string                       // Config file to look for
    Parse(content []byte) (*Tool, error)    // Extract commands
}
```

**Registered Parsers:**

| Parser | Config File | Tool Name |
|--------|-------------|-----------|
| `MakefileParser` | `Makefile` | Makefile |
| `PackageJsonParser` | `package.json` | npm |
| `MiseParser` | `mise.toml` | mise |
| `JustfileParser` | `Justfile`, `justfile` | just |
| `TaskfileParser` | `Taskfile.yml` | task |
| `CargoParser` | `Cargo.toml` | cargo |
| `PyprojectParser` | `pyproject.toml` | Python |
| `DockerComposeParser` | `docker-compose.yml` | docker-compose |

### Terminal Package (`internal/terminal/`)

```go
// CaptureContext gets tmux scrollback buffer
func CaptureContext(lines int) (context string, warning string, err error)

// GetTmuxInfo extracts current tmux session info
func GetTmuxInfo() TmuxInfo

// InTmux checks if running inside tmux
func InTmux() bool
```

### Logging Package (`internal/logging/`)

```go
// NewLogger creates a session logger
func NewLogger(
    query string,
    claudeMd string,
    termCtx string,
    docsCtx string,
    model string,
    tmuxInfo terminal.TmuxInfo,
) *Logger

// AddIteration logs a generation attempt
func (l *Logger) AddIteration(
    feedback string,
    sysPrompt string,
    userPrompt string,
    rawOutput string,
    cmd string,
    explanation string,
)

// Finalize marks session as complete and writes to disk
func (l *Logger) Finalize(status FinalStatus, feedback string) error

// ListLogs returns all log summaries
func ListLogs() ([]LogSummary, error)

// SearchLogs filters logs by query string
func SearchLogs(query string) ([]LogSummary, error)

// GetLogDir returns the log directory path
func GetLogDir() (string, error)
```

### Config Package (`internal/config/`)

```go
// Load creates a Config with the specified model
func Load(model string) *Config

// LoadClaudeMd reads user preferences file
func LoadClaudeMd() (string, error)

// EnsureClaudeMd creates default preferences if missing
func EnsureClaudeMd() error

// GetConfigDir returns ~/.config/cmd
func GetConfigDir() (string, error)

// GetClaudeMdPath returns full path to claude.md
func GetClaudeMdPath() (string, error)
```

### Docs Package (`internal/docs/`)

```go
// Detect finds documentation files in a directory
func Detect(dir string) *DetectionResult

// FormatForPrompt converts docs to prompt-ready text
func (r *DetectionResult) FormatForPrompt() string
```

### Clipboard Package (`internal/clipboard/`)

```go
// Copy writes text to system clipboard
func Copy(text string) error
```

### Server Package (`internal/server/`)

```go
// NewServer creates an HTTP server with embedded assets
func NewServer(port int, assets fs.FS) *Server

// Start begins serving HTTP requests
func (s *Server) Start() error

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error

// OpenBrowser opens the default browser to the server URL
func OpenBrowser(url string) error

// URL returns the server's base URL
func (s *Server) URL() string
```

---

## Frontend API Client (`web/src/api/logs.ts`)

```typescript
// Fetch paginated log list with filters
async function fetchLogs(params: FilterParams): Promise<LogListResponse>

// Fetch full log details by ID
async function fetchLogById(id: string): Promise<SessionLog>
```

---

## CLI Interface

```
cmd [options] <natural language query>

Options:
  --model <model>        Claude model to use (default: claude-sonnet-4-20250514)
  --context-lines <n>    Lines of tmux scrollback (default: 100)
  --logs                 Open log viewer in browser
  --help                 Show usage information

Interactive Commands:
  A - Accept command (copies to clipboard)
  R - Reject with feedback (refine command)
  Q - Quit without accepting
```

---

## External Dependencies

### Claude CLI

The tool calls `claude` as a subprocess:

```bash
claude \
  -p \                                    # Piped input
  --model <model> \                       # Model selection
  --output-format json \                  # JSON output
  --append-system-prompt <prompt> \       # System instructions
  --json-schema <schema>                  # Structured output schema
```

The `claude` CLI must be installed and configured with `ANTHROPIC_API_KEY`.
