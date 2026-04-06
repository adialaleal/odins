package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/i18n"
	"github.com/spf13/cobra"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "Onboarding guide for ODINS",
	Long:  `Shows the interactive introduction guide. Can be run at any time.`,
	RunE:  runWelcome,
}

func runWelcome(cmd *cobra.Command, args []string) error {
	return showWelcome(false)
}

// showWelcome displays the onboarding guide.
//   - firstRun = true  → full guide, offers to run `odins init`
//   - firstRun = false → short "getting started in this project" guide
//     when in a folder without .odins; full guide available via `odins welcome`
func showWelcome(firstRun bool) error {
	violet := "\033[38;5;141m"
	dim    := "\033[38;5;245m"
	bold   := "\033[1m"
	reset  := "\033[0m"
	green  := "\033[38;5;114m"

	// Detect whether we are in a no-project folder (for short welcome path).
	cwd, _ := os.Getwd()
	inNoProject := !config.ExistsProject(cwd)

	clear()

	// ── Logo ─────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(violet + bold + `   ██████╗ ██████╗ ██╗███╗   ██╗███████╗` + reset)
	fmt.Println(violet +        `  ██╔═══██╗██╔══██╗██║████╗  ██║██╔════╝` + reset)
	fmt.Println(violet +        `  ██║   ██║██║  ██║██║██╔██╗ ██║███████╗` + reset)
	fmt.Println(violet +        `  ██║   ██║██║  ██║██║██║╚██╗██║╚════██║` + reset)
	fmt.Println(violet +        `  ╚██████╔╝██████╔╝██║██║ ╚████║███████║` + reset)
	fmt.Println(violet +        `   ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═══╝╚══════╝` + reset)
	fmt.Println()
	fmt.Println(dim + "  ᚦ ᚢ ᚱ ᛋ ᛏ ᚨ ᛉ ᚾ   " + i18n.T("welcome.tagline") + reset)
	fmt.Println()

	// ── Short path: already onboarded, running in a new project folder ───
	if !firstRun && inNoProject {
		return showProjectWelcome(cwd, bold, dim, violet, green, reset)
	}

	// ── Full guide ────────────────────────────────────────────────────────
	pause(dim, reset)

	// Section 1
	section(i18n.T("welcome.section.what"), violet, bold, dim, reset)
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.elimina"))
	fmt.Println()
	fmt.Println(dim + "  " + i18n.T("welcome.sem_odins") + "                      " + i18n.T("welcome.com_odins") + reset)
	fmt.Println("  http://localhost:3000           https://web.<projeto>.odins")
	fmt.Println("  http://localhost:4000           https://api.<projeto>.odins")
	fmt.Println("  http://localhost:5173           https://admin.<projeto>.odins")
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.https_auto"))
	fmt.Println()
	pause(dim, reset)

	// Section 2
	section(i18n.T("welcome.section.how"), violet, bold, dim, reset)
	fmt.Println()
	fmt.Println("  1. " + bold + "DNS" + reset + "   — " + i18n.T("welcome.how_dns"))
	fmt.Println("  2. " + bold + "Proxy" + reset + " — " + i18n.T("welcome.how_proxy"))
	fmt.Println("  3. " + bold + "HTTPS" + reset + " — " + i18n.T("welcome.how_https"))
	fmt.Println()
	pause(dim, reset)

	// Section 3
	section(i18n.T("welcome.section.domains"), violet, bold, dim, reset)
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.domain_is"))
	fmt.Println()
	fmt.Println("    odins domain add <projeto>")
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.domain_landing"))
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.subdomain_is"))
	fmt.Println()
	fmt.Println("    web.<projeto>.odins   → seu Next.js na porta 3000")
	fmt.Println("    api.<projeto>.odins   → sua API na porta 4000")
	fmt.Println("    admin.<projeto>.odins → painel admin na porta 5173")
	fmt.Println()
	pause(dim, reset)

	// Section 4
	section(i18n.T("welcome.section.config"), violet, bold, dim, reset)
	fmt.Println()
	fmt.Println("  " + bold + ".odins" + reset + ":")
	fmt.Println()
	fmt.Println(dim + "  [project]" + reset)
	fmt.Println(dim + `  name    = "meu-projeto"` + reset)
	fmt.Println(dim + `  domain  = "meu-projeto"   # workspace pai` + reset)
	fmt.Println(dim + `  runtime = "node"` + reset)
	fmt.Println()
	fmt.Println(dim + "  [[routes]]" + reset)
	fmt.Println(dim + `  subdomain = "web"         # → web.meu-projeto.odins` + reset)
	fmt.Println(dim + "  port      = 3000" + reset)
	fmt.Println(dim + "  https     = true" + reset)
	fmt.Println()
	fmt.Println("  " + i18n.T("welcome.auto_detect"))
	fmt.Println()
	pause(dim, reset)

	// Section 5
	section(i18n.T("welcome.section.commands"), violet, bold, dim, reset)
	fmt.Println()
	printCmd("odins init",              i18n.T("cmd.init_desc"), bold, dim, reset)
	printCmd("odins domain add <proj>", i18n.T("cmd.domain_add_desc"), bold, dim, reset)
	printCmd("odins up",                i18n.T("cmd.up_desc"), bold, dim, reset)
	printCmd("odins ls",                i18n.T("cmd.ls_desc"), bold, dim, reset)
	printCmd("odins kill <fqdn>",       i18n.T("cmd.kill_desc"), bold, dim, reset)
	printCmd("odins down",              i18n.T("cmd.down_desc"), bold, dim, reset)
	printCmd("odins",                   i18n.T("cmd.tui_desc"), bold, dim, reset)
	printCmd("odins welcome",           i18n.T("cmd.welcome_desc"), bold, dim, reset)
	fmt.Println()
	pause(dim, reset)

	// Section 6: Next steps
	section(i18n.T("welcome.section.next"), violet, bold, dim, reset)
	fmt.Println()

	cfg, _ := config.LoadGlobal()
	if !cfg.OnboardingDone && cfg.TLD == "" {
		// Fresh install — offer to run init
		fmt.Println("  " + i18n.T("welcome.not_configured"))
		fmt.Println()
		fmt.Print("  " + i18n.T("welcome.run_init") + " " + dim + i18n.T("welcome.run_init_prompt") + " " + reset)
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		yes := i18n.T("welcome.run_init_yes")
		if ans == "" || ans == yes || ans == "y" || ans == "s" {
			fmt.Println()
			return runInit(nil, nil)
		}
		fmt.Println()
		fmt.Println("  " + i18n.T("welcome.ok"))
	} else {
		fmt.Println(green + "  " + i18n.T("welcome.already_configured") + reset)
		fmt.Println()
		fmt.Println("  " + i18n.T("welcome.create_domain"))
		fmt.Println("    " + bold + "odins domain add meu-projeto" + reset)
		fmt.Println()
		fmt.Println("  " + i18n.T("welcome.then_project"))
		fmt.Println("    " + bold + "odins up" + reset)
	}

	// Mark onboarding as done
	cfg.OnboardingDone = true
	_ = config.SaveGlobal(cfg)

	fmt.Println()
	return nil
}

// showProjectWelcome is the short welcome shown when running `odins` in a
// folder without .odins, for users already globally onboarded.
func showProjectWelcome(cwd, bold, dim, violet, green, reset string) error {
	section(i18n.T("welcome.new_folder.title"), violet, bold, dim, reset)
	fmt.Println()

	d := detect.Project(cwd)
	if d.Runtime != "unknown" {
		fmt.Println(green + "  " + i18n.Tf("welcome.new_folder.detected",
			d.Name, d.Runtime, d.Framework, d.Port) + reset)
		fmt.Println()
		fmt.Println("  " + i18n.T("welcome.new_folder.activate"))
		fmt.Println()
		fmt.Println("    " + bold + "odins up" + reset)
	} else {
		fmt.Println("  " + i18n.T("welcome.new_folder.manual"))
		fmt.Println()
		fmt.Println("    " + bold + "odins add <subdominio> --port <porta>" + reset)
		fmt.Println()
		fmt.Println("  " + i18n.T("welcome.new_folder.or_add"))
		fmt.Println()
		fmt.Println("    " + bold + "odins domain add <projeto>" + reset)
		fmt.Println("    " + bold + "odins up" + reset)
	}

	fmt.Println()
	fmt.Println(dim + "  " + i18n.T("welcome.new_folder.see_guide") + reset)
	fmt.Println()
	return nil
}

func section(title, violet, bold, dim, reset string) {
	line := strings.Repeat("─", 50)
	fmt.Println(violet + bold + "  " + title + reset)
	fmt.Println(dim + "  " + line + reset)
}

func printCmd(cmd, desc, bold, dim, reset string) {
	fmt.Printf("  %-34s %s%s%s\n", bold+cmd+reset, dim, desc, reset)
}

func pause(dim, reset string) {
	fmt.Print(dim + "  " + i18n.T("welcome.enter") + reset)
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
}

func clear() {
	fmt.Print("\033[2J\033[H")
}
