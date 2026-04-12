package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/service"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adrg/xdg"
)

func TestDetectJSONGolden(t *testing.T) {
	withTestXDG(t)

	fixtureDir := copyFixtureToTemp(t, filepath.Join("..", "testdata", "fixtures", "node-vite"), "node-vite")
	stdout, stderr, code := runCLI(t, "detect", "--json", "--dir", fixtureDir)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	normalized := sanitizeJSON(t, stdout, map[string]string{
		fixtureDir: "$FIXTURE_DIR",
	})
	assertGolden(t, filepath.Join("testdata", "golden", "detect-node-vite.json"), normalized)
}

func TestLSJSONGolden(t *testing.T) {
	withTestXDG(t)

	if err := config.SaveGlobal(config.DefaultGlobalConfig()); err != nil {
		t.Fatalf("SaveGlobal() error = %v", err)
	}

	store := &state.Store{
		Routes: []state.Route{
			{
				ID:        "odins-api.rankly.odin",
				Subdomain: "api.rankly.odin",
				Port:      1,
				Project:   "rankly",
				Runtime:   "node",
				HTTPS:     true,
				CreatedAt: time.Date(2026, 4, 6, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	if err := store.Save(); err != nil {
		t.Fatalf("store.Save() error = %v", err)
	}

	stdout, stderr, code := runCLI(t, "ls", "--json")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	assertGolden(t, filepath.Join("testdata", "golden", "ls.json"), sanitizeJSON(t, stdout, nil))
}

func TestDoctorJSONGolden(t *testing.T) {
	tmpRoot := withTestXDG(t)

	cfg := config.DefaultGlobalConfig()
	if err := config.SaveGlobal(cfg); err != nil {
		t.Fatalf("SaveGlobal() error = %v", err)
	}

	oldFactory := serviceFactory
	serviceFactory = func() *service.Manager {
		return service.New(service.Runtime{
			BrewInstalled:        func() bool { return true },
			BrewInstall:          func(string) error { return nil },
			BrewFormulaInstalled: func(string) bool { return true },
			BrewServiceStart:     func(string) error { return nil },
			BrewServiceRestart:   func(string) error { return nil },
			BrewServiceRunning:   func(string) bool { return true },
			GenerateDNSConfig:    func([]string, int) error { return nil },
			LinkDNSConfig:        func() error { return nil },
			SudoWriteResolver:    func(string, int) error { return nil },
			SudoTrustCA:          func(string) error { return nil },
			CaddyCAPath:          func() string { return filepath.Join(tmpRoot, "caddy", "root.crt") },
			InstallMkcertCA:      func() error { return nil },
			IssueMkcert:          func(string) error { return nil },
			ResolverPath:         func(string) string { return "/tmp/odins-test/resolver/odin" },
			FileExists: func(path string) bool {
				switch {
				case strings.Contains(path, "config.toml"):
					return true
				case strings.Contains(path, "dnsmasq.conf"):
					return true
				case strings.Contains(path, "/tmp/odins-test/resolver/odin"):
					return true
				case strings.Contains(path, "root.crt"):
					return true
				case strings.Contains(path, "certs"):
					return true
				default:
					return false
				}
			},
			ReadFile:          os.ReadFile,
			Sleep:             func(time.Duration) {},
			CaddyInit:         func(string) error { return nil },
			CaddyIsRunning:    func() bool { return true },
			CaddyAddRoute:     func(state.Route) error { return nil },
			CaddyRemoveRoute:  func(string) error { return nil },
			CaddyAddDomain:    func(string, string) error { return nil },
			CaddyRemoveDomain: func(string) error { return nil },
			NginxInit:         func() error { return nil },
			NginxIsRunning:    func() bool { return true },
			NginxAddRoute:     func(state.Route) error { return nil },
			NginxRemoveRoute:  func(string) error { return nil },
			ApacheIsRunning:   func() bool { return true },
			ApacheAddRoute:    func(state.Route) error { return nil },
			ApacheRemoveRoute: func(string) error { return nil },
		})
	}
	t.Cleanup(func() { serviceFactory = oldFactory })

	stdout, stderr, code := runCLI(t, "doctor", "--json")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	normalized := sanitizeJSON(t, stdout, map[string]string{
		tmpRoot: "$TMP",
	})
	assertGolden(t, filepath.Join("testdata", "golden", "doctor.json"), normalized)
}

func TestErrorJSONGolden(t *testing.T) {
	withTestXDG(t)

	stdout, stderr, code := runCLI(t, "kill", "missing.rankly.odin", "--json")
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	assertGolden(t, filepath.Join("testdata", "golden", "kill-missing-error.json"), sanitizeJSON(t, stdout, nil))
}

func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := ExecuteWithArgs(args, &stdout, &stderr)
	return stdout.String(), stderr.String(), code
}

func withTestXDG(t *testing.T) string {
	t.Helper()

	tmpRoot := t.TempDir()
	oldConfig := xdg.ConfigHome
	oldData := xdg.DataHome
	oldCache := xdg.CacheHome
	xdg.ConfigHome = filepath.Join(tmpRoot, "config")
	xdg.DataHome = filepath.Join(tmpRoot, "data")
	xdg.CacheHome = filepath.Join(tmpRoot, "cache")
	t.Cleanup(func() {
		xdg.ConfigHome = oldConfig
		xdg.DataHome = oldData
		xdg.CacheHome = oldCache
	})
	return tmpRoot
}

func copyFixtureToTemp(t *testing.T, src, name string) string {
	t.Helper()

	dst := filepath.Join(t.TempDir(), name)
	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	}); err != nil {
		t.Fatalf("copyFixtureToTemp() error = %v", err)
	}
	return dst
}

func sanitizeJSON(t *testing.T, raw string, replacements map[string]string) string {
	t.Helper()
	for from, to := range replacements {
		raw = strings.ReplaceAll(raw, from, to)
	}

	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, raw)
	}
	normalized, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}
	return string(normalized) + "\n"
}

func assertGolden(t *testing.T, path, got string) {
	t.Helper()

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(want) != got {
		t.Fatalf("golden mismatch for %s\nwant:\n%s\ngot:\n%s", path, string(want), got)
	}
}
