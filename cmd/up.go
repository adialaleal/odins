package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/i18n"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply routes from .odins config in the current directory",
	Long: `Read the .odins file in the current directory and apply all routes.

If no .odins exists in the current directory, ODINS checks $HOME/.odins as
a global config. If neither exists, it auto-detects the project type.

Examples:
  cd ~/Projects/rankly && odins up
  odins up --dir ~/Projects/api
  odins up --global     # explicitly read $HOME/.odins`,
	RunE: runUp,
}

var (
	upDir    string
	upGlobal bool
)

func init() {
	upCmd.Flags().StringVar(&upDir, "dir", "", "Directory to use instead of current")
	upCmd.Flags().BoolVar(&upGlobal, "global", false, "Read $HOME/.odins as global config")
}

func runUp(cmd *cobra.Command, args []string) error {
	dir := upDir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	// Resolve which directory to use for the .odins file:
	// 1. --global flag → always use $HOME
	// 2. No .odins in CWD → fall back to $HOME/.odins if it exists
	// 3. Otherwise → use CWD (existing behaviour)
	resolvedDir := dir
	if upGlobal {
		home, _ := os.UserHomeDir()
		resolvedDir = home
	} else if !config.ExistsProject(dir) {
		if home, err := os.UserHomeDir(); err == nil && config.ExistsProject(home) {
			resolvedDir = home
			fmt.Printf("  → Usando config global em %s\n", filepath.Join(home, config.ProjectConfigFile))
		}
	}

	projectCfgPath := filepath.Join(resolvedDir, config.ProjectConfigFile)

	var projCfg config.ProjectConfig

	if config.ExistsProject(resolvedDir) {
		// Load existing .odins
		projCfg, err = config.LoadProject(projectCfgPath)
		if err != nil {
			return fmt.Errorf("read .odins: %w", err)
		}
		fmt.Println("  " + i18n.Tf("up.reading", projCfg.Project.Name))
	} else {
		// Auto-detect project in the original CWD (not home)
		fmt.Println("  " + i18n.T("up.detecting"))
		d := detect.Project(dir)

		if d.Runtime == "unknown" {
			return fmt.Errorf("%s", i18n.Tf("up.not_detected", dir))
		}

		fmt.Println("  " + i18n.Tf("up.detected", d.Runtime, d.Framework, d.Port))
		fmt.Println("  " + i18n.Tf("up.start_cmd", d.StartCmd))

		// Build default project config — project name IS the domain
		projCfg = config.ProjectConfig{
			Project: config.ProjectInfo{
				Name:      d.Name,
				Runtime:   d.Runtime,
				Framework: d.Framework,
				Domain:    d.Name,
			},
			Routes: []config.RouteConfig{
				{
					Subdomain: d.Name,
					Port:      d.Port,
					HTTPS:     true,
				},
			},
		}

		// Save the generated .odins to original CWD
		savePath := filepath.Join(dir, config.ProjectConfigFile)
		if err := config.SaveProject(savePath, projCfg); err != nil {
			fmt.Println("  " + i18n.Tf("up.save_warn", err))
		} else {
			fmt.Println("  " + i18n.Tf("up.created", savePath))
		}
	}

	// Apply routes
	store, err := state.Load()
	if err != nil {
		return err
	}

	domain := projCfg.Project.Domain
	// If no explicit domain, use project name (project.odins style)
	if domain == "" {
		domain = projCfg.Project.Name
	}

	if domain != "" {
		fmt.Println("  " + i18n.Tf("up.domain", domain, cfg.TLD))
	}

	applied := 0
	for _, rc := range projCfg.Routes {
		fqdn := buildFQDN(rc.Subdomain, domain, projCfg.Project.Name, cfg.TLD)

		r := state.Route{
			ID:              "odins-" + fqdn,
			Subdomain:       fqdn,
			Port:            rc.Port,
			Project:         projCfg.Project.Name,
			Runtime:         projCfg.Project.Runtime,
			Domain:          domain,
			DockerContainer: rc.DockerContainer,
			HTTPS:           rc.HTTPS,
			CreatedAt:       time.Now(),
		}

		if err := proxyAdd(cfg, r); err != nil {
			fmt.Println("  " + i18n.Tf("up.route_error", fqdn, err))
			continue
		}

		store.Add(r)

		proto := "https"
		if !r.HTTPS {
			proto = "http"
		}
		fmt.Println("  " + i18n.Tf("up.route_ok", proto, fqdn, rc.Port))
		applied++
	}

	if err := store.Save(); err != nil {
		return err
	}

	// Regenerate landing page if attached to a domain
	if domain != "" {
		regeneratePageForDomain(cfg, store, domain)
		fmt.Println("  " + i18n.Tf("up.page_updated", domain, cfg.TLD))
	}

	fmt.Println()
	fmt.Println("  " + i18n.Tf("up.applied", applied, projCfg.Project.Name))
	return nil
}

// buildFQDN constructs the full FQDN for a route.
//   - domain set:          subdomain.domain.tld   (e.g. web.project.odins)
//   - subdomain has dot:   subdomain.tld           (e.g. api.rankly.odins)
//   - otherwise:           subdomain.project.tld
func buildFQDN(subdomain, domain, project, tld string) string {
	if domain != "" {
		return subdomain + "." + domain + "." + tld
	}
	for _, c := range subdomain {
		if c == '.' {
			return subdomain + "." + tld
		}
	}
	return subdomain + "." + project + "." + tld
}
