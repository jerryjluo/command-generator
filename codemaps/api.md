# API Reference

> Last updated: 2026-02-13

## CLI Interface

```
cmd [options] [query]
cmd --logs

Options:
  --model <model>        Claude model to use (default: opus)
  --context-lines <n>    Lines of tmux scrollback (default: 100)
  --output <file>        Write accepted command to file instead of clipboard
  --logs                 Launch TUI log viewer
  --help                 Show usage information

Interactive Commands:
  A - Accept command (copies to clipboard or writes to --output file)
  R - Reject with feedback (refine command)
  Q - Quit without accepting
```

## TUI Log Viewer Keybindings

### Global

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit |
| `?` | Toggle help |

### List View

| Key | Action |
|-----|--------|
| `enter` | View log details |
| `/` | Search logs |
| `s` | Cycle status filter (all → accepted → rejected → quit) |
| `c` | Copy selected log's command |
| `esc` | Clear search |
| Arrow keys / PgUp / PgDn | Navigate |

### Detail View

| Key | Action |
|-----|--------|
| `tab` / `l` | Next tab |
| `shift+tab` / `h` | Previous tab |
| `1-8` | Jump to specific tab |
| `c` | Copy active tab content |
| `esc` / `backspace` | Back to list |
| Arrow keys | Scroll content |

---

## Internal Go APIs

### Claude Package (`internal/claude/`)

```go
// GenerateCommand calls Claude CLI to generate a shell command
func GenerateCommand(
    model string,              // Claude model to use (e.g., "opus")
    claudeMdContent string,    // User preferences from ~/.config/cmd/claude.md
    terminalContext string,     // Captured tmux scrollback
    buildToolsContext string,   // Detected build commands
    docsContext string,         // Documentation sections
    userQuery string,           // Natural language request
    feedback string,            // Optional refinement feedback
) (*GenerateResult, error)

// CheckClaudeCLI verifies the claude CLI is installed
func CheckClaudeCLI() error
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
| `MakefileParser` | `Makefile` | make |
| `PackageJSONParser` | `package.json` | npm |
| `MiseParser` | `mise.toml` | mise |
| `JustfileParser` | `justfile` | just |
| `TaskfileParser` | `Taskfile.yml` | task |
| `CargoParser` | `Cargo.toml` | cargo |
| `PyprojectParser` | `pyproject.toml` | python |
| `DockerComposeParser` | `docker-compose.yml` | docker-compose |

### Terminal Package (`internal/terminal/`)

```go
// CaptureContext gets tmux scrollback buffer
func CaptureContext(lines int) (context string, warning string, err error)

// GetTmuxInfo extracts current tmux session info
func GetTmuxInfo() TmuxInfo

// InTmux checks if running inside tmux
func InTmux() bool

const ScrollbackLines = 100 // Default scrollback capture
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
func (l *Logger) Finalize(status FinalStatus, feedback string)

// ListLogs returns all log summaries (newest first)
func ListLogs() ([]LogSummary, error)

// SearchLogs filters logs by query string (searches query, command, explanation)
func SearchLogs(query string) ([]LogSummary, error)

// ReadLog reads a single log file by ID
func ReadLog(id string) (*SessionLog, error)

// ReadLogWithID reads a log and wraps it with its ID
func ReadLogWithID(id string) (*SessionLogWithID, error)

// GetLogDir returns the log directory path
func GetLogDir() (string, error)
```

### Config Package (`internal/config/`)

```go
// Load creates a Config with the specified model (default: "opus")
func Load(model string) *Config

// LoadClaudeMd reads user preferences file
func LoadClaudeMd() (string, error)

// EnsureClaudeMd creates default preferences if missing
func EnsureClaudeMd() error

// EnsureConfigDir creates the config directory if needed
func EnsureConfigDir() error

// GetConfigDir returns ~/.config/cmd
func GetConfigDir() (string, error)

// GetClaudeMdPath returns full path to claude.md
func GetClaudeMdPath() (string, error)

// Constants
const DefaultModel  = "opus"
const ConfigDirName = "cmd"
const ClaudeMdName  = "claude.md"
```

### Docs Package (`internal/docs/`)

```go
// Detect reads documentation files and extracts command-related sections
func Detect(dir string) *Result

// FormatForPrompt converts docs to prompt-ready text
func (r *Result) FormatForPrompt() string

// DocFiles defines which files to scan (in priority order)
var DocFiles = []string{"README.md", "CLAUDE.md", "AGENTS.md"}

// RelevantHeadingPatterns defines heading keywords to match
var RelevantHeadingPatterns = []string{
    "build", "development", "dev", "installation", "install",
    "usage", "commands", "cli", "getting started", "quick start",
    "quickstart", "running", "run", "setup", "prerequisites", "requirements",
}
```

### Clipboard Package (`internal/clipboard/`)

```go
// Copy writes text to system clipboard
func Copy(text string) error
```

### TUI Package (`internal/tui/`)

```go
// Run starts the TUI log viewer (alternate screen mode)
func Run()
```

---

## External Dependencies

### Claude CLI

The tool calls `claude` as a subprocess:

```bash
claude \
  -p \                                    # Piped input (non-interactive)
  --model <model> \                       # Model selection
  --output-format json \                  # JSON output
  --append-system-prompt <prompt> \       # System instructions
  --json-schema <schema> \               # Structured output schema
  <prompt>                                # User prompt as positional arg
```

The `claude` CLI must be installed and available on PATH.
