# Contributing to ODINS

First off — thank you for taking the time to contribute! 🎉

ODINS is a small open-source project built for macOS developers. Every bug report, suggestion, translation, and pull request helps make it better.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Adding a Language](#adding-a-language)
- [Commit Convention](#commit-convention)

---

## Code of Conduct

Be respectful. Be constructive. We're all here to build something useful together.

---

## How Can I Contribute?

### 🐛 Report a Bug

Open an [issue](https://github.com/adialaleal/odins/issues/new?template=bug_report.md) and include:

- macOS version
- ODINS version (`odins --version`)
- Steps to reproduce
- What you expected vs what happened
- Relevant output from `odins ls`, `brew services info dnsmasq`, or `curl http://localhost:2019/config/`

### 💡 Suggest a Feature

Open an [issue](https://github.com/adialaleal/odins/issues/new?template=feature_request.md) describing:

- The problem you're trying to solve
- Your proposed solution
- Any alternatives you considered

### 🌐 Add or Improve a Translation

ODINS supports PT, EN, and ES. Translation strings live in:

```
internal/i18n/
├── pt.go   ← Portuguese (reference)
├── en.go   ← English
└── es.go   ← Spanish
```

To add a new language, see [Adding a Language](#adding-a-language).

### 🔧 Fix a Bug or Implement a Feature

Check the [open issues](https://github.com/adialaleal/odins/issues) — especially ones labeled `good first issue` or `help wanted`.

---

## Development Setup

**Requirements:**
- macOS (ODINS is macOS-only)
- Go 1.22+
- Homebrew

```bash
# 1. Fork and clone
git clone https://github.com/<your-username>/odins
cd odins

# 2. Install dependencies
go mod tidy

# 3. Build
go build -o odins .

# 4. Run locally (replaces installed binary temporarily)
./odins --help

# 5. Run tests
go test -v -race ./...

# 6. Lint
go vet ./...
```

---

## Project Structure

```
odins/
├── cmd/                    # Cobra CLI commands
│   ├── root.go             # Entry point, Caddy sync, welcome trigger
│   ├── init.go             # odins init — one-time setup
│   ├── up.go               # odins up — apply routes from .odins
│   ├── domain.go           # odins domain add/rm
│   ├── add.go              # odins add <sub> --port <n>
│   ├── kill.go             # odins kill <sub>
│   ├── down.go             # odins down
│   ├── ls.go               # odins ls
│   └── welcome.go          # Onboarding guide
│
├── internal/
│   ├── config/             # Global + project config (TOML)
│   ├── detect/             # Project auto-detection (Node/Go/Python)
│   ├── dns/                # dnsmasq config generation
│   ├── helper/             # Privileged ops via macOS auth dialog
│   ├── i18n/               # Multi-language strings (PT/EN/ES)
│   ├── page/               # HTML landing page generator
│   ├── proxy/
│   │   ├── caddy/          # Caddy Admin API client
│   │   └── nginx/          # Nginx config manager
│   ├── state/              # Route + domain persistence (JSON)
│   └── tui/                # Bubble Tea TUI
│       ├── app.go          # Main TUI model
│       ├── components/     # Reusable UI components
│       ├── screens/        # Dashboard, AddRoute, Settings, Logs
│       └── styles/         # Lipgloss style definitions
│
├── pkg/
│   └── brew/               # Homebrew helper (install, services)
│
├── main.go                 # Binary entry point + language detection
├── .goreleaser.yaml        # Release pipeline (GoReleaser)
└── install.sh              # curl | bash installer
```

---

## Submitting a Pull Request

1. **Fork** the repository and create a branch from `main`:
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make your changes.** Follow the existing code style.

3. **Test** your changes:
   ```bash
   go test -v -race ./...
   go vet ./...
   go build ./...
   ```

4. **Commit** using [Conventional Commits](#commit-convention).

5. **Push** and open a Pull Request against `main`.

6. Fill in the PR template — describe what changed and why.

**PR checklist:**
- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] New user-facing strings added to all three i18n catalogs (`pt.go`, `en.go`, `es.go`)
- [ ] CHANGELOG.md updated under `[Unreleased]`

---

## Adding a Language

1. Copy `internal/i18n/en.go` to `internal/i18n/<lang>.go`
2. Translate all string values (keep the keys identical)
3. Register it in `internal/i18n/i18n.go`:
   ```go
   const XX Lang = "xx"

   var catalogs = map[Lang]map[string]string{
       PT: ptStrings,
       EN: enStrings,
       ES: esStrings,
       XX: xxStrings,  // ← add here
   }
   ```
4. Add detection in `main.go` `initLang()`:
   ```go
   if strings.HasPrefix(v, "xx") { return i18n.XX }
   ```
5. Open a PR — translations are always welcome!

---

## Commit Convention

ODINS uses [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add Docker Compose auto-detection
fix: dnsmasq not binding on port 5300
docs: update README troubleshooting section
refactor: extract FQDN builder to shared helper
test: add unit tests for detect package
i18n: add French translation
```

Types: `feat` · `fix` · `docs` · `refactor` · `test` · `chore` · `i18n` · `style`

---

## Questions?

Open a [Discussion](https://github.com/adialaleal/odins/discussions) or reach out via [issues](https://github.com/adialaleal/odins/issues).

Thank you for contributing to ODINS! ⚡
