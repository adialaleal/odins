package detect

import (
	"os"
	"path/filepath"
)

// DetectedProject contains auto-detected project information.
type DetectedProject struct {
	Name        string
	Runtime     string // "node", "go", "python", "unknown"
	Framework   string // "nextjs", "gin", "fastapi", etc.
	Port        int
	StartCmd    string
	HasDocker   bool
	HasCompose  bool
}

// Project detects the project type and configuration in dir.
func Project(dir string) DetectedProject {
	name := filepath.Base(dir)
	if name == "." || name == "" {
		if abs, err := filepath.Abs(dir); err == nil {
			name = filepath.Base(abs)
		}
	}

	d := DetectedProject{
		Name:    name,
		Runtime: "unknown",
		Port:    8080,
	}

	d.HasDocker = fileExists(filepath.Join(dir, "Dockerfile"))
	d.HasCompose = fileExists(filepath.Join(dir, "docker-compose.yml")) ||
		fileExists(filepath.Join(dir, "docker-compose.yaml")) ||
		fileExists(filepath.Join(dir, "compose.yml"))

	if nd := detectNode(dir); nd != nil {
		d.Runtime = "node"
		d.Framework = nd.Framework
		d.Port = nd.Port
		d.StartCmd = nd.StartCmd
		return d
	}

	if gd := detectGo(dir); gd != nil {
		d.Runtime = "go"
		d.Framework = gd.Framework
		d.Port = gd.Port
		d.StartCmd = gd.StartCmd
		return d
	}

	if pd := detectPython(dir); pd != nil {
		d.Runtime = "python"
		d.Framework = pd.Framework
		d.Port = pd.Port
		d.StartCmd = pd.StartCmd
		return d
	}

	return d
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
