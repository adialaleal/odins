package mcp

import (
	"github.com/adialaleal/odins/internal/service"
	"github.com/mark3labs/mcp-go/server"
)

const serverInstructions = `ODINS — Local DNS + Reverse Proxy manager for macOS developers.

Operational rules:
- Inspect first, change later. Use odins_detect before creating or modifying routes.
- Use odins_doctor to diagnose problems before manual troubleshooting.
- odins_init may require sudo for DNS resolver and certificate trust.
- When a project has a .odins file, prefer odins_up over individual odins_add calls.
- Domain workspaces group related services: create with odins_domain_add, then use domain field in routes.
- Use odins_scan to discover all projects under a directory tree.`

// NewServer creates a configured MCP server wrapping the ODINS service layer.
func NewServer(mgr *service.Manager, version string) *server.MCPServer {
	s := server.NewMCPServer("odins", version,
		server.WithToolCapabilities(false),
		server.WithResourceCapabilities(true, false),
		server.WithInstructions(serverInstructions),
		server.WithRecovery(),
	)

	registerTools(s, mgr)
	registerResources(s)

	return s
}
