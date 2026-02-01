package buildtools

import (
	"github.com/BurntSushi/toml"
)

// PyprojectParser parses pyproject.toml scripts
type PyprojectParser struct{}

// FileName returns the pyproject.toml filename
func (p *PyprojectParser) FileName() string {
	return "pyproject.toml"
}

type pyprojectConfig struct {
	Project struct {
		Scripts map[string]string `toml:"scripts"`
	} `toml:"project"`
	Tool struct {
		Poetry struct {
			Scripts map[string]string `toml:"scripts"`
		} `toml:"poetry"`
		PDM struct {
			Scripts map[string]interface{} `toml:"scripts"`
		} `toml:"pdm"`
	} `toml:"tool"`
}

// Parse extracts scripts from pyproject.toml
func (p *PyprojectParser) Parse(content []byte) (*Tool, error) {
	var cfg pyprojectConfig
	if err := toml.Unmarshal(content, &cfg); err != nil {
		return nil, nil // Graceful degradation
	}

	tool := &Tool{
		Name:     "python",
		File:     "pyproject.toml",
		Commands: []Command{},
	}

	// Check PEP 621 project.scripts
	for name := range cfg.Project.Scripts {
		tool.Commands = append(tool.Commands, Command{Name: name})
	}

	// Check tool.poetry.scripts
	for name := range cfg.Tool.Poetry.Scripts {
		tool.Commands = append(tool.Commands, Command{Name: name})
	}

	// Check tool.pdm.scripts
	for name := range cfg.Tool.PDM.Scripts {
		tool.Commands = append(tool.Commands, Command{Name: name})
	}

	if len(tool.Commands) == 0 {
		return nil, nil
	}

	return tool, nil
}
