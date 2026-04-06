package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/adialaleal/odins/internal/cert"
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/dns"
	"github.com/adialaleal/odins/internal/helper"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/proxy/nginx"
	"github.com/adialaleal/odins/pkg/brew"
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
  5. Start all services`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println("  ____  ____  ___ _   _ ____")
	fmt.Println(" / __ \\|  _ \\|_ _| \\ | / ___|")
	fmt.Println("| |  | | | | || ||  \\| \\___ \\")
	fmt.Println("| |__| | |_| || || |\\  |___) |")
	fmt.Println(" \\____/|____/|___|_| \\_|____/ ")
	fmt.Println()
	fmt.Println("  The All-Father of Local DNS — Setup")
	fmt.Println()

	// Step 1: Check Homebrew
	step(1, "Verificando Homebrew")
	if !brew.IsInstalled() {
		return fmt.Errorf("Homebrew não encontrado. Instale em https://brew.sh")
	}
	ok()

	// Step 2: Choose TLD
	step(2, "Escolha o TLD para seus domínios locais")
	tld := chooseTLD()
	fmt.Printf("  → TLD escolhido: .%s\n\n", tld)

	// Step 3: Choose proxy backend
	step(3, "Escolha o reverse proxy")
	backend := chooseBackend()
	fmt.Printf("  → Proxy escolhido: %s\n\n", backend)

	// Step 4: Install dnsmasq
	step(4, "Instalando dnsmasq")
	if err := brew.Install("dnsmasq"); err != nil {
		return err
	}
	if err := dns.GenerateConfig([]string{tld}, 5353); err != nil {
		return err
	}
	if err := dns.LinkConfig(); err != nil {
		fmt.Printf("  ⚠  Não foi possível linkar config dnsmasq: %v\n", err)
	}
	if err := brew.ServiceRestart("dnsmasq"); err != nil {
		return err
	}
	ok()

	// Step 5: Install proxy
	step(5, fmt.Sprintf("Instalando %s", backend))
	var proxyFormula string
	switch backend {
	case "nginx":
		proxyFormula = "nginx"
	case "apache":
		proxyFormula = "httpd"
	default:
		proxyFormula = "caddy"
	}
	if err := brew.Install(proxyFormula); err != nil {
		return err
	}
	if err := brew.ServiceStart(proxyFormula); err != nil {
		fmt.Printf("  ⚠  %s start: %v\n", proxyFormula, err)
	}
	ok()

	// Step 6: Write /etc/resolver/<tld> via sudo
	step(6, fmt.Sprintf("Configurando /etc/resolver/%s (requer sudo)", tld))
	fmt.Printf("  → ODINS precisa de acesso root para criar /etc/resolver/%s\n", tld)
	fmt.Println("  → Isso permite que seu Mac resolva *."+tld+" para 127.0.0.1")
	fmt.Println()
	if err := helper.SudoWriteResolver(tld, 5353); err != nil {
		return fmt.Errorf("write resolver: %w", err)
	}
	ok()

	// Step 7: Trust HTTPS certificate
	step(7, "Configurando HTTPS local")
	if backend == "caddy" {
		// Caddy needs to be running first to generate its CA
		fmt.Println("  → Iniciando Caddy e aguardando CA ser gerada...")
		waitForCaddyCA()
		caPath := cert.CaddyCAPath()
		if caPath != "" {
			fmt.Printf("  → CA encontrada em %s\n", caPath)
			fmt.Println("  → Adicionando ao keychain do macOS (requer sudo)")
			if err := helper.SudoTrustCA(caPath); err != nil {
				fmt.Printf("  ⚠  %v\n", err)
			} else {
				ok()
			}
		} else {
			fmt.Println("  ⚠  CA do Caddy ainda não gerada — abra um domínio no browser para ativá-la")
		}

		// Init Caddy with base config
		caddyClient := caddy.New()
		if err := caddyClient.Init(tld); err != nil {
			fmt.Printf("  ⚠  Caddy config init: %v\n", err)
		}
	} else {
		// Use mkcert for nginx/apache
		if brew.IsFormulaInstalled("mkcert") || installMkcert() {
			if err := cert.InstallMkcertCA(); err != nil {
				fmt.Printf("  ⚠  mkcert install: %v\n", err)
			} else {
				ok()
			}
		}

		// Init nginx include directory
		if backend == "nginx" {
			nginx.New().Init()
		} else {
			os.MkdirAll(apacheConfDir(), 0755)
		}
	}

	// Step 8: Save config
	step(8, "Salvando configuração")
	cfg := config.GlobalConfig{
		TLD:          tld,
		ProxyBackend: config.ProxyBackend(backend),
		DnsmasqPort:  5353,
		CaddyAdmin:   "http://localhost:2019",
		HTTPPort:     80,
		HTTPSPort:    443,
	}
	if err := config.SaveGlobal(cfg); err != nil {
		return err
	}
	ok()

	// Done
	fmt.Println()
	fmt.Println("  ✓ ODINS configurado com sucesso!")
	fmt.Println()
	fmt.Printf("  Domínios disponíveis: https://<projeto>.%s\n", tld)
	fmt.Println()
	fmt.Println("  Próximos passos:")
	fmt.Printf("    cd meu-projeto && odins up     # detecta e cria as rotas\n")
	fmt.Printf("    odins add api.projeto.%s --port 3000\n", tld)
	fmt.Printf("    odins                           # abrir TUI\n")
	fmt.Println()

	return nil
}

func chooseTLD() string {
	fmt.Println()
	for i, t := range config.SupportedTLDs {
		warn := ""
		if t.Warning != "" {
			warn = " ⚠"
		}
		fmt.Printf("  [%d] %s%s\n", i+1, t.Label, warn)
	}
	fmt.Println()
	fmt.Print("  Escolha (1): ")

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
	backendValues := []string{"caddy", "nginx", "apache"}

	fmt.Println()
	for i, b := range backends {
		fmt.Printf("  [%d] %s\n", i+1, b)
	}
	fmt.Println()
	fmt.Print("  Escolha (1): ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	if input == "" {
		return backendValues[0]
	}
	for i := range backendValues {
		if fmt.Sprintf("%d", i+1) == input {
			return backendValues[i]
		}
	}
	return backendValues[0]
}

func step(n int, desc string) {
	fmt.Printf("  [%d] %s...\n", n, desc)
}

func ok() {
	fmt.Println("      ✓ OK")
}

func waitForCaddyCA() {
	// Wait up to 10s for Caddy to generate its CA
	for i := 0; i < 20; i++ {
		if cert.CaddyCAPath() != "" {
			return
		}
		exec.Command("sleep", "0.5").Run()
	}
}

func installMkcert() bool {
	fmt.Println("  → Instalando mkcert...")
	return brew.Install("mkcert") == nil
}

func apacheConfDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.config/odins/apache/vhosts"
}
