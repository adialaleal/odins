package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
)

// ScanOptions configures a directory scan.
type ScanOptions struct {
	Directory   string // root dir (default ~/Projects)
	MaxDepth    int    // max depth to walk (default 3)
	CreateOdins bool   // create .odins files where missing
}

// ScannedProject represents a detected project in the scan results.
type ScannedProject struct {
	Path      string             `json:"path"`
	Name      string             `json:"name"`
	Runtime   string             `json:"runtime"`
	Framework string             `json:"framework"`
	Port      int                `json:"port"`
	StartCmd  string             `json:"start_cmd,omitempty"`
	HasOdins  bool               `json:"has_odins"`
	HasDocker bool               `json:"has_docker"`
	Domain    string             `json:"domain,omitempty"`
	Routes    []config.RouteConfig `json:"routes,omitempty"`
}

// ScanResult is returned by Scan.
type ScanResult struct {
	RootDirectory string           `json:"root_directory"`
	Projects      []ScannedProject `json:"projects"`
	Created       int              `json:"created"`
}

// skipDirs are directories to skip during walk.
var skipDirs = map[string]bool{
	"node_modules":  true,
	".git":          true,
	"vendor":        true,
	".venv":         true,
	"__pycache__":   true,
	"dist":          true,
	"build":         true,
	".next":         true,
	".nuxt":         true,
	".cache":        true,
	".output":       true,
	"target":        true,
	".terraform":    true,
	".gradle":       true,
	".idea":         true,
	".vscode":       true,
}

// projectIndicators are files that signal a project root.
var projectIndicators = []string{
	"package.json",
	"go.mod",
	"pyproject.toml",
	"requirements.txt",
	"Cargo.toml",
	"Dockerfile",
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	".odins",
}

// Scan walks a directory tree and detects projects.
func Scan(opts ScanOptions) (ScanResult, error) {
	dir := opts.Directory
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ScanResult{}, err
		}
		dir = filepath.Join(home, "Projects")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return ScanResult{}, err
	}

	maxDepth := opts.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 3
	}

	result := ScanResult{RootDirectory: dir}

	// Track detected project paths so we skip their subtrees.
	detectedPaths := make(map[string]bool)

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible dirs
		}

		if !d.IsDir() {
			return nil
		}

		name := d.Name()

		// Skip hidden dirs (except the root itself) and known junk dirs.
		if path != dir {
			if skipDirs[name] || (strings.HasPrefix(name, ".") && name != ".") {
				return filepath.SkipDir
			}
		}

		// Enforce max depth.
		rel, _ := filepath.Rel(dir, path)
		depth := 0
		if rel != "." {
			depth = strings.Count(rel, string(os.PathSeparator)) + 1
		}
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Skip subtrees of already-detected projects.
		for detected := range detectedPaths {
			if path != detected && strings.HasPrefix(path, detected+string(os.PathSeparator)) {
				return filepath.SkipDir
			}
		}

		// Check if this dir looks like a project root.
		if !isProjectRoot(path) {
			return nil
		}

		detected := detect.Project(path)
		if detected.Runtime == "unknown" && !detected.HasDocker && !detected.HasCompose {
			return nil
		}

		detectedPaths[path] = true

		hasOdins := config.ExistsProject(path)
		project := ScannedProject{
			Path:      path,
			Name:      detected.Name,
			Runtime:   detected.Runtime,
			Framework: detected.Framework,
			Port:      detected.Port,
			StartCmd:  detected.StartCmd,
			HasOdins:  hasOdins,
			HasDocker: detected.HasDocker,
		}

		// If .odins exists, load it for domain/routes info.
		if hasOdins {
			cfg, err := config.LoadProject(filepath.Join(path, config.ProjectConfigFile))
			if err == nil {
				project.Domain = cfg.Project.Domain
				project.Routes = cfg.Routes
			}
		}

		// Optionally create .odins where missing.
		if opts.CreateOdins && !hasOdins && detected.Runtime != "unknown" {
			cfg := recommendedConfig(detected)
			cfgPath := filepath.Join(path, config.ProjectConfigFile)
			if err := config.SaveProject(cfgPath, cfg); err == nil {
				project.HasOdins = true
				project.Domain = cfg.Project.Domain
				project.Routes = cfg.Routes
				result.Created++
			}
		}

		result.Projects = append(result.Projects, project)
		return filepath.SkipDir // skip subtree of detected project
	})

	return result, nil
}

func isProjectRoot(dir string) bool {
	for _, indicator := range projectIndicators {
		if _, err := os.Stat(filepath.Join(dir, indicator)); err == nil {
			return true
		}
	}
	return false
}

// recommendedConfig builds a default .odins config from detection results.
func recommendedConfig(d detect.DetectedProject) config.ProjectConfig {
	https := true
	subdomain := d.Name
	// Normalize subdomain: replace underscores with hyphens for DNS compatibility.
	subdomain = strings.ReplaceAll(subdomain, "_", "-")

	return config.ProjectConfig{
		Project: config.ProjectInfo{
			Name:      d.Name,
			Runtime:   d.Runtime,
			Framework: d.Framework,
		},
		Routes: []config.RouteConfig{
			{
				Subdomain: subdomain,
				Port:      d.Port,
				HTTPS:     https,
			},
		},
	}
}
