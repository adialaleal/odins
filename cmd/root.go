package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
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
  odins domain add tatoh  — criar workspace tatoh.odins
  odins up                — read .odins config and apply routes
  odins add <domain>      — add a single route
  odins ls                — list active routes
  odins kill <domain>     — remove a route
  odins welcome           — guia de onboarding
  odins                   — open the TUI dashboard`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputJSON {
			return service.InvalidInput("modo TUI indisponível com --json; use um subcomando como `odins detect --json` ou `odins ls --json`")
		}
		if !isInteractiveIO() {
			return cmd.Help()
		}

		// First-run detection: show welcome before TUI if never onboarded.
		cfg, err := config.LoadGlobal()
		if err == nil && !cfg.OnboardingDone {
			if err := showWelcome(true); err != nil {
				return err
			}
		}
		return tui.Run()
	},
}

// SetVersion wires version info injected by GoReleaser into the Cobra root command.
// Cobra automatically adds --version / -v when Version is non-empty.
func SetVersion(v, c, d string) {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", v, c, d)
}

// Execute is the entry point called from main.go.
func Execute() int {
	return ExecuteWithArgs(os.Args[1:], os.Stdout, os.Stderr)
}

// ExecuteWithArgs runs the CLI with explicit args and writers. Useful for tests.
func ExecuteWithArgs(args []string, stdout, stderr io.Writer) int {
	resetCLIState()
	rootCmd.SetArgs(args)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	if err := rootCmd.Execute(); err != nil {
		cmdName := commandNameFromArgs(args)
		if jsonRequested(args) {
			if jsonErr := writeJSONError(stdout, cmdName, normalizeCobraError(err)); jsonErr != nil {
				fmt.Fprintln(stderr, jsonErr)
			}
		} else {
			fmt.Fprintln(stderr, normalizeCobraError(err))
		}
		return service.ExitCodeForError(normalizeCobraError(err))
	}

	return 0
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Render structured JSON output")
	rootCmd.AddCommand(
		initCmd,
		addCmd,
		upCmd,
		downCmd,
		lsCmd,
		killCmd,
		domainCmd,
		detectCmd,
		doctorCmd,
		welcomeCmd,
	)
}

func normalizeCobraError(err error) error {
	if err == nil {
		return nil
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "required flag"):
		return service.InvalidInput(message)
	case strings.Contains(message, "unknown command"):
		return service.InvalidInput(message)
	case strings.Contains(message, "accepts"):
		return service.InvalidInput(message)
	case strings.Contains(message, "unknown shorthand flag"):
		return service.InvalidInput(message)
	default:
		return err
	}
}

func resetCLIState() {
	outputJSON = false
	initNonInteractive = false
	initTLD = ""
	initBackend = ""
	addPort = 0
	addDocker = ""
	addProject = ""
	addDomain = ""
	addNoHTTPS = false
	upDir = ""
	domainTitle = ""
	domainDesc = ""
	detectDir = ""
}
