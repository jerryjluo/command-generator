package buildtools

import (
	"gopkg.in/yaml.v3"
)

// DockerComposeParser parses docker-compose.yml services
type DockerComposeParser struct{}

// FileName returns the docker-compose.yml filename
func (p *DockerComposeParser) FileName() string {
	return "docker-compose.yml"
}

type dockerCompose struct {
	Services map[string]interface{} `yaml:"services"`
}

// Parse extracts services from docker-compose.yml
func (p *DockerComposeParser) Parse(content []byte) (*Tool, error) {
	var dc dockerCompose
	if err := yaml.Unmarshal(content, &dc); err != nil {
		return nil, nil // Graceful degradation
	}

	if len(dc.Services) == 0 {
		return nil, nil
	}

	tool := &Tool{
		Name:     "docker-compose",
		File:     "docker-compose.yml",
		Commands: []Command{},
	}

	// Add general docker-compose commands
	tool.Commands = append(tool.Commands, Command{
		Name:        "up",
		Description: "Start all services",
	})
	tool.Commands = append(tool.Commands, Command{
		Name:        "down",
		Description: "Stop all services",
	})
	tool.Commands = append(tool.Commands, Command{
		Name:        "build",
		Description: "Build all services",
	})
	tool.Commands = append(tool.Commands, Command{
		Name:        "logs",
		Description: "View logs",
	})

	// Add per-service commands
	for serviceName := range dc.Services {
		tool.Commands = append(tool.Commands, Command{
			Name:        "up " + serviceName,
			Description: "Start " + serviceName + " service",
		})
	}

	return tool, nil
}
