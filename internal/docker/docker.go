package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
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

// IsRunning returns true if the Docker daemon is accessible.
func IsRunning() bool {
	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: &http.Transport{},
	}
	resp, err := client.Get("http://localhost:2375/_ping")
	if err == nil {
		resp.Body.Close()
		return resp.StatusCode == 200
	}
	return dockerSocketExists()
}

// dockerSocketExists returns true if the Docker CLI can reach the daemon.
func dockerSocketExists() bool {
	return exec.Command("docker", "info").Run() == nil
}

// FindContainerPort looks up the host port mapped to a container's internal port.
// Uses "docker port <container> <port>" and parses the output.
func FindContainerPort(ctx context.Context, containerName string, internalPort int) (string, error) {
	out, err := exec.CommandContext(ctx, "docker", "port", containerName, fmt.Sprintf("%d", internalPort)).Output()
	if err != nil {
		// Docker not available or container not running — fall back to direct port
		return fmt.Sprintf("127.0.0.1:%d", internalPort), nil
	}
	// Output format: "0.0.0.0:32768\n" or ":::32768\n"
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	if line == "" {
		return fmt.Sprintf("127.0.0.1:%d", internalPort), nil
	}
	// Extract port from "0.0.0.0:PORT" or ":::PORT"
	_, port, err := net.SplitHostPort(line)
	if err != nil {
		return fmt.Sprintf("127.0.0.1:%d", internalPort), nil
	}
	return "127.0.0.1:" + port, nil
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
