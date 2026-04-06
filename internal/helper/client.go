// Package helper provides a client for the odins-helper privileged binary.
// The helper runs as root via launchd and handles operations that require
// elevated privileges: writing /etc/resolver files and trusting CA certificates.
package helper

import (
	"encoding/json"
	"fmt"
	"net"
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
	// The helper binary is embedded or downloaded alongside odins
	// For now, use sudo to write the resolver directly
	return nil
}

// SudoWriteResolver is a fallback that writes /etc/resolver/<tld> via sudo.
func SudoWriteResolver(tld string, port int) error {
	content := fmt.Sprintf("nameserver 127.0.0.1\nport %d\n", port)
	tmpFile := fmt.Sprintf("/tmp/odins-resolver-%s", tld)

	if err := writeTemp(tmpFile, content); err != nil {
		return err
	}

	out, err := exec.Command(
		"sudo", "-p",
		fmt.Sprintf("[ODINS] Autorização para criar /etc/resolver/%s (DNS local): ", tld),
		"cp", tmpFile,
		fmt.Sprintf("/etc/resolver/%s", tld),
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("sudo write resolver: %w\n%s", err, string(out))
	}

	return nil
}

// SudoTrustCA adds a CA cert to the system keychain via sudo.
func SudoTrustCA(certPath string) error {
	out, err := exec.Command(
		"sudo", "-p",
		"[ODINS] Autorização para confiar no certificado HTTPS local: ",
		"security", "add-trusted-cert",
		"-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain",
		certPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("sudo trust CA: %w\n%s", err, string(out))
	}
	return nil
}

func writeTemp(path, content string) error {
	return exec.Command("bash", "-c", fmt.Sprintf("echo %q > %s", content, path)).Run()
}
