package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
	"github.com/adialaleal/odins/internal/state"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [subdomain]",
	Short: "Open a local domain in the default browser",
	Long: `Open a local ODINS domain in the default browser.

Examples:
  odins open              # open the current project's domain landing page
  odins open api          # open https://api.<domain>.<tld> (inferred from .odins)
  odins open api.tatoh    # open https://api.tatoh.<tld>`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadGlobal()
		if err != nil {
			return fmt.Errorf("carregar configuração: %w", err)
		}

		tld := cfg.TLD

		var fqdn string

		if len(args) == 0 {
			// No arg: open the domain landing page inferred from .odins in cwd
			fqdn, err = domainLandingPage(tld)
			if err != nil {
				return err
			}
		} else {
			arg := args[0]
			// If it already contains a dot, treat as "subdomain.domain"
			if strings.Contains(arg, ".") {
				fqdn = arg + "." + tld
			} else {
				// Single token: try to resolve from store using .odins context
				fqdn, err = resolveSubdomain(arg, tld)
				if err != nil {
					return err
				}
			}
		}

		url := "https://" + fqdn
		fmt.Fprintf(cmd.OutOrStdout(), "abrindo %s\n", url)
		return exec.Command("open", url).Run()
	},
}

// domainLandingPage returns the landing page FQDN for the current project's domain.
func domainLandingPage(tld string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("obter diretório atual: %w", err)
	}

	if !config.ExistsProject(cwd) {
		return "", service.InvalidInput("nenhum arquivo .odins encontrado no diretório atual; especifique um subdomínio")
	}

	proj, err := config.LoadProject(filepath.Join(cwd, config.ProjectConfigFile))
	if err != nil {
		return "", fmt.Errorf("carregar .odins: %w", err)
	}

	domain := proj.Project.Domain
	if domain == "" {
		domain = proj.Project.Name
	}

	return domain + "." + tld, nil
}

// resolveSubdomain finds the FQDN for a bare subdomain token using .odins + store.
func resolveSubdomain(subdomain, tld string) (string, error) {
	cwd, _ := os.Getwd()

	// Try to infer domain from .odins file in cwd
	if config.ExistsProject(cwd) {
		proj, err := config.LoadProject(filepath.Join(cwd, config.ProjectConfigFile))
		if err == nil && proj.Project.Domain != "" {
			fqdn := subdomain + "." + proj.Project.Domain + "." + tld
			// Confirm the route exists in state
			if routeExists(fqdn) {
				return fqdn, nil
			}
			// Return the constructed FQDN anyway — user may not have run odins up yet
			return fqdn, nil
		}
	}

	// Fallback: search state for any route whose subdomain starts with the token
	store, err := state.Load()
	if err != nil {
		return "", fmt.Errorf("carregar estado: %w", err)
	}

	for _, r := range store.Routes {
		// r.Subdomain is the full FQDN (e.g. "api.tatoh.odin")
		if strings.HasPrefix(r.Subdomain, subdomain+".") {
			return r.Subdomain, nil
		}
	}

	return "", service.InvalidInput(fmt.Sprintf("subdomínio %q não encontrado; verifique com `odins ls`", subdomain))
}

func routeExists(fqdn string) bool {
	store, err := state.Load()
	if err != nil {
		return false
	}
	for _, r := range store.Routes {
		if r.Subdomain == fqdn {
			return true
		}
	}
	return false
}
