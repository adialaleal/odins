package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type nodeResult struct {
	Framework string
	Port      int
	StartCmd  string
}

type packageJSON struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
	Deps    map[string]string `json:"dependencies"`
	DevDeps map[string]string `json:"devDependencies"`
}

func detectNode(dir string) *nodeResult {
	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	r := &nodeResult{Port: 3000, StartCmd: "npm run dev"}

	// Determine start command
	if _, ok := pkg.Scripts["dev"]; ok {
		r.StartCmd = "npm run dev"
	} else if _, ok := pkg.Scripts["start"]; ok {
		r.StartCmd = "npm start"
	}

	// Detect framework
	allDeps := mergeMaps(pkg.Deps, pkg.DevDeps)

	switch {
	case fileExists(filepath.Join(dir, "next.config.js")) ||
		fileExists(filepath.Join(dir, "next.config.ts")) ||
		fileExists(filepath.Join(dir, "next.config.mjs")) ||
		hasDep(allDeps, "next"):
		r.Framework = "nextjs"
		r.Port = 3000
		r.StartCmd = "npm run dev"

	case fileExists(filepath.Join(dir, "nuxt.config.js")) ||
		fileExists(filepath.Join(dir, "nuxt.config.ts")) ||
		hasDep(allDeps, "nuxt"):
		r.Framework = "nuxt"
		r.Port = 3000

	case fileExists(filepath.Join(dir, "vite.config.js")) ||
		fileExists(filepath.Join(dir, "vite.config.ts")) ||
		fileExists(filepath.Join(dir, "vite.config.mjs")) ||
		hasDep(allDeps, "vite"):
		r.Framework = "vite"
		r.Port = 5173
		r.StartCmd = "npm run dev"

	case hasDep(allDeps, "remix") || hasDep(allDeps, "@remix-run/node"):
		r.Framework = "remix"
		r.Port = 3000

	case hasDep(allDeps, "fastify"):
		r.Framework = "fastify"
		r.Port = 3000

	case hasDep(allDeps, "express"):
		r.Framework = "express"
		r.Port = 3000

	case hasDep(allDeps, "hapi") || hasDep(allDeps, "@hapi/hapi"):
		r.Framework = "hapi"
		r.Port = 3000

	case hasDep(allDeps, "koa"):
		r.Framework = "koa"
		r.Port = 3000

	case hasDep(allDeps, "nestjs") || hasDep(allDeps, "@nestjs/core"):
		r.Framework = "nestjs"
		r.Port = 3000

	default:
		r.Framework = "node"
	}

	// Override port from .env files
	if p := portFromEnvFiles(dir); p > 0 {
		r.Port = p
	}

	return r
}

func portFromEnvFiles(dir string) int {
	envFiles := []string{".env.local", ".env.development", ".env"}
	for _, name := range envFiles {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "PORT=") {
				val := strings.TrimPrefix(line, "PORT=")
				val = strings.Trim(val, "\"'")
				if p, err := strconv.Atoi(val); err == nil && p > 0 {
					return p
				}
			}
		}
	}
	return 0
}

func hasDep(deps map[string]string, name string) bool {
	_, ok := deps[name]
	return ok
}

func mergeMaps(a, b map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
