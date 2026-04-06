# Changelog

All notable changes to ODINS are documented here.

## [Unreleased]

### Added
- `odins detect` for read-only project inspection and `.odins` recommendation
- `odins doctor` for local environment diagnostics
- `--json` structured output for operational commands
- `--non-interactive`, `--tld`, and `--backend` flags for `odins init`
- AI-friendly docs under `docs/ai/`
- Published adapters for Codex, Claude Code, and Antigravity
- Detection fixtures, JSON golden tests, and AI pack sync checks

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
- 7 TLD options: `.odin`, `.odins`, `.test`, `.dev`, `.lan`, `.internal`, `.local`
- macOS DNS resolver via `/etc/resolver/<tld>` (one-time sudo)
- GitHub Actions CI and GoReleaser release pipeline
