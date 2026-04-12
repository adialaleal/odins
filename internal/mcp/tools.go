package mcp

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/adialaleal/odins/internal/scanner"
	"github.com/adialaleal/odins/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer, mgr *service.Manager) {
	registerDetect(s, mgr)
	registerDoctor(s, mgr)
	registerLs(s, mgr)
	registerDomainLs(s, mgr)
	registerAdd(s, mgr)
	registerUp(s, mgr)
	registerDown(s, mgr)
	registerKill(s, mgr)
	registerDomainAdd(s, mgr)
	registerDomainRm(s, mgr)
	registerScan(s)
	registerOpen(s)
	registerInit(s, mgr)
}

// serviceResultToMCP converts a service result to an MCP tool result.
func serviceResultToMCP(data any, warnings []string, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(service.ErrorMessageForError(err)), nil
	}
	if warnings == nil {
		warnings = []string{}
	}
	envelope := map[string]any{"data": data, "warnings": warnings}
	jsonBytes, _ := json.MarshalIndent(envelope, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// --- odins_detect ---

func registerDetect(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_detect",
		mcp.WithDescription("Detect project runtime, framework, port, and recommended .odins config for a directory"),
		mcp.WithString("directory", mcp.Required(), mcp.Description("Absolute path to the project directory")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir, err := req.RequireString("directory")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, warnings, svcErr := mgr.Detect(dir)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_doctor ---

func registerDoctor(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_doctor",
		mcp.WithDescription("Run environment diagnostics: DNS, proxy, HTTPS, store health"),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, warnings, svcErr := mgr.Doctor()
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_ls ---

func registerLs(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_ls",
		mcp.WithDescription("List all active routes with live upstream status checks"),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, warnings, svcErr := mgr.ListRoutes()
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_domain_ls ---

func registerDomainLs(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_domain_ls",
		mcp.WithDescription("List all registered domain workspaces with service counts"),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, warnings, svcErr := mgr.DomainList()
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_add ---

func registerAdd(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_add",
		mcp.WithDescription("Add a reverse proxy route from a local subdomain to a running service"),
		mcp.WithString("subdomain", mcp.Required(), mcp.Description("Subdomain FQDN, e.g. api.tatoh.odin")),
		mcp.WithNumber("port", mcp.Required(), mcp.Description("Local port to proxy to (1-65535)")),
		mcp.WithString("docker", mcp.Description("Docker container name (optional)")),
		mcp.WithString("project", mcp.Description("Project name (inferred from subdomain if omitted)")),
		mcp.WithString("domain", mcp.Description("Domain workspace name, e.g. tatoh")),
		mcp.WithBoolean("https", mcp.Description("Enable HTTPS (default: true)")),
		mcp.WithIdempotentHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		subdomain, err := req.RequireString("subdomain")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		port, err := req.RequireInt("port")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		opts := service.AddRouteOptions{
			Subdomain: subdomain,
			Port:      port,
			Docker:    req.GetString("docker", ""),
			Project:   req.GetString("project", ""),
			Domain:    req.GetString("domain", ""),
			HTTPS:     req.GetBool("https", true),
		}

		result, warnings, svcErr := mgr.AddRoute(opts)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_up ---

func registerUp(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_up",
		mcp.WithDescription("Apply routes from the .odins config in a directory. Creates .odins from auto-detection if missing."),
		mcp.WithString("directory", mcp.Required(), mcp.Description("Absolute path to the project directory")),
		mcp.WithIdempotentHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir, err := req.RequireString("directory")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, warnings, svcErr := mgr.Up(dir)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_down ---

func registerDown(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_down",
		mcp.WithDescription("Remove all routes declared in the .odins config of a directory"),
		mcp.WithString("directory", mcp.Required(), mcp.Description("Absolute path to the project directory")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir, err := req.RequireString("directory")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, warnings, svcErr := mgr.Down(dir)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_kill ---

func registerKill(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_kill",
		mcp.WithDescription("Remove a single route by its FQDN subdomain"),
		mcp.WithString("subdomain", mcp.Required(), mcp.Description("Full FQDN to remove, e.g. api.tatoh.odin")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		subdomain, err := req.RequireString("subdomain")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, warnings, svcErr := mgr.Kill(subdomain)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_domain_add ---

func registerDomainAdd(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_domain_add",
		mcp.WithDescription("Create a domain workspace with an auto-generated landing page"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Domain name, e.g. tatoh")),
		mcp.WithString("title", mcp.Description("Title for the landing page (defaults to name)")),
		mcp.WithString("description", mcp.Description("Description for the domain workspace")),
		mcp.WithIdempotentHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		title := req.GetString("title", "")
		desc := req.GetString("description", "")
		result, warnings, svcErr := mgr.DomainAdd(name, title, desc)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_domain_rm ---

func registerDomainRm(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_domain_rm",
		mcp.WithDescription("Remove a domain workspace and all its attached routes"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Domain name to remove")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, warnings, svcErr := mgr.DomainRemove(name)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}

// --- odins_scan ---

func registerScan(s *server.MCPServer) {
	tool := mcp.NewTool("odins_scan",
		mcp.WithDescription("Scan a directory tree for projects, detect their stacks, and update the global project index"),
		mcp.WithString("directory", mcp.Description("Root directory to scan (default: ~/Projects)")),
		mcp.WithNumber("max_depth", mcp.Description("Maximum directory depth to walk (default: 3)")),
		mcp.WithBoolean("create_odins", mcp.Description("Create .odins files where missing (default: false)")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		maxDepth := 3
		if v, err := req.RequireInt("max_depth"); err == nil {
			maxDepth = v
		}

		opts := scanner.ScanOptions{
			Directory:   req.GetString("directory", ""),
			MaxDepth:    maxDepth,
			CreateOdins: req.GetBool("create_odins", false),
		}

		result, err := scanner.Scan(opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := scanner.UpdateIndex(result); err != nil {
			return mcp.NewToolResultError("scan concluído mas falha ao atualizar índice: " + err.Error()), nil
		}

		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(jsonBytes)), nil
	})
}

// --- odins_open ---

func registerOpen(s *server.MCPServer) {
	tool := mcp.NewTool("odins_open",
		mcp.WithDescription("Open a local FQDN in the default browser"),
		mcp.WithString("fqdn", mcp.Required(), mcp.Description("Full FQDN to open, e.g. api.tatoh.odin")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fqdn, err := req.RequireString("fqdn")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		url := "https://" + fqdn
		if err := exec.Command("open", url).Run(); err != nil {
			return mcp.NewToolResultError("falha ao abrir browser: " + err.Error()), nil
		}
		return mcp.NewToolResultText("aberto: " + url), nil
	})
}

// --- odins_init ---

func registerInit(s *server.MCPServer, mgr *service.Manager) {
	tool := mcp.NewTool("odins_init",
		mcp.WithDescription("One-time ODINS setup: install DNS, proxy, HTTPS. May require sudo. Only needed once per machine."),
		mcp.WithString("tld", mcp.Description("TLD for local domains (default: odin)")),
		mcp.WithString("backend", mcp.Description("Reverse proxy backend: caddy, nginx, or apache (default: caddy)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := service.InitOptions{
			NonInteractive: true, // MCP runs without terminal
			TLD:            req.GetString("tld", ""),
			Backend:        req.GetString("backend", ""),
		}
		result, warnings, svcErr := mgr.Init(opts)
		return serviceResultToMCP(result, warnings, svcErr)
	})
}
