package cmd

import (
	"fmt"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/page"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Gerenciar domínios locais (workspaces de projeto)",
	Long: `Cria e gerencia domínios locais. Um domínio é um workspace que agrupa serviços.

Exemplo:
  odins domain add tatoh                  # cria tatoh.odins (landing page)
  odins domain add tatoh --title "Tatoh"  # com título personalizado
  odins domain ls                          # lista todos os domínios
  odins domain rm tatoh                    # remove domínio e seus serviços`,
}

var (
	domainTitle string
	domainDesc  string
)

func init() {
	domainCmd.AddCommand(domainAddCmd, domainLsCmd, domainRmCmd)
	domainAddCmd.Flags().StringVar(&domainTitle, "title", "", "Título exibido na landing page")
	domainAddCmd.Flags().StringVar(&domainDesc, "desc", "", "Descrição do domínio")
}

// ─── odins domain add ────────────────────────────────────────────────────────

var domainAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Criar um novo domínio local com landing page",
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
		return fmt.Errorf("domínio '%s' já existe", name)
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
		fmt.Printf("  ⚠  landing page: %v\n", err)
	}

	// Register with Caddy
	caddyClient := caddy.New()
	if err := caddyClient.AddDomain(hostname, pageDir); err != nil {
		fmt.Printf("  ⚠  caddy domain route: %v\n", err)
		fmt.Printf("     (rode 'odins init' se o Caddy ainda não foi configurado)\n")
	}

	fmt.Println()
	fmt.Printf("  ✓ Domínio criado: https://%s\n", hostname)
	fmt.Printf("  → Landing page gerada em %s\n", pageDir)
	fmt.Println()
	fmt.Printf("  Para adicionar serviços, crie um .odins com:\n")
	fmt.Printf("    [project]\n")
	fmt.Printf("    domain = \"%s\"\n", name)
	fmt.Println()
	fmt.Printf("  Depois rode: odins up\n")
	fmt.Println()
	return nil
}

// ─── odins domain ls ─────────────────────────────────────────────────────────

var domainLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Listar domínios cadastrados",
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
		fmt.Println("  Nenhum domínio cadastrado.")
		fmt.Println("  Use: odins domain add <nome>")
		return nil
	}

	fmt.Println()
	fmt.Printf("  %-20s  %-30s  %s\n", "DOMÍNIO", "FQDN", "SERVIÇOS")
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
	Short: "Remover um domínio e todos os seus serviços",
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
		return fmt.Errorf("domínio '%s' não encontrado", name)
	}

	routes := store.ByDomain(name)
	hostname := name + "." + cfg.TLD

	// Remove each route from proxy
	for _, r := range routes {
		if err := proxyRemove(cfg, r.Subdomain); err != nil {
			fmt.Printf("  ⚠  proxy remove %s: %v\n", r.Subdomain, err)
		}
	}

	// Remove domain route from Caddy
	caddyClient := caddy.New()
	if err := caddyClient.RemoveDomain(hostname); err != nil {
		fmt.Printf("  ⚠  caddy domain remove: %v\n", err)
	}

	// Remove from store (also removes all routes in that domain)
	store.RemoveDomain(name)
	if err := store.Save(); err != nil {
		return err
	}

	fmt.Printf("  ✓ Domínio '%s' e %d serviço(s) removidos.\n", name, len(routes))
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
// e.g. "web.tatoh.odins" with domain="tatoh" tld="odins" → "web"
func extractSubdomain(fqdn, domain, tld string) string {
	suffix := "." + domain + "." + tld
	if len(fqdn) > len(suffix) && fqdn[len(fqdn)-len(suffix):] == suffix {
		return fqdn[:len(fqdn)-len(suffix)]
	}
	return fqdn
}
