package detect

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type goResult struct {
	Framework string
	Port      int
	StartCmd  string
}

var portRegex = regexp.MustCompile(`[":]([\d]{4,5})[")\s]`)

func detectGo(dir string) *goResult {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); os.IsNotExist(err) {
		return nil
	}

	r := &goResult{Port: 8080, StartCmd: "go run ."}

	// Read go.mod to identify frameworks
	modData, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return r
	}
	modStr := string(modData)

	switch {
	case strings.Contains(modStr, "github.com/gin-gonic/gin"):
		r.Framework = "gin"
		r.Port = 8080

	case strings.Contains(modStr, "github.com/labstack/echo"):
		r.Framework = "echo"
		r.Port = 1323

	case strings.Contains(modStr, "github.com/gofiber/fiber"):
		r.Framework = "fiber"
		r.Port = 3000

	case strings.Contains(modStr, "github.com/gorilla/mux"):
		r.Framework = "gorilla"
		r.Port = 8080

	case strings.Contains(modStr, "github.com/go-chi/chi"):
		r.Framework = "chi"
		r.Port = 8080

	default:
		r.Framework = "go"
	}

	// Try to find port in main.go
	if p := portFromGoFile(filepath.Join(dir, "main.go")); p > 0 {
		r.Port = p
	}

	// Override with .env
	if p := portFromEnvFiles(dir); p > 0 {
		r.Port = p
	}

	return r
}

func portFromGoFile(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	content := string(data)

	// Look for common port patterns: ":8080", Listen(":3000"), etc.
	patterns := []string{
		`ListenAndServe\(":(\d+)"`,
		`\.Listen\(":(\d+)"`,
		`\.Run\(":(\d+)"`,
		`\.Start\(":(\d+)"`,
		`"PORT",\s*"(\d+)"`,
	}

	for _, pat := range patterns {
		re := regexp.MustCompile(pat)
		if m := re.FindStringSubmatch(content); m != nil {
			if p, err := strconv.Atoi(m[1]); err == nil {
				return p
			}
		}
	}

	// Fallback: find any 4-5 digit port in strings
	if m := portRegex.FindStringSubmatch(content); m != nil {
		if p, err := strconv.Atoi(m[1]); err == nil && p >= 1024 && p <= 65535 {
			return p
		}
	}

	return 0
}
