package i18n

var enStrings = map[string]string{
	// ── Welcome / onboarding ─────────────────────────────────────────────
	"welcome.section.what":     "What is ODINS?",
	"welcome.section.how":      "How does it work?",
	"welcome.section.domains":  "Domains and Subdomains",
	"welcome.section.config":   "Per-Project Config (.odins)",
	"welcome.section.commands": "Main Commands",
	"welcome.section.next":     "Next Steps",
	"welcome.tagline":          "The All-Father of Local DNS",
	"welcome.elimina":          "ODINS eliminates port conflicts in your local development.",
	"welcome.sem_odins":        "Without ODINS:",
	"welcome.com_odins":        "With ODINS:",
	"welcome.https_auto":       "Each project gets a clean subdomain with automatic HTTPS.",
	"welcome.how_dns":          "DNS   — dnsmasq resolves *.<project>.odins → 127.0.0.1",
	"welcome.how_proxy":        "Proxy — Caddy routes web.<project>.odins → localhost:3000",
	"welcome.how_https":        "HTTPS — Caddy manages certificates automatically",
	"welcome.domain_is":        "A domain is the central workspace for your projects:",
	"welcome.domain_landing":   "creates <project>.odins — a landing page listing all workspace services with real-time status.",
	"welcome.subdomain_is":     "Each subdomain is a project/service:",
	"welcome.auto_detect":      "ODINS auto-detects Node.js, Go, and Python projects.",
	"welcome.not_configured":   "It looks like ODINS has not been configured on this machine yet.",
	"welcome.run_init":         "Run odins init now?",
	"welcome.run_init_prompt":  "[Y/n]",
	"welcome.run_init_yes":     "y",
	"welcome.ok":               "No problem! When ready, run: odins init",
	"welcome.already_configured": "✓ ODINS is already configured.",
	"welcome.create_domain":    "Start by creating a domain:",
	"welcome.then_project":     "Then inside a project:",
	"welcome.enter":            "[Press Enter to continue...]",
	// New-folder welcome (short)
	"welcome.new_folder.title":     "Getting started in this project",
	"welcome.new_folder.detected":  "Project detected: %s (%s/%s, port %d)",
	"welcome.new_folder.activate":  "Run odins up to activate routes automatically.",
	"welcome.new_folder.manual":    "No project detected. Create a .odins file with:",
	"welcome.new_folder.or_add":    "Or add a route manually:",
	"welcome.new_folder.see_guide": "Full guide: odins welcome",

	// ── Command descriptions ──────────────────────────────────────────────
	"cmd.init_desc":       "One-time setup: DNS, proxy, HTTPS",
	"cmd.domain_add_desc": "Create workspace <project>.odins",
	"cmd.up_desc":         "Activate routes for the current project",
	"cmd.ls_desc":         "List active routes",
	"cmd.kill_desc":       "Remove a route",
	"cmd.down_desc":       "Remove all routes for the current project",
	"cmd.tui_desc":        "Open TUI dashboard",
	"cmd.welcome_desc":    "Show this guide again",

	// ── TUI app messages ─────────────────────────────────────────────────
	"tui.detected":         "Project '%s' (%s/%s, port %d) detected — press [u] to activate routes",
	"tui.no_project":       "No .odins found — press [a] to add a route manually",
	"tui.activating":       "Activating routes...",
	"tui.error":            "Error: %s",
	"tui.no_routes":        "No routes applied",
	"tui.routes_activated": "✓ %d route(s) activated!",
	"tui.config_save_err":  "Error saving config: %s",
	"tui.config_saved":     "Settings saved!",
	"tui.save_error":       "Error saving: %s",
	"tui.proxy_error":      "Proxy error: %s",
	"tui.route_added":      "✓ %s → :%d added!",
	"tui.route_removed":    "%s removed.",
	"tui.read_project_err": "read .odins: %s",
	"tui.no_detect":        "project not detected in %s",

	// ── Footer hints ─────────────────────────────────────────────────────
	"hint.add":        "add",
	"hint.up":         "odins up",
	"hint.remove":     "remove",
	"hint.settings":   "settings",
	"hint.logs":       "logs",
	"hint.quit":       "quit",
	"hint.back":       "back",
	"hint.next_field": "next field",
	"hint.confirm":    "confirm",
	"hint.cancel":     "cancel",
	"hint.select":     "select",
	"hint.navigate":   "navigate",
	"hint.save":       "save",
	"hint.scroll":     "scroll",

	// ── Dashboard ────────────────────────────────────────────────────────
	"dash.title":          "Dashboard — Active Routes",
	"dash.routes":         "%d routes",
	"dash.remove_confirm": "Remove %s?",

	// ── Add route form ───────────────────────────────────────────────────
	"add.title":             "Add Route",
	"add.field.subdomain":   "Subdomain",
	"add.field.port":        "Port",
	"add.field.project":     "Project",
	"add.field.docker":      "Docker",
	"add.detected":          "✦ Detected: %s/%s (port %d)",
	"add.err.subdomain":     "subdomain is required",
	"add.err.port_required": "port is required",
	"add.err.port_invalid":  "invalid port",

	// ── Settings screen ──────────────────────────────────────────────────
	"settings.title":       "Settings",
	"settings.tld_label":   "TLD:",
	"settings.proxy_label": "Proxy:",
	"settings.warn_local":  "⚠  .local conflicts with mDNS/Bonjour on macOS — use with caution",
	"settings.info":        "Domains: *%s → 127.0.0.1  |  Proxy: %s",

	// ── Logs screen ──────────────────────────────────────────────────────
	"logs.title": "Proxy Logs",

	// ── Confirm modal ────────────────────────────────────────────────────
	"modal.confirm": "confirm",
	"modal.cancel":  "cancel",

	// ── CLI: odins up ────────────────────────────────────────────────────
	"up.reading":      "→ Reading .odins for project '%s'",
	"up.detecting":    "→ .odins not found, detecting project...",
	"up.not_detected": "could not detect project type in %s\nCreate a .odins manually or use: odins add <subdomain> --port <port>",
	"up.detected":     "→ Detected: %s/%s (port %d)",
	"up.start_cmd":    "→ Start command: %s",
	"up.save_warn":    "⚠  Could not save .odins: %v",
	"up.created":      "→ .odins created at %s",
	"up.domain":       "→ Domain: %s.%s",
	"up.route_error":  "✗ %s: %v",
	"up.route_ok":     "✓ %s://%s → :%d",
	"up.page_updated": "→ Landing page updated: https://%s.%s",
	"up.applied":      "%d route(s) activated for '%s'",

	// ── CLI: odins down ──────────────────────────────────────────────────
	"down.no_project":   ".odins not found in %s",
	"down.proxy_warn":   "⚠  %s: %v",
	"down.removed":      "✓ %s removed",
	"down.page_updated": "→ Landing page updated: https://%s.%s",
	"down.total":        "%d route(s) removed for '%s'",

	// ── CLI: odins kill ──────────────────────────────────────────────────
	"kill.not_found":  "route '%s' not found",
	"kill.proxy_warn": "⚠  proxy remove: %v",
	"kill.removed":    "✓ %s removed",

	// ── CLI: odins ls ────────────────────────────────────────────────────
	"ls.empty": "No active routes. Use 'odins add' or 'odins up' to add one.",

	// ── CLI: odins domain ────────────────────────────────────────────────
	"domain.exists":          "domain '%s' already exists",
	"domain.page_warn":       "⚠  landing page: %v",
	"domain.caddy_warn":      "⚠  caddy domain route: %v",
	"domain.caddy_hint":      "(run 'odins init' if Caddy has not been configured yet)",
	"domain.created":         "✓ Domain created: https://%s",
	"domain.page_generated":  "→ Landing page generated at %s",
	"domain.add_service":     "To add services, create a .odins file with:",
	"domain.then_up":         "Then run: odins up",
	"domain.empty":           "No domains registered.",
	"domain.empty_hint":      "Use: odins domain add <name>",
	"domain.header.domain":   "DOMAIN",
	"domain.header.fqdn":     "FQDN",
	"domain.header.services": "SERVICES",
	"domain.not_found":       "domain '%s' not found",
	"domain.removed":         "✓ Domain '%s' and %d service(s) removed.",
	"domain.proxy_warn":      "⚠  proxy remove %s: %v",
	"domain.caddy_rm_warn":   "⚠  caddy domain remove: %v",

	// ── Landing page HTML ────────────────────────────────────────────────
	"page.tagline":     "The All-Father of Local DNS",
	"page.services":    "service",
	"page.services_pl": "services",
	"page.empty_state": "No services yet.\nAdd projects with <code>odins up</code>",
	"page.generated_by":"Generated by",

	// ── TLD labels ───────────────────────────────────────────────────────
	"tld.odin.label":     ".odin — thematic, no conflicts (default)",
	"tld.odins.label":    ".odins — thematic variant",
	"tld.test.label":     ".test — IANA reserved, no HSTS",
	"tld.dev.label":      ".dev — popular (requires HTTPS, HSTS preloaded)",
	"tld.lan.label":      ".lan — common on local networks",
	"tld.internal.label": ".internal — corporate use",
	"tld.local.label":    ".local — ⚠️  conflicts with mDNS/Bonjour",
}
