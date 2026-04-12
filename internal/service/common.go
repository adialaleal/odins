package service

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/state"
)

func buildFQDN(subdomain, domain, project, tld string) string {
	if domain != "" {
		return subdomain + "." + domain + "." + tld
	}
	for _, c := range subdomain {
		if c == '.' {
			return subdomain + "." + tld
		}
	}
	return subdomain + "." + project + "." + tld
}

func normalizeDir(dir string) (string, error) {
	if dir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func validateBackend(value string) (config.ProxyBackend, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", string(config.BackendCaddy):
		return config.BackendCaddy, nil
	case string(config.BackendNginx):
		return config.BackendNginx, nil
	case string(config.BackendApache):
		return config.BackendApache, nil
	default:
		return "", invalidInput("proxy backend inválido: %s", value)
	}
}

func validateTLD(value string) (string, error) {
	if value == "" {
		return config.DefaultGlobalConfig().TLD, nil
	}
	for _, item := range config.SupportedTLDs {
		if item.TLD == value {
			return value, nil
		}
	}
	return "", invalidInput("TLD inválido: %s", value)
}

func projectConfigPath(dir string) string {
	return filepath.Join(dir, config.ProjectConfigFile)
}

func recommendedConfig(name string, detectedRuntime, detectedFramework string, port int) config.ProjectConfig {
	return config.ProjectConfig{
		Project: config.ProjectInfo{
			Name:      name,
			Runtime:   detectedRuntime,
			Framework: detectedFramework,
		},
		Routes: []config.RouteConfig{
			{
				Subdomain: name,
				Port:      port,
				HTTPS:     true,
			},
		},
	}
}

func defaultRouteProtocol(route state.Route) string {
	if route.HTTPS {
		return "https"
	}
	return "http"
}

func currentPlatformSupported() bool {
	return runtime.GOOS == "darwin"
}
