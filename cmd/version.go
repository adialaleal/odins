package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print detailed version information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputJSON {
			data := map[string]string{
				"version": buildVersion,
				"commit":  buildCommit,
				"built":   buildDate,
				"go":      runtime.Version(),
				"os":      runtime.GOOS,
				"arch":    runtime.GOARCH,
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(data)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "odins %s\n", buildVersion)
		fmt.Fprintf(cmd.OutOrStdout(), "  commit  : %s\n", buildCommit)
		fmt.Fprintf(cmd.OutOrStdout(), "  built   : %s\n", buildDate)
		fmt.Fprintf(cmd.OutOrStdout(), "  go      : %s\n", runtime.Version())
		fmt.Fprintf(cmd.OutOrStdout(), "  os/arch : %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return nil
	},
}
