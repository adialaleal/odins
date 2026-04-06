package cmd

import (
	"fmt"
	"os"

	"github.com/adialaleal/odins/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odins",
	Short: "The All-Father of Local DNS",
	Long: `ODINS — Local DNS + Reverse Proxy manager for macOS developers.

Stop fighting with ports. Route your local projects to beautiful subdomains
with automatic HTTPS. Works with Node.js, Go, Python, Docker, and more.

  odins init           — one-time setup (DNS, proxy, HTTPS)
  odins up             — read .odins config and apply routes
  odins add <domain>   — add a single route
  odins ls             — list active routes
  odins kill <domain>  — remove a route
  odins                — open the TUI dashboard`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
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
	)
}
