package aiclient

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// MCPEntry is the MCP server config entry written to AI client configs.
type MCPEntry struct {
	Command string   `json:"command" toml:"command"`
	Args    []string `json:"args" toml:"args"`
}

// OdinsBinaryPath returns the path to the odins binary.
func OdinsBinaryPath() string {
	if p, err := exec.LookPath("odins"); err == nil {
		return p
	}
	if p, err := os.Executable(); err == nil {
		return p
	}
	return "odins"
}

// ConfigureClient writes the ODINS MCP server entry into a client's config.
func ConfigureClient(client AIClient, odinsBinary string) error {
	if client.ConfigPath == "" {
		return fmt.Errorf("sem path de configuração para %s", client.Name)
	}

	if client.Format == "toml" {
		return configureTOML(client, odinsBinary)
	}
	return configureJSON(client, odinsBinary)
}

func configureJSON(client AIClient, binary string) error {
	// Ensure parent directory exists.
	dir := filepath.Dir(client.ConfigPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Load existing config or start fresh.
	data := make(map[string]any)
	if raw, err := os.ReadFile(client.ConfigPath); err == nil {
		json.Unmarshal(raw, &data) // ignore parse errors, start fresh on corrupt
	}

	// Navigate to or create the root key.
	rootKey := client.RootKey
	servers, ok := data[rootKey].(map[string]any)
	if !ok {
		servers = make(map[string]any)
	}

	// Don't overwrite existing odins entry.
	if _, exists := servers["odins"]; exists {
		return nil
	}

	entry := map[string]any{
		"command": binary,
		"args":    []string{"mcp"},
	}

	// VS Code and Cursor need a "type" field.
	if client.RootKey == "servers" || client.Name == "Cursor" {
		entry["type"] = "stdio"
	}

	servers["odins"] = entry
	data[rootKey] = servers

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(client.ConfigPath, out, 0o644)
}

func configureTOML(client AIClient, binary string) error {
	// Ensure parent directory exists.
	dir := filepath.Dir(client.ConfigPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Load existing config or start fresh.
	data := make(map[string]any)
	if raw, err := os.ReadFile(client.ConfigPath); err == nil {
		toml.Unmarshal(raw, &data)
	}

	// Navigate to mcp_servers section.
	rootKey := client.RootKey
	servers, ok := data[rootKey].(map[string]any)
	if !ok {
		servers = make(map[string]any)
	}

	// Don't overwrite existing odins entry.
	if _, exists := servers["odins"]; exists {
		return nil
	}

	servers["odins"] = MCPEntry{
		Command: binary,
		Args:    []string{"mcp"},
	}
	data[rootKey] = servers

	f, err := os.Create(client.ConfigPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(data)
}

// ConfiguredResult describes the outcome of configuring a client.
type ConfiguredResult struct {
	Client     string `json:"client"`
	ConfigPath string `json:"config_path"`
	Configured bool   `json:"configured"`
	Skipped    bool   `json:"skipped"`
	Error      string `json:"error,omitempty"`
}

// ConfigureAll configures all detected clients and returns results.
func ConfigureAll(clients []AIClient) []ConfiguredResult {
	binary := OdinsBinaryPath()
	var results []ConfiguredResult

	for _, c := range clients {
		r := ConfiguredResult{
			Client:     c.Name,
			ConfigPath: c.ConfigPath,
		}

		err := ConfigureClient(c, binary)
		if err != nil {
			r.Error = err.Error()
		} else {
			r.Configured = true
		}

		results = append(results, r)
	}

	return results
}
