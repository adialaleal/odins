package cmd

import (
	"fmt"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/i18n"
	"github.com/adialaleal/odins/internal/page"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage local domains (project workspaces)",
	Long: `Create and manage local domains. A domain is a workspace that groups services.

Example:
  odins domain add myproject                  # creates myproject.odins (landing page)
  odins domain add myproject --title "MyApp"  # with custom title
  odins domain ls                              # list all domains
  odins domain rm myproject                    # remove domain and its services`,
}

var (
	domainTitle string
	domainDesc  string
)

func init() {
	domainCmd.AddCommand(domainAddCmd, domainLsCmd, domainRmCmd)
	domainAddCmd.Flags().StringVar(&domainTitle, "title", "", "Title displayed on the landing page")
	domainAddCmd.Flags().StringVar(&domainDesc, "desc", "", "Domain description")
}

// ─── odins domain add ────────────────────────────────────────────────────────

var domainAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new local domain with landing page",
	Args:  cobra.ExactArgs(1),
	RunE:  runDomainAdd,
}

func runDomainAdd(cmd *cobra.Command, args []string) error {
	name := strings.ToLower(args[0])

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	store, err := state.Load()
	if err != nil {
		return err
	}

	if _, exists := store.GetDomain(name); exists {
		return fmt.Errorf("%s", i18n.Tf("domain.exists", name))
	}

	title := domainTitle
	if title == "" {
		title = name
	}

	d := state.Domain{
		Name:        name,
		Title:       title,
		Description: domainDesc,
	}
	store.AddDomain(d)
	if err := store.Save(); err != nil {
		return err
	}

	hostname := name + "." + cfg.TLD
	pageDir := page.PageDir(name)

	// Generate landing page (empty, no routes yet)
	if err := page.Generate(page.PageData{
		Domain:      name,
		TLD:         cfg.TLD,
		Title:       title,
		Description: domainDesc,
	}); err != nil {
		fmt.Println("  " + i18n.Tf("domain.page_warn", err))
	}

	// Register with Caddy
	caddyClient := caddy.New()
	if err := caddyClient.AddDomain(hostname, pageDir); err != nil {
		fmt.Println("  " + i18n.Tf("domain.caddy_warn", err))
		fmt.Println("     " + i18n.T("domain.caddy_hint"))
	}

	fmt.Println()
	fmt.Println("  " + i18n.Tf("domain.created", hostname))
	fmt.Println("  " + i18n.Tf("domain.page_generated", pageDir))
	fmt.Println()
	fmt.Println("  " + i18n.T("domain.add_service"))
	fmt.Printf("    [project]\n")
	fmt.Printf("    domain = \"%s\"\n", name)
	fmt.Println()
	fmt.Println("  " + i18n.T("domain.then_up"))
	fmt.Println()
	return nil
}

// ─── odins domain ls ─────────────────────────────────────────────────────────

var domainLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List registered domains",
	RunE:  runDomainLs,
}

func runDomainLs(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	store, err := state.Load()
	if err != nil {
		return err
	}

	if len(store.Domains) == 0 {
		fmt.Println("  " + i18n.T("domain.empty"))
		fmt.Println("  " + i18n.T("domain.empty_hint"))
		return nil
	}

	fmt.Println()
	fmt.Printf("  %-20s  %-30s  %s\n",
		i18n.T("domain.header.domain"),
		i18n.T("domain.header.fqdn"),
		i18n.T("domain.header.services"))
	fmt.Printf("  %s\n", strings.Repeat("─", 70))

	for _, d := range store.Domains {
		hostname := d.Name + "." + cfg.TLD
		routes := store.ByDomain(d.Name)
		fmt.Printf("  %-20s  %-30s  %d\n", d.Name, hostname, len(routes))
	}
	fmt.Println()
	return nil
}

// ─── odins domain rm ─────────────────────────────────────────────────────────

var domainRmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a domain and all its services",
	Args:  cobra.ExactArgs(1),
	RunE:  runDomainRm,
}

func runDomainRm(cmd *cobra.Command, args []string) error {
	name := strings.ToLower(args[0])

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	store, err := state.Load()
	if err != nil {
		return err
	}

	if _, exists := store.GetDomain(name); !exists {
		return fmt.Errorf("%s", i18n.Tf("domain.not_found", name))
	}

	routes := store.ByDomain(name)
	hostname := name + "." + cfg.TLD

	// Remove each route from proxy
	for _, r := range routes {
		if err := proxyRemove(cfg, r.Subdomain); err != nil {
			fmt.Println("  " + i18n.Tf("domain.proxy_warn", r.Subdomain, err))
		}
	}

	// Remove domain route from Caddy
	caddyClient := caddy.New()
	if err := caddyClient.RemoveDomain(hostname); err != nil {
		fmt.Println("  " + i18n.Tf("domain.caddy_rm_warn", err))
	}

	// Remove from store (also removes all routes in that domain)
	store.RemoveDomain(name)
	if err := store.Save(); err != nil {
		return err
	}

	fmt.Println("  " + i18n.Tf("domain.removed", name, len(routes)))
	return nil
}

// regeneratePageForDomain regenerates the landing page for a domain after route changes.
func regeneratePageForDomain(cfg config.GlobalConfig, store *state.Store, domainName string) {
	d, ok := store.GetDomain(domainName)
	if !ok {
		return
	}

	routes := store.ByDomain(domainName)
	var routeInfos []page.RouteInfo
	for _, r := range routes {
		routeInfos = append(routeInfos, page.RouteInfo{
			Subdomain: extractSubdomain(r.Subdomain, domainName, cfg.TLD),
			FQDN:      r.Subdomain,
			Port:      r.Port,
			Runtime:   r.Runtime,
			Project:   r.Project,
		})
	}

	_ = page.Generate(page.PageData{
		Domain:      domainName,
		TLD:         cfg.TLD,
		Title:       d.Title,
		Description: d.Description,
		Routes:      routeInfos,
	})
}

// extractSubdomain strips the domain+tld suffix to get just the subdomain part.
// e.g. "web.project.odins" with domain="project" tld="odins" → "web"
func extractSubdomain(fqdn, domain, tld string) string {
	suffix := "." + domain + "." + tld
	if len(fqdn) > len(suffix) && fqdn[len(fqdn)-len(suffix):] == suffix {
		return fqdn[:len(fqdn)-len(suffix)]
	}
	return fqdn
}
