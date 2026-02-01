# cmd

A CLI tool that generates shell commands from natural language using Claude AI.

Describe what you want to do in plain English, and `cmd` generates the appropriate shell command with an explanation. Refine it iteratively with feedback until you get exactly what you need.

## Features

- **Natural language to shell commands** - Just describe what you want
- **Context-aware** - Automatically detects your terminal history (tmux) and available build tools
- **Iterative refinement** - Provide feedback to adjust the generated command
- **Build tool detection** - Recognizes Makefile, package.json, mise, just, task, cargo, pyproject.toml, and docker-compose
- **Session logging** - All generations are logged for review
- **Web log viewer** - Browse your generation history in a web interface

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [mise](https://mise.jdx.dev/) (task runner)
- [Claude Code CLI](https://claude.ai/code) (`claude` command must be available)
- [Node.js 18+](https://nodejs.org/) (for web frontend)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/command_generator.git
cd command_generator

# Install CLI only
mise run install-cli

# Or install everything (CLI + web viewer)
npm install --prefix web
mise run install
```

The binary is installed to `~/.local/bin/cmd`. Make sure this is in your PATH.

## Usage

### Basic Usage

```bash
cmd "find all Python files modified in the last week"
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
cmd [options] "your request"

Options:
  --model <model>         Claude model to use (default: claude-sonnet-4-20250514)
  --context-lines <n>     Lines of terminal history to include (default: 100)
  --logs                  Open the web log viewer
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

### Web Log Viewer

View your command generation history:

```bash
cmd --logs
```

Opens a web interface at `http://localhost:8765` where you can:
- Browse all generation sessions
- Filter by status, model, or search query
- View full context (terminal history, build tools, prompts)

## Configuration

User preferences are stored in `~/.config/cmd/claude.md`. This file is automatically created on first run and can be customized to influence command generation.

Example preferences:
```markdown
- Target macOS with zsh
- Prefer modern CLI tools (ripgrep, fd, bat) over traditional ones
- Use verbose flags for clarity
```

## How It Works

1. Captures your recent terminal history (requires tmux)
2. Detects build tools in your current directory
3. Sends context + your request to Claude with a JSON schema
4. Displays the generated command and explanation
5. Loops for feedback until you accept or quit
6. Copies accepted command to clipboard and logs the session

## Development

```bash
# Build the CLI
mise run build

# Run frontend dev server
mise run web-dev

# Build everything
mise run build-all

# Lint frontend
cd web && npm run lint
```
