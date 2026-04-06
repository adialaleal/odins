package cmd

import (
	"fmt"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <subdomain>",
	Short: "Remove a route",
	Long: `Remove a route by its subdomain.

Examples:
  odins kill api.rankly.odin
  odins kill app.rankly.odin`,
	Args: cobra.ExactArgs(1),
	RunE: runKill,
}

func runKill(cmd *cobra.Command, args []string) error {
	subdomain := args[0]

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	store, err := state.Load()
	if err != nil {
		return err
	}

	route, ok := store.Get(subdomain)
	if !ok {
		return fmt.Errorf("rota '%s' não encontrada", subdomain)
	}

	domainName := route.Domain

	if err := proxyRemove(cfg, subdomain); err != nil {
		fmt.Printf("  ⚠  proxy remove: %v\n", err)
	}

	store.Remove(subdomain)
	if err := store.Save(); err != nil {
		return err
	}

	// Regenerate landing page if the route belonged to a domain
	if domainName != "" {
		regeneratePageForDomain(cfg, store, domainName)
	}

	fmt.Printf("  ✓ %s removido\n", subdomain)
	return nil
}
