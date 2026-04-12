package cmd

import (
	"log"
	"os"

	odinsmcp "github.com/adialaleal/odins/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server for AI coding tools (stdio)",
	Long: `Start the ODINS MCP (Model Context Protocol) server over stdio.

This exposes all ODINS operations as MCP tools that AI coding tools
can discover and use natively. Supported clients include Claude Code,
VS Code Copilot, Cursor, Codex CLI, and Gemini CLI.

Configuration example for Claude Code (.mcp.json):
  { "mcpServers": { "odins": { "command": "odins", "args": ["mcp"] } } }`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// All logging must go to stderr — stdout is the JSON-RPC channel.
		log.SetOutput(os.Stderr)

		mgr := serviceFactory()
		version := buildVersion
		if version == "" {
			version = "dev"
		}

		s := odinsmcp.NewServer(mgr, version)
		return server.ServeStdio(s)
	},
}
