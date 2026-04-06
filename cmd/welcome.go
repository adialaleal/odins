package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/spf13/cobra"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "Mostrar guia de boas-vindas e onboarding do ODINS",
	Long:  `Exibe o guia interativo de introdução ao ODINS. Pode ser rodado a qualquer momento.`,
	RunE:  runWelcome,
}

func runWelcome(cmd *cobra.Command, args []string) error {
	return showWelcome(false)
}

// showWelcome displays the onboarding guide.
// If firstRun is true it offers to run `odins init` right after.
func showWelcome(firstRun bool) error {
	violet := "\033[38;5;141m"
	dim    := "\033[38;5;245m"
	bold   := "\033[1m"
	reset  := "\033[0m"
	green  := "\033[38;5;114m"

	clear()

	fmt.Println()
	fmt.Println(violet + bold + `   ██████╗ ██████╗ ██╗███╗   ██╗███████╗` + reset)
	fmt.Println(violet +         `  ██╔═══██╗██╔══██╗██║████╗  ██║██╔════╝` + reset)
	fmt.Println(violet +         `  ██║   ██║██║  ██║██║██╔██╗ ██║███████╗` + reset)
	fmt.Println(violet +         `  ██║   ██║██║  ██║██║██║╚██╗██║╚════██║` + reset)
	fmt.Println(violet +         `  ╚██████╔╝██████╔╝██║██║ ╚████║███████║` + reset)
	fmt.Println(violet +         `   ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═══╝╚══════╝` + reset)
	fmt.Println()
	fmt.Println(dim + "  ᚦ ᚢ ᚱ ᛋ ᛏ ᚨ ᛉ ᚾ   The All-Father of Local DNS" + reset)
	fmt.Println()
	pause()

	// ── O que é o ODINS ──────────────────────────────────────────────
	section("O que é o ODINS?")
	fmt.Println()
	fmt.Println("  ODINS elimina a guerra de portas no seu desenvolvimento local.")
	fmt.Println()
	fmt.Println(dim + "  Sem ODINS:                      Com ODINS:" + reset)
	fmt.Println("  http://localhost:3000           https://web.tatoh.odins")
	fmt.Println("  http://localhost:4000           https://api.tatoh.odins")
	fmt.Println("  http://localhost:5173           https://admin.tatoh.odins")
	fmt.Println()
	fmt.Println("  Cada projeto ganha um subdomínio bonito com HTTPS automático.")
	fmt.Println()
	pause()

	// ── Como funciona ─────────────────────────────────────────────────
	section("Como funciona?")
	fmt.Println()
	fmt.Println("  1. " + bold + "DNS" + reset + "   — dnsmasq resolve *.tatoh.odins → 127.0.0.1")
	fmt.Println("  2. " + bold + "Proxy" + reset + " — Caddy roteia web.tatoh.odins → localhost:3000")
	fmt.Println("  3. " + bold + "HTTPS" + reset + " — Caddy gerencia certificados automaticamente")
	fmt.Println()
	pause()

	// ── Domínios e subdomínios ────────────────────────────────────────
	section("Domínios e Subdomínios")
	fmt.Println()
	fmt.Println("  Um " + bold + "domínio" + reset + " é o workspace central dos seus projetos:")
	fmt.Println()
	fmt.Println("    odins domain add tatoh")
	fmt.Println()
	fmt.Println("  Isso cria " + violet + "tatoh.odins" + reset + " — uma landing page que lista")
	fmt.Println("  todos os serviços do workspace com status em tempo real.")
	fmt.Println()
	fmt.Println("  Cada " + bold + "subdomínio" + reset + " é um projeto/serviço:")
	fmt.Println()
	fmt.Println("    web.tatoh.odins   → seu Next.js na porta 3000")
	fmt.Println("    api.tatoh.odins   → sua API na porta 4000")
	fmt.Println("    admin.tatoh.odins → painel admin na porta 5173")
	fmt.Println()
	pause()

	// ── Arquivo .odins ────────────────────────────────────────────────
	section("Configuração por Projeto (.odins)")
	fmt.Println()
	fmt.Println("  Cada projeto tem um arquivo " + bold + ".odins" + reset + " na raiz:")
	fmt.Println()
	fmt.Println(dim + "  [project]" + reset)
	fmt.Println(dim + `  name   = "tatoh_web"` + reset)
	fmt.Println(dim + `  domain = "tatoh"        # workspace pai` + reset)
	fmt.Println(dim + `  runtime = "node"` + reset)
	fmt.Println()
	fmt.Println(dim + "  [[routes]]" + reset)
	fmt.Println(dim + `  subdomain = "web"       # → web.tatoh.odins` + reset)
	fmt.Println(dim + "  port      = 3000" + reset)
	fmt.Println(dim + "  https     = true" + reset)
	fmt.Println()
	fmt.Println("  O ODINS detecta Node.js, Go e Python automaticamente.")
	fmt.Println()
	pause()

	// ── Comandos ─────────────────────────────────────────────────────
	section("Comandos Principais")
	fmt.Println()
	printCmd("odins init",             "Setup único: DNS, proxy, HTTPS")
	printCmd("odins domain add tatoh", "Criar workspace tatoh.odins")
	printCmd("odins up",               "Ativar rotas do projeto atual")
	printCmd("odins ls",               "Listar rotas ativas")
	printCmd("odins kill <fqdn>",      "Remover uma rota")
	printCmd("odins down",             "Remover todas as rotas do projeto")
	printCmd("odins",                  "Abrir painel TUI")
	printCmd("odins welcome",          "Ver este guia novamente")
	fmt.Println()
	pause()

	// ── Próximos passos ───────────────────────────────────────────────
	section("Próximos Passos")
	fmt.Println()

	cfg, _ := config.LoadGlobal()
	if !cfg.OnboardingDone && cfg.TLD == "" {
		// Fresh install — offer to run init
		fmt.Println("  Parece que o ODINS ainda não foi configurado nesta máquina.")
		fmt.Println()
		fmt.Print("  Rodar " + bold + "odins init" + reset + " agora? " + dim + "[S/n] " + reset)
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		if ans == "" || ans == "s" || ans == "y" {
			fmt.Println()
			return runInit(nil, nil)
		}
		fmt.Println()
		fmt.Println("  Tudo bem! Quando estiver pronto, rode: " + bold + "odins init" + reset)
	} else {
		fmt.Println(green + "  ✓ ODINS já está configurado." + reset)
		fmt.Println()
		fmt.Println("  Comece criando um domínio:")
		fmt.Println("    " + bold + "odins domain add meu-workspace" + reset)
		fmt.Println()
		fmt.Println("  Depois num projeto:")
		fmt.Println("    " + bold + "odins up" + reset)
	}

	// Mark onboarding as done
	cfg.OnboardingDone = true
	_ = config.SaveGlobal(cfg)

	fmt.Println()
	return nil
}

func section(title string) {
	bold  := "\033[1m"
	violet := "\033[38;5;141m"
	reset  := "\033[0m"
	dim    := "\033[38;5;245m"
	line   := strings.Repeat("─", 50)
	fmt.Println(violet + bold + "  " + title + reset)
	fmt.Println(dim + "  " + line + reset)
}

func printCmd(cmd, desc string) {
	bold  := "\033[1m"
	dim   := "\033[38;5;245m"
	reset := "\033[0m"
	fmt.Printf("  %-32s %s%s%s\n", bold+cmd+reset, dim, desc, reset)
}

func pause() {
	dim   := "\033[38;5;245m"
	reset := "\033[0m"
	fmt.Print(dim + "  [Enter para continuar...]" + reset)
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
}

func clear() {
	fmt.Print("\033[2J\033[H")
}
