package service

import (
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/docker"
	"github.com/adialaleal/odins/internal/page"
	"github.com/adialaleal/odins/internal/state"
)

// RouteStatus contains the route plus its live availability.
type RouteStatus struct {
	Route   state.Route `json:"route"`
	Up      bool        `json:"up"`
	Proto   string      `json:"proto"`
	Runtime string      `json:"runtime"`
}

// ListRoutesResult is returned by ls.
type ListRoutesResult struct {
	Routes []RouteStatus `json:"routes"`
	Count  int           `json:"count"`
}

// AddRouteOptions configures a route addition.
type AddRouteOptions struct {
	Subdomain string
	Port      int
	Docker    string
	Project   string
	Domain    string
	HTTPS     bool
}

// AddRouteResult is returned by add.
type AddRouteResult struct {
	Route         state.Route `json:"route"`
	DomainPageURL string      `json:"domain_page_url,omitempty"`
}

// AppliedRoute captures a route applied by odins up.
type AppliedRoute struct {
	Route state.Route `json:"route"`
	Proto string      `json:"proto"`
}

// UpResult is returned by odins up.
type UpResult struct {
	Directory         string             `json:"directory"`
	ProjectConfigPath string             `json:"project_config_path"`
	GeneratedConfig   bool               `json:"generated_config"`
	AutoDetected      bool               `json:"auto_detected"`
	Project           config.ProjectInfo `json:"project"`
	Routes            []AppliedRoute     `json:"routes"`
	DomainPageURL     string             `json:"domain_page_url,omitempty"`
}

// RemovedRoute captures a route removed by down.
type RemovedRoute struct {
	Subdomain string `json:"subdomain"`
}

// DownResult is returned by odins down.
type DownResult struct {
	Directory     string         `json:"directory"`
	Project       string         `json:"project"`
	RemovedRoutes []RemovedRoute `json:"removed_routes"`
	DomainPageURL string         `json:"domain_page_url,omitempty"`
}

// KillResult is returned by odins kill.
type KillResult struct {
	Subdomain     string `json:"subdomain"`
	DomainPageURL string `json:"domain_page_url,omitempty"`
}

// DomainSummary is the shape exposed by domain ls.
type DomainSummary struct {
	Domain   state.Domain `json:"domain"`
	Hostname string       `json:"hostname"`
	Services int          `json:"services"`
}

// DomainListResult is returned by domain ls.
type DomainListResult struct {
	Domains []DomainSummary `json:"domains"`
	Count   int             `json:"count"`
}

// DomainAddResult is returned by domain add.
type DomainAddResult struct {
	Domain   state.Domain `json:"domain"`
	Hostname string       `json:"hostname"`
	PageDir  string       `json:"page_dir"`
}

// DomainRemoveResult is returned by domain rm.
type DomainRemoveResult struct {
	Domain        string         `json:"domain"`
	Hostname      string         `json:"hostname"`
	RemovedRoutes []RemovedRoute `json:"removed_routes"`
}

// ListRoutes lists active routes with live status checks.
func (m *Manager) ListRoutes() (ListRoutesResult, []string, error) {
	store, err := state.Load()
	if err != nil {
		return ListRoutesResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}

	result := ListRoutesResult{}
	for _, route := range store.Routes {
		runtimeLabel := route.Runtime
		if route.DockerContainer != "" {
			runtimeLabel = "docker"
		}
		if runtimeLabel == "" {
			runtimeLabel = "unknown"
		}
		result.Routes = append(result.Routes, RouteStatus{
			Route:   route,
			Up:      docker.CheckSubdomain(route.Port),
			Proto:   stringsUpper(defaultRouteProtocol(route)),
			Runtime: runtimeLabel,
		})
	}
	result.Count = len(result.Routes)

	return result, nil, nil
}

// AddRoute creates a single route.
func (m *Manager) AddRoute(opts AddRouteOptions) (AddRouteResult, []string, error) {
	if opts.Subdomain == "" {
		return AddRouteResult{}, nil, invalidInput("subdomínio é obrigatório")
	}
	if opts.Port <= 0 || opts.Port > 65535 {
		return AddRouteResult{}, nil, invalidInput("porta inválida: %d", opts.Port)
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return AddRouteResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	project := opts.Project
	if project == "" {
		parts := splitDomainParts(opts.Subdomain)
		if len(parts) >= 2 {
			project = parts[1]
		} else if len(parts) == 1 {
			project = parts[0]
		}
	}

	route := state.Route{
		ID:              "odins-" + opts.Subdomain,
		Subdomain:       opts.Subdomain,
		Port:            opts.Port,
		Project:         project,
		Domain:          opts.Domain,
		DockerContainer: opts.Docker,
		HTTPS:           opts.HTTPS,
		CreatedAt:       time.Now(),
	}

	if err := m.addProxyRoute(cfg, route); err != nil {
		return AddRouteResult{}, nil, runtimeFailure(err, "não foi possível adicionar a rota %s", opts.Subdomain)
	}

	store, err := state.Load()
	if err != nil {
		return AddRouteResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}

	store.Add(route)
	if err := store.Save(); err != nil {
		return AddRouteResult{}, nil, runtimeFailure(err, "não foi possível persistir a rota %s", opts.Subdomain)
	}

	var warnings []string
	result := AddRouteResult{Route: route}
	if opts.Domain != "" {
		if err := regeneratePageForDomain(cfg, store, opts.Domain); err != nil {
			warnings = append(warnings, err.Error())
		} else {
			result.DomainPageURL = "https://" + opts.Domain + "." + cfg.TLD
		}
	}

	return result, warnings, nil
}

// Up applies routes from .odins or from auto-detection.
func (m *Manager) Up(dir string) (UpResult, []string, error) {
	normalizedDir, err := normalizeDir(dir)
	if err != nil {
		return UpResult{}, nil, runtimeFailure(err, "não foi possível resolver o diretório do projeto")
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return UpResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	projectCfgPath := projectConfigPath(normalizedDir)
	var projectCfg config.ProjectConfig
	result := UpResult{
		Directory:         normalizedDir,
		ProjectConfigPath: projectCfgPath,
	}
	var warnings []string

	if config.ExistsProject(normalizedDir) {
		projectCfg, err = config.LoadProject(projectCfgPath)
		if err != nil {
			return UpResult{}, nil, configurationError(err, "não foi possível ler o arquivo .odins do projeto")
		}
	} else {
		detected := detect.Project(normalizedDir)
		if detected.Runtime == "unknown" {
			return UpResult{}, nil, configurationError(nil, "não foi possível detectar o tipo de projeto em %s", normalizedDir)
		}

		projectCfg = recommendedConfig(detected.Name, detected.Runtime, detected.Framework, detected.Port)
		result.AutoDetected = true
		if err := config.SaveProject(projectCfgPath, projectCfg); err != nil {
			warnings = append(warnings, "ODINS detectou o projeto, mas não conseguiu salvar o arquivo .odins automaticamente.")
		} else {
			result.GeneratedConfig = true
		}
	}

	store, err := state.Load()
	if err != nil {
		return UpResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}

	effectiveDomain := projectCfg.Project.Domain
	if effectiveDomain == "" {
		effectiveDomain = projectCfg.Project.Name
	}

	for _, routeCfg := range projectCfg.Routes {
		fqdn := buildFQDN(routeCfg.Subdomain, effectiveDomain, projectCfg.Project.Name, cfg.TLD)
		route := state.Route{
			ID:              "odins-" + fqdn,
			Subdomain:       fqdn,
			Port:            routeCfg.Port,
			Project:         projectCfg.Project.Name,
			Runtime:         projectCfg.Project.Runtime,
			Domain:          effectiveDomain,
			DockerContainer: routeCfg.DockerContainer,
			HTTPS:           routeCfg.HTTPS,
			CreatedAt:       time.Now(),
		}
		if err := m.addProxyRoute(cfg, route); err != nil {
			warnings = append(warnings, "Falha ao aplicar "+fqdn+": "+err.Error())
			continue
		}
		store.Add(route)
		result.Routes = append(result.Routes, AppliedRoute{
			Route: route,
			Proto: defaultRouteProtocol(route),
		})
	}

	if err := store.Save(); err != nil {
		return UpResult{}, nil, runtimeFailure(err, "não foi possível persistir as rotas do projeto")
	}

	if effectiveDomain != "" {
		if err := regeneratePageForDomain(cfg, store, effectiveDomain); err != nil {
			warnings = append(warnings, err.Error())
		} else {
			result.DomainPageURL = "https://" + effectiveDomain + "." + cfg.TLD
		}
	}

	projectCfg.Project.Domain = effectiveDomain
	result.Project = projectCfg.Project
	return result, warnings, nil
}

// Down removes all routes declared in the current project .odins file.
func (m *Manager) Down(dir string) (DownResult, []string, error) {
	normalizedDir, err := normalizeDir(dir)
	if err != nil {
		return DownResult{}, nil, runtimeFailure(err, "não foi possível resolver o diretório do projeto")
	}
	if !config.ExistsProject(normalizedDir) {
		return DownResult{}, nil, configurationError(nil, ".odins não encontrado em %s", normalizedDir)
	}

	projectCfg, err := config.LoadProject(projectConfigPath(normalizedDir))
	if err != nil {
		return DownResult{}, nil, configurationError(err, "não foi possível ler o arquivo .odins do projeto")
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return DownResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	store, err := state.Load()
	if err != nil {
		return DownResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}

	var warnings []string
	result := DownResult{
		Directory: normalizedDir,
		Project:   projectCfg.Project.Name,
	}
	effectiveDomain := projectCfg.Project.Domain
	if effectiveDomain == "" {
		effectiveDomain = projectCfg.Project.Name
	}
	for _, routeCfg := range projectCfg.Routes {
		fqdn := buildFQDN(routeCfg.Subdomain, effectiveDomain, projectCfg.Project.Name, cfg.TLD)
		if err := m.removeProxyRoute(cfg, fqdn); err != nil {
			warnings = append(warnings, "Falha ao remover "+fqdn+": "+err.Error())
		}
		store.Remove(fqdn)
		result.RemovedRoutes = append(result.RemovedRoutes, RemovedRoute{Subdomain: fqdn})
	}

	if err := store.Save(); err != nil {
		return DownResult{}, nil, runtimeFailure(err, "não foi possível persistir a remoção das rotas")
	}

	if effectiveDomain != "" {
		if err := regeneratePageForDomain(cfg, store, effectiveDomain); err != nil {
			warnings = append(warnings, err.Error())
		} else {
			result.DomainPageURL = "https://" + effectiveDomain + "." + cfg.TLD
		}
	}

	return result, warnings, nil
}

// Kill removes a single route by FQDN.
func (m *Manager) Kill(subdomain string) (KillResult, []string, error) {
	if subdomain == "" {
		return KillResult{}, nil, invalidInput("subdomínio é obrigatório")
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return KillResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	store, err := state.Load()
	if err != nil {
		return KillResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}

	route, ok := store.Get(subdomain)
	if !ok {
		return KillResult{}, nil, configurationError(nil, "rota '%s' não encontrada", subdomain)
	}

	var warnings []string
	if err := m.removeProxyRoute(cfg, subdomain); err != nil {
		warnings = append(warnings, "Falha ao remover "+subdomain+" do proxy: "+err.Error())
	}

	store.Remove(subdomain)
	if err := store.Save(); err != nil {
		return KillResult{}, nil, runtimeFailure(err, "não foi possível persistir a remoção da rota")
	}

	result := KillResult{Subdomain: subdomain}
	if route.Domain != "" {
		if err := regeneratePageForDomain(cfg, store, route.Domain); err != nil {
			warnings = append(warnings, err.Error())
		} else {
			result.DomainPageURL = "https://" + route.Domain + "." + cfg.TLD
		}
	}

	return result, warnings, nil
}

// DomainAdd creates a new domain workspace and its landing page.
func (m *Manager) DomainAdd(name, title, description string) (DomainAddResult, []string, error) {
	if name == "" {
		return DomainAddResult{}, nil, invalidInput("nome do domínio é obrigatório")
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return DomainAddResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	store, err := state.Load()
	if err != nil {
		return DomainAddResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de rotas")
	}
	if _, exists := store.GetDomain(name); exists {
		return DomainAddResult{}, nil, configurationError(nil, "domínio '%s' já existe", name)
	}

	if title == "" {
		title = name
	}

	domain := state.Domain{
		Name:        name,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}
	store.AddDomain(domain)
	if err := store.Save(); err != nil {
		return DomainAddResult{}, nil, runtimeFailure(err, "não foi possível persistir o domínio '%s'", name)
	}

	hostname := name + "." + cfg.TLD
	pageDir := page.PageDir(name)
	var warnings []string
	if err := page.Generate(page.PageData{
		Domain:      name,
		TLD:         cfg.TLD,
		Title:       title,
		Description: description,
	}); err != nil {
		warnings = append(warnings, "Falha ao gerar a landing page do domínio: "+err.Error())
	}
	if cfg.ProxyBackend == config.BackendCaddy {
		if err := m.rt.CaddyAddDomain(hostname, pageDir); err != nil {
			warnings = append(warnings, "Falha ao registrar a landing page no Caddy: "+err.Error())
		}
	} else {
		warnings = append(warnings, "Landing pages de domínio são registradas automaticamente apenas com o backend Caddy nesta versão.")
	}

	return DomainAddResult{
		Domain:   domain,
		Hostname: hostname,
		PageDir:  pageDir,
	}, warnings, nil
}

// DomainList lists all configured domains.
func (m *Manager) DomainList() (DomainListResult, []string, error) {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return DomainListResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	store, err := state.Load()
	if err != nil {
		return DomainListResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de domínios")
	}

	result := DomainListResult{}
	for _, domain := range store.Domains {
		result.Domains = append(result.Domains, DomainSummary{
			Domain:   domain,
			Hostname: domain.Name + "." + cfg.TLD,
			Services: len(store.ByDomain(domain.Name)),
		})
	}
	result.Count = len(result.Domains)
	return result, nil, nil
}

// DomainRemove removes a domain workspace and all attached routes.
func (m *Manager) DomainRemove(name string) (DomainRemoveResult, []string, error) {
	if name == "" {
		return DomainRemoveResult{}, nil, invalidInput("nome do domínio é obrigatório")
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return DomainRemoveResult{}, nil, configurationError(err, "não foi possível carregar a configuração global")
	}

	store, err := state.Load()
	if err != nil {
		return DomainRemoveResult{}, nil, runtimeFailure(err, "não foi possível carregar o registro de domínios")
	}
	if _, exists := store.GetDomain(name); !exists {
		return DomainRemoveResult{}, nil, configurationError(nil, "domínio '%s' não encontrado", name)
	}

	var warnings []string
	routes := store.ByDomain(name)
	for _, route := range routes {
		if err := m.removeProxyRoute(cfg, route.Subdomain); err != nil {
			warnings = append(warnings, "Falha ao remover "+route.Subdomain+": "+err.Error())
		}
	}

	hostname := name + "." + cfg.TLD
	if cfg.ProxyBackend == config.BackendCaddy {
		if err := m.rt.CaddyRemoveDomain(hostname); err != nil {
			warnings = append(warnings, "Falha ao remover a landing page do domínio do Caddy: "+err.Error())
		}
	}

	store.RemoveDomain(name)
	if err := store.Save(); err != nil {
		return DomainRemoveResult{}, nil, runtimeFailure(err, "não foi possível persistir a remoção do domínio '%s'", name)
	}

	result := DomainRemoveResult{
		Domain:   name,
		Hostname: hostname,
	}
	for _, route := range routes {
		result.RemovedRoutes = append(result.RemovedRoutes, RemovedRoute{Subdomain: route.Subdomain})
	}
	return result, warnings, nil
}

func (m *Manager) addProxyRoute(cfg config.GlobalConfig, route state.Route) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return m.rt.NginxAddRoute(route)
	case config.BackendApache:
		return m.rt.ApacheAddRoute(route)
	default:
		return m.rt.CaddyAddRoute(route)
	}
}

func (m *Manager) removeProxyRoute(cfg config.GlobalConfig, subdomain string) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return m.rt.NginxRemoveRoute(subdomain)
	case config.BackendApache:
		return m.rt.ApacheRemoveRoute(subdomain)
	default:
		return m.rt.CaddyRemoveRoute(subdomain)
	}
}

func regeneratePageForDomain(cfg config.GlobalConfig, store *state.Store, domainName string) error {
	domain, ok := store.GetDomain(domainName)
	if !ok {
		return nil
	}

	var routes []page.RouteInfo
	for _, route := range store.ByDomain(domainName) {
		routes = append(routes, page.RouteInfo{
			Subdomain: extractSubdomain(route.Subdomain, domainName, cfg.TLD),
			FQDN:      route.Subdomain,
			Port:      route.Port,
			Runtime:   route.Runtime,
			Project:   route.Project,
		})
	}

	if err := page.Generate(page.PageData{
		Domain:      domainName,
		TLD:         cfg.TLD,
		Title:       domain.Title,
		Description: domain.Description,
		Routes:      routes,
	}); err != nil {
		return runtimeFailure(err, "não foi possível regenerar a landing page do domínio '%s'", domainName)
	}

	return nil
}

func extractSubdomain(fqdn, domain, tld string) string {
	suffix := "." + domain + "." + tld
	if len(fqdn) > len(suffix) && fqdn[len(fqdn)-len(suffix):] == suffix {
		return fqdn[:len(fqdn)-len(suffix)]
	}
	return fqdn
}

func splitDomainParts(value string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(value); i++ {
		if value[i] == '.' {
			if i > start {
				parts = append(parts, value[start:i])
			}
			start = i + 1
		}
	}
	if start < len(value) {
		parts = append(parts, value[start:])
	}
	return parts
}

func stringsUpper(value string) string {
	if value == "" {
		return value
	}
	if len(value) == 1 {
		if value[0] >= 'a' && value[0] <= 'z' {
			return string(value[0] - 32)
		}
		return value
	}
	var out []byte
	for i := 0; i < len(value); i++ {
		b := value[i]
		if b >= 'a' && b <= 'z' {
			b -= 32
		}
		out = append(out, b)
	}
	return string(out)
}
