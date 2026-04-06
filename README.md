<div align="center">

```
  ____  ____  ___ _   _ ____
 / __ \|  _ \|_ _| \ | / ___|
| |  | | | | || ||  \| \___ \
| |__| | |_| || || |\  |___) |
 \____/|____/|___|_| \_|____/
```

**The All-Father of Local DNS**

[![CI](https://github.com/adialaleal/odins/actions/workflows/ci.yml/badge.svg)](https://github.com/adialaleal/odins/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/adialaleal/odins?color=7c3aed&label=latest)](https://github.com/adialaleal/odins/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/license-MIT-6d28d9)](LICENSE)
[![macOS](https://img.shields.io/badge/macOS-only-lightgrey?logo=apple)](https://www.apple.com/macos/)

*Stop fighting with ports. Route your local projects to beautiful domains with automatic HTTPS.*

[Install](#install) · [Quick Start](#quick-start) · [Commands](#commands) · [Contributing](#contributing)

</div>

---

## What is ODINS?

ODINS is a **local DNS + reverse proxy manager** for macOS developers.

Instead of juggling `localhost:3000`, `localhost:5173`, `localhost:8080`... you get:

```
https://api.rankly.odins   →  localhost:3000   (Node.js / Express)
https://app.rankly.odins   →  localhost:5173   (Vite / React)
https://jobs.rankly.odins  →  localhost:8080   (Go / Gin)
```

Zero config. One command. Automatic HTTPS. Beautiful TUI dashboard.

---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/adialaleal/odins/main/install.sh | bash
```

**Requirements:** macOS · [Homebrew](https://brew.sh)

> **Homebrew tap** *(coming soon)*
> ```bash
> brew install adialaleal/odins/odins
> ```

---

## Quick Start

```bash
# 1. One-time setup — installs dnsmasq + Caddy, configures DNS and HTTPS
odins init

# 2. Go to any project directory
cd ~/Projects/my-api

# 3. Apply routes — auto-detects Node.js / Go / Python
odins up
#  → Detectado: node/nextjs (porta 3000)
#  → ✓ https://my-api.my-api.odins → :3000

# 4. Open in browser
open https://my-api.odins

# 5. Or manage everything from the TUI
odins
```

---

## TUI Dashboard

```
┌─ ODINS ─────────────────────────────── 4 rotas ─────────────────────────────┐
│                                                                               │
│  STATUS  SUBDOMAIN              PORT   PROTO  RUNTIME   PROJECT              │
│  ──────  ─────────              ────   ─────  ───────   ───────              │
│  ●       app.rankly.odins       3000   HTTPS  node      rankly               │
│  ●       api.rankly.odins       4000   HTTPS  node      rankly               │
│  ●       jobs.rankly.odins      8080   HTTPS  go        rankly               │
│  ○       worker.rankly.odins    5000   HTTPS  docker    rankly               │
│                                                                               │
│  [a] adicionar  [u] odins up  [x] remover  [s] settings  [l] logs  [q] sair │
└───────────────────────────────────────────────────────────────────────────────┘
```

---

## Commands

| Command | Description |
|---|---|
| `odins` | Open TUI dashboard |
| `odins init` | One-time setup: DNS, proxy, HTTPS |
| `odins up` | Apply routes from `.odins` in current directory |
| `odins down` | Remove all routes for current project |
| `odins domain add <name>` | Create a domain workspace (landing page) |
| `odins add <sub> --port <n>` | Add a single route |
| `odins kill <subdomain>` | Remove a specific route |
| `odins ls` | List all active routes |

### Flags

```bash
# Route with Docker container
odins add worker.rankly --port 5000 --docker rankly_worker_1

# HTTP only (no HTTPS)
odins add admin.rankly --port 8080 --no-https

# Use global ~/.odins config
odins up --global

# Use a specific directory
odins up --dir ~/Projects/my-api
```

---

## `.odins` Project Config

Place a `.odins` file in your project root — or let `odins up` generate one automatically:

```toml
[project]
name      = "rankly"
runtime   = "node"       # auto-detected: node | go | python
framework = "nextjs"     # auto-detected
domain    = "rankly"     # → routes become sub.rankly.<tld>

[[routes]]
subdomain = "app"        # → app.rankly.odins
port      = 3000
https     = true

[[routes]]
subdomain = "api"        # → api.rankly.odins
port      = 4000
https     = true

[[routes]]
subdomain = "worker"     # Docker container
port      = 5000
docker_container = "rankly_worker_1"
https     = true
```

---

## Project Auto-Detection

`odins up` detects your stack automatically — no config needed:

| Runtime | Detected via | Frameworks |
|---|---|---|
| **Node.js** | `package.json` | Next.js, Nuxt, Vite, Express, Fastify, NestJS, Remix, Hapi, Koa |
| **Go** | `go.mod` | Gin, Echo, Fiber, Chi, Gorilla |
| **Python** | `requirements.txt`, `pyproject.toml`, `manage.py` | FastAPI, Django, Flask, Sanic, Tornado |

Port is read from `.env` / `.env.local` / `.env.development` when present.

---

## Proxy Backends

Choose during `odins init`:

| Backend | Notes |
|---|---|
| **Caddy** *(recommended)* | Auto HTTPS via internal TLS · hot-reload via Admin API |
| **Nginx** | Familiar config · uses mkcert for HTTPS |
| **Apache** | VirtualHost-based · uses mkcert for HTTPS |

---

## TLD Options

| TLD | Notes |
|---|---|
| `.odins` | Thematic, no conflicts — **default** |
| `.odin` | Shorter variant |
| `.test` | IANA reserved for testing |
| `.dev` | Popular · HTTPS required (Caddy handles it) |
| `.lan` | Common in local networks |
| `.internal` | Enterprise-style |
| `.local` | ⚠️ mDNS conflict on macOS — use with caution |

---

## Multi-Language Support

ODINS auto-detects your system language:

```bash
LANG=en_US.UTF-8 odins ls   # English
LANG=es_ES.UTF-8 odins ls   # Español
LANG=pt_BR.UTF-8 odins ls   # Português
```

Override in `~/.config/odins/config.toml`:

```toml
language = "en"   # "pt" | "en" | "es"
```

---

## Architecture

```
Browser → *.rankly.odins
               ↓
         dnsmasq :5300   (wildcard *.odins → 127.0.0.1)
               ↓
         /etc/resolver/odins  (macOS resolver)
               ↓
         Caddy :443 / :80  (reverse proxy + HTTPS)
               ↓
         localhost:<port>  (your app)
```

- **DNS** — dnsmasq via Homebrew on port `5300`, `/etc/resolver/<tld>` for macOS routing
- **Proxy** — Caddy / Nginx / Apache managed entirely by ODINS
- **HTTPS** — Caddy internal TLS (auto-trusted) or mkcert for Nginx/Apache
- **State** — `~/Library/Application Support/odins/routes.json`
- **Config** — `~/.config/odins/config.toml`

---

## Troubleshooting

**DNS not resolving?**
```bash
brew services restart dnsmasq
scutil --dns | grep odins   # should show nameserver 127.0.0.1, port 5300
```

**Caddy not starting?**
```bash
brew services restart caddy
curl http://localhost:2019/config/
```

**Certificate not trusted in browser?**
```bash
# Run this once — opens macOS auth dialog
odins init
```

**Route missing after Caddy restart?**
Routes are automatically re-synced from state on every `odins` command. Just run any `odins` command to restore.

**Port 80/443 already in use?**
```bash
sudo lsof -i :80 -i :443
```

---

## Contributing

ODINS is open source and contributions are very welcome! See **[CONTRIBUTING.md](CONTRIBUTING.md)** for the full guide.

**Quick start for contributors:**

```bash
git clone https://github.com/adialaleal/odins
cd odins
go mod tidy
go build -o odins .
./odins --help
```

**Good first issues** are labeled [`good first issue`](https://github.com/adialaleal/odins/issues?q=label%3A%22good+first+issue%22) on GitHub.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for the full history.

---

<div align="center">

MIT © [Adialá Leal](https://github.com/adialaleal) · Built with [Go](https://go.dev) + [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Caddy](https://caddyserver.com)

*If ODINS saves you time, consider giving it a ⭐*

</div>
