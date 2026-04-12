package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
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
	resolvedDir, dirWarnings, err := resolveUpDir(upDir, upGlobal)
	if err != nil {
		return err
	}

	manager := serviceFactory()
	result, warnings, err := manager.Up(resolvedDir)
	if err != nil {
		return err
	}
	warnings = append(dirWarnings, warnings...)

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "up", result, warnings)
	}

	if result.GeneratedConfig {
		writeTextLine(cmd.OutOrStdout(), "  → .odins criado em %s", result.ProjectConfigPath)
	} else if result.AutoDetected {
		writeTextLine(cmd.OutOrStdout(), "  → Projeto detectado automaticamente para '%s'", result.Project.Name)
	} else {
		writeTextLine(cmd.OutOrStdout(), "  → Lendo .odins do projeto '%s'", result.Project.Name)
	}
	for _, route := range result.Routes {
		writeTextLine(cmd.OutOrStdout(), "  ✓ %s://%s → :%d", route.Proto, route.Route.Subdomain, route.Route.Port)
	}
	if result.DomainPageURL != "" {
		writeTextLine(cmd.OutOrStdout(), "  → Landing page atualizada: %s", result.DomainPageURL)
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  %d rota(s) ativada(s) para '%s'", len(result.Routes), result.Project.Name)
	return nil
}

func resolveUpDir(dir string, useGlobal bool) (string, []string, error) {
	if useGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", nil, service.InvalidInput("não foi possível resolver o diretório HOME para `odins up --global`")
		}
		return home, nil, nil
	}

	baseDir := dir
	if baseDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", nil, service.InvalidInput("não foi possível resolver o diretório atual")
		}
		baseDir = cwd
	}

	resolvedDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", nil, service.InvalidInput(fmt.Sprintf("não foi possível resolver o diretório %q", baseDir))
	}

	if config.ExistsProject(resolvedDir) {
		return resolvedDir, nil, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return resolvedDir, nil, nil
	}
	if config.ExistsProject(home) {
		return home, []string{"Usando config global em " + filepath.Join(home, config.ProjectConfigFile)}, nil
	}

	return resolvedDir, nil, nil
}
