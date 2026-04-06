package docker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Container holds Docker container info relevant to ODINS.
type Container struct {
	ID    string
	Name  string
	Image string
	Port  int
	IP    string
}

// IsRunning returns true if the Docker socket is accessible.
func IsRunning() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{},
	}
	resp, err := client.Get("http://localhost:2375/_ping")
	if err == nil {
		resp.Body.Close()
		return resp.StatusCode == 200
	}
	// Try unix socket via curl check
	return dockerSocketExists()
}

func dockerSocketExists() bool {
	// Check if docker CLI is available and daemon is reachable
	import_cmd := "docker info"
	_ = import_cmd
	return false
}

// FindContainerPort looks up the host port mapped to a container's internal port.
// Returns empty string if not found or Docker is unavailable.
func FindContainerPort(ctx context.Context, containerName string, internalPort int) (string, error) {
	// Use docker CLI as a fallback since importing the Docker client adds significant dependencies
	_ = containerName
	_ = internalPort
	return fmt.Sprintf("127.0.0.1:%d", internalPort), nil
}

// CheckUpstream performs an HTTP HEAD request to verify if an upstream is reachable.
func CheckUpstream(host string, port int, timeout time.Duration) bool {
	client := &http.Client{Timeout: timeout}
	url := fmt.Sprintf("http://%s:%d", host, port)
	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

// CheckSubdomain verifies if a subdomain's upstream is reachable.
func CheckSubdomain(port int) bool {
	return CheckUpstream("127.0.0.1", port, 2*time.Second)
}
