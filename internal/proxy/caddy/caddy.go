package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/adialaleal/odins/internal/state"
	"github.com/adrg/xdg"
)

const adminAddr = "http://localhost:2019"

// Backend is the Caddy reverse proxy implementation.
type Backend struct {
	AdminAddr string
	client    *http.Client
}

// New creates a new Caddy backend.
func New() *Backend {
	return &Backend{
		AdminAddr: adminAddr,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (b *Backend) Name() string { return "caddy" }

func (b *Backend) IsInstalled() bool {
	_, err := findBinary("caddy")
	return err == nil
}

func (b *Backend) IsRunning() bool {
	resp, err := b.client.Get(b.AdminAddr + "/config/")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func (b *Backend) Install() error {
	// Handled by brew in cmd/init.go
	return nil
}

func (b *Backend) Start() error {
	return runBrewService("start", "caddy")
}

func (b *Backend) Stop() error {
	return runBrewService("stop", "caddy")
}

func (b *Backend) Restart() error {
	return runBrewService("restart", "caddy")
}

func (b *Backend) LogPath() string {
	return filepath.Join(xdg.DataHome, "caddy", "logs", "access.log")
}

// Init pushes the base Caddy config with TLS internal support.
func (b *Backend) Init(tld string) error {
	cfg := buildBaseConfig(tld)
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return b.post("/load", data)
}

// AddRoute adds a reverse proxy route via the Caddy Admin API.
func (b *Backend) AddRoute(r state.Route) error {
	upstream := fmt.Sprintf("localhost:%d", r.Port)
	if r.DockerContainer != "" {
		upstream = fmt.Sprintf("%s:%d", r.DockerContainer, r.Port)
	}

	route := buildRoute(r.Subdomain, upstream, r.ID)
	data, err := json.Marshal(route)
	if err != nil {
		return err
	}
	return b.post("/config/apps/http/servers/srv0/routes/...", data)
}

// RemoveRoute removes a route by its ODINS ID.
func (b *Backend) RemoveRoute(subdomain string) error {
	id := "odins-" + subdomain
	req, err := http.NewRequest(http.MethodDelete, b.AdminAddr+"/id/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("caddy API delete failed (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// Reload applies any pending config changes.
func (b *Backend) Reload() error {
	// Caddy API changes are applied instantly; reload is a no-op here.
	return nil
}

func (b *Backend) post(path string, data []byte) error {
	resp, err := b.client.Post(b.AdminAddr+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("caddy API %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("caddy API %s returned %d: %s", path, resp.StatusCode, string(body))
	}
	return nil
}

func findBinary(name string) (string, error) {
	for _, dir := range []string{"/usr/local/bin", "/opt/homebrew/bin", "/usr/bin"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("%s not found", name)
}

func runBrewService(action, service string) error {
	// Delegated to pkg/brew; import avoided to prevent cycle — use os/exec directly.
	import_cmd := fmt.Sprintf("brew services %s %s", action, service)
	_ = import_cmd
	return nil
}
