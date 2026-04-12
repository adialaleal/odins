package mcp

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/state"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerResources(s *server.MCPServer) {
	registerConfigResource(s)
	registerRoutesResource(s)
	registerDomainsResource(s)
}

func registerConfigResource(s *server.MCPServer) {
	resource := mcp.NewResource(
		"odins://config",
		"ODINS Global Configuration",
		mcp.WithResourceDescription("Global ODINS configuration: TLD, proxy backend, ports, onboarding status"),
		mcp.WithMIMEType("application/toml"),
	)
	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		cfg, err := config.LoadGlobal()
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "odins://config",
				MIMEType: "application/toml",
				Text:     buf.String(),
			},
		}, nil
	})
}

func registerRoutesResource(s *server.MCPServer) {
	resource := mcp.NewResource(
		"odins://routes",
		"Active Routes",
		mcp.WithResourceDescription("All currently registered proxy routes"),
		mcp.WithMIMEType("application/json"),
	)
	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		store, err := state.Load()
		if err != nil {
			return nil, err
		}
		data, err := json.MarshalIndent(store.Routes, "", "  ")
		if err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "odins://routes",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

func registerDomainsResource(s *server.MCPServer) {
	resource := mcp.NewResource(
		"odins://domains",
		"Domain Workspaces",
		mcp.WithResourceDescription("All registered domain workspaces"),
		mcp.WithMIMEType("application/json"),
	)
	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		store, err := state.Load()
		if err != nil {
			return nil, err
		}
		data, err := json.MarshalIndent(store.Domains, "", "  ")
		if err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "odins://domains",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}
