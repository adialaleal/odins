package apache

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/adialaleal/odins/internal/state"
	"github.com/adrg/xdg"
	"path/filepath"
)

// Backend is the Apache httpd reverse proxy implementation.
type Backend struct{}

func New() *Backend { return &Backend{} }

func (b *Backend) Name() string { return "apache" }

func (b *Backend) IsInstalled() bool {
	_, err := exec.LookPath("apachectl")
	return err == nil
}

func (b *Backend) IsRunning() bool {
	out, _ := exec.Command("apachectl", "status").Output()
	return len(out) > 0
}

func (b *Backend) Install() error { return nil }

func (b *Backend) Start() error {
	return runBrewService("start", "httpd")
}

func (b *Backend) Stop() error {
	return runBrewService("stop", "httpd")
}

func (b *Backend) Restart() error {
	return runBrewService("restart", "httpd")
}

func (b *Backend) Reload() error {
	out, err := exec.Command("apachectl", "graceful").CombinedOutput()
	if err != nil {
		return fmt.Errorf("apachectl graceful: %w\n%s", err, string(out))
	}
	return nil
}

func (b *Backend) LogPath() string {
	for _, path := range []string{
		"/usr/local/var/log/httpd/access_log",
		"/opt/homebrew/var/log/httpd/access_log",
	} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return filepath.Join(xdg.DataHome, "odins", "apache-access.log")
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

func runBrewService(action, service string) error {
	out, err := exec.Command("brew", "services", action, service).CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew services %s %s: %w\n%s", action, service, err, string(out))
	}
	return nil
}
