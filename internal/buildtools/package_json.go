package buildtools

import (
	"encoding/json"
)

// PackageJSONParser parses package.json scripts
type PackageJSONParser struct{}

// FileName returns the package.json filename
func (p *PackageJSONParser) FileName() string {
	return "package.json"
}

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// Parse extracts npm scripts from package.json
func (p *PackageJSONParser) Parse(content []byte) (*Tool, error) {
	var pkg packageJSON
	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, nil // Graceful degradation
	}

	if len(pkg.Scripts) == 0 {
		return nil, nil
	}

	tool := &Tool{
		Name:     "npm",
		File:     "package.json",
		Commands: []Command{},
	}

	for name, script := range pkg.Scripts {
		desc := script
		// Truncate long scripts for readability
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		tool.Commands = append(tool.Commands, Command{
			Name:        name,
			Description: desc,
		})
	}

	return tool, nil
}
