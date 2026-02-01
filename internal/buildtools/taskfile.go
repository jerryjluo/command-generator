package buildtools

import (
	"gopkg.in/yaml.v3"
)

// TaskfileParser parses Taskfile.yml tasks
type TaskfileParser struct{}

// FileName returns the Taskfile.yml filename
func (p *TaskfileParser) FileName() string {
	return "Taskfile.yml"
}

type taskfile struct {
	Tasks map[string]taskDef `yaml:"tasks"`
}

type taskDef struct {
	Desc string `yaml:"desc"`
}

// Parse extracts tasks from Taskfile.yml
func (p *TaskfileParser) Parse(content []byte) (*Tool, error) {
	var tf taskfile
	if err := yaml.Unmarshal(content, &tf); err != nil {
		return nil, nil // Graceful degradation
	}

	if len(tf.Tasks) == 0 {
		return nil, nil
	}

	tool := &Tool{
		Name:     "task",
		File:     "Taskfile.yml",
		Commands: []Command{},
	}

	for name, task := range tf.Tasks {
		cmd := Command{Name: name}
		if task.Desc != "" {
			cmd.Description = task.Desc
		}
		tool.Commands = append(tool.Commands, cmd)
	}

	return tool, nil
}
