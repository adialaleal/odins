package service

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adialaleal/odins/internal/cert"
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/dns"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adrg/xdg"
)

// DoctorCheck is a single environment diagnostic.
type DoctorCheck struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Status  string `json:"status"`
	Details string `json:"details"`
	Action  string `json:"action,omitempty"`
}

// DoctorResult is returned by odins doctor.
type DoctorResult struct {
	Healthy bool                `json:"healthy"`
	Config  config.GlobalConfig `json:"config"`
	Checks  []DoctorCheck       `json:"checks"`
}

// Doctor inspects the local ODINS environment.
func (m *Manager) Doctor() (DoctorResult, []string, error) {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return DoctorResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	result := DoctorResult{Config: cfg}
	addCheck := func(check DoctorCheck) {
		result.Checks = append(result.Checks, check)
		if !check.OK {
			result.Healthy = false
		}
	}

	result.Healthy = true
	addCheck(DoctorCheck{
		Name:    "platform",
		OK:      runtime.GOOS == "darwin",
		Status:  mapBool(runtime.GOOS == "darwin", "supported", "unsupported"),
		Details: "Sistema operacional atual: " + runtime.GOOS,
		Action:  "Use o ODINS em uma máquina macOS.",
	})

	homebrewOK := m.rt.BrewInstalled()
	addCheck(DoctorCheck{
		Name:    "homebrew",
		OK:      homebrewOK,
		Status:  mapBool(homebrewOK, "installed", "missing"),
		Details: "Homebrew é usado para instalar dnsmasq e o proxy local.",
		Action:  "Instale o Homebrew em https://brew.sh e rode `odins init` novamente.",
	})

	globalConfigExists := m.rt.FileExists(config.ConfigPath())
	addCheck(DoctorCheck{
		Name:    "global_config",
		OK:      globalConfigExists,
		Status:  mapBool(globalConfigExists, "present", "missing"),
		Details: "Arquivo esperado em " + config.ConfigPath(),
		Action:  "Rode `odins init` para gerar a configuração global.",
	})

	dnsmasqInstalled := m.rt.BrewFormulaInstalled("dnsmasq")
	addCheck(DoctorCheck{
		Name:    "dnsmasq_formula",
		OK:      dnsmasqInstalled,
		Status:  mapBool(dnsmasqInstalled, "installed", "missing"),
		Details: "dnsmasq fornece a resolução wildcard local.",
		Action:  "Rode `odins init` para instalar o dnsmasq.",
	})

	dnsmasqRunning := m.rt.BrewServiceRunning("dnsmasq")
	addCheck(DoctorCheck{
		Name:    "dnsmasq_service",
		OK:      dnsmasqRunning,
		Status:  mapBool(dnsmasqRunning, "running", "stopped"),
		Details: "Serviço Homebrew `dnsmasq`.",
		Action:  "Rode `brew services restart dnsmasq` ou `odins init`.",
	})

	dnsConfigExists := m.rt.FileExists(dns.DnsmasqConfPath())
	addCheck(DoctorCheck{
		Name:    "dnsmasq_config",
		OK:      dnsConfigExists,
		Status:  mapBool(dnsConfigExists, "present", "missing"),
		Details: "Arquivo esperado em " + dns.DnsmasqConfPath(),
		Action:  "Rode `odins init` para gerar a configuração do dnsmasq.",
	})

	resolverPath := m.rt.ResolverPath(cfg.TLD)
	resolverExists := m.rt.FileExists(resolverPath)
	addCheck(DoctorCheck{
		Name:    "resolver",
		OK:      resolverExists,
		Status:  mapBool(resolverExists, "present", "missing"),
		Details: "Arquivo esperado em " + resolverPath,
		Action:  "Rode `odins init` para criar o resolver local.",
	})

	proxyFormula := string(cfg.ProxyBackend)
	proxyInstalled := m.rt.BrewFormulaInstalled(proxyFormulaToFormula(proxyFormula))
	addCheck(DoctorCheck{
		Name:    "proxy_formula",
		OK:      proxyInstalled,
		Status:  mapBool(proxyInstalled, "installed", "missing"),
		Details: "Backend configurado: " + proxyFormula,
		Action:  "Rode `odins init --backend " + proxyFormula + "` para instalar o proxy configurado.",
	})

	proxyRunning := m.proxyRunning(cfg.ProxyBackend)
	addCheck(DoctorCheck{
		Name:    "proxy_service",
		OK:      proxyRunning,
		Status:  mapBool(proxyRunning, "running", "stopped"),
		Details: "Serviço do backend " + proxyFormula,
		Action:  "Inicie ou reinicie o proxy configurado e rode `odins doctor` novamente.",
	})

	storePath := filepath.Join(xdg.DataHome, "odins", "routes.json")
	store, loadErr := state.Load()

	certOK, certDetails, certAction := m.certificateCheck(cfg, store)
	addCheck(DoctorCheck{
		Name:    "certificates",
		OK:      certOK,
		Status:  mapBool(certOK, "ready", "attention"),
		Details: certDetails,
		Action:  certAction,
	})
	storeOK := loadErr == nil
	storeDetails := "Store carregado com sucesso."
	if !m.rt.FileExists(storePath) {
		storeDetails = "Store ainda não existe; ele será criado na primeira rota salva."
	}
	if loadErr != nil {
		storeDetails = loadErr.Error()
	}
	addCheck(DoctorCheck{
		Name:    "store",
		OK:      storeOK,
		Status:  mapBool(storeOK, "ready", "error"),
		Details: storeDetails,
		Action:  "Apague ou corrija o arquivo de store se ele estiver corrompido.",
	})

	var warnings []string
	if storeOK && len(store.Routes) == 0 {
		warnings = append(warnings, "Nenhuma rota ativa encontrada; isso é normal se você ainda não rodou `odins up` ou `odins add`.")
	}

	return result, warnings, nil
}

func (m *Manager) proxyRunning(backend config.ProxyBackend) bool {
	switch backend {
	case config.BackendNginx:
		return m.rt.NginxIsRunning()
	case config.BackendApache:
		return m.rt.ApacheIsRunning()
	default:
		return m.rt.CaddyIsRunning()
	}
}

func (m *Manager) certificateCheck(cfg config.GlobalConfig, store *state.Store) (bool, string, string) {
	switch cfg.ProxyBackend {
	case config.BackendCaddy:
		path := m.rt.CaddyCAPath()
		if path == "" {
			return false,
				"A CA local do Caddy ainda não foi encontrada no disco.",
				"Abra um domínio local HTTPS ou rode `odins init` novamente para forçar a geração da CA."
		}
		return true,
			"CA local do Caddy encontrada em " + path,
			""
	default:
		certDir := cert.CertDir()
		if !m.rt.FileExists(certDir) {
			return false,
				"O diretório de certificados do mkcert ainda não existe em " + certDir + ".",
				"Rode `odins init` ou gere um certificado mkcert para o domínio desejado."
		}
		if store == nil || len(store.Routes) == 0 {
			return true,
				"Diretório de certificados encontrado em " + certDir + ". Nenhuma rota ativa para verificar.",
				""
		}
		var missing []string
		for _, r := range store.Routes {
			certFile := filepath.Join(certDir, r.Subdomain+".pem")
			keyFile := filepath.Join(certDir, r.Subdomain+"-key.pem")
			if !m.rt.FileExists(certFile) || !m.rt.FileExists(keyFile) {
				missing = append(missing, r.Subdomain)
			}
		}
		if len(missing) > 0 {
			return false,
				"Certificados ausentes para: " + strings.Join(missing, ", "),
				"Rode `odins up` para gerar os certificados faltantes ou `mkcert` manualmente."
		}
		return true,
			"Todos os certificados das rotas ativas encontrados em " + certDir,
			""
	}
}

func mapBool(value bool, whenTrue, whenFalse string) string {
	if value {
		return whenTrue
	}
	return whenFalse
}

func proxyFormulaToFormula(backend string) string {
	switch strings.ToLower(backend) {
	case "apache":
		return "httpd"
	default:
		return backend
	}
}
