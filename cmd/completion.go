package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion script",
	Long: `Generate an autocompletion script for odins in the specified shell.

Bash:
  odins completion bash > /opt/homebrew/etc/bash_completion.d/odins
  # or: source <(odins completion bash)

Zsh:
  odins completion zsh > "${fpath[1]}/_odins"
  # or add to ~/.zshrc: source <(odins completion zsh)

Fish:
  odins completion fish > ~/.config/fish/completions/odins.fish`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(out)
		case "zsh":
			return rootCmd.GenZshCompletion(out)
		case "fish":
			return rootCmd.GenFishCompletion(out, true)
		default:
			return fmt.Errorf("shell não suportado: %s (use bash, zsh ou fish)", args[0])
		}
	},
}
