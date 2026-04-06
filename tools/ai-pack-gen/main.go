package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	check := flag.Bool("check", false, "Verify generated AI packs are up to date")
	flag.Parse()

	root, err := os.Getwd()
	if err != nil {
		fail(err)
	}

	docs, err := loadDocs(root)
	if err != nil {
		fail(err)
	}

	files := generatedFiles(docs)
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var mismatches []string
	for _, relativePath := range paths {
		targetPath := filepath.Join(root, relativePath)
		content := files[relativePath]
		existing, err := os.ReadFile(targetPath)
		if err == nil && bytes.Equal(existing, []byte(content)) {
			continue
		}
		if *check {
			mismatches = append(mismatches, relativePath)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			fail(err)
		}
		if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
			fail(err)
		}
		fmt.Println("updated", relativePath)
	}

	if *check && len(mismatches) > 0 {
		fail(fmt.Errorf("AI packs desatualizados: %s", strings.Join(mismatches, ", ")))
	}
}

func loadDocs(root string) (map[string]string, error) {
	files := []string{
		"setup-local.md",
		"apply-to-project.md",
		"workspace-multi-service.md",
		"doctor-troubleshoot.md",
	}

	out := make(map[string]string, len(files))
	for _, name := range files {
		path := filepath.Join(root, "docs", "ai", name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		out[name] = strings.TrimSpace(string(data))
	}
	return out, nil
}

func generatedFiles(docs map[string]string) map[string]string {
	setup := docs["setup-local.md"]
	apply := docs["apply-to-project.md"]
	workspace := docs["workspace-multi-service.md"]
	doctor := docs["doctor-troubleshoot.md"]

	return map[string]string{
		"AGENTS.md":                       renderAgents(setup, apply, workspace, doctor),
		"CLAUDE.md":                       renderClaude(setup, apply, workspace, doctor),
		"ai/README.md":                    renderAIReadme(),
		"ai/antigravity/PROJECT_RULES.md": renderAntigravity(setup, apply, workspace, doctor),
		"ai/codex/skills/odins-setup/SKILL.md": renderCodexSkill(
			"odins-setup",
			"Use quando o usuário quiser instalar, inicializar ou validar o ambiente local do ODINS no macOS.",
			setup,
		),
		"ai/codex/skills/odins-apply/SKILL.md": renderCodexSkill(
			"odins-apply",
			"Use quando o usuário quiser aplicar o ODINS a um projeto, detectar a stack, gerar `.odins` ou organizar múltiplos serviços em um workspace.",
			apply+"\n\n---\n\n"+workspace,
		),
		"ai/codex/skills/odins-troubleshooting/SKILL.md": renderCodexSkill(
			"odins-troubleshooting",
			"Use quando o usuário relatar problemas de DNS, HTTPS, proxy, rotas ou configuração local do ODINS.",
			doctor,
		),
	}
}

func renderAgents(setup, apply, workspace, doctor string) string {
	return joinLines(
		"# AGENTS",
		"",
		"Guia neutro para agentes e desenvolvedores que precisam operar o ODINS neste repositório ou em projetos que usam o ODINS.",
		"",
		"## Visão rápida",
		"",
		"- ODINS é um gerenciador local de DNS + reverse proxy para macOS.",
		"- Use `odins detect` para inspecionar um projeto antes de mudar qualquer arquivo.",
		"- Use `odins doctor` para diagnosticar o ambiente antes de troubleshooting manual.",
		"- Use `--json` sempre que a saída precisar ser consumida por automação ou agentes.",
		"",
		"## Comandos seguros de leitura",
		"",
		"- `odins detect --json --dir <path>`",
		"- `odins doctor --json`",
		"- `odins ls --json`",
		"- `odins domain ls --json`",
		"",
		"## Comandos com side effects",
		"",
		"- `odins init`",
		"- `odins up`",
		"- `odins add <subdomain> --port <port>`",
		"- `odins kill <subdomain>`",
		"- `odins down`",
		"- `odins domain add <name>`",
		"- `odins domain rm <name>`",
		"",
		"## Cuidados obrigatórios",
		"",
		"- `odins init` pode pedir `sudo` para criar `/etc/resolver/<tld>` e confiar no certificado local.",
		"- O suporte oficial desta versão é macOS.",
		"- Prefira `caddy` como backend padrão.",
		"- Em automação, prefira `--non-interactive` e `--json`.",
		"",
		"## Fonte canônica",
		"",
		"- [docs/ai/setup-local.md](docs/ai/setup-local.md)",
		"- [docs/ai/apply-to-project.md](docs/ai/apply-to-project.md)",
		"- [docs/ai/workspace-multi-service.md](docs/ai/workspace-multi-service.md)",
		"- [docs/ai/doctor-troubleshoot.md](docs/ai/doctor-troubleshoot.md)",
		"",
		"## Resumos rápidos",
		"",
		setup,
		"",
		"---",
		"",
		apply,
		"",
		"---",
		"",
		workspace,
		"",
		"---",
		"",
		doctor,
	)
}

func renderClaude(setup, apply, workspace, doctor string) string {
	return joinLines(
		"# CLAUDE.md",
		"",
		"Este projeto expõe uma camada AI Friendly para Claude Code via CLI estruturado, docs canônicas e regras operacionais curtas.",
		"",
		"## Regra principal",
		"",
		"- Inspecione primeiro, mude depois.",
		"- Prefira `odins detect --json` antes de propor ou escrever `.odins`.",
		"- Prefira `odins doctor --json` antes de troubleshooting manual.",
		"- Explique quando um comando pode pedir `sudo`.",
		"- Em automação, prefira `--json` e `--non-interactive`.",
		"",
		"## Documentos base",
		"",
		"- [docs/ai/setup-local.md](docs/ai/setup-local.md)",
		"- [docs/ai/apply-to-project.md](docs/ai/apply-to-project.md)",
		"- [docs/ai/workspace-multi-service.md](docs/ai/workspace-multi-service.md)",
		"- [docs/ai/doctor-troubleshoot.md](docs/ai/doctor-troubleshoot.md)",
		"",
		"## Resumo operacional",
		"",
		"### Setup",
		"",
		setup,
		"",
		"### Aplicação em projeto",
		"",
		apply,
		"",
		"### Workspaces",
		"",
		workspace,
		"",
		"### Troubleshooting",
		"",
		doctor,
	)
}

func renderAIReadme() string {
	return joinLines(
		"# AI Packs",
		"",
		"Este diretório contém adapters gerados a partir da fonte canônica em `docs/ai/`.",
		"",
		"- `codex/skills/`: skills PT-BR para uso com Codex.",
		"- `antigravity/PROJECT_RULES.md`: rules pack portátil para Antigravity.",
		"- `../CLAUDE.md`: adapter de projeto para Claude Code.",
		"",
		"Para regenerar:",
		"",
		"```bash",
		"go run ./tools/ai-pack-gen",
		"```",
		"",
		"Para validar sincronismo:",
		"",
		"```bash",
		"go run ./tools/ai-pack-gen --check",
		"```",
	)
}

func renderAntigravity(setup, apply, workspace, doctor string) string {
	return joinLines(
		"# ODINS Project Rules For Antigravity",
		"",
		"Copie estas regras para o Project Settings ou para o arquivo de regras usado no projeto.",
		"",
		"## Core",
		"",
		"- Você está trabalhando com ODINS, um gerenciador local de DNS + reverse proxy para macOS.",
		"- Sempre explore o repositório local antes de propor mudanças.",
		"- Use `odins detect --json` antes de criar ou alterar `.odins`.",
		"- Use `odins doctor --json` antes de troubleshooting manual.",
		"- Em automação, prefira `--json` e `--non-interactive`.",
		"- Avise o usuário antes de qualquer ação que possa pedir `sudo`.",
		"",
		"## Safe Reads",
		"",
		"- `odins detect --json --dir <path>`",
		"- `odins doctor --json`",
		"- `odins ls --json`",
		"- `odins domain ls --json`",
		"",
		"## Side Effects",
		"",
		"- `odins init`",
		"- `odins up`",
		"- `odins add <subdomain> --port <port>`",
		"- `odins kill <subdomain>`",
		"- `odins down`",
		"- `odins domain add <name>`",
		"- `odins domain rm <name>`",
		"",
		"## Playbooks",
		"",
		setup,
		"",
		"---",
		"",
		apply,
		"",
		"---",
		"",
		workspace,
		"",
		"---",
		"",
		doctor,
	)
}

func renderCodexSkill(name, description, content string) string {
	return joinLines(
		"# "+name,
		"",
		description,
		"",
		"## Workflow",
		"",
		content,
	)
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
