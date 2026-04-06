# ODINS

```
  ____  ____  ___ _   _ ____
 / __ \|  _ \|_ _| \ | / ___|
| |  | | | | || ||  \| \___ \
| |__| | |_| || || |\  |___) |
 \____/|____/|___|_| \_|____/

  The All-Father of Local DNS
```

**ODINS** is a local DNS + reverse proxy manager for macOS developers.

Stop fighting with ports. Route your local projects to memorable subdomains with automatic HTTPS — and never type `localhost:3000` again.

```
https://api.rankly.odin   →  localhost:3000   (Node.js/Express)
https://app.rankly.odin   →  localhost:5173   (Vite/React)
https://api.rankly.odin   →  localhost:8080   (Go/Gin)
```

---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/adialaleal/odins/main/install.sh | bash
```

**Requirements:** macOS, [Homebrew](https://brew.sh)

---

## Quick Start

```bash
# 1. One-time setup (installs dnsmasq + caddy, configures DNS and HTTPS)
odins init

# 2. In any project directory
cd ~/Projects/rankly/api
odins up   # auto-detects Node.js/Go/Python and creates .odins

# 3. Open your browser
open https://api.rankly.odin

# 4. Manage routes
odins ls
odins kill api.rankly.odin

# 5. TUI dashboard
odins
```

---

## TLD Options

Choose your preferred TLD during `odins init`. All are supported:

| TLD | Notes |
|-----|-------|
| `.odin` | Thematic, no conflicts — **default** |
| `.odins` | Thematic variant |
| `.test` | IANA reserved for testing, no HSTS |
| `.dev` | Popular, but requires HTTPS (Chrome HSTS) — Caddy handles it |
| `.lan` | Common in local networks |
| `.internal` | Enterprise-style |
| `.local` | ⚠️ mDNS conflict on macOS — use with caution |

---

## Commands

| Command | Description |
|---------|-------------|
| `odins` | Open TUI dashboard |
| `odins init` | One-time setup: DNS, proxy, HTTPS |
| `odins up` | Apply routes from `.odins` in current directory |
| `odins down` | Remove all routes for current project |
| `odins add <subdomain> --port <port>` | Add a single route |
| `odins kill <subdomain>` | Remove a route |
| `odins ls` | List all active routes |

### `odins add` flags

```bash
odins add api.rankly.odin --port 3000
odins add worker.rankly.odin --port 5000 --docker rankly_worker_1
odins add admin.rankly.odin --port 8080 --no-https
```

---

## `.odins` Project Config

Place a `.odins` file in your project root:

```toml
[project]
name = "rankly"
runtime = "node"       # auto-detected
framework = "nextjs"   # auto-detected

[[routes]]
subdomain = "app"      # → app.rankly.<tld>
port = 3000
https = true

[[routes]]
subdomain = "api"
port = 4000
https = true

# Docker support
[[routes]]
subdomain = "worker"
port = 5000
docker_container = "rankly_worker_1"
https = true
```

Then run `odins up` from any directory in your project.

---

## Project Detection

When no `.odins` exists, `odins up` auto-detects your project:

### Node.js
- Detects via `package.json`
- Frameworks: **Next.js**, **Nuxt**, **Vite**, **Express**, **Fastify**, **NestJS**, **Remix**, **Hapi**, **Koa**
- Port: reads `PORT` from `.env`, `.env.local`, `.env.development`

### Go
- Detects via `go.mod`
- Frameworks: **Gin** (8080), **Echo** (1323), **Fiber** (3000), **Chi** (8080), **Gorilla** (8080)
- Port: scans `main.go` for `ListenAndServe`, `:Listen`, `:Run`

### Python
- Detects via `manage.py`, `pyproject.toml`, `requirements.txt`
- Frameworks: **Django** (8000), **FastAPI** (8000), **Flask** (5000), **Sanic** (8000), **Tornado** (8888)
- Port: reads `PORT` from `.env`

---

## Proxy Backends

Choose during `odins init`:

| Backend | Notes |
|---------|-------|
| **Caddy** *(recommended)* | Auto HTTPS via internal TLS, hot-reload via Admin API |
| **Nginx** | Familiar config, uses mkcert for HTTPS |
| **Apache** | VirtualHost-based, uses mkcert for HTTPS |

---

## Docker Support

ODINS works seamlessly with Docker projects:

```toml
[[routes]]
subdomain = "api"
port = 3000
docker_container = "my_api_container"
```

If the container exposes port 3000 to the host, ODINS proxies to `localhost:3000`. The TUI dashboard shows a `○` dot when the container is stopped.

---

## HTTPS

- **Caddy**: Uses Caddy's built-in TLS (automatically trusted after `odins init`)
- **Nginx/Apache**: Uses [mkcert](https://github.com/FiloSottile/mkcert) — installed automatically

After `odins init`, all `*.odin` (or your chosen TLD) domains are HTTPS by default — no browser warnings.

---

## Architecture

```
Browser → *.rankly.odin
             ↓ (dnsmasq wildcard → 127.0.0.1)
         Reverse Proxy :443
         (Caddy / Nginx / Apache)
             ↓
         localhost:3000  (your app)
```

- **DNS**: dnsmasq via Homebrew on port 5353, `/etc/resolver/<tld>` for macOS
- **Proxy**: Caddy/Nginx/Apache managed by ODINS
- **HTTPS**: Caddy internal TLS or mkcert, trusted in macOS keychain
- **State**: `~/.local/share/odins/routes.json`
- **Config**: `~/.config/odins/config.toml`

---

## Troubleshooting

**DNS not resolving?**
```bash
brew services restart dnsmasq
scutil --dns | grep odin   # should show nameserver 127.0.0.1
```

**Caddy not starting?**
```bash
brew services restart caddy
curl http://localhost:2019/config/   # Caddy Admin API
```

**Certificate not trusted?**
```bash
# For Caddy:
sudo security add-trusted-cert -d -r trustRoot \
  -k /Library/Keychains/System.keychain \
  ~/.local/share/caddy/pki/authorities/local/root.crt

# For mkcert:
mkcert -install
```

**Port 80/443 already in use?**
```bash
sudo lsof -i :80 -i :443
```

---

## License

MIT © [Adialá Leal](https://github.com/adialaleal)
