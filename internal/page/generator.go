// Package page generates the landing page HTML for ODINS domains.
package page

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/adrg/xdg"
)

// RouteInfo is a simplified view of a route for HTML templating.
type RouteInfo struct {
	Subdomain string
	FQDN      string
	Port      int
	Runtime   string
	Project   string
}

// PageData is passed to the HTML template.
type PageData struct {
	Domain      string
	TLD         string
	FQDN        string // domain.tld
	Title       string
	Description string
	Routes      []RouteInfo
	GeneratedAt string
}

// PagesDir returns the base directory for all domain landing pages.
func PagesDir() string {
	return filepath.Join(xdg.DataHome, "odins", "pages")
}

// PageDir returns the directory for a specific domain's landing page.
func PageDir(domain string) string {
	return filepath.Join(PagesDir(), domain)
}

// Generate writes the index.html for the given domain to disk.
func Generate(data PageData) error {
	data.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")
	if data.FQDN == "" {
		data.FQDN = data.Domain + "." + data.TLD
	}
	if data.Title == "" {
		data.Title = data.FQDN
	}

	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	dir := PageDir(data.Domain)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir pages: %w", err)
	}

	path := filepath.Join(dir, "index.html")
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write page: %w", err)
	}

	return nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}} — ODINS</title>
  <style>
    :root {
      --bg:      #0f0f0f;
      --surface: #1a1a2e;
      --border:  #3730a3;
      --violet:  #7c3aed;
      --vlight:  #a78bfa;
      --vdim:    #6d28d9;
      --text:    #e5e7eb;
      --muted:   #9ca3af;
      --green:   #10b981;
      --gray:    #6b7280;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: var(--bg);
      color: var(--text);
      font-family: ui-monospace, 'JetBrains Mono', 'Fira Code', monospace;
      min-height: 100vh;
    }

    /* ── Header ── */
    header {
      padding: 2rem 2.5rem 1.5rem;
      background: var(--surface);
      border-bottom: 1px solid var(--border);
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }
    .logo-block {}
    .logo {
      font-size: 1.6rem;
      font-weight: 700;
      letter-spacing: .05em;
      background: linear-gradient(135deg, #e9d5ff, #a78bfa, #7c3aed);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
    .domain-fqdn {
      color: var(--vdim);
      font-size: .95rem;
      margin-top: .35rem;
    }
    .domain-desc {
      color: var(--muted);
      font-size: .82rem;
      margin-top: .3rem;
    }
    .header-meta {
      text-align: right;
      font-size: .75rem;
      color: var(--muted);
      padding-top: .25rem;
    }
    .badge-count {
      display: inline-block;
      background: var(--border);
      color: var(--vlight);
      border-radius: 999px;
      padding: .1rem .6rem;
      font-size: .72rem;
      margin-top: .4rem;
    }

    /* ── Grid ── */
    .services {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 1rem;
      padding: 2rem 2.5rem;
    }
    .empty {
      grid-column: 1/-1;
      text-align: center;
      color: var(--muted);
      padding: 3rem 0;
      font-size: .9rem;
    }

    /* ── Cards ── */
    .card {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 10px;
      padding: 1.25rem 1.5rem;
      transition: border-color .25s, transform .15s, box-shadow .25s;
      cursor: default;
      text-decoration: none;
      color: inherit;
      display: block;
    }
    .card:hover { transform: translateY(-2px); box-shadow: 0 4px 24px #7c3aed22; }
    .card.up    { border-color: var(--green); }
    .card.down  { border-color: var(--gray); opacity: .72; }

    .card-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: .6rem;
    }
    .svc-name { font-size: 1.05rem; color: var(--vlight); font-weight: 600; }
    .dot {
      width: 10px; height: 10px;
      border-radius: 50%;
      background: var(--gray);
      flex-shrink: 0;
      transition: background .3s, box-shadow .3s;
    }
    .dot.up { background: var(--green); box-shadow: 0 0 8px var(--green); }

    .svc-url {
      color: var(--vdim);
      font-size: .82rem;
      word-break: break-all;
      text-decoration: none;
    }
    .svc-url:hover { color: var(--vlight); }
    .svc-port    { color: var(--muted); font-size: .78rem; margin-top: .4rem; }
    .svc-runtime {
      color: var(--vdim);
      font-size: .7rem;
      text-transform: uppercase;
      letter-spacing: .08em;
      margin-top: .5rem;
      font-weight: 600;
    }

    /* ── Footer ── */
    footer {
      text-align: center;
      padding: 1.5rem;
      font-size: .72rem;
      color: var(--muted);
      border-top: 1px solid var(--surface);
    }
    footer a { color: var(--vdim); text-decoration: none; }
    footer a:hover { color: var(--vlight); }
  </style>
</head>
<body>

<header>
  <div class="logo-block">
    <div class="logo">ODINS</div>
    <div class="domain-fqdn">{{.FQDN}}</div>
    {{if .Description}}<div class="domain-desc">{{.Description}}</div>{{end}}
  </div>
  <div class="header-meta">
    The All-Father of Local DNS<br>
    <span class="badge-count" id="count">{{len .Routes}} service{{if ne (len .Routes) 1}}s{{end}}</span>
  </div>
</header>

<div class="services">
  {{if not .Routes}}
  <div class="empty">
    Nenhum serviço ainda.<br>
    Configure <code>domain = "{{.Domain}}"</code> no <code>.odins</code><br>
    e depois rode <code>odins up</code>
  </div>
  {{else}}
  {{range .Routes}}
  <a class="card" id="card-{{.Subdomain}}" href="https://{{.FQDN}}" target="_blank">
    <div class="card-top">
      <span class="svc-name">{{.Subdomain}}</span>
      <div class="dot" id="dot-{{.Subdomain}}"></div>
    </div>
    <div class="svc-url">{{.FQDN}}</div>
    <div class="svc-port">:{{.Port}}</div>
    {{if .Runtime}}<div class="svc-runtime">{{.Runtime}}</div>{{end}}
  </a>
  {{end}}
  {{end}}
</div>

<footer>
  Gerado por <a href="https://github.com/adialaleal/odins">ODINS</a>
  &nbsp;·&nbsp; {{.GeneratedAt}}
</footer>

<script>
const services = [{{range .Routes}}
  { subdomain: "{{.Subdomain}}", url: "https://{{.FQDN}}" },{{end}}
];

let upCount = 0;

async function check(svc) {
  try {
    const ctrl = new AbortController();
    const tid  = setTimeout(() => ctrl.abort(), 2500);
    await fetch(svc.url, { mode: 'no-cors', signal: ctrl.signal });
    clearTimeout(tid);
    return true;
  } catch { return false; }
}

async function poll() {
  upCount = 0;
  for (const s of services) {
    const up   = await check(s);
    const dot  = document.getElementById('dot-'  + s.subdomain);
    const card = document.getElementById('card-' + s.subdomain);
    if (up) upCount++;
    if (dot)  dot.className  = 'dot'  + (up ? ' up' : '');
    if (card) card.className = 'card' + (up ? ' up' : ' down');
  }
  const el = document.getElementById('count');
  if (el) el.textContent = upCount + '/' + services.length + ' online';
}

poll();
setInterval(poll, 5000);
</script>
</body>
</html>`
