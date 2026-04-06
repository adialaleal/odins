package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const ProjectConfigFile = ".odins"

// RouteConfig represents a single route in a project .odins file.
type RouteConfig struct {
	Subdomain       string `toml:"subdomain"`
	Port            int    `toml:"port"`
	HTTPS           bool   `toml:"https"`
	DockerContainer string `toml:"docker_container,omitempty"`
}

// ProjectInfo holds metadata about the project.
type ProjectInfo struct {
	Name      string `toml:"name"`
	Runtime   string `toml:"runtime,omitempty"`
	Framework string `toml:"framework,omitempty"`
}

// ProjectConfig is the structure of a .odins file.
type ProjectConfig struct {
	Project ProjectInfo   `toml:"project"`
	Routes  []RouteConfig `toml:"routes"`
}

// LoadProject reads a .odins file from the given path.
func LoadProject(path string) (ProjectConfig, error) {
	var cfg ProjectConfig
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}

// SaveProject writes a ProjectConfig to a .odins file.
func SaveProject(path string, cfg ProjectConfig) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// ExistsProject returns true if .odins exists in the given directory.
func ExistsProject(dir string) bool {
	_, err := os.Stat(dir + "/" + ProjectConfigFile)
	return err == nil
}
