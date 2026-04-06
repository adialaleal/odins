package cmd

import (
	"fmt"

	"github.com/adialaleal/odins/internal/service"
	"github.com/spf13/cobra"
)

func exactArgs(expected int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != expected {
			return service.InvalidInput(fmt.Sprintf("%q aceita exatamente %d argumento(s)", cmd.CommandPath(), expected))
		}
		return nil
	}
}
