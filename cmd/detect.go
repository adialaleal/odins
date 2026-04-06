package cmd

import (
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect project runtime, framework and recommended .odins config",
	RunE:  runDetect,
}

var detectDir string

func init() {
	detectCmd.Flags().StringVar(&detectDir, "dir", "", "Directory to inspect instead of current")
}

func runDetect(cmd *cobra.Command, args []string) error {
	manager := serviceFactory()
	result, warnings, err := manager.Detect(detectDir)
	if err != nil {
		return err
	}

	if outputJSON {
		return writeJSONSuccess(cmd.OutOrStdout(), "detect", result, warnings)
	}

	writeTextLine(cmd.OutOrStdout(), "  Diretório: %s", result.Directory)
	writeTextLine(cmd.OutOrStdout(), "  Runtime: %s", result.Detected.Runtime)
	writeTextLine(cmd.OutOrStdout(), "  Framework: %s", result.Detected.Framework)
	writeTextLine(cmd.OutOrStdout(), "  Porta: %d", result.Detected.Port)
	if result.Detected.StartCmd != "" {
		writeTextLine(cmd.OutOrStdout(), "  Start: %s", result.Detected.StartCmd)
	}
	writeTextLine(cmd.OutOrStdout(), "  .odins existente: %t", result.ProjectConfigExists)
	if !result.ProjectConfigExists {
		writeTextLine(cmd.OutOrStdout(), "  .odins recomendado: %s", result.ProjectConfigPath)
	}
	for _, warning := range warnings {
		writeTextLine(cmd.OutOrStdout(), "  ⚠  %s", warning)
	}
	return nil
}
