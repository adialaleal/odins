package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply routes from .odins config in the current directory",
	Long: `Read the .odins file in the current directory and apply all routes.

If no .odins file exists, ODINS will auto-detect the project type (Node.js,
Go, Python) and create one interactively.

Examples:
  cd ~/Projects/rankly && odins up
  odins up --dir ~/Projects/api`,
	RunE: runUp,
}

var upDir string

func init() {
	upCmd.Flags().StringVar(&upDir, "dir", "", "Directory to use instead of current")
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

	projectCfgPath := filepath.Join(dir, config.ProjectConfigFile)

	var projCfg config.ProjectConfig

	if config.ExistsProject(dir) {
		// Load existing .odins
		projCfg, err = config.LoadProject(projectCfgPath)
		if err != nil {
			return fmt.Errorf("read .odins: %w", err)
		}
		fmt.Printf("  → Lendo .odins do projeto '%s'\n", projCfg.Project.Name)
	} else {
		// Auto-detect project
		fmt.Println("  → .odins não encontrado, detectando projeto...")
		d := detect.Project(dir)

		if d.Runtime == "unknown" {
			return fmt.Errorf("não foi possível detectar o tipo de projeto em %s\n"+
				"Crie um .odins manualmente ou use: odins add <subdomain> --port <port>", dir)
		}

		fmt.Printf("  → Detectado: %s/%s (porta %d)\n", d.Runtime, d.Framework, d.Port)
		fmt.Printf("  → Comando de start: %s\n", d.StartCmd)

		// Build default project config
		projCfg = config.ProjectConfig{
			Project: config.ProjectInfo{
				Name:      d.Name,
				Runtime:   d.Runtime,
				Framework: d.Framework,
			},
			Routes: []config.RouteConfig{
				{
					Subdomain: d.Name,
					Port:      d.Port,
					HTTPS:     true,
				},
			},
		}

		// Save the generated .odins
		if err := config.SaveProject(projectCfgPath, projCfg); err != nil {
			fmt.Printf("  ⚠  Não foi possível salvar .odins: %v\n", err)
		} else {
			fmt.Printf("  → .odins criado em %s\n", projectCfgPath)
		}
	}

	// Apply routes
	store, err := state.Load()
	if err != nil {
		return err
	}

	applied := 0
	for _, rc := range projCfg.Routes {
		fqdn := buildFQDN(rc.Subdomain, projCfg.Project.Name, cfg.TLD)

		r := state.Route{
			Subdomain:       fqdn,
			Port:            rc.Port,
			Project:         projCfg.Project.Name,
			Runtime:         projCfg.Project.Runtime,
			DockerContainer: rc.DockerContainer,
			HTTPS:           rc.HTTPS,
			CreatedAt:       time.Now(),
		}

		if err := proxyAdd(cfg, r); err != nil {
			fmt.Printf("  ✗ %s: %v\n", fqdn, err)
			continue
		}

		store.Add(r)

		proto := "https"
		if !r.HTTPS {
			proto = "http"
		}
		fmt.Printf("  ✓ %s://%s → :%d\n", proto, fqdn, rc.Port)
		applied++
	}

	if err := store.Save(); err != nil {
		return err
	}

	fmt.Printf("\n  %d rota(s) ativada(s) para '%s'\n", applied, projCfg.Project.Name)
	return nil
}

// buildFQDN constructs the full subdomain.
// If the route subdomain already contains a dot (e.g. "api.rankly"),
// it is used as-is with TLD appended. Otherwise: subdomain.project.tld.
func buildFQDN(subdomain, project, tld string) string {
	for _, c := range subdomain {
		if c == '.' {
			return subdomain + "." + tld
		}
	}
	return subdomain + "." + project + "." + tld
}
