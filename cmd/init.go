package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/aiclient"
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "One-time setup: DNS, reverse proxy, and HTTPS",
	Long: `odins init configures your Mac for local domain routing.

It will:
  1. Install dnsmasq, caddy (or nginx/apache) via Homebrew
  2. Ask which TLD and proxy backend you want
  3. Configure DNS wildcard resolution (one sudo prompt)
  4. Set up HTTPS with a trusted local certificate
  5. Start all services

Run again at any time to repair or reconfigure.`,
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
	out := commandWriter(cmd)
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
		return writeJSONSuccess(out, "init", result, warnings)
	}

	writeTextLine(out, "")
	writeTextLine(out, "  ✓ ODINS configurado com sucesso!")
	writeTextLine(out, "  → TLD: .%s", result.TLD)
	writeTextLine(out, "  → Proxy: %s", result.Backend)
	for _, step := range result.Steps {
		if step.OK {
			writeTextLine(out, "  ✓ %s: %s", step.Name, step.Detail)
			continue
		}
		if step.Warning != "" {
			writeTextLine(out, "  ⚠  %s: %s", step.Name, step.Warning)
		}
	}
	for _, warning := range warnings {
		writeTextLine(out, "  ⚠  %s", warning)
	}
	writeTextLine(out, "")
	writeTextLine(out, "  Domínios disponíveis: https://<projeto>.%s", result.TLD)
	writeTextLine(out, "")
	writeTextLine(out, "  Próximos passos:")
	writeTextLine(out, "    cd meu-projeto && odins detect --json")
	writeTextLine(out, "    cd meu-projeto && odins up")
	writeTextLine(out, "    odins doctor")
	writeTextLine(out, "")

	// AI client auto-configuration.
	if !initNonInteractive && isInteractiveIO() {
		configureAIClients(out)
	}

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

func configureAIClients(out io.Writer) {
	clients := aiclient.DetectClients()
	if len(clients) == 0 {
		return
	}

	writeTextLine(out, "  AI Tools detectados:")
	for _, c := range clients {
		writeTextLine(out, "    • %s", c.Name)
	}
	writeTextLine(out, "")
	fmt.Fprint(os.Stdout, "  Configurar integração MCP para essas ferramentas? [Y/n] ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "" && input != "y" && input != "s" && input != "sim" && input != "yes" {
		return
	}

	results := aiclient.ConfigureAll(clients)
	for _, r := range results {
		if r.Error != "" {
			writeTextLine(out, "    ⚠  %s: %s", r.Client, r.Error)
		} else if r.Configured {
			writeTextLine(out, "    ✓ %s → %s", r.Client, r.ConfigPath)
		}
	}
	writeTextLine(out, "")
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
