package i18n

var ptStrings = map[string]string{
	// ── Welcome / onboarding ─────────────────────────────────────────────
	"welcome.section.what":     "O que é o ODINS?",
	"welcome.section.how":      "Como funciona?",
	"welcome.section.domains":  "Domínios e Subdomínios",
	"welcome.section.config":   "Configuração por Projeto (.odins)",
	"welcome.section.commands": "Comandos Principais",
	"welcome.section.next":     "Próximos Passos",
	"welcome.tagline":          "The All-Father of Local DNS",
	"welcome.elimina":          "ODINS elimina a guerra de portas no seu desenvolvimento local.",
	"welcome.sem_odins":        "Sem ODINS:",
	"welcome.com_odins":        "Com ODINS:",
	"welcome.https_auto":       "Cada projeto ganha um subdomínio bonito com HTTPS automático.",
	"welcome.how_dns":          "DNS   — dnsmasq resolve *.<projeto>.odins → 127.0.0.1",
	"welcome.how_proxy":        "Proxy — Caddy roteia web.<projeto>.odins → localhost:3000",
	"welcome.how_https":        "HTTPS — Caddy gerencia certificados automaticamente",
	"welcome.domain_is":        "Um domínio é o workspace central dos seus projetos:",
	"welcome.domain_landing":   "cria <projeto>.odins — uma landing page que lista todos os serviços do workspace com status em tempo real.",
	"welcome.subdomain_is":     "Cada subdomínio é um projeto/serviço:",
	"welcome.auto_detect":      "O ODINS detecta Node.js, Go e Python automaticamente.",
	"welcome.not_configured":   "Parece que o ODINS ainda não foi configurado nesta máquina.",
	"welcome.run_init":         "Rodar odins init agora?",
	"welcome.run_init_prompt":  "[S/n]",
	"welcome.run_init_yes":     "s",
	"welcome.ok":               "Tudo bem! Quando estiver pronto, rode: odins init",
	"welcome.already_configured": "✓ ODINS já está configurado.",
	"welcome.create_domain":    "Comece criando um domínio:",
	"welcome.then_project":     "Depois num projeto:",
	"welcome.enter":            "[Enter para continuar...]",
	// New-folder welcome (short)
	"welcome.new_folder.title":     "Começando neste projeto",
	"welcome.new_folder.detected":  "Projeto detectado: %s (%s/%s, porta %d)",
	"welcome.new_folder.activate":  "Rode odins up para ativar as rotas automaticamente.",
	"welcome.new_folder.manual":    "Nenhum projeto detectado. Crie um .odins com:",
	"welcome.new_folder.or_add":    "Ou adicione uma rota manualmente:",
	"welcome.new_folder.see_guide": "Ver guia completo: odins welcome",

	// ── Command descriptions (shown in welcome) ───────────────────────────
	"cmd.init_desc":       "Setup único: DNS, proxy, HTTPS",
	"cmd.domain_add_desc": "Criar workspace <projeto>.odins",
	"cmd.up_desc":         "Ativar rotas do projeto atual",
	"cmd.ls_desc":         "Listar rotas ativas",
	"cmd.kill_desc":       "Remover uma rota",
	"cmd.down_desc":       "Remover todas as rotas do projeto",
	"cmd.tui_desc":        "Abrir painel TUI",
	"cmd.welcome_desc":    "Ver este guia novamente",

	// ── TUI app messages ─────────────────────────────────────────────────
	"tui.detected":        "Projeto '%s' (%s/%s, porta %d) detectado — pressione [u] para ativar as rotas",
	"tui.no_project":      "Nenhum .odins encontrado — pressione [a] para adicionar uma rota manualmente",
	"tui.activating":      "Ativando rotas...",
	"tui.error":           "Erro: %s",
	"tui.no_routes":       "Nenhuma rota aplicada",
	"tui.routes_activated":"✓ %d rota(s) ativada(s)!",
	"tui.config_save_err": "Erro ao salvar config: %s",
	"tui.config_saved":    "Configurações salvas!",
	"tui.save_error":      "Erro ao salvar: %s",
	"tui.proxy_error":     "Erro no proxy: %s",
	"tui.route_added":     "✓ %s → :%d adicionado!",
	"tui.route_removed":   "%s removido.",
	"tui.read_project_err":"ler .odins: %s",
	"tui.no_detect":       "projeto não detectado em %s",

	// ── Footer hints ─────────────────────────────────────────────────────
	"hint.add":        "adicionar",
	"hint.up":         "odins up",
	"hint.remove":     "remover",
	"hint.settings":   "settings",
	"hint.logs":       "logs",
	"hint.quit":       "sair",
	"hint.back":       "voltar",
	"hint.next_field": "próximo campo",
	"hint.confirm":    "confirmar",
	"hint.cancel":     "cancelar",
	"hint.select":     "selecionar",
	"hint.navigate":   "navegar",
	"hint.save":       "salvar",
	"hint.scroll":     "scroll",

	// ── Dashboard ────────────────────────────────────────────────────────
	"dash.title":          "Dashboard — Rotas Ativas",
	"dash.routes":         "%d rotas",
	"dash.remove_confirm": "Remover %s?",

	// ── Add route form ───────────────────────────────────────────────────
	"add.title":              "Adicionar Rota",
	"add.field.subdomain":    "Subdomínio",
	"add.field.port":         "Porta",
	"add.field.project":      "Projeto",
	"add.field.docker":       "Docker",
	"add.detected":           "✦ Detectado: %s/%s (porta %d)",
	"add.err.subdomain":      "subdomínio é obrigatório",
	"add.err.port_required":  "porta é obrigatória",
	"add.err.port_invalid":   "porta inválida",

	// ── Settings screen ──────────────────────────────────────────────────
	"settings.title":      "Configurações",
	"settings.tld_label":  "TLD:",
	"settings.proxy_label":"Proxy:",
	"settings.warn_local": "⚠  .local conflita com mDNS/Bonjour no macOS — use com cuidado",
	"settings.info":       "Domínios: *%s → 127.0.0.1  |  Proxy: %s",

	// ── Logs screen ──────────────────────────────────────────────────────
	"logs.title": "Logs do Proxy",

	// ── Confirm modal ────────────────────────────────────────────────────
	"modal.confirm": "confirmar",
	"modal.cancel":  "cancelar",

	// ── CLI: odins up ────────────────────────────────────────────────────
	"up.reading":      "→ Lendo .odins do projeto '%s'",
	"up.detecting":    "→ .odins não encontrado, detectando projeto...",
	"up.not_detected": "não foi possível detectar o tipo de projeto em %s\nCrie um .odins manualmente ou use: odins add <subdomain> --port <port>",
	"up.detected":     "→ Detectado: %s/%s (porta %d)",
	"up.start_cmd":    "→ Comando de start: %s",
	"up.save_warn":    "⚠  Não foi possível salvar .odins: %v",
	"up.created":      "→ .odins criado em %s",
	"up.domain":       "→ Domínio: %s.%s",
	"up.route_error":  "✗ %s: %v",
	"up.route_ok":     "✓ %s://%s → :%d",
	"up.page_updated": "→ Landing page atualizada: https://%s.%s",
	"up.applied":      "%d rota(s) ativada(s) para '%s'",

	// ── CLI: odins down ──────────────────────────────────────────────────
	"down.no_project":  ".odins não encontrado em %s",
	"down.proxy_warn":  "⚠  %s: %v",
	"down.removed":     "✓ %s removido",
	"down.page_updated":"→ Landing page atualizada: https://%s.%s",
	"down.total":       "%d rota(s) removida(s) para '%s'",

	// ── CLI: odins kill ──────────────────────────────────────────────────
	"kill.not_found":  "rota '%s' não encontrada",
	"kill.proxy_warn": "⚠  proxy remove: %v",
	"kill.removed":    "✓ %s removido",

	// ── CLI: odins ls ────────────────────────────────────────────────────
	"ls.empty": "Nenhuma rota ativa. Use 'odins add' ou 'odins up' para adicionar.",

	// ── CLI: odins domain ────────────────────────────────────────────────
	"domain.exists":          "domínio '%s' já existe",
	"domain.page_warn":       "⚠  landing page: %v",
	"domain.caddy_warn":      "⚠  caddy domain route: %v",
	"domain.caddy_hint":      "(rode 'odins init' se o Caddy ainda não foi configurado)",
	"domain.created":         "✓ Domínio criado: https://%s",
	"domain.page_generated":  "→ Landing page gerada em %s",
	"domain.add_service":     "Para adicionar serviços, crie um .odins com:",
	"domain.then_up":         "Depois rode: odins up",
	"domain.empty":           "Nenhum domínio cadastrado.",
	"domain.empty_hint":      "Use: odins domain add <nome>",
	"domain.header.domain":   "DOMÍNIO",
	"domain.header.fqdn":     "FQDN",
	"domain.header.services": "SERVIÇOS",
	"domain.not_found":       "domínio '%s' não encontrado",
	"domain.removed":         "✓ Domínio '%s' e %d serviço(s) removidos.",
	"domain.proxy_warn":      "⚠  proxy remove %s: %v",
	"domain.caddy_rm_warn":   "⚠  caddy domain remove: %v",

	// ── Landing page HTML ────────────────────────────────────────────────
	"page.tagline":       "The All-Father of Local DNS",
	"page.services":      "service",
	"page.services_pl":   "services",
	"page.empty_state":   "Nenhum serviço ainda.\nAdicione projetos com <code>odins up</code>",
	"page.generated_by":  "Gerado por",

	// ── TLD labels (settings screen) ─────────────────────────────────────
	"tld.odin.label":     ".odin — temático, sem conflitos (padrão)",
	"tld.odins.label":    ".odins — variante temática",
	"tld.test.label":     ".test — reservado IANA, sem HSTS",
	"tld.dev.label":      ".dev — popular (requer HTTPS, HSTS preloaded)",
	"tld.lan.label":      ".lan — comum em redes locais",
	"tld.internal.label": ".internal — uso corporativo",
	"tld.local.label":    ".local — ⚠️  conflito com mDNS/Bonjour",
}
