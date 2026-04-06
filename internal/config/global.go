package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// SupportedTLDs lists all TLDs ODINS can configure.
var SupportedTLDs = []struct {
	TLD     string
	Label   string
	Warning string
}{
	{TLD: "odin", Label: ".odin — temático, sem conflitos (padrão)"},
	{TLD: "odins", Label: ".odins — variante temática"},
	{TLD: "test", Label: ".test — reservado IANA, sem HSTS"},
	{TLD: "dev", Label: ".dev — popular (requer HTTPS, HSTS preloaded)"},
	{TLD: "lan", Label: ".lan — comum em redes locais"},
	{TLD: "internal", Label: ".internal — uso corporativo"},
	{TLD: "local", Label: ".local — ⚠️  conflito com mDNS/Bonjour", Warning: "mDNS conflict"},
}

// ProxyBackend represents a supported reverse proxy.
type ProxyBackend string

const (
	BackendCaddy  ProxyBackend = "caddy"
	BackendNginx  ProxyBackend = "nginx"
	BackendApache ProxyBackend = "apache"
)

// GlobalConfig holds ODINS global settings stored at ~/.config/odins/config.toml.
type GlobalConfig struct {
	TLD            string       `toml:"tld" json:"tld"`
	ProxyBackend   ProxyBackend `toml:"proxy_backend" json:"proxy_backend"`
	DnsmasqPort    int          `toml:"dnsmasq_port" json:"dnsmasq_port"`
	CaddyAdmin     string       `toml:"caddy_admin" json:"caddy_admin"`
	HTTPPort       int          `toml:"http_port" json:"http_port"`
	HTTPSPort      int          `toml:"https_port" json:"https_port"`
	OnboardingDone bool         `toml:"onboarding_done" json:"onboarding_done"`
	// Language overrides auto-detection. Values: "pt", "en", "es". Empty = auto.
	Language string `toml:"language,omitempty" json:"language,omitempty"`
}

// DefaultGlobalConfig returns the default global configuration.
func DefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		TLD:          "odin",
		ProxyBackend: BackendCaddy,
		DnsmasqPort:  5300,
		CaddyAdmin:   "http://localhost:2019",
		HTTPPort:     80,
		HTTPSPort:    443,
	}
}

// ConfigDir returns the ODINS config directory path.
func ConfigDir() string {
	return filepath.Join(xdg.ConfigHome, "odins")
}

// ConfigPath returns the path to the global config file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

// LoadGlobal loads the global config from disk. Returns defaults if not found.
func LoadGlobal() (GlobalConfig, error) {
	cfg := DefaultGlobalConfig()
	path := ConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// SaveGlobal writes the global config to disk.
func SaveGlobal(cfg GlobalConfig) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(ConfigPath())
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
