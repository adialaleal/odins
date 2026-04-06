package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/page"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adialaleal/odins/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odins",
	Short: "The All-Father of Local DNS",
	Long: `ODINS — Local DNS + Reverse Proxy manager for macOS developers.

Stop fighting with ports. Route your local projects to beautiful subdomains
with automatic HTTPS. Works with Node.js, Go, Python, Docker, and more.

  odins init              — one-time setup (DNS, proxy, HTTPS)
  odins domain add <proj> — create workspace <proj>.odins
  odins up                — read .odins config and apply routes
  odins add <domain>      — add a single route
  odins ls                — list active routes
  odins kill <domain>     — remove a route
  odins welcome           — onboarding guide
  odins                   — open the TUI dashboard`,
	// Sync Caddy from persisted state before every command so routes
	// survive Caddy restarts without requiring the user to re-run odins up.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		syncCaddyFromState()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadGlobal()
		cwd, _ := os.Getwd()
		noProject := !config.ExistsProject(cwd)

		// Show welcome when:
		//   a) Global first-run (OnboardingDone == false), OR
		//   b) Running in a folder without .odins (per-project first-run)
		if err == nil && (!cfg.OnboardingDone || noProject) {
			firstRun := !cfg.OnboardingDone
			if err := showWelcome(firstRun); err != nil {
				return err
			}
			// After the per-project welcome, don't open TUI —
			// user was guided to run odins up / odins init.
			if noProject && !config.ExistsProject(cwd) {
				return nil
			}
		}

		return tui.Run()
	},
}

// syncCaddyFromState re-applies all routes and domain landing pages from the
// persisted store into Caddy. This is a no-op when Caddy is not running or
// the proxy backend is not Caddy. It ensures routes survive Caddy restarts.
func syncCaddyFromState() {
	cfg, err := config.LoadGlobal()
	if err != nil || cfg.ProxyBackend != config.BackendCaddy {
		return
	}
	c := caddy.New()
	if !c.IsRunning() {
		return
	}
	store, err := state.Load()
	if err != nil {
		return
	}
	// Build domain → page-dir map
	domainPages := make(map[string]string)
	for _, d := range store.Domains {
		dir := filepath.Join(page.PagesDir(), d.Name)
		hostname := d.Name + "." + cfg.TLD
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			domainPages[hostname] = dir
		}
	}
	// Fire-and-forget: sync failures are non-fatal (Caddy may have routes already)
	_ = c.SyncRoutes(store.Routes, domainPages)
}

// SetVersion wires version info injected by GoReleaser into the Cobra root command.
// Cobra automatically adds --version / -v when Version is non-empty.
func SetVersion(v, c, d string) {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", v, c, d)
}

// Execute is the entry point called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		initCmd,
		addCmd,
		upCmd,
		downCmd,
		lsCmd,
		killCmd,
		domainCmd,
		welcomeCmd,
	)
}
