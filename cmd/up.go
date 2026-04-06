package cmd

import (
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply routes from .odins config in the current directory",
	Long: `Read the .odins file in the current directory and apply all routes.

If no .odins file exists, ODINS will auto-detect the project type (Node.js,
Go, Python) and create one automatically.

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
	manager := serviceFactory()
	result, warnings, err := manager.Up(upDir)
	if err != nil {
		return err
	}

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
