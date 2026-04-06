package cmd

import (
	"fmt"
	"text/tabwriter"
	"os"

	"github.com/adialaleal/odins/internal/docker"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all active routes",
	RunE:    runLs,
}

func runLs(cmd *cobra.Command, args []string) error {
	store, err := state.Load()
	if err != nil {
		return err
	}

	if len(store.Routes) == 0 {
		fmt.Println("  Nenhuma rota ativa. Use 'odins add' ou 'odins up' para adicionar.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\n  STATUS\tSUBDOMAIN\tPORT\tPROTO\tRUNTIME\tPROJECT")
	fmt.Fprintln(w, "  ──────\t─────────\t────\t─────\t───────\t───────")

	for _, r := range store.Routes {
		status := "○"
		if docker.CheckSubdomain(r.Port) {
			status = "●"
		}

		proto := "HTTP"
		if r.HTTPS {
			proto = "HTTPS"
		}

		runtime := r.Runtime
		if r.DockerContainer != "" {
			runtime = "docker"
		}
		if runtime == "" {
			runtime = "—"
		}

		fmt.Fprintf(w, "  %s\t%s\t%d\t%s\t%s\t%s\n",
			status, r.Subdomain, r.Port, proto, runtime, r.Project)
	}

	w.Flush()
	fmt.Println()
	return nil
}
