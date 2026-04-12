package detect

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ComposeService represents a service detected from a docker-compose file.
type ComposeService struct {
	Name string `json:"name"`
	Port int    `json:"port"` // host-side port from the ports mapping
}

// ParseComposeServices reads a docker-compose file in dir and returns
// detected services with their host-mapped ports.
// Supports "host:container" and bare "port" formats under the ports key.
func ParseComposeServices(dir string) []ComposeService {
	for _, name := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		return parseComposeYAML(string(data))
	}
	return nil
}

// parseComposeYAML performs a minimal line-by-line parse of docker-compose YAML
// to extract service names and host ports without pulling in a YAML dependency.
func parseComposeYAML(content string) []ComposeService {
	var services []ComposeService
	lines := strings.Split(content, "\n")

	var currentService string
	inServices := false
	inPorts := false
	serviceIndent := -1
	portsIndent := -1

	for _, raw := range lines {
		line := strings.TrimRight(raw, " \t\r")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		indent := leadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		// Top-level "services:" key
		if trimmed == "services:" {
			inServices = true
			serviceIndent = -1
			inPorts = false
			continue
		}

		if !inServices {
			continue
		}

		// Reset if we leave the services block (back to top level)
		if indent == 0 && trimmed != "services:" {
			inServices = false
			continue
		}

		// Detect service name (first indented key inside services)
		if serviceIndent < 0 && indent > 0 && strings.HasSuffix(trimmed, ":") {
			serviceIndent = indent
		}

		if indent == serviceIndent && strings.HasSuffix(trimmed, ":") {
			currentService = strings.TrimSuffix(trimmed, ":")
			inPorts = false
			portsIndent = -1
			continue
		}

		if currentService == "" {
			continue
		}

		// Detect "ports:" key under a service
		if trimmed == "ports:" {
			inPorts = true
			portsIndent = indent
			continue
		}

		// Leave ports block if we dedent
		if inPorts && indent <= portsIndent && trimmed != "ports:" {
			inPorts = false
		}

		// Parse port entries (list items starting with "- ")
		if inPorts && strings.HasPrefix(trimmed, "- ") {
			portStr := strings.TrimPrefix(trimmed, "- ")
			portStr = strings.Trim(portStr, "\"'")
			port := parseHostPort(portStr)
			if port > 0 {
				services = appendOrUpdate(services, currentService, port)
			}
		}
	}

	return services
}

// parseHostPort extracts the host port from formats like "8080:8080", "8080", "127.0.0.1:8080:8080".
func parseHostPort(s string) int {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		// bare port number
		p, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		return p
	case 2:
		// host:container
		p, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		return p
	case 3:
		// ip:host:container
		p, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		return p
	}
	return 0
}

// appendOrUpdate adds a service or skips if already present (first port wins).
func appendOrUpdate(services []ComposeService, name string, port int) []ComposeService {
	for _, s := range services {
		if s.Name == name {
			return services
		}
	}
	return append(services, ComposeService{Name: name, Port: port})
}

func leadingSpaces(s string) int {
	count := 0
	for _, c := range s {
		if c == ' ' {
			count++
		} else if c == '\t' {
			count += 2
		} else {
			break
		}
	}
	return count
}
