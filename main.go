package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/term"

	"github.com/jerryluo/cmd/internal/buildtools"
	"github.com/jerryluo/cmd/internal/claude"
	"github.com/jerryluo/cmd/internal/clipboard"
	"github.com/jerryluo/cmd/internal/config"
	"github.com/jerryluo/cmd/internal/docs"
	"github.com/jerryluo/cmd/internal/logging"
	"github.com/jerryluo/cmd/internal/terminal"
	"github.com/jerryluo/cmd/internal/tui"
)

func main() {
	// Parse flags
	model := flag.String("model", "", "Claude model to use (default: opus)")
	contextLines := flag.Int("context-lines", terminal.ScrollbackLines, "Number of tmux scrollback lines to capture")
	help := flag.Bool("help", false, "Show help")
	logs := flag.Bool("logs", false, "Launch log viewer")
	output := flag.String("output", "", "Write accepted command to file instead of clipboard")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	// Handle --logs flag
	if *logs {
		tui.Run()
		return
	}

	// Handle Ctrl+C for clean exit (especially when launched from shell key bindings)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println()
		os.Exit(130)
	}()

	// Get the query from arguments or interactive prompt
	reader := bufio.NewReader(os.Stdin)
	args := flag.Args()
	var query string
	if len(args) == 0 {
		fmt.Print("What do you need? ")
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		query = strings.TrimSpace(line)
		if query == "" {
			os.Exit(0)
		}
	} else {
		query = strings.Join(args, " ")
	}

	// Check claude CLI is available
	if err := claude.CheckClaudeCLI(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Load config and ensure claude.md exists
	cfg := config.Load(*model)
	if err := config.EnsureClaudeMd(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create claude.md: %v\n", err)
	}

	// Load claude.md content
	claudeMdContent, err := config.LoadClaudeMd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load claude.md: %v\n", err)
	}

	// Capture terminal context
	terminalContext, warning, err := terminal.CaptureContext(*contextLines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}
	if warning != "" {
		fmt.Fprintln(os.Stderr, warning)
	}

	// Get tmux info for display
	tmuxInfo := terminal.GetTmuxInfo()

	// Detect build tools in current directory
	buildToolsResult := buildtools.Detect(".")
	buildToolsContext := ""
	if buildToolsResult != nil {
		buildToolsContext = buildToolsResult.FormatForPrompt()
	}

	// Detect documentation files
	docsResult := docs.Detect(".")
	docsContext := ""
	if docsResult != nil {
		docsContext = docsResult.FormatForPrompt()
	}

	// Initialize request logger
	logger := logging.NewLogger(query, claudeMdContent, terminalContext, docsContext, cfg.Model, tmuxInfo)

	// Interactive loop
	feedback := ""

	for {
		// Display generation message with model and tmux context
		var tmuxContext string
		if tmuxInfo.InTmux {
			tmuxContext = fmt.Sprintf("tmux: %s/%s/%s", tmuxInfo.Session, tmuxInfo.Window, tmuxInfo.Pane)
		} else {
			tmuxContext = "no tmux context"
		}
		fmt.Printf("\nGenerating command using %s (%s)...\n", cfg.Model, tmuxContext)

		result, err := claude.GenerateCommand(cfg.Model, claudeMdContent, terminalContext, buildToolsContext, docsContext, query, feedback)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Log this iteration
		logger.AddIteration(feedback, result.SystemPrompt, result.UserPrompt,
			result.RawOutput, result.Response.Command, result.Response.Explanation)

		// Display the command and explanation
		fmt.Println()
		fmt.Printf("\033[1mCommand:\033[0m %s\n", result.Response.Command)
		fmt.Println()
		fmt.Println("\033[1mExplanation:\033[0m")
		printExplanation(result.Response.Explanation)
		fmt.Println()

		// Prompt for action
		fmt.Print("\033[1m[A]\033[0mccept  \033[1m[R]\033[0meject with feedback  \033[1m[Q]\033[0muit: ")

		key, err := readSingleKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			os.Exit(1)
		}
		fmt.Println() // Move to next line after keypress

		switch key {
		case 'a', 'A':
			logger.Finalize(logging.StatusAccepted, "")
			if *output != "" {
				if err := os.WriteFile(*output, []byte(result.Response.Command), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to %s: %v\n", *output, err)
					os.Exit(1)
				}
			} else {
				// Copy to clipboard
				if err := clipboard.Copy(result.Response.Command); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not copy to clipboard: %v\n", err)
					fmt.Printf("Command: %s\n", result.Response.Command)
				} else {
					fmt.Println("Command copied to clipboard!")
				}
			}
			os.Exit(0)

		case 'q', 'Q':
			logger.Finalize(logging.StatusQuit, "")
			fmt.Println("Exiting without copying.")
			os.Exit(0)

		case 'r', 'R':
			// Get feedback using normal buffered input
			fmt.Print("Enter feedback: ")
			feedback, err = reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading feedback: %v\n", err)
				os.Exit(1)
			}
			feedback = strings.TrimSpace(feedback)
			if feedback == "" {
				fmt.Println("No feedback provided, please try again.")
				continue
			}
			// Loop continues with new feedback

		case 3: // Ctrl+C
			logger.Finalize(logging.StatusQuit, "")
			fmt.Println("^C")
			os.Exit(0)

		default:
			fmt.Println("Invalid option. Please enter A, R, or Q.")
		}
	}
}

func printUsage() {
	fmt.Println("cmd - Generate CLI commands from natural language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cmd [options] [query]")
	fmt.Println("  cmd --logs")
	fmt.Println()
	fmt.Println("If no query is provided, an interactive prompt is shown.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --model <model>       Claude model to use (default: opus)")
	fmt.Println("  --context-lines <n>   Number of tmux scrollback lines to capture (default: 100)")
	fmt.Println("  --output <file>       Write accepted command to file instead of clipboard")
	fmt.Println("  --logs                Launch log viewer")
	fmt.Println("  --help                Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cmd \"find all large files modified today\"")
	fmt.Println("  cmd --model sonnet \"compress all images in current directory\"")
	fmt.Println("  cmd --output /tmp/cmd.txt")
	fmt.Println("  cmd --logs")
	fmt.Println()
	fmt.Println("Shell integration:")
	fmt.Println("  Fish: Press Ctrl+G to generate a command directly on your prompt")
	fmt.Println("  Install: mise run install (includes fish integration)")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  ~/.config/cmd/claude.md - Customize command generation preferences")
}

func printExplanation(explanation string) {
	lines := strings.Split(explanation, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Add bullet if not already present
		if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "•") {
			fmt.Printf("  • %s\n", line)
		} else {
			// Replace - or * with • for consistency
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")
			line = strings.TrimPrefix(line, "• ")
			fmt.Printf("  • %s\n", line)
		}
	}
}

// readSingleKey reads a single keypress without requiring Enter
func readSingleKey() (byte, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	return b[0], err
}
