package buildtools

import (
	"os"
)

// CargoParser parses Cargo.toml and returns standard cargo commands
type CargoParser struct{}

// FileName returns the Cargo.toml filename
func (p *CargoParser) FileName() string {
	return "Cargo.toml"
}

// Standard cargo commands available for all Rust projects
var cargoCommands = []Command{
	{Name: "build", Description: "Compile the current package"},
	{Name: "run", Description: "Run the current package"},
	{Name: "test", Description: "Run the tests"},
	{Name: "check", Description: "Analyze the current package"},
	{Name: "clean", Description: "Remove generated artifacts"},
	{Name: "doc", Description: "Build documentation"},
	{Name: "clippy", Description: "Run clippy lints"},
	{Name: "fmt", Description: "Format code with rustfmt"},
}

// Parse returns standard cargo commands if Cargo.toml exists
func (p *CargoParser) Parse(content []byte) (*Tool, error) {
	// Just verify the file can be read (is valid)
	// We don't need to parse the TOML since we return standard commands
	if len(content) == 0 {
		return nil, nil
	}

	// Check if file actually exists and is readable
	if _, err := os.Stat("Cargo.toml"); os.IsNotExist(err) {
		return nil, nil
	}

	tool := &Tool{
		Name:     "cargo",
		File:     "Cargo.toml",
		Commands: cargoCommands,
	}

	return tool, nil
}
