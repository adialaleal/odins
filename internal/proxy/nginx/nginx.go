package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adialaleal/odins/internal/state"
	"github.com/adrg/xdg"
)

// Backend is the Nginx reverse proxy implementation.
type Backend struct{}

func New() *Backend { return &Backend{} }

func (b *Backend) Name() string { return "nginx" }

func (b *Backend) IsInstalled() bool {
	_, err := exec.LookPath("nginx")
	return err == nil
}

func (b *Backend) IsRunning() bool {
	out, err := exec.Command("brew", "services", "list").Output()
	if err != nil {
		return false
	}
	for _, line := range splitLines(string(out)) {
		if len(line) > 0 && line[0] == "nginx" && line[1] == "started" {
			return true
		}
	}
	return false
}

func (b *Backend) Install() error { return nil }

func (b *Backend) Start() error {
	return runBrewService("start", "nginx")
}

func (b *Backend) Stop() error {
	return runBrewService("stop", "nginx")
}

func (b *Backend) Restart() error {
	return runBrewService("restart", "nginx")
}

func (b *Backend) Reload() error {
	out, err := exec.Command("nginx", "-s", "reload").CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx reload: %w\n%s", err, string(out))
	}
	return nil
}

func (b *Backend) LogPath() string {
	// macOS Homebrew nginx default log location
	for _, path := range []string{
		"/usr/local/var/log/nginx/access.log",
		"/opt/homebrew/var/log/nginx/access.log",
	} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return filepath.Join(xdg.DataHome, "odins", "nginx-access.log")
}

func (b *Backend) AddRoute(r state.Route) error {
	if err := writeVhost(r); err != nil {
		return err
	}
	return b.Reload()
}

func (b *Backend) RemoveRoute(subdomain string) error {
	if err := removeVhost(subdomain); err != nil {
		return err
	}
	return b.Reload()
}

// Init writes the nginx include directive config.
func (b *Backend) Init() error {
	return os.MkdirAll(confDir(), 0755)
}

func runBrewService(action, service string) error {
	out, err := exec.Command("brew", "services", action, service).CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew services %s %s: %w\n%s", action, service, err, string(out))
	}
	return nil
}

func splitLines(s string) [][]string {
	var result [][]string
	for _, line := range splitStr(s, "\n") {
		fields := splitFields(line)
		if len(fields) >= 2 {
			result = append(result, fields)
		}
	}
	return result
}

func splitStr(s, sep string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep[0] {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

func splitFields(s string) []string {
	var out []string
	start := -1
	for i, c := range s {
		if c != ' ' && c != '\t' {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				out = append(out, s[start:i])
				start = -1
			}
		}
	}
	if start != -1 {
		out = append(out, s[start:])
	}
	return out
}
