package service

import (
	"testing"
	"time"

	"github.com/adialaleal/odins/internal/state"
)

func TestInitNonInteractiveDefaults(t *testing.T) {
	t.Parallel()

	manager := New(Runtime{
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
		CaddyCAPath:          func() string { return "/tmp/caddy/root.crt" },
		InstallMkcertCA:      func() error { return nil },
		IssueMkcert:          func(string) error { return nil },
		ResolverPath:         func(string) string { return "/tmp/resolver/odin" },
		FileExists:           func(string) bool { return true },
		ReadFile:             func(string) ([]byte, error) { return nil, nil },
		Sleep:                func(_ time.Duration) {},
		CaddyInit:            func(string) error { return nil },
		CaddyIsRunning:       func() bool { return true },
		CaddyAddRoute:        func(state.Route) error { return nil },
		CaddyRemoveRoute:     func(string) error { return nil },
		CaddyAddDomain:       func(string, string) error { return nil },
		CaddyRemoveDomain:    func(string) error { return nil },
		NginxInit:            func() error { return nil },
		NginxIsRunning:       func() bool { return true },
		NginxAddRoute:        func(state.Route) error { return nil },
		NginxRemoveRoute:     func(string) error { return nil },
		ApacheIsRunning:      func() bool { return true },
		ApacheAddRoute:       func(state.Route) error { return nil },
		ApacheRemoveRoute:    func(string) error { return nil },
	})

	result, _, err := manager.Init(InitOptions{NonInteractive: true})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if result.TLD != "odin" {
		t.Fatalf("TLD = %q, want %q", result.TLD, "odin")
	}
	if result.Backend != "caddy" {
		t.Fatalf("Backend = %q, want %q", result.Backend, "caddy")
	}
}

func TestInitInvalidBackend(t *testing.T) {
	t.Parallel()

	manager := New(DefaultRuntime())
	_, _, err := manager.Init(InitOptions{NonInteractive: true, Backend: "bad-backend"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ErrorCodeForError(err) != CodeInvalidInput {
		t.Fatalf("error code = %q, want %q", ErrorCodeForError(err), CodeInvalidInput)
	}
	if ExitCodeForError(err) != 2 {
		t.Fatalf("exit code = %d, want 2", ExitCodeForError(err))
	}
}
