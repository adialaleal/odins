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
	return exec.Command("pgrep", "-x", "nginx").Run() == nil
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

