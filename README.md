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
https://api.rankly.odin   →  localhost:3000   (Node.js / Express)
https://app.rankly.odin   →  localhost:5173   (Vite / React)
https://jobs.rankly.odin  →  localhost:8080   (Go / Gin)
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

## AI Friendly

ODINS agora expõe uma superfície oficial para agentes via `CLI + JSON + docs`.

- Use `odins detect --json` para inspecionar um projeto sem alterar arquivos.
- Use `odins doctor --json` para diagnosticar DNS, proxy, HTTPS e store.
- Use `--json` nos comandos operacionais para saídas estáveis em automação.
- Consulte [AGENTS.md](AGENTS.md), [CLAUDE.md](CLAUDE.md) e [`ai/`](ai/) para os adapters publicados.

Documentação canônica:

- [docs/ai/setup-local.md](docs/ai/setup-local.md)
- [docs/ai/apply-to-project.md](docs/ai/apply-to-project.md)
- [docs/ai/workspace-multi-service.md](docs/ai/workspace-multi-service.md)
- [docs/ai/doctor-troubleshoot.md](docs/ai/doctor-troubleshoot.md)

---

## Quick Start

```bash
# 1. One-time setup — installs dnsmasq + Caddy, configures DNS and HTTPS
odins init

# 2. Inspect any project first
cd ~/Projects/my-api
odins detect --json

# 3. Apply routes in the project directory
odins up
#  → Detectado: node/vite (porta 5173)
#  → ✓ https://my-api.my-api.odin → :5173

# 4. Validate routes and the local environment
odins ls
odins doctor

# 5. Open the TUI dashboard
odins
```

---

## TUI Dashboard

```
┌─ ODINS ─────────────────────────────── 4 rotas ─────────────────────────────┐
│                                                                               │
│  STATUS  SUBDOMAIN              PORT   PROTO  RUNTIME   PROJECT              │
│  ──────  ─────────              ────   ─────  ───────   ───────              │
│  ●       app.rankly.odin        3000   HTTPS  node      rankly               │
│  ●       api.rankly.odin        4000   HTTPS  node      rankly               │
│  ●       jobs.rankly.odin       8080   HTTPS  go        rankly               │
│  ○       worker.rankly.odin     5000   HTTPS  docker    rankly               │
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
| `odins detect` | Inspect project runtime/framework/port and recommend a `.odins` |
| `odins up` | Apply routes from `.odins` in current directory |
| `odins down` | Remove all routes for current project |
| `odins add <sub> --port <n>` | Add a single route |
| `odins kill <subdomain>` | Remove a specific route |
| `odins ls` | List all active routes |
| `odins doctor` | Diagnose Homebrew, DNS, proxy, HTTPS and store health |
| `odins domain add <name>` | Create a workspace domain landing page |
| `odins domain ls` | List workspace domains |
| `odins domain rm <name>` | Remove a workspace domain |

### Machine-readable mode

All operational commands support `--json`.

```bash
odins detect --json
odins ls --json
odins doctor --json
odins init --json --non-interactive --tld odin --backend caddy
```

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
subdomain = "app"        # → app.rankly.odin
port      = 3000
https     = true

[[routes]]
subdomain = "api"        # → api.rankly.odin
port      = 4000
https     = true

[[routes]]
subdomain = "worker"     # Docker container
port      = 5000
docker_container = "rankly_worker_1"
https     = true
```

---

## Prompt Recipes

Copy and paste one of these prompts into your AI coding tool:

### Detect a project

```text
Explore este repositório primeiro. Depois rode `odins detect --json` na raiz do projeto, resuma runtime/framework/porta detectados e proponha o `.odins` recomendado sem aplicar mudanças ainda.
```

### Propose and apply `.odins`

```text
Explore este projeto, rode `odins detect --json`, proponha o `.odins` ideal e só depois aplique com `odins up`. Se houver qualquer ação com sudo, me avise antes.
```

### Connect to a workspace domain

```text
Explore o projeto, sugira como conectá-lo a um workspace ODINS existente ou novo usando `odins domain add` e o campo `domain` no `.odins`. Prefira `odins detect --json` antes de qualquer mudança.
```

### Validate the environment

```text
Rode `odins doctor --json`, explique cada check com problema e proponha a próxima ação mínima para deixar o ambiente saudável.
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
| `.odin` | Thematic, no conflicts — **default** |
| `.odins` | Thematic variant |
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
Browser → *.rankly.odin
               ↓
         dnsmasq :5300   (wildcard *.odin → 127.0.0.1)
               ↓
         /etc/resolver/odin  (macOS resolver)
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
odins doctor
brew services restart dnsmasq
scutil --dns | grep odin   # should show nameserver 127.0.0.1, port 5300
```

**Caddy not starting?**
```bash
odins doctor --json
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
