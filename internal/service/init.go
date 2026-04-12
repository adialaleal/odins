package service

import (
	"time"

	"github.com/adialaleal/odins/internal/config"
)

// InitOptions configures odins init.
type InitOptions struct {
	NonInteractive bool
	TLD            string
	Backend        string
}

// InitStep captures the outcome of each setup step.
type InitStep struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Detail  string `json:"detail,omitempty"`
	Warning string `json:"warning,omitempty"`
}

// InitResult is returned by odins init.
type InitResult struct {
	TLD     string              `json:"tld"`
	Backend string              `json:"backend"`
	Config  config.GlobalConfig `json:"config"`
	Steps   []InitStep          `json:"steps"`
}

// ResolveInitOptions resolves validated defaults for init.
func (m *Manager) ResolveInitOptions(opts InitOptions) (InitOptions, error) {
	tld, err := validateTLD(opts.TLD)
	if err != nil {
		return InitOptions{}, err
	}

	backend, err := validateBackend(opts.Backend)
	if err != nil {
		return InitOptions{}, err
	}

	return InitOptions{
		NonInteractive: opts.NonInteractive,
		TLD:            tld,
		Backend:        string(backend),
	}, nil
}

// Init configures the local ODINS environment.
func (m *Manager) Init(opts InitOptions) (InitResult, []string, error) {
	resolved, err := m.ResolveInitOptions(opts)
	if err != nil {
		return InitResult{}, nil, err
	}

	if !currentPlatformSupported() {
		return InitResult{}, nil, environmentNotReady("ODINS atualmente suporta apenas macOS")
	}

	result := InitResult{
		TLD:     resolved.TLD,
		Backend: resolved.Backend,
	}
	var warnings []string

	if !m.rt.BrewInstalled() {
		return InitResult{}, nil, environmentNotReady("Homebrew não encontrado. Instale em https://brew.sh")
	}
	result.Steps = append(result.Steps, InitStep{Name: "homebrew", OK: true, Detail: "Homebrew detectado"})

	if err := m.rt.BrewInstall("dnsmasq"); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível instalar dnsmasq")
	}
	if err := m.rt.GenerateDNSConfig([]string{resolved.TLD}, 5300); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível gerar a configuração do dnsmasq")
	}
	if err := m.rt.LinkDNSConfig(); err != nil {
		warnings = append(warnings, "Não foi possível linkar a configuração do dnsmasq: "+err.Error())
		result.Steps = append(result.Steps, InitStep{Name: "dnsmasq-link", OK: false, Warning: err.Error()})
	} else {
		result.Steps = append(result.Steps, InitStep{Name: "dnsmasq-link", OK: true, Detail: "Configuração do dnsmasq vinculada"})
	}
	if err := m.rt.BrewServiceRestart("dnsmasq"); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível reiniciar o dnsmasq")
	}
	result.Steps = append(result.Steps, InitStep{Name: "dnsmasq", OK: true, Detail: "dnsmasq instalado e reiniciado"})

	proxyFormula := "caddy"
	switch resolved.Backend {
	case string(config.BackendNginx):
		proxyFormula = "nginx"
	case string(config.BackendApache):
		proxyFormula = "httpd"
	}
	if err := m.rt.BrewInstall(proxyFormula); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível instalar %s", proxyFormula)
	}
	if config.ProxyBackend(resolved.Backend) == config.BackendCaddy && m.rt.CaddyEnsureConfig != nil {
		if err := m.rt.CaddyEnsureConfig(); err != nil {
			warnings = append(warnings, "Não foi possível preparar o Caddyfile base: "+err.Error())
			result.Steps = append(result.Steps, InitStep{Name: "caddyfile", OK: false, Warning: err.Error()})
		} else {
			result.Steps = append(result.Steps, InitStep{Name: "caddyfile", OK: true, Detail: "Caddyfile base preparado"})
		}
	}
	if err := m.rt.BrewServiceRestart(proxyFormula); err != nil {
		warnings = append(warnings, "Não foi possível iniciar o serviço "+proxyFormula+": "+err.Error())
		result.Steps = append(result.Steps, InitStep{Name: "proxy-service", OK: false, Warning: err.Error()})
	} else {
		result.Steps = append(result.Steps, InitStep{Name: "proxy-service", OK: true, Detail: "Serviço " + proxyFormula + " reiniciado"})
	}
	if config.ProxyBackend(resolved.Backend) == config.BackendCaddy {
		m.waitForCaddyAPI()
	}

	if err := m.rt.SudoWriteResolver(resolved.TLD, 5300); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível escrever /etc/resolver/%s", resolved.TLD)
	}
	if m.rt.SudoFlushDNS != nil {
		m.rt.SudoFlushDNS()
	}
	result.Steps = append(result.Steps, InitStep{Name: "resolver", OK: true, Detail: "/etc/resolver/" + resolved.TLD + " configurado"})

	switch config.ProxyBackend(resolved.Backend) {
	case config.BackendCaddy:
		if err := m.rt.CaddyInit(resolved.TLD); err != nil {
			warnings = append(warnings, "Não foi possível inicializar a configuração base do Caddy: "+err.Error())
			result.Steps = append(result.Steps, InitStep{Name: "caddy-config", OK: false, Warning: err.Error()})
		} else {
			result.Steps = append(result.Steps, InitStep{Name: "caddy-config", OK: true, Detail: "Configuração base do Caddy carregada"})
		}
		m.waitForCaddyCA()
		caPath := m.rt.CaddyCAPath()
		if caPath == "" {
			warnings = append(warnings, "A CA local do Caddy ainda não foi gerada. Abra um domínio local ou rode odins doctor depois do primeiro acesso HTTPS.")
			result.Steps = append(result.Steps, InitStep{Name: "certificates", OK: false, Warning: "CA do Caddy ainda não gerada"})
		} else {
			if err := m.rt.SudoTrustCA(caPath); err != nil {
				warnings = append(warnings, "Não foi possível confiar na CA local do Caddy: "+err.Error())
				result.Steps = append(result.Steps, InitStep{Name: "certificates", OK: false, Warning: err.Error()})
			} else {
				result.Steps = append(result.Steps, InitStep{Name: "certificates", OK: true, Detail: "CA local do Caddy adicionada ao keychain do macOS"})
			}
		}
	default:
		if !m.rt.BrewFormulaInstalled("mkcert") {
			if err := m.rt.BrewInstall("mkcert"); err != nil {
				warnings = append(warnings, "Não foi possível instalar mkcert automaticamente: "+err.Error())
			}
		}
		if err := m.rt.InstallMkcertCA(); err != nil {
			warnings = append(warnings, "Não foi possível instalar a CA do mkcert: "+err.Error())
			result.Steps = append(result.Steps, InitStep{Name: "certificates", OK: false, Warning: err.Error()})
		} else {
			result.Steps = append(result.Steps, InitStep{Name: "certificates", OK: true, Detail: "CA do mkcert instalada"})
		}
		if config.ProxyBackend(resolved.Backend) == config.BackendNginx {
			if err := m.rt.NginxInit(); err != nil {
				warnings = append(warnings, "Não foi possível preparar a configuração do Nginx: "+err.Error())
				result.Steps = append(result.Steps, InitStep{Name: "nginx-config", OK: false, Warning: err.Error()})
			} else {
				result.Steps = append(result.Steps, InitStep{Name: "nginx-config", OK: true, Detail: "Diretório de configuração do Nginx preparado"})
			}
		}
	}

	result.Config = config.GlobalConfig{
		TLD:            resolved.TLD,
		ProxyBackend:   config.ProxyBackend(resolved.Backend),
		DnsmasqPort:    5300,
		CaddyAdmin:     "http://localhost:2019",
		HTTPPort:       80,
		HTTPSPort:      443,
		OnboardingDone: true,
	}
	if err := config.SaveGlobal(result.Config); err != nil {
		return InitResult{}, nil, runtimeFailure(err, "não foi possível salvar a configuração global")
	}
	result.Steps = append(result.Steps, InitStep{Name: "global-config", OK: true, Detail: "Configuração global salva"})

	return result, warnings, nil
}

func (m *Manager) waitForCaddyCA() {
	for i := 0; i < 20; i++ {
		if m.rt.CaddyCAPath() != "" {
			return
		}
		m.rt.Sleep(500 * time.Millisecond)
	}
}

func (m *Manager) waitForCaddyAPI() {
	for i := 0; i < 20; i++ {
		if m.rt.CaddyIsRunning != nil && m.rt.CaddyIsRunning() {
			return
		}
		m.rt.Sleep(500 * time.Millisecond)
	}
}
