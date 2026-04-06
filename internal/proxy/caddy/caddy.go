package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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
// If the server config doesn't exist yet (Caddy freshly started), it initialises
// a base config with an empty srv0 first, then appends the route.
func (b *Backend) AddRoute(r state.Route) error {
	upstream := fmt.Sprintf("localhost:%d", r.Port)
	if r.DockerContainer != "" {
		upstream = fmt.Sprintf("%s:%d", r.DockerContainer, r.Port)
	}

	routeID := r.ID
	if routeID == "" {
		routeID = "odins-" + r.Subdomain
	}
	route := buildRoute(r.Subdomain, upstream, routeID)
	data, err := json.Marshal(route)
	if err != nil {
		return err
	}

	// Try to append to existing srv0 routes array.
	if err := b.post("/config/apps/http/servers/srv0/routes", data); err != nil {
		// srv0 likely doesn't exist — bootstrap a base config and retry.
		if initErr := b.initBase(); initErr != nil {
			return fmt.Errorf("caddy init base: %w (original: %v)", initErr, err)
		}
		return b.post("/config/apps/http/servers/srv0/routes", data)
	}
	return nil
}

// initBase pushes a minimal Caddy config that creates srv0.
func (b *Backend) initBase() error {
	cfg := buildBaseConfig("")
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return b.post("/load", data)
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

// AddDomain registers a static landing page route for a domain in Caddy.
// pageDir is the directory containing index.html (e.g. ~/.local/share/odins/pages/tatoh).
func (b *Backend) AddDomain(hostname, pageDir string) error {
	route := buildDomainRoute(hostname, pageDir)
	data, err := json.Marshal(route)
	if err != nil {
		return err
	}
	if err := b.post("/config/apps/http/servers/srv0/routes", data); err != nil {
		if initErr := b.initBase(); initErr != nil {
			return fmt.Errorf("caddy init base: %w (original: %v)", initErr, err)
		}
		return b.post("/config/apps/http/servers/srv0/routes", data)
	}
	return nil
}

// RemoveDomain deletes the landing page route for a domain from Caddy.
func (b *Backend) RemoveDomain(hostname string) error {
	id := "odins-domain-" + hostname
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
		return fmt.Errorf("caddy API delete domain failed (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// Reload applies any pending config changes.
func (b *Backend) Reload() error {
	// Caddy API changes are applied instantly; reload is a no-op here.
	return nil
}

// SyncRoutes rebuilds the entire Caddy config from the provided routes and
// domain page directories. Call this after Caddy restarts to restore state.
// domainPages maps domain name → landing page directory path.
func (b *Backend) SyncRoutes(routes []state.Route, domainPages map[string]string) error {
	allRoutes := make([]interface{}, 0, len(routes)+len(domainPages))

	// Domain landing page routes first (lower priority than service routes)
	for domain, pageDir := range domainPages {
		allRoutes = append(allRoutes, buildDomainRoute(domain, pageDir))
	}
	// Service reverse-proxy routes
	for _, r := range routes {
		upstream := fmt.Sprintf("localhost:%d", r.Port)
		if r.DockerContainer != "" {
			upstream = fmt.Sprintf("%s:%d", r.DockerContainer, r.Port)
		}
		id := r.ID
		if id == "" {
			id = "odins-" + r.Subdomain
		}
		allRoutes = append(allRoutes, buildRoute(r.Subdomain, upstream, id))
	}

	cfg := buildBaseConfig("")
	apps := cfg["apps"].(map[string]interface{})
	http := apps["http"].(map[string]interface{})
	servers := http["servers"].(map[string]interface{})
	srv0 := servers["srv0"].(map[string]interface{})
	srv0["routes"] = allRoutes

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return b.post("/load", data)
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
	cmd := exec.Command(brewBin(), "services", action, service)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func brewBin() string {
	for _, p := range []string{"/opt/homebrew/bin/brew", "/usr/local/bin/brew"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "brew"
}
