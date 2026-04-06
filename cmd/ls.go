package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all active routes",
	RunE:    runLs,
}

func runLs(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.ListRoutes()
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "ls", result, warnings)
	}

	if result.Count == 0 {
		writeTextLine(cmd.OutOrStdout(), "  Nenhuma rota ativa. Use 'odins add' ou 'odins up' para adicionar.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\n  STATUS\tSUBDOMAIN\tPORT\tPROTO\tRUNTIME\tPROJECT")
	fmt.Fprintln(w, "  ──────\t─────────\t────\t─────\t───────\t───────")
	for _, route := range result.Routes {
		status := "○"
		if route.Up {
			status = "●"
		}
		fmt.Fprintf(w, "  %s\t%s\t%d\t%s\t%s\t%s\n",
			status,
			route.Route.Subdomain,
			route.Route.Port,
			route.Proto,
			route.Runtime,
			route.Route.Project,
		)
	}
	w.Flush()
	writeTextLine(cmd.OutOrStdout(), "")
	return nil
}
