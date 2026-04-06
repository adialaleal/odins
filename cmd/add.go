package cmd

import (
	"github.com/adialaleal/odins/internal/service"
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
	Args: exactArgs(1),
	RunE: runAdd,
}

var (
	addPort    int
	addDocker  string
	addProject string
	addDomain  string
	addNoHTTPS bool
)

func init() {
	addCmd.Flags().IntVarP(&addPort, "port", "p", 0, "Local port to proxy (required)")
	addCmd.Flags().StringVarP(&addDocker, "docker", "d", "", "Docker container name")
	addCmd.Flags().StringVar(&addProject, "project", "", "Project name (inferred from subdomain if not set)")
	addCmd.Flags().StringVar(&addDomain, "domain", "", "Domain workspace (e.g. tatoh → tatoh.odins)")
	addCmd.Flags().BoolVar(&addNoHTTPS, "no-https", false, "Disable HTTPS for this route")
	addCmd.MarkFlagRequired("port")
}

func runAdd(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.AddRoute(service.AddRouteOptions{
		Subdomain: args[0],
		Port:      addPort,
		Docker:    addDocker,
		Project:   addProject,
		Domain:    addDomain,
		HTTPS:     !addNoHTTPS,
	})
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "add", result, warnings)
	}

	proto := "https"
	if !result.Route.HTTPS {
		proto = "http"
	}
	writeTextLine(cmd.OutOrStdout(), "  ✓ %s://%s → localhost:%d", proto, result.Route.Subdomain, result.Route.Port)
	if result.DomainPageURL != "" {
		writeTextLine(cmd.OutOrStdout(), "  → Landing page atualizada: %s", result.DomainPageURL)
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	return nil
}
