package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adialaleal/odins/internal/service"
	"github.com/spf13/cobra"
)

var outputJSON bool

type jsonSuccessEnvelope struct {
	OK       bool     `json:"ok"`
	Command  string   `json:"command"`
	Data     any      `json:"data"`
	Warnings []string `json:"warnings"`
}

type jsonErrorEnvelope struct {
	OK      bool          `json:"ok"`
	Command string        `json:"command"`
	Error   jsonErrorBody `json:"error"`
}

type jsonErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func isInteractiveIO() bool {
	stdoutInfo, stdoutErr := os.Stdout.Stat()
	stdinInfo, stdinErr := os.Stdin.Stat()
	if stdoutErr != nil || stdinErr != nil {
		return false
	}
	return (stdoutInfo.Mode()&os.ModeCharDevice) != 0 && (stdinInfo.Mode()&os.ModeCharDevice) != 0
}

func writeJSONSuccess(w io.Writer, command string, data any, warnings []string) error {
	if warnings == nil {
		warnings = []string{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonSuccessEnvelope{
		OK:       true,
		Command:  command,
		Data:     data,
		Warnings: warnings,
	})
}

func writeJSONError(w io.Writer, command string, err error) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonErrorEnvelope{
		OK:      false,
		Command: command,
		Error: jsonErrorBody{
			Code:    service.ErrorCodeForError(err),
			Message: service.ErrorMessageForError(err),
		},
	})
}

func jsonRequested(args []string) bool {
	for _, arg := range args {
		if arg == "--json" {
			return true
		}
		if strings.HasPrefix(arg, "--json=") {
			return !strings.HasSuffix(arg, "=false")
		}
	}
	return false
}

func commandNameFromArgs(args []string) string {
	if len(args) == 0 {
		return "odins"
	}

	var parts []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if len(parts) > 0 {
				break
			}
			continue
		}
		parts = append(parts, arg)
		if parts[0] != "domain" {
			break
		}
		if len(parts) == 2 {
			break
		}
	}

	if len(parts) == 0 {
		return "odins"
	}

	return strings.Join(parts, " ")
}

func renderJSONMaybe(cmdName string, out io.Writer, data any, warnings []string) error {
	if !outputJSON {
		return nil
	}
	return writeJSONSuccess(out, cmdName, data, warnings)
}

func commandWriter(cmd *cobra.Command) io.Writer {
	if cmd == nil {
		return os.Stdout
	}
	return cmd.OutOrStdout()
}

func writeTextLine(w io.Writer, format string, args ...any) {
	fmt.Fprintf(w, format+"\n", args...)
}
