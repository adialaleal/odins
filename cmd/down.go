package cmd

import (
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
	manager := serviceFactory()
	result, warnings, err := manager.Down("")
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "down", result, warnings)
	}

	for _, route := range result.RemovedRoutes {
		writeTextLine(cmd.OutOrStdout(), "  ✓ %s removido", route.Subdomain)
	}
	if result.DomainPageURL != "" {
		writeTextLine(cmd.OutOrStdout(), "  → Landing page atualizada: %s", result.DomainPageURL)
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  %d rota(s) removida(s) para '%s'", len(result.RemovedRoutes), result.Project)
	return nil
}
