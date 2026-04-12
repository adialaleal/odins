package cmd

import (
	"fmt"
	"io"

	"github.com/adialaleal/odins/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	scanMaxDepth    int
	scanCreateOdins bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [dir]",
	Short: "Scan directory tree for projects and build a global index",
	Long: `Scan a directory tree, detect projects, and register them in a global index.

Default directory is ~/Projects. Detected projects are recorded in
~/.config/odins/projects.json for AI tool discovery.

Examples:
  odins scan                          # scan ~/Projects
  odins scan ~/Code --max-depth 2     # scan with limited depth
  odins scan --create-odins           # also create .odins where missing
  odins scan --json                   # structured JSON output`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := ""
		if len(args) > 0 {
			dir = args[0]
		}

		opts := scanner.ScanOptions{
			Directory:   dir,
			MaxDepth:    scanMaxDepth,
			CreateOdins: scanCreateOdins,
		}

		result, err := scanner.Scan(opts)
		if err != nil {
			return fmt.Errorf("falha no scan: %w", err)
		}

		if err := scanner.UpdateIndex(result); err != nil {
			return fmt.Errorf("falha ao atualizar índice: %w", err)
		}

		out := cmd.OutOrStdout()

		if outputJSON {
			return writeJSONSuccess(out, "scan", result, nil)
		}

		printScanTable(out, result)
		return nil
	},
}

func init() {
	scanCmd.Flags().IntVar(&scanMaxDepth, "max-depth", 3, "Maximum directory depth to scan")
	scanCmd.Flags().BoolVar(&scanCreateOdins, "create-odins", false, "Create .odins files where missing")
}

func printScanTable(w io.Writer, result scanner.ScanResult) {
	fmt.Fprintf(w, "Diretório: %s\n", result.RootDirectory)
	fmt.Fprintf(w, "Projetos encontrados: %d\n\n", len(result.Projects))

	if len(result.Projects) == 0 {
		fmt.Fprintln(w, "Nenhum projeto detectado.")
		return
	}

	fmt.Fprintf(w, "  %-25s %-10s %-12s %5s  %s\n", "NOME", "RUNTIME", "FRAMEWORK", "PORTA", "STATUS")
	fmt.Fprintf(w, "  %-25s %-10s %-12s %5s  %s\n", "────", "───────", "─────────", "─────", "──────")

	for _, p := range result.Projects {
		status := "novo"
		if p.HasOdins {
			status = ".odins"
		}
		fmt.Fprintf(w, "  %-25s %-10s %-12s %5d  %s\n",
			truncate(p.Name, 25),
			p.Runtime,
			truncate(p.Framework, 12),
			p.Port,
			status,
		)
	}

	if result.Created > 0 {
		fmt.Fprintf(w, "\n%d arquivo(s) .odins criado(s).\n", result.Created)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
