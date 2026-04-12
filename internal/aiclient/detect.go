package aiclient

import (
	"os"
	"os/exec"
	"path/filepath"
)

// AIClient represents a detected AI coding tool.
type AIClient struct {
	Name       string `json:"name"`
	ConfigPath string `json:"config_path"`
	Format     string `json:"format"`   // "json" or "toml"
	RootKey    string `json:"root_key"` // "mcpServers", "servers", or "mcp_servers"
	Detected   bool   `json:"detected"`
}

// DetectClients checks for known AI coding tools on the machine.
func DetectClients() []AIClient {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	clients := []AIClient{
		detectClaudeCode(home),
		detectVSCode(home),
		detectCursor(home),
		detectCodex(home),
		detectGemini(home),
	}

	var found []AIClient
	for _, c := range clients {
		if c.Detected {
			found = append(found, c)
		}
	}
	return found
}

// DetectAllClients returns all clients (detected or not) for display.
func DetectAllClients() []AIClient {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []AIClient{
		detectClaudeCode(home),
		detectVSCode(home),
		detectCursor(home),
		detectCodex(home),
		detectGemini(home),
	}
}

func detectClaudeCode(home string) AIClient {
	c := AIClient{
		Name:    "Claude Code",
		Format:  "json",
		RootKey: "mcpServers",
	}

	// Check for claude binary on PATH.
	if _, err := exec.LookPath("claude"); err == nil {
		c.Detected = true
	}

	// Config at ~/.claude.json (global).
	path := filepath.Join(home, ".claude.json")
	if _, err := os.Stat(path); err == nil {
		c.Detected = true
		c.ConfigPath = path
		return c
	}

	// If binary found but no config yet, set the target path.
	if c.Detected {
		c.ConfigPath = path
	}

	return c
}

func detectVSCode(home string) AIClient {
	c := AIClient{
		Name:    "VS Code",
		Format:  "json",
		RootKey: "servers",
	}

	dir := filepath.Join(home, ".vscode")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		c.Detected = true
		c.ConfigPath = filepath.Join(dir, "mcp.json")
	}

	return c
}

func detectCursor(home string) AIClient {
	c := AIClient{
		Name:    "Cursor",
		Format:  "json",
		RootKey: "mcpServers",
	}

	dir := filepath.Join(home, ".cursor")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		c.Detected = true
		c.ConfigPath = filepath.Join(dir, "mcp.json")
	}

	return c
}

func detectCodex(home string) AIClient {
	c := AIClient{
		Name:    "Codex CLI",
		Format:  "toml",
		RootKey: "mcp_servers",
	}

	path := filepath.Join(home, ".codex", "config.toml")
	if _, err := os.Stat(path); err == nil {
		c.Detected = true
		c.ConfigPath = path
		return c
	}

	if _, err := exec.LookPath("codex"); err == nil {
		c.Detected = true
		c.ConfigPath = path
	}

	return c
}

func detectGemini(home string) AIClient {
	c := AIClient{
		Name:    "Gemini CLI",
		Format:  "json",
		RootKey: "mcpServers",
	}

	path := filepath.Join(home, ".gemini", "settings.json")
	if _, err := os.Stat(path); err == nil {
		c.Detected = true
		c.ConfigPath = path
		return c
	}

	if _, err := exec.LookPath("gemini"); err == nil {
		c.Detected = true
		c.ConfigPath = path
	}

	return c
}
