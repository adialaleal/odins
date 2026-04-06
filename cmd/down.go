package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/i18n"
	"github.com/adialaleal/odins/internal/proxy/apache"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/proxy/nginx"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Remove all routes for the current project",
	Long: `Read the .odins file in the current directory and remove all its routes.

Examples:
  cd ~/Projects/rankly && odins down`,
	RunE: runDown,
}

func runDown(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	projectCfgPath := filepath.Join(dir, config.ProjectConfigFile)
	if !config.ExistsProject(dir) {
		return fmt.Errorf("%s", i18n.Tf("down.no_project", dir))
	}

	projCfg, err := config.LoadProject(projectCfgPath)
	if err != nil {
		return err
	}

	globalCfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	store, err := state.Load()
	if err != nil {
		return err
	}

	domain := projCfg.Project.Domain
	if domain == "" {
		domain = projCfg.Project.Name
	}

	removed := 0
	for _, rc := range projCfg.Routes {
		fqdn := buildFQDN(rc.Subdomain, domain, projCfg.Project.Name, globalCfg.TLD)
		if err := proxyRemove(globalCfg, fqdn); err != nil {
			fmt.Println("  " + i18n.Tf("down.proxy_warn", fqdn, err))
		}
		store.Remove(fqdn)
		fmt.Println("  " + i18n.Tf("down.removed", fqdn))
		removed++
	}

	if err := store.Save(); err != nil {
		return err
	}

	// Regenerate landing page if project belonged to a domain
	if domain != "" {
		regeneratePageForDomain(globalCfg, store, domain)
		fmt.Println("  " + i18n.Tf("down.page_updated", domain, globalCfg.TLD))
	}

	fmt.Println()
	fmt.Println("  " + i18n.Tf("down.total", removed, projCfg.Project.Name))
	return nil
}

func proxyRemove(cfg config.GlobalConfig, subdomain string) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return nginx.New().RemoveRoute(subdomain)
	case config.BackendApache:
		return apache.New().RemoveRoute(subdomain)
	default:
		return caddy.New().RemoveRoute(subdomain)
	}
}
