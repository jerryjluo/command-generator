# cmd

> Last updated: 2026-02-13

Generate shell commands from natural language using Claude AI. Press **Ctrl+G** in fish shell and describe what you need — the accepted command is placed directly on your prompt, ready to edit or execute.

`cmd` can also be used standalone: pass a query as an argument or run it with no arguments for an interactive prompt. It copies the accepted command to your clipboard.

## Features

- **Inline fish shell integration** - Press Ctrl+G to generate a command that lands directly on your prompt line — the primary way to use `cmd`
- **Standalone CLI** - Also works as `cmd "your query"` or `cmd` with an interactive prompt, copying to clipboard
- **Context-aware** - Automatically detects your terminal history (tmux) and available build tools
- **Iterative refinement** - Provide feedback to adjust the generated command
- **Build tool detection** - Recognizes Makefile, package.json, mise, just, task, cargo, pyproject.toml, and docker-compose
- **Documentation detection** - Includes README, CONTRIBUTING, and other docs as context
- **Session logging** - All generations are logged for review
- **TUI log viewer** - Browse your generation history in a terminal interface

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [mise](https://mise.jdx.dev/) (task runner)
- [Claude Code CLI](https://claude.ai/code) (`claude` command must be available)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/command_generator.git
cd command_generator

# Build and install (CLI + fish shell integration)
mise run install
```

The binary is installed to `~/.local/bin/cmd`. Make sure this is in your PATH.

Fish shell integration is automatically installed to `~/.config/fish/conf.d/cmd.fish`. Restart fish or run `source ~/.config/fish/conf.d/cmd.fish` to activate.

## Usage

### Fish Shell (Ctrl+G)

The primary way to use `cmd`. Press **Ctrl+G** anywhere in fish to describe what you need. Once you accept, the generated command is placed directly on your prompt line — ready to review, edit, or execute. Press Ctrl+C at any point to cancel and return to your prompt.

```bash
# Activate (already done if you ran mise run install)
source shell/cmd.fish

# Or manually
cp shell/cmd.fish ~/.config/fish/conf.d/cmd.fish
```

### Standalone CLI

```bash
# Pass a query as an argument
cmd "find all Python files modified in the last week"

# Or run with no arguments for an interactive prompt
cmd
# What do you need? find all Python files modified in the last week
```

Output:
```
Command: find . -name "*.py" -mtime -7

Explanation: This finds all files ending in .py that were modified within
the last 7 days, starting from the current directory.

[A]ccept, [R]eject with feedback, [Q]uit:
```

- Press **A** to accept (copies command to clipboard)
- Press **R** to provide feedback and regenerate
- Press **Q** to quit

### Options

```bash
cmd [options] [query]

Options:
  --model <model>         Claude model to use (default: opus)
  --context-lines <n>     Lines of terminal history to include (default: 100)
  --output <file>         Write accepted command to file instead of clipboard
  --logs                  Open the log viewer
  --help                  Show help
```

### Examples

```bash
# Git operations
cmd "squash the last 3 commits"

# File operations
cmd "compress all images in this folder to 80% quality"

# System administration
cmd "show disk usage sorted by size, human readable"

# Docker
cmd "run a postgres container with persistent storage"

# Using detected build tools
cmd "run the tests"  # Uses your Makefile, package.json, etc.
```

### Log Viewer

View your command generation history:

```bash
cmd --logs
```

Opens a terminal UI where you can:
- Browse all generation sessions
- Filter by status or search query
- View full context (terminal history, build tools, prompts)
- Copy commands to clipboard

## Configuration

User preferences are stored in `~/.config/cmd/claude.md`. This file is automatically created on first run and can be customized to influence command generation.

Example preferences:
```markdown
- Target macOS with zsh
- Prefer modern CLI tools (ripgrep, fd, bat) over traditional ones
- Use verbose flags for clarity
```

## How It Works

1. Gets your query (from arguments or interactive prompt)
2. Captures your recent terminal history (requires tmux)
3. Detects build tools in your current directory
4. Detects documentation files (README, CONTRIBUTING, etc.)
5. Sends context + your request to Claude with a JSON schema
6. Displays the generated command and explanation
7. Loops for feedback until you accept or quit
8. Copies accepted command to clipboard (or writes to `--output` file) and logs the session

## Development

```bash
# Build the CLI
mise run build
```
