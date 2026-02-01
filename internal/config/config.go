package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultModel    = "opus"
	ConfigDirName   = "cmd"
	ClaudeMdName    = "claude.md"
	DefaultClaudeMd = `# Command Generation Preferences

- Generate commands for macOS/zsh unless context suggests otherwise
- Prefer modern CLI tools when available (ripgrep over grep, fd over find, etc.)
- Use safe defaults (e.g., prefer interactive flags like -i for destructive operations)
`
)

type Config struct {
	Model       string
	ClaudeMdDir string
}

// GetConfigDir returns the path to ~/.config/cmd
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", ConfigDirName), nil
}

// GetClaudeMdPath returns the path to the claude.md file
func GetClaudeMdPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ClaudeMdName), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}

// EnsureClaudeMd creates the claude.md file with defaults if it doesn't exist
func EnsureClaudeMd() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	claudeMdPath, err := GetClaudeMdPath()
	if err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(claudeMdPath); os.IsNotExist(err) {
		// Create with default content
		return os.WriteFile(claudeMdPath, []byte(DefaultClaudeMd), 0644)
	}

	return nil
}

// LoadClaudeMd reads the claude.md file content
func LoadClaudeMd() (string, error) {
	claudeMdPath, err := GetClaudeMdPath()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(claudeMdPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return string(content), nil
}

// Load returns a Config with the specified model or default
func Load(model string) *Config {
	if model == "" {
		model = DefaultModel
	}

	configDir, _ := GetConfigDir()

	return &Config{
		Model:       model,
		ClaudeMdDir: configDir,
	}
}
