package detect

import (
	"os"
	"path/filepath"
	"strings"
)

type pythonResult struct {
	Framework string
	Port      int
	StartCmd  string
}

func detectPython(dir string) *pythonResult {
	hasDjango := fileExists(filepath.Join(dir, "manage.py"))
	hasPyproject := fileExists(filepath.Join(dir, "pyproject.toml"))
	hasRequirements := fileExists(filepath.Join(dir, "requirements.txt"))
	hasSetup := fileExists(filepath.Join(dir, "setup.py"))

	if !hasDjango && !hasPyproject && !hasRequirements && !hasSetup {
		return nil
	}

	r := &pythonResult{Port: 8000, StartCmd: "python -m uvicorn main:app --reload"}

	if hasDjango {
		r.Framework = "django"
		r.Port = 8000
		r.StartCmd = "python manage.py runserver"
		return r
	}

	// Check requirements.txt or pyproject.toml for known frameworks
	deps := readPythonDeps(dir)

	switch {
	case containsAny(deps, "fastapi", "uvicorn"):
		r.Framework = "fastapi"
		r.Port = 8000
		r.StartCmd = "uvicorn main:app --reload"

	case containsAny(deps, "flask"):
		r.Framework = "flask"
		r.Port = 5000
		r.StartCmd = "flask run"

	case containsAny(deps, "starlette"):
		r.Framework = "starlette"
		r.Port = 8000

	case containsAny(deps, "tornado"):
		r.Framework = "tornado"
		r.Port = 8888

	case containsAny(deps, "sanic"):
		r.Framework = "sanic"
		r.Port = 8000

	default:
		r.Framework = "python"
	}

	// Override with .env
	if p := portFromEnvFiles(dir); p > 0 {
		r.Port = p
	}

	return r
}

func readPythonDeps(dir string) []string {
	var deps []string

	// requirements.txt
	if data, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			pkg := strings.Split(strings.TrimSpace(line), "=")[0]
			pkg = strings.Split(pkg, ">")[0]
			pkg = strings.Split(pkg, "<")[0]
			pkg = strings.ToLower(strings.TrimSpace(pkg))
			if pkg != "" && !strings.HasPrefix(pkg, "#") {
				deps = append(deps, pkg)
			}
		}
	}

	// pyproject.toml (simple scan)
	if data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			// Look for dependency entries
			if strings.HasPrefix(line, "\"") || strings.HasPrefix(line, "'") {
				pkg := strings.Split(line, "\"")[1]
				pkg = strings.Split(pkg, "=")[0]
				pkg = strings.Split(pkg, ">")[0]
				pkg = strings.ToLower(strings.TrimSpace(pkg))
				if pkg != "" {
					deps = append(deps, pkg)
				}
			}
		}
	}

	return deps
}

func containsAny(list []string, items ...string) bool {
	for _, item := range items {
		for _, l := range list {
			if strings.Contains(l, item) {
				return true
			}
		}
	}
	return false
}
