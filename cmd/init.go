package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "One-time setup: DNS, reverse proxy, and HTTPS",
	Long: `odins init configures your Mac for local domain routing.

It will:
  1. Install dnsmasq and a reverse proxy via Homebrew
  2. Configure DNS wildcard resolution
  3. Set up HTTPS with a trusted local certificate
  4. Save the ODINS global config`,
	RunE: runInit,
}

var (
	initNonInteractive bool
	initTLD            string
	initBackend        string
)

func init() {
	initCmd.Flags().BoolVar(&initNonInteractive, "non-interactive", false, "Use defaults and flags without interactive prompts")
	initCmd.Flags().StringVar(&initTLD, "tld", "", "TLD for local domains")
	initCmd.Flags().StringVar(&initBackend, "backend", "", "Reverse proxy backend: caddy, nginx, apache")
}

func runInit(cmd *cobra.Command, args []string) error {
	selectedTLD := strings.TrimSpace(initTLD)
	selectedBackend := strings.TrimSpace(initBackend)

	if !outputJSON && !initNonInteractive && isInteractiveIO() {
		if selectedTLD == "" {
			selectedTLD = chooseTLD()
		}
		if selectedBackend == "" {
			selectedBackend = chooseBackend()
		}
	}

	manager := serviceFactory()
	result, warnings, err := manager.Init(service.InitOptions{
		NonInteractive: initNonInteractive || outputJSON || !isInteractiveIO(),
		TLD:            selectedTLD,
		Backend:        selectedBackend,
	})
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "init", result, warnings)
	}

	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  ✓ ODINS configurado com sucesso!")
	writeTextLine(cmd.OutOrStdout(), "  → TLD: .%s", result.TLD)
	writeTextLine(cmd.OutOrStdout(), "  → Proxy: %s", result.Backend)
	for _, step := range result.Steps {
		if step.OK {
			writeTextLine(cmd.OutOrStdout(), "  ✓ %s: %s", step.Name, step.Detail)
			continue
		}
		if step.Warning != "" {
			writeTextLine(cmd.OutOrStdout(), "  ⚠  %s: %s", step.Name, step.Warning)
		}
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  Domínios disponíveis: https://<projeto>.%s", result.TLD)
	writeTextLine(cmd.OutOrStdout(), "")
	writeTextLine(cmd.OutOrStdout(), "  Próximos passos:")
	writeTextLine(cmd.OutOrStdout(), "    cd meu-projeto && odins detect --json")
	writeTextLine(cmd.OutOrStdout(), "    cd meu-projeto && odins up")
	writeTextLine(cmd.OutOrStdout(), "    odins doctor")
	writeTextLine(cmd.OutOrStdout(), "")
	return nil
}

func chooseTLD() string {
	writeTextLine(os.Stdout, "")
	for i, t := range config.SupportedTLDs {
		warn := ""
		if t.Warning != "" {
			warn = " ⚠"
		}
		writeTextLine(os.Stdout, "  [%d] %s%s", i+1, t.Label, warn)
	}
	writeTextLine(os.Stdout, "")
	fmt.Fprint(os.Stdout, "  Escolha (1): ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	if input == "" {
		return config.SupportedTLDs[0].TLD
	}
	for i, t := range config.SupportedTLDs {
		if fmt.Sprintf("%d", i+1) == input {
			return t.TLD
		}
	}
	return config.SupportedTLDs[0].TLD
}

func chooseBackend() string {
	backends := []string{"caddy (recomendado — HTTPS automático)", "nginx", "apache"}
	values := []string{"caddy", "nginx", "apache"}

	writeTextLine(os.Stdout, "")
	for i, backend := range backends {
		writeTextLine(os.Stdout, "  [%d] %s", i+1, backend)
	}
	writeTextLine(os.Stdout, "")
	fmt.Fprint(os.Stdout, "  Escolha (1): ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	if input == "" {
		return values[0]
	}
	for i, value := range values {
		if fmt.Sprintf("%d", i+1) == input {
			return value
		}
	}
	return values[0]
}
