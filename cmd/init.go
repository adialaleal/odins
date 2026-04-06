package cmd

import (
	"fmt"
	"os"
	pathpkg "path/filepath"
	"strings"
	"time"

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
  5. Start all services

Run again at any time to repair or reconfigure.`,
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

	// Step 4: Install and configure dnsmasq
	step(4, "Configurando dnsmasq")
	if err := brew.Install("dnsmasq"); err != nil {
		return err
	}
	if err := dns.GenerateConfig([]string{tld}, 5300); err != nil {
		return err
	}
	if err := brew.ServiceRestart("dnsmasq"); err != nil {
		fmt.Printf("  ⚠  dnsmasq restart: %v\n", err)
	}
	ok()

	// Step 5: Install and configure proxy
	step(5, fmt.Sprintf("Configurando %s", backend))
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

	// For Caddy: create the Caddyfile BEFORE starting the brew service.
	// The plist runs: caddy run --config /opt/homebrew/etc/Caddyfile
	// Without this file caddy crash-loops immediately.
	if proxyFormula == "caddy" {
		if err := ensureCaddyfile(); err != nil {
			fmt.Printf("  ⚠  criar Caddyfile: %v\n", err)
		}
	}

	if err := brew.ServiceRestart(proxyFormula); err != nil {
		fmt.Printf("  ⚠  %s restart: %v\n", proxyFormula, err)
	}

	// Wait for Caddy admin API to be ready (up to 10s)
	if proxyFormula == "caddy" {
		waitForCaddyAPI()
	}
	ok()

	// Step 6: Write /etc/resolver/<tld> via sudo
	step(6, fmt.Sprintf("Configurando /etc/resolver/%s", tld))
	fmt.Printf("  → Uma janela de autenticação será aberta para criar /etc/resolver/%s\n", tld)
	fmt.Println("  → Isso permite que seu Mac resolva *."+tld+" para 127.0.0.1")
	fmt.Println()
	if err := helper.SudoWriteResolver(tld, 5300); err != nil {
		return fmt.Errorf("write resolver: %w", err)
	}
	// Flush macOS DNS cache so the new resolver takes effect immediately
	helper.SudoFlushDNS()
	ok()

	// Step 7: Trust HTTPS certificate
	step(7, "Configurando HTTPS local")
	if backend == "caddy" {
		// Push ODINS base config to Caddy (TLS internal)
		caddyClient := caddy.New()
		if err := caddyClient.Init(tld); err != nil {
			fmt.Printf("  ⚠  Caddy config init: %v\n", err)
		}

		// Wait for Caddy to generate its internal CA
		fmt.Println("  → Aguardando Caddy gerar CA local...")
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
			ok()
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
		if backend == "nginx" {
			nginx.New().Init()
		} else {
			os.MkdirAll(apacheConfDir(), 0755)
		}
	}

	// Step 8: Save config
	step(8, "Salvando configuração")
	cfg := config.GlobalConfig{
		TLD:            tld,
		ProxyBackend:   config.ProxyBackend(backend),
		DnsmasqPort:    5300,
		CaddyAdmin:     "http://localhost:2019",
		HTTPPort:       80,
		HTTPSPort:      443,
		OnboardingDone: true,
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
	fmt.Printf("    odins domain add meu-projeto   # criar workspace\n")
	fmt.Printf("    cd meu-projeto && odins up      # detecta e cria as rotas\n")
	fmt.Printf("    odins                           # abrir TUI\n")
	fmt.Println()

	return nil
}

// ensureCaddyfile creates a minimal Caddyfile so that `brew services start caddy`
// does not crash-loop. The file just enables the admin API on localhost:2019.
// ODINS manages all routes via the API — the Caddyfile itself stays minimal.
func ensureCaddyfile() error {
	candidates := []string{
		"/opt/homebrew/etc/Caddyfile",
		"/usr/local/etc/Caddyfile",
	}
	for _, path := range candidates {
		if _, err := os.Stat(pathpkg.Dir(path)); err == nil {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				content := "{\n\tadmin localhost:2019\n}\n"
				return os.WriteFile(path, []byte(content), 0644)
			}
			return nil // already exists — don't overwrite user customisation
		}
	}
	return fmt.Errorf("Homebrew etc directory not found")
}

// waitForCaddyAPI blocks until Caddy's admin API is responsive (max 10s).
func waitForCaddyAPI() {
	caddyClient := caddy.New()
	for i := 0; i < 20; i++ {
		if caddyClient.IsRunning() {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
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
	for i := 0; i < 20; i++ {
		if cert.CaddyCAPath() != "" {
			return
		}
		time.Sleep(500 * time.Millisecond)
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
