package service

import (
	"os"
	"path/filepath"
	"time"

	"github.com/adialaleal/odins/internal/cert"
	"github.com/adialaleal/odins/internal/dns"
	"github.com/adialaleal/odins/internal/helper"
	"github.com/adialaleal/odins/internal/proxy/apache"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/proxy/nginx"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adialaleal/odins/pkg/brew"
)

// Runtime holds the side-effecting operations used by the service layer.
type Runtime struct {
	BrewInstalled        func() bool
	BrewInstall          func(string) error
	BrewFormulaInstalled func(string) bool
	BrewServiceStart     func(string) error
	BrewServiceRestart   func(string) error
	BrewServiceRunning   func(string) bool
	GenerateDNSConfig    func([]string, int) error
	LinkDNSConfig        func() error
	SudoWriteResolver    func(string, int) error
	SudoFlushDNS         func()
	SudoTrustCA          func(string) error
	CaddyCAPath          func() string
	InstallMkcertCA      func() error
	ResolverPath         func(string) string
	FileExists           func(string) bool
	ReadFile             func(string) ([]byte, error)
	Sleep                func(time.Duration)
	CaddyEnsureConfig    func() error
	CaddyInit            func(string) error
	CaddyIsRunning       func() bool
	CaddyAddRoute        func(state.Route) error
	CaddyRemoveRoute     func(string) error
	CaddyAddDomain       func(string, string) error
	CaddyRemoveDomain    func(string) error
	NginxInit            func() error
	NginxIsRunning       func() bool
	NginxAddRoute        func(state.Route) error
	NginxRemoveRoute     func(string) error
	ApacheIsRunning      func() bool
	ApacheAddRoute       func(state.Route) error
	ApacheRemoveRoute    func(string) error
}

// DefaultRuntime returns the real runtime used by the CLI.
func DefaultRuntime() Runtime {
	caddyClient := caddy.New()
	nginxClient := nginx.New()
	apacheClient := apache.New()

	return Runtime{
		BrewInstalled:        brew.IsInstalled,
		BrewInstall:          brew.Install,
		BrewFormulaInstalled: brew.IsFormulaInstalled,
		BrewServiceStart:     brew.ServiceStart,
		BrewServiceRestart:   brew.ServiceRestart,
		BrewServiceRunning:   brew.ServiceRunning,
		GenerateDNSConfig:    dns.GenerateConfig,
		LinkDNSConfig:        dns.LinkConfig,
		SudoWriteResolver:    helper.SudoWriteResolver,
		SudoFlushDNS:         helper.SudoFlushDNS,
		SudoTrustCA:          helper.SudoTrustCA,
		CaddyCAPath:          cert.CaddyCAPath,
		InstallMkcertCA:      cert.InstallMkcertCA,
		ResolverPath:         dns.ResolverPath,
		FileExists: func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		},
		ReadFile:          os.ReadFile,
		Sleep:             time.Sleep,
		CaddyEnsureConfig: ensureCaddyfile,
		CaddyInit:         caddyClient.Init,
		CaddyIsRunning:    caddyClient.IsRunning,
		CaddyAddRoute:     caddyClient.AddRoute,
		CaddyRemoveRoute:  caddyClient.RemoveRoute,
		CaddyAddDomain:    caddyClient.AddDomain,
		CaddyRemoveDomain: caddyClient.RemoveDomain,
		NginxInit:         nginxClient.Init,
		NginxIsRunning:    nginxClient.IsRunning,
		NginxAddRoute:     nginxClient.AddRoute,
		NginxRemoveRoute:  nginxClient.RemoveRoute,
		ApacheIsRunning:   apacheClient.IsRunning,
		ApacheAddRoute:    apacheClient.AddRoute,
		ApacheRemoveRoute: apacheClient.RemoveRoute,
	}
}

// Manager coordinates ODINS operations using a pluggable runtime.
type Manager struct {
	rt Runtime
}

// New creates a service manager.
func New(rt Runtime) *Manager {
	return &Manager{rt: rt}
}

func ensureCaddyfile() error {
	candidates := []string{
		"/opt/homebrew/etc/Caddyfile",
		"/usr/local/etc/Caddyfile",
	}

	for _, path := range candidates {
		if _, err := os.Stat(filepath.Dir(path)); err == nil {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				content := "{\n\tadmin localhost:2019\n}\n"
				return os.WriteFile(path, []byte(content), 0o644)
			}
			return nil
		}
	}

	return nil
}
