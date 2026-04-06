package cert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
)

// CertDir returns the directory where ODINS stores certificates.
func CertDir() string {
	return filepath.Join(xdg.ConfigHome, "odins", "certs")
}

// CaddyCAPath returns the path to Caddy's internal root CA cert.
func CaddyCAPath() string {
	// Caddy stores its internal CA here by default on macOS
	candidates := []string{
		filepath.Join(xdg.DataHome, "caddy", "pki", "authorities", "local", "root.crt"),
		filepath.Join(os.Getenv("HOME"), ".local", "share", "caddy", "pki", "authorities", "local", "root.crt"),
		"/usr/local/var/lib/caddy/pki/authorities/local/root.crt",
		"/opt/homebrew/var/lib/caddy/pki/authorities/local/root.crt",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// TrustCaddyCA adds Caddy's internal root CA to the macOS system keychain.
// Must be called with elevated privileges (via helper).
func TrustCaddyCA() error {
	caPath := CaddyCAPath()
	if caPath == "" {
		return fmt.Errorf("caddy CA not found — run odins init with caddy started first")
	}
	out, err := exec.Command(
		"security", "add-trusted-cert",
		"-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain",
		caPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("trust caddy CA: %w\n%s", err, string(out))
	}
	return nil
}

// IssueMkcert issues a certificate for a domain using mkcert.
// The cert is placed in CertDir().
func IssueMkcert(domain string) error {
	if err := os.MkdirAll(CertDir(), 0755); err != nil {
		return err
	}
	certFile := filepath.Join(CertDir(), domain+".pem")
	keyFile := filepath.Join(CertDir(), domain+"-key.pem")

	out, err := exec.Command(
		"mkcert",
		"-cert-file", certFile,
		"-key-file", keyFile,
		domain,
		"*."+domain,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkcert %s: %w\n%s", domain, err, string(out))
	}
	return nil
}

// InstallMkcertCA installs the mkcert CA into the system trust store.
func InstallMkcertCA() error {
	out, err := exec.Command("mkcert", "-install").CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkcert -install: %w\n%s", err, string(out))
	}
	return nil
}
