# Backend Structure

> Last updated: 2026-02-13

## Directory Structure

```
/
├── main.go                     # CLI entry point (~270 lines)
├── go.mod                      # Module: github.com/jerryluo/cmd
├── go.sum                      # Dependency lock
├── mise.toml                   # Task runner config
├── shell/
│   └── cmd.fish                # Fish shell integration (Ctrl+G)
└── internal/
    ├── buildtools/             # Build tool detection
    │   ├── buildtools.go       # Detection orchestration + types
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
    │   ├── docs.go             # Documentation detection + types
    │   ├── parser.go           # Markdown parsing logic
    │   └── docs_test.go        # Tests
    ├── logging/
    │   └── logging.go          # Session logging + log querying
    ├── terminal/
    │   └── context.go          # tmux context capture
    └── tui/
        ├── tui.go              # TUI entry point + main model
        ├── list.go             # Log list view (table + search + filters)
        ├── detail.go           # Log detail view (8 tabs + viewport)
        ├── keys.go             # Keybinding definitions
        └── styles.go           # Lipgloss styling
```

## Entry Point (`main.go`)

### Responsibilities

1. **Flag Parsing**: `--model`, `--context-lines`, `--output`, `--logs`, `--help`
2. **Mode Selection**: TUI log viewer (`--logs`) vs command generation
3. **Context Gathering**: Combines config, tmux scrollback, build tools, docs
4. **Interactive Loop**: Accept/Reject/Quit handling with single-key input
5. **Output**: Clipboard copy or file write via `--output`

### Key Functions

| Function | Purpose |
|----------|---------|
| `main()` | Entry point, orchestrates entire flow |
| `printUsage()` | Displays help message |
| `printExplanation()` | Formats explanation with bullet points |
| `readSingleKey()` | Raw terminal input for A/R/Q (via `golang.org/x/term`) |

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
| `Makefile` | MakefileParser | Make targets (with preceding comment descriptions) |
| `package.json` | PackageJSONParser | npm scripts (with script body as description) |
| `mise.toml` | MiseParser | mise tasks (description or run command) |
| `justfile` | JustfileParser | just recipes (with preceding comment descriptions) |
| `Taskfile.yml` | TaskfileParser | task commands (with `desc` field) |
| `Cargo.toml` | CargoParser | Standard cargo commands (build, run, test, etc.) |
| `pyproject.toml` | PyprojectParser | PEP 621 + Poetry + PDM scripts |
| `docker-compose.yml` | DockerComposeParser | Standard commands + per-service up |

**Key Function:**
```go
func Detect(dir string) *DetectionResult
func (r *DetectionResult) FormatForPrompt() string
```

### `internal/claude/`

Communicates with Claude CLI using JSON schema for structured output.

**Key Functions:**
```go
func GenerateCommand(
    model, claudeMdContent, terminalContext,
    buildToolsContext, docsContext, userQuery, feedback string,
) (*GenerateResult, error)

func CheckClaudeCLI() error
```

**Claude CLI Invocation:**
```bash
claude -p --model <model> --output-format json \
    --append-system-prompt <prompt> \
    --json-schema <schema> \
    <user_prompt>
```

**Response Handling:**
1. Parse outer `ClaudeResponse` JSON (with `result`, `structured_output`, `is_error`)
2. Prefer `structured_output` (from `--json-schema`)
3. Fallback: extract JSON from `result` field (handles markdown code blocks)

**JSON Extraction (`extractJSON`):**
- Tries parsing as-is
- Regex extraction from markdown code blocks
- Brace-matching with string escape awareness

### `internal/clipboard/`

Platform-specific clipboard operations.

| Platform | Tool Used |
|----------|-----------|
| macOS | `pbcopy` |
| Linux | `xclip` (primary) → `xsel` (fallback) |

### `internal/config/`

Manages user preferences.

**Paths:**
- Config dir: `~/.config/cmd/`
- Preferences: `~/.config/cmd/claude.md`

**Constants:**
```go
const DefaultModel  = "opus"
const ConfigDirName = "cmd"
const ClaudeMdName  = "claude.md"
```

**Default claude.md:**
```markdown
# Command Generation Preferences

- Generate commands for macOS/zsh unless context suggests otherwise
- Prefer modern CLI tools when available (ripgrep over grep, fd over find, etc.)
- Use safe defaults (e.g., prefer interactive flags like -i for destructive operations)
```

### `internal/docs/`

Detects documentation files and extracts command-related sections.

**Scanned Files:** `README.md`, `CLAUDE.md`, `AGENTS.md`

**Parsing Logic (`parser.go`):**
- Finds headings matching relevant keywords (build, install, usage, setup, etc.)
- Extracts full section content under matching headings
- Detects standalone shell code blocks outside relevant sections
- Tracks code fences to avoid false heading matches

**Relevant Heading Keywords:**
`build`, `development`, `dev`, `installation`, `install`, `usage`, `commands`, `cli`, `getting started`, `quick start`, `quickstart`, `running`, `run`, `setup`, `prerequisites`, `requirements`

### `internal/logging/`

JSON session logging with atomic writes.

**Log Location:** `~/.local/share/cmd/logs/`
**Filename Format:** `2006-01-02T15-04-05Z.json` (UTC)

**Key Features:**
- Atomic writes (temp file + `os.Rename`)
- Thread-safe with `sync.Mutex`
- Nil-safe methods (no-op if logger creation failed)
- Supports search across query, command, and explanation fields
- Log listing sorted by timestamp descending

**Key Functions:**
```go
func NewLogger(...) *Logger
func (l *Logger) AddIteration(...)
func (l *Logger) Finalize(status FinalStatus, feedback string)
func ListLogs() ([]LogSummary, error)
func SearchLogs(query string) ([]LogSummary, error)
func ReadLog(id string) (*SessionLog, error)
func ReadLogWithID(id string) (*SessionLogWithID, error)
```

### `internal/tui/`

Terminal UI log viewer built with Charm's Bubbletea framework.

**Architecture:** Two-view state machine (list ↔ detail)

| File | Purpose |
|------|---------|
| `tui.go` | Main model, view routing, clipboard ops, `Run()` entry point |
| `list.go` | Log list table with search and status filtering |
| `detail.go` | 8-tab detail view with scrollable viewport |
| `keys.go` | `keyMap` struct with all keybindings |
| `styles.go` | Lipgloss styles (colors, borders, layout) |

**View States:** `listView`, `detailView`

**List View Features:**
- Table: Query, Status, Model, Time, Command columns
- Search: `/` activates search input
- Status filter: `s` cycles all → accepted → rejected → quit → all
- Copy: `c` copies selected log's command

**Detail View Features:**
- 8 tabs: Response, System Prompt, User Prompt, User Query, Tmux Context, Documentation, Build Tools, Preferences
- Tab navigation: `tab`/`l` (next), `shift+tab`/`h` (prev), `1-8` (jump)
- Scrollable viewport for long content
- `c` copies content of active tab

**Key Types:**
```go
type viewState int // listView, detailView

type model struct {
    state, list, detail, keys, help, width, height, ready, showHelp, statusMessage
}

type clipboardCopyMsg struct { err error }
type clearStatusMsg struct{}
```

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
- `tmux capture-pane -p -S -N` (scrollback capture)
- `tmux display-message -p #S` (session name)
- `tmux display-message -p #W` (window name)
- `tmux display-message -p #P` (pane index)

---

## Shell Integration (`shell/cmd.fish`)

Fish shell function bound to Ctrl+G:

```fish
function cmd-generate --description "Generate a command with AI"
    set -l tmpfile (mktemp /tmp/cmd-output.XXXXXX)
    command stty sane </dev/tty 2>/dev/null
    command cmd --output $tmpfile </dev/tty
    if test $status -eq 0 -a -s $tmpfile
        commandline -r (cat $tmpfile)
    end
    rm -f $tmpfile
    commandline -f repaint
end

bind \cg cmd-generate
```

**Key Details:**
- `stty sane` resets terminal from fish's raw mode
- Reads from `/dev/tty` for proper terminal I/O
- Uses `--output` to write to temp file, then `commandline -r` to place on prompt
- Installed to `~/.config/fish/conf.d/cmd.fish`

---

## Dependencies

```
github.com/jerryluo/cmd
go 1.25.3

require:
    github.com/BurntSushi/toml v1.6.0          # TOML parsing (mise, pyproject)
    github.com/charmbracelet/bubbles v1.0.0     # TUI components (table, help, viewport)
    github.com/charmbracelet/bubbletea v1.3.10  # TUI framework
    github.com/charmbracelet/lipgloss v1.1.0    # TUI styling
    golang.org/x/term v0.39.0                   # Terminal raw mode
    gopkg.in/yaml.v3 v3.0.1                     # YAML parsing (taskfile, docker-compose)
```

---

## Build Tasks (`mise.toml`)

| Task | Description |
|------|-------------|
| `build` | Build binary (`go build -o cmd .`) |
| `install` | Build + install to `~/.local/bin/cmd` + fish integration |
| `uninstall` | Remove binary and fish integration |
