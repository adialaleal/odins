package i18n

var esStrings = map[string]string{
	// ── Welcome / onboarding ─────────────────────────────────────────────
	"welcome.section.what":     "¿Qué es ODINS?",
	"welcome.section.how":      "¿Cómo funciona?",
	"welcome.section.domains":  "Dominios y Subdominios",
	"welcome.section.config":   "Configuración por Proyecto (.odins)",
	"welcome.section.commands": "Comandos Principales",
	"welcome.section.next":     "Próximos Pasos",
	"welcome.tagline":          "The All-Father of Local DNS",
	"welcome.elimina":          "ODINS elimina los conflictos de puertos en tu desarrollo local.",
	"welcome.sem_odins":        "Sin ODINS:",
	"welcome.com_odins":        "Con ODINS:",
	"welcome.https_auto":       "Cada proyecto obtiene un subdominio elegante con HTTPS automático.",
	"welcome.how_dns":          "DNS   — dnsmasq resuelve *.<proyecto>.odins → 127.0.0.1",
	"welcome.how_proxy":        "Proxy — Caddy enruta web.<proyecto>.odins → localhost:3000",
	"welcome.how_https":        "HTTPS — Caddy gestiona los certificados automáticamente",
	"welcome.domain_is":        "Un dominio es el workspace central de tus proyectos:",
	"welcome.domain_landing":   "crea <proyecto>.odins — una landing page que lista todos los servicios del workspace con estado en tiempo real.",
	"welcome.subdomain_is":     "Cada subdominio es un proyecto/servicio:",
	"welcome.auto_detect":      "ODINS detecta automáticamente proyectos Node.js, Go y Python.",
	"welcome.not_configured":   "Parece que ODINS aún no ha sido configurado en esta máquina.",
	"welcome.run_init":         "¿Ejecutar odins init ahora?",
	"welcome.run_init_prompt":  "[S/n]",
	"welcome.run_init_yes":     "s",
	"welcome.ok":               "¡Sin problema! Cuando estés listo, ejecuta: odins init",
	"welcome.already_configured": "✓ ODINS ya está configurado.",
	"welcome.create_domain":    "Empieza creando un dominio:",
	"welcome.then_project":     "Luego en un proyecto:",
	"welcome.enter":            "[Presiona Enter para continuar...]",
	// New-folder welcome (short)
	"welcome.new_folder.title":     "Empezando en este proyecto",
	"welcome.new_folder.detected":  "Proyecto detectado: %s (%s/%s, puerto %d)",
	"welcome.new_folder.activate":  "Ejecuta odins up para activar las rutas automáticamente.",
	"welcome.new_folder.manual":    "Ningún proyecto detectado. Crea un .odins con:",
	"welcome.new_folder.or_add":    "O agrega una ruta manualmente:",
	"welcome.new_folder.see_guide": "Guía completa: odins welcome",

	// ── Command descriptions ──────────────────────────────────────────────
	"cmd.init_desc":       "Configuración única: DNS, proxy, HTTPS",
	"cmd.domain_add_desc": "Crear workspace <proyecto>.odins",
	"cmd.up_desc":         "Activar rutas del proyecto actual",
	"cmd.ls_desc":         "Listar rutas activas",
	"cmd.kill_desc":       "Eliminar una ruta",
	"cmd.down_desc":       "Eliminar todas las rutas del proyecto",
	"cmd.tui_desc":        "Abrir panel TUI",
	"cmd.welcome_desc":    "Ver esta guía de nuevo",

	// ── TUI app messages ─────────────────────────────────────────────────
	"tui.detected":         "Proyecto '%s' (%s/%s, puerto %d) detectado — presiona [u] para activar las rutas",
	"tui.no_project":       "No se encontró .odins — presiona [a] para agregar una ruta manualmente",
	"tui.activating":       "Activando rutas...",
	"tui.error":            "Error: %s",
	"tui.no_routes":        "Ninguna ruta aplicada",
	"tui.routes_activated": "✓ %d ruta(s) activada(s)!",
	"tui.config_save_err":  "Error al guardar configuración: %s",
	"tui.config_saved":     "¡Configuración guardada!",
	"tui.save_error":       "Error al guardar: %s",
	"tui.proxy_error":      "Error en el proxy: %s",
	"tui.route_added":      "✓ %s → :%d agregado!",
	"tui.route_removed":    "%s eliminado.",
	"tui.read_project_err": "leer .odins: %s",
	"tui.no_detect":        "proyecto no detectado en %s",

	// ── Footer hints ─────────────────────────────────────────────────────
	"hint.add":        "agregar",
	"hint.up":         "odins up",
	"hint.remove":     "eliminar",
	"hint.settings":   "ajustes",
	"hint.logs":       "logs",
	"hint.quit":       "salir",
	"hint.back":       "volver",
	"hint.next_field": "siguiente campo",
	"hint.confirm":    "confirmar",
	"hint.cancel":     "cancelar",
	"hint.select":     "seleccionar",
	"hint.navigate":   "navegar",
	"hint.save":       "guardar",
	"hint.scroll":     "scroll",

	// ── Dashboard ────────────────────────────────────────────────────────
	"dash.title":          "Dashboard — Rutas Activas",
	"dash.routes":         "%d rutas",
	"dash.remove_confirm": "¿Eliminar %s?",

	// ── Add route form ───────────────────────────────────────────────────
	"add.title":             "Agregar Ruta",
	"add.field.subdomain":   "Subdominio",
	"add.field.port":        "Puerto",
	"add.field.project":     "Proyecto",
	"add.field.docker":      "Docker",
	"add.detected":          "✦ Detectado: %s/%s (puerto %d)",
	"add.err.subdomain":     "el subdominio es obligatorio",
	"add.err.port_required": "el puerto es obligatorio",
	"add.err.port_invalid":  "puerto inválido",

	// ── Settings screen ──────────────────────────────────────────────────
	"settings.title":       "Configuración",
	"settings.tld_label":   "TLD:",
	"settings.proxy_label": "Proxy:",
	"settings.warn_local":  "⚠  .local conflicta con mDNS/Bonjour en macOS — usar con precaución",
	"settings.info":        "Dominios: *%s → 127.0.0.1  |  Proxy: %s",

	// ── Logs screen ──────────────────────────────────────────────────────
	"logs.title": "Logs del Proxy",

	// ── Confirm modal ────────────────────────────────────────────────────
	"modal.confirm": "confirmar",
	"modal.cancel":  "cancelar",

	// ── CLI: odins up ────────────────────────────────────────────────────
	"up.reading":      "→ Leyendo .odins del proyecto '%s'",
	"up.detecting":    "→ .odins no encontrado, detectando proyecto...",
	"up.not_detected": "no se pudo detectar el tipo de proyecto en %s\nCrea un .odins manualmente o usa: odins add <subdomain> --port <port>",
	"up.detected":     "→ Detectado: %s/%s (puerto %d)",
	"up.start_cmd":    "→ Comando de inicio: %s",
	"up.save_warn":    "⚠  No se pudo guardar .odins: %v",
	"up.created":      "→ .odins creado en %s",
	"up.domain":       "→ Dominio: %s.%s",
	"up.route_error":  "✗ %s: %v",
	"up.route_ok":     "✓ %s://%s → :%d",
	"up.page_updated": "→ Landing page actualizada: https://%s.%s",
	"up.applied":      "%d ruta(s) activada(s) para '%s'",

	// ── CLI: odins down ──────────────────────────────────────────────────
	"down.no_project":   ".odins no encontrado en %s",
	"down.proxy_warn":   "⚠  %s: %v",
	"down.removed":      "✓ %s eliminado",
	"down.page_updated": "→ Landing page actualizada: https://%s.%s",
	"down.total":        "%d ruta(s) eliminada(s) para '%s'",

	// ── CLI: odins kill ──────────────────────────────────────────────────
	"kill.not_found":  "ruta '%s' no encontrada",
	"kill.proxy_warn": "⚠  proxy remove: %v",
	"kill.removed":    "✓ %s eliminado",

	// ── CLI: odins ls ────────────────────────────────────────────────────
	"ls.empty": "Sin rutas activas. Usa 'odins add' o 'odins up' para agregar una.",

	// ── CLI: odins domain ────────────────────────────────────────────────
	"domain.exists":          "el dominio '%s' ya existe",
	"domain.page_warn":       "⚠  landing page: %v",
	"domain.caddy_warn":      "⚠  caddy domain route: %v",
	"domain.caddy_hint":      "(ejecuta 'odins init' si Caddy aún no ha sido configurado)",
	"domain.created":         "✓ Dominio creado: https://%s",
	"domain.page_generated":  "→ Landing page generada en %s",
	"domain.add_service":     "Para agregar servicios, crea un .odins con:",
	"domain.then_up":         "Luego ejecuta: odins up",
	"domain.empty":           "Ningún dominio registrado.",
	"domain.empty_hint":      "Usa: odins domain add <nombre>",
	"domain.header.domain":   "DOMINIO",
	"domain.header.fqdn":     "FQDN",
	"domain.header.services": "SERVICIOS",
	"domain.not_found":       "dominio '%s' no encontrado",
	"domain.removed":         "✓ Dominio '%s' y %d servicio(s) eliminados.",
	"domain.proxy_warn":      "⚠  proxy remove %s: %v",
	"domain.caddy_rm_warn":   "⚠  caddy domain remove: %v",

	// ── Landing page HTML ────────────────────────────────────────────────
	"page.tagline":     "The All-Father of Local DNS",
	"page.services":    "servicio",
	"page.services_pl": "servicios",
	"page.empty_state": "Sin servicios aún.\nAgrega proyectos con <code>odins up</code>",
	"page.generated_by":"Generado por",

	// ── TLD labels ───────────────────────────────────────────────────────
	"tld.odin.label":     ".odin — temático, sin conflictos (predeterminado)",
	"tld.odins.label":    ".odins — variante temática",
	"tld.test.label":     ".test — reservado IANA, sin HSTS",
	"tld.dev.label":      ".dev — popular (requiere HTTPS, HSTS preloaded)",
	"tld.lan.label":      ".lan — común en redes locales",
	"tld.internal.label": ".internal — uso corporativo",
	"tld.local.label":    ".local — ⚠️  conflicto con mDNS/Bonjour",
}
