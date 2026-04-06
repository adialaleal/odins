package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

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

var domainAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Criar um novo domínio local com landing page",
	Args:  exactArgs(1),
	RunE:  runDomainAdd,
}

func runDomainAdd(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.DomainAdd(strings.ToLower(args[0]), domainTitle, domainDesc)
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "domain add", result, warnings)
	}

	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  ✓ Domínio criado: https://%s", result.Hostname)
	writeTextLine(cmd.OutOrStdout(), "  → Landing page gerada em %s", result.PageDir)
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  Para adicionar serviços, crie um .odins com:")
	writeTextLine(cmd.OutOrStdout(), "    [project]")
	writeTextLine(cmd.OutOrStdout(), "    domain = %q", result.Domain.Name)
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  Depois rode: odins up")
	writeTextLine(cmd.OutOrStdout(), "")
	return nil
}

var domainLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List registered domains",
	RunE:  runDomainLs,
}

func runDomainLs(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.DomainList()
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "domain ls", result, warnings)
	}

	if result.Count == 0 {
		writeTextLine(cmd.OutOrStdout(), "  Nenhum domínio cadastrado.")
		writeTextLine(cmd.OutOrStdout(), "  Use: odins domain add <nome>")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\n  DOMÍNIO\tFQDN\tSERVIÇOS")
	fmt.Fprintln(w, "  ───────\t────\t────────")
	for _, item := range result.Domains {
		fmt.Fprintf(w, "  %s\t%s\t%d\n", item.Domain.Name, item.Hostname, item.Services)
	}
	w.Flush()
	writeTextLine(cmd.OutOrStdout(), "")
	return nil
}

var domainRmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remover um domínio e todos os seus serviços",
	Args:  exactArgs(1),
	RunE:  runDomainRm,
}

func runDomainRm(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.DomainRemove(strings.ToLower(args[0]))
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "domain rm", result, warnings)
	}

	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	writeTextLine(cmd.OutOrStdout(), "  ✓ Domínio '%s' e %d serviço(s) removidos.", result.Domain, len(result.RemovedRoutes))
	return nil
}
