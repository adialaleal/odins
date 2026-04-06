package dns

import (
	"fmt"
	"os"
	"path/filepath"
)

const resolverDir = "/etc/resolver"

// ResolverContent returns the content for /etc/resolver/<tld>.
func ResolverContent(port int) string {
	return fmt.Sprintf("nameserver 127.0.0.1\nport %d\n", port)
}

// ResolverPath returns the path for a given TLD resolver file.
func ResolverPath(tld string) string {
	return filepath.Join(resolverDir, tld)
}

// NeedsPrivilegedWrite returns true — /etc/resolver requires root.
func NeedsPrivilegedWrite() bool {
	info, err := os.Stat(resolverDir)
	if err != nil {
		return true
	}
	return info.Mode().Perm()&0200 == 0 // not writable by current user
}

// WriteResolver writes /etc/resolver/<tld> directly (requires root, called from helper).
func WriteResolver(tld string, port int) error {
	if err := os.MkdirAll(resolverDir, 0755); err != nil {
		return fmt.Errorf("mkdir /etc/resolver: %w", err)
	}
	path := ResolverPath(tld)
	return os.WriteFile(path, []byte(ResolverContent(port)), 0644)
}

// RemoveResolver deletes /etc/resolver/<tld> (requires root, called from helper).
func RemoveResolver(tld string) error {
	path := ResolverPath(tld)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
