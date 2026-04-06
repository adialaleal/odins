package cmd

import (
	"fmt"
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/proxy/nginx"
	"github.com/adialaleal/odins/internal/proxy/apache"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <subdomain>",
	Short: "Add a reverse proxy route",
	Long: `Add a route from a local subdomain to a running service.

Examples:
  odins add api.rankly.odin --port 3000
  odins add app.rankly.odin --port 5173
  odins add worker.rankly.odin --port 4000 --docker rankly_worker_1`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

var (
	addPort    int
	addDocker  string
	addProject string
	addNoHTTPS bool
)

func init() {
	addCmd.Flags().IntVarP(&addPort, "port", "p", 0, "Local port to proxy (required)")
	addCmd.Flags().StringVarP(&addDocker, "docker", "d", "", "Docker container name")
	addCmd.Flags().StringVar(&addProject, "project", "", "Project name (inferred from subdomain if not set)")
	addCmd.Flags().BoolVar(&addNoHTTPS, "no-https", false, "Disable HTTPS for this route")
	addCmd.MarkFlagRequired("port")
}

func runAdd(cmd *cobra.Command, args []string) error {
	subdomain := args[0]

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	project := addProject
	if project == "" {
		// Infer project from subdomain (second segment)
		parts := splitDomain(subdomain)
		if len(parts) >= 2 {
			project = parts[1]
		} else {
			project = parts[0]
		}
	}

	r := state.Route{
		Subdomain:       subdomain,
		Port:            addPort,
		Project:         project,
		DockerContainer: addDocker,
		HTTPS:           !addNoHTTPS,
		CreatedAt:       time.Now(),
	}

	// Add to proxy
	if err := proxyAdd(cfg, r); err != nil {
		return fmt.Errorf("proxy add: %w", err)
	}

	// Persist to state
	store, err := state.Load()
	if err != nil {
		return err
	}
	store.Add(r)
	if err := store.Save(); err != nil {
		return err
	}

	proto := "https"
	if !r.HTTPS {
		proto = "http"
	}
	fmt.Printf("  ✓ %s://%s → localhost:%d\n", proto, subdomain, addPort)
	return nil
}

func proxyAdd(cfg config.GlobalConfig, r state.Route) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return nginx.New().AddRoute(r)
	case config.BackendApache:
		return apache.New().AddRoute(r)
	default:
		return caddy.New().AddRoute(r)
	}
}

func splitDomain(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			if i > start {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}
