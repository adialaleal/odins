// Package helper provides a client for the odins-helper privileged binary.
// The helper runs as root via launchd and handles operations that require
// elevated privileges: writing /etc/resolver files and trusting CA certificates.
package helper

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const socketPath = "/var/run/odins-helper.sock"

// Op is the type of privileged operation to perform.
type Op string

const (
	OpWriteResolver  Op = "write_resolver"
	OpRemoveResolver Op = "remove_resolver"
	OpTrustCA        Op = "trust_ca"
)

// Request is sent to the helper over the Unix socket.
type Request struct {
	Op      Op                `json:"op"`
	Payload map[string]string `json:"payload"`
}

// Response is returned by the helper.
type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// IsRunning returns true if the helper daemon is listening.
func IsRunning() bool {
	conn, err := net.DialTimeout("unix", socketPath, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// WriteResolver instructs the helper to write /etc/resolver/<tld>.
func WriteResolver(tld string, port int) error {
	return call(Request{
		Op: OpWriteResolver,
		Payload: map[string]string{
			"tld":  tld,
			"port": fmt.Sprintf("%d", port),
		},
	})
}

// RemoveResolver instructs the helper to delete /etc/resolver/<tld>.
func RemoveResolver(tld string) error {
	return call(Request{
		Op:      OpRemoveResolver,
		Payload: map[string]string{"tld": tld},
	})
}

// TrustCA instructs the helper to add a CA cert to the system keychain.
func TrustCA(certPath string) error {
	return call(Request{
		Op:      OpTrustCA,
		Payload: map[string]string{"cert_path": certPath},
	})
}

func call(req Request) error {
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return fmt.Errorf("helper not running: %w", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return fmt.Errorf("helper send: %w", err)
	}

	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return fmt.Errorf("helper recv: %w", err)
	}

	if !resp.OK {
		return fmt.Errorf("helper error: %s", resp.Error)
	}
	return nil
}

// InstallHelper installs the odins-helper via sudo.
// This is called once during odins init.
func InstallHelper() error {
	return nil
}

// osascriptRun executes a shell command with macOS administrator privileges
// by showing the native system authentication dialog (like any proper Mac app).
// Falls back to sudo on the terminal if osascript fails.
func osascriptRun(shellCmd string) error {
	// Escape double-quotes inside the shell command for the AppleScript string
	escaped := ""
	for _, c := range shellCmd {
		if c == '"' || c == '\\' {
			escaped += "\\"
		}
		escaped += string(c)
	}
	script := fmt.Sprintf(`do shell script "%s" with administrator privileges`, escaped)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SudoWriteResolver writes /etc/resolver/<tld> using the macOS authentication
// dialog (osascript). Creates the /etc/resolver/ directory if needed.
func SudoWriteResolver(tld string, port int) error {
	content := fmt.Sprintf("nameserver 127.0.0.1\\nport %d\\n", port)
	resolverPath := fmt.Sprintf("/etc/resolver/%s", tld)
	// Build a single shell command: create dir, write file via printf
	shellCmd := fmt.Sprintf(
		"mkdir -p /etc/resolver && printf '%s' > %s",
		content, resolverPath,
	)
	if err := osascriptRun(shellCmd); err != nil {
		// Fallback: interactive sudo on the terminal
		prompt := fmt.Sprintf("\n[ODINS] Senha para criar /etc/resolver/%s (DNS local): ", tld)
		fileContent := fmt.Sprintf("nameserver 127.0.0.1\nport %d\n", port)
		tmpFile := fmt.Sprintf("/tmp/odins-resolver-%s", tld)
		if werr := os.WriteFile(tmpFile, []byte(fileContent), 0644); werr != nil {
			return fmt.Errorf("write temp resolver: %w", werr)
		}
		shell := fmt.Sprintf("mkdir -p /etc/resolver && cp %s /etc/resolver/%s", tmpFile, tld)
		cmd2 := exec.Command("sudo", "-p", prompt, "bash", "-c", shell)
		cmd2.Stdin = os.Stdin
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		if err2 := cmd2.Run(); err2 != nil {
			return fmt.Errorf("write resolver: %w", err2)
		}
	}
	return nil
}

// SudoTrustCA adds a CA cert to the system keychain using the macOS
// authentication dialog (osascript), falling back to sudo.
func SudoTrustCA(certPath string) error {
	shellCmd := fmt.Sprintf(
		"security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s",
		certPath,
	)
	if err := osascriptRun(shellCmd); err != nil {
		// Fallback: interactive sudo on the terminal
		cmd2 := exec.Command(
			"sudo", "-p",
			"\n[ODINS] Senha para confiar no certificado HTTPS local: ",
			"security", "add-trusted-cert",
			"-d", "-r", "trustRoot",
			"-k", "/Library/Keychains/System.keychain",
			certPath,
		)
		cmd2.Stdin = os.Stdin
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		if err2 := cmd2.Run(); err2 != nil {
			return fmt.Errorf("trust CA: %w", err2)
		}
	}
	return nil
}

// SudoFlushDNS flushes the macOS DNS cache and restarts mDNSResponder,
// using the macOS authentication dialog for the privileged killall call.
func SudoFlushDNS() {
	// dscacheutil doesn't need root
	exec.Command("dscacheutil", "-flushcache").Run()
	// killall -HUP mDNSResponder needs root on macOS 10.10+
	osascriptRun("killall -HUP mDNSResponder")
}
