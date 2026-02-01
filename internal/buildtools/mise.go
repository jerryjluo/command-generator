package buildtools

import (
	"github.com/BurntSushi/toml"
)

// MiseParser parses mise.toml tasks
type MiseParser struct{}

// FileName returns the mise.toml filename
func (p *MiseParser) FileName() string {
	return "mise.toml"
}

// miseTask represents a task definition in mise.toml
// Tasks can be defined as [tasks.taskname] sections
type miseTask struct {
	Description string      `toml:"description"`
	Run         interface{} `toml:"run"` // Can be string or array
}

type miseConfig struct {
	Tasks map[string]miseTask `toml:"tasks"`
}

// Parse extracts mise tasks from mise.toml
func (p *MiseParser) Parse(content []byte) (*Tool, error) {
	var cfg miseConfig
	if err := toml.Unmarshal(content, &cfg); err != nil {
		return nil, nil // Graceful degradation
	}

	if len(cfg.Tasks) == 0 {
		return nil, nil
	}

	tool := &Tool{
		Name:     "mise",
		File:     "mise.toml",
		Commands: []Command{},
	}

	for name, task := range cfg.Tasks {
		cmd := Command{Name: name}
		if task.Description != "" {
			cmd.Description = task.Description
		} else if task.Run != nil {
			// Use run command as description if no description provided
			var desc string
			switch v := task.Run.(type) {
			case string:
				desc = v
			}
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			if desc != "" {
				cmd.Description = desc
			}
		}
		tool.Commands = append(tool.Commands, cmd)
	}

	return tool, nil
}
