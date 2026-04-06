package cmd

import (
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <subdomain>",
	Short: "Remove a route",
	Long: `Remove a route by its subdomain.

Examples:
  odins kill api.rankly.odin
  odins kill app.rankly.odin`,
	Args: exactArgs(1),
	RunE: runKill,
}

func runKill(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.Kill(args[0])
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "kill", result, warnings)
	}

	writeTextLine(cmd.OutOrStdout(), "  ✓ %s removido", result.Subdomain)
	if result.DomainPageURL != "" {
		writeTextLine(cmd.OutOrStdout(), "  → Landing page atualizada: %s", result.DomainPageURL)
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	return nil
}
