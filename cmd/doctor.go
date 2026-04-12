package cmd

import (
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Inspect local ODINS health and suggest next actions",
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.Doctor()
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "doctor", result, warnings)
	}

	if result.Healthy {
		writeTextLine(cmd.OutOrStdout(), "  ✓ Ambiente ODINS saudável")
	} else {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  Ambiente ODINS requer atenção")
	}
	for _, check := range result.Checks {
		if check.OK {
			writeTextLine(cmd.OutOrStdout(), "  ✓ %-16s %s", check.Name, check.Details)
			continue
		}
		writeTextLine(cmd.OutOrStdout(), "  ✗ %-16s %s", check.Name, check.Details)
		if check.Action != "" {
			writeTextLine(cmd.OutOrStdout(), "    → %s", check.Action)
		}
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	return nil
}
