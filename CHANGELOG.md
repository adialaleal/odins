# Changelog

All notable changes to ODINS are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

---

## [0.7.0] - 2026-04-06

### Added
- **i18n** — full PT / EN / ES support across all user-facing strings (`internal/i18n/`)
- **Language auto-detection** from `LANG` / `LC_ALL` / `LC_MESSAGES` env vars; override via `language` field in config
- **`$HOME/.odins` global config** — `odins up` falls back to `~/.odins` when no project config exists; `odins up --global` forces it
- **Welcome trigger for new folders** — running `odins` in a directory without `.odins` now shows the onboarding guide even after first-run
- **Project-name domain** — `odins up` auto-derives the domain from `project.name` when no explicit `domain` is set in `.odins`
- **`odins domain add`** — creates a domain workspace with a landing page served by Caddy
- **Auto-sync Caddy on startup** — routes and domain landing pages are re-applied from state on every `odins` command; routes survive Caddy restarts without re-running `odins up`
- **macOS auth dialog** for privileged operations (`/etc/resolver`, CA trust) — uses `osascript` instead of terminal `sudo`
- **`SudoFlushDNS`** helper — flushes macOS DNS cache and restarts `mDNSResponder` via auth dialog

### Fixed
- **dnsmasq port changed 5353 → 5300** — port 5353 is held by macOS `mDNSResponder` (Bonjour) at kernel level; dnsmasq silently failed to bind
- **Caddy crash-loop on first start** — `odins init` now creates a minimal `Caddyfile` before starting the Caddy brew service
- **Caddy route append API** — fixed invalid `POST .../routes/...` path (not supported in Caddy v2.11); corrected to `POST .../routes`
- **dnsmasq config** — writes directly to Homebrew path instead of symlink to prevent broken-link failures on macOS
- **`settings.go`** — removed unused `"fmt"` import after i18n migration

### Changed
- All user-facing strings in CLI and TUI are now routed through `i18n.T()` / `i18n.Tf()`
- Welcome guide examples updated from hardcoded `tatoh.odins` → generic `<projeto>.odins`
- `odins init` flushes macOS DNS cache after writing `/etc/resolver/<tld>`

---

## [0.1.0] - 2026-04-06

### Added
- Initial release
- `odins init` — one-time setup: dnsmasq, Caddy/Nginx/Apache, HTTPS, DNS resolver
- `odins up` — reads `.odins` project config and applies all routes
- `odins down` — removes all routes for the current project
- `odins add <subdomain> --port <n>` — add a single reverse proxy route
- `odins kill <subdomain>` — remove a route
- `odins ls` — list all active routes with live status
- TUI dashboard with Bubble Tea (slide-up screen transitions)
- Auto-detection of Node.js, Go, Python projects and their frameworks
- Support for Caddy (default), Nginx, and Apache as reverse proxy backends
- HTTPS via Caddy internal TLS or mkcert
- Docker container routing via `docker_container` field in `.odins`
- 7 TLD options: `.odins`, `.odin`, `.test`, `.dev`, `.lan`, `.internal`, `.local`
- macOS DNS resolver via `/etc/resolver/<tld>` (one-time sudo)
- GitHub Actions CI and GoReleaser release pipeline
