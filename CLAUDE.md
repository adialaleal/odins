# CLAUDE.md

Este projeto expõe uma camada AI Friendly para Claude Code via CLI estruturado, docs canônicas e regras operacionais curtas.

## Regra principal

- Inspecione primeiro, mude depois.
- Prefira `odins detect --json` antes de propor ou escrever `.odins`.
- Prefira `odins doctor --json` antes de troubleshooting manual.
- Explique quando um comando pode pedir `sudo`.
- Em automação, prefira `--json` e `--non-interactive`.

## Documentos base

- [docs/ai/setup-local.md](docs/ai/setup-local.md)
- [docs/ai/apply-to-project.md](docs/ai/apply-to-project.md)
- [docs/ai/workspace-multi-service.md](docs/ai/workspace-multi-service.md)
- [docs/ai/doctor-troubleshoot.md](docs/ai/doctor-troubleshoot.md)

## Resumo operacional

### Setup

# Setup Local

Use este guia quando o objetivo for instalar e preparar o ODINS na máquina do usuário.

## Objetivo

Deixar o macOS apto a resolver domínios locais via ODINS com DNS wildcard, proxy reverso e HTTPS local.

## Checklist rápido

1. Confirmar que a máquina é macOS.
2. Confirmar que o Homebrew está instalado.
3. Explicar que `odins init` pode pedir `sudo` para configurar `/etc/resolver` e confiar no certificado local.
4. Preferir `caddy` como backend padrão, a menos que o usuário peça `nginx` ou `apache`.
5. Preferir `.odin` como TLD padrão, a menos que o usuário peça outro TLD suportado.

## Comandos recomendados

```bash
odins init
odins init --non-interactive --tld odin --backend caddy
odins init --json --non-interactive --tld odin --backend caddy
odins doctor
odins doctor --json
```

## Regras para agentes

- Explorar primeiro o ambiente e só depois propor mudanças.
- Explicar claramente quando um comando vai pedir `sudo`.
- Em automação, preferir `--non-interactive`.
- Em automação ou integrações, preferir `--json`.
- Se o usuário estiver em contexto não interativo, não assumir que o `sudo` vai funcionar sem confirmação.

## Saída esperada

- Configuração global salva em `~/.config/odins/config.toml`.
- `dnsmasq` instalado e reiniciado.
- Backend escolhido instalado.
- Resolver local criado em `/etc/resolver/<tld>`.
- Certificado local preparado para HTTPS.

## Troubleshooting curto

- Se `odins init` falhar, rodar `odins doctor`.
- Se o problema for Homebrew ausente, orientar instalação do Homebrew antes de continuar.
- Se o problema for certificado ou DNS, verificar as ações sugeridas por `odins doctor`.

### Aplicação em projeto

# Apply To Project

Use este guia quando o objetivo for aplicar o ODINS a um projeto local.

## Objetivo

Detectar a stack do projeto, propor uma configuração `.odins`, aplicar as rotas e validar o acesso via domínio local.

## Fluxo recomendado

1. Inspecionar o repositório local.
2. Rodar `odins detect` antes de mudar qualquer arquivo.
3. Revisar o `.odins` recomendado.
4. Rodar `odins up` para aplicar as rotas.
5. Validar com `odins ls` e `odins doctor` se necessário.

## Comandos recomendados

```bash
odins detect --json
odins up
odins up --json
odins ls
odins ls --json
```

## O que o agente deve observar

- `runtime`, `framework`, `port` e `start_cmd` detectados.
- Se já existe um `.odins` no projeto.
- Se o projeto parece usar Docker ou Docker Compose.
- Se a porta detectada faz sentido com o framework.

## Regras para agentes

- Preferir `odins detect` antes de escrever ou sobrescrever `.odins`.
- Só perguntar ao usuário o mínimo necessário: domínio/workspace, subdomínios desejados e confirmação para operações sensíveis.
- Se a detecção falhar, explicar o motivo e propor um `.odins` manual.
- Ao aplicar rotas, preferir `odins up` em vez de várias chamadas `odins add`, salvo quando o usuário quer uma rota isolada.

## Resultado esperado

- `.odins` presente ou revisado.
- Rotas aplicadas.
- Projeto acessível por um FQDN local como `https://app.projeto.odin`.

### Workspaces

# Workspace Multi Service

Use este guia quando o projeto faz parte de um workspace com vários serviços.

## Objetivo

Organizar múltiplos serviços sob um domínio local de workspace, com landing page e subdomínios previsíveis.

## Fluxo recomendado

1. Criar ou validar um domínio de workspace.
2. Garantir que cada projeto tenha `domain = "<workspace>"` no `.odins`.
3. Aplicar as rotas de cada projeto com `odins up`.
4. Validar a landing page do workspace e os subdomínios individuais.

## Comandos recomendados

```bash
odins domain add tatoh
odins domain add tatoh --json
odins domain ls
odins domain ls --json
odins domain rm tatoh
odins domain rm tatoh --json
```

## Convenções sugeridas

- Usar `web`, `api`, `admin`, `worker` ou nomes curtos equivalentes como `subdomain`.
- Usar um workspace por contexto de produto, time ou suite de serviços.
- Evitar FQDNs longos demais quando um workspace já agrupa os projetos.

## Regras para agentes

- Quando houver mais de um serviço, sugerir `odins domain add` antes de sair criando FQDNs soltos.
- Revisar a configuração `domain` do `.odins` antes de rodar `odins up`.
- Avisar que a landing page do domínio é melhor suportada com backend `caddy`.

### Troubleshooting

# Doctor Troubleshoot

Use este guia quando o usuário relatar falhas de DNS, HTTPS, proxy, rotas ou comportamento inesperado do ODINS.

## Objetivo

Diagnosticar rapidamente o ambiente e indicar a próxima ação mais segura.

## Fluxo recomendado

1. Rodar `odins doctor --json`.
2. Ler os checks e as `action` sugeridas.
3. Só partir para troubleshooting manual se o `doctor` não for suficiente.
4. Se necessário, confirmar `odins ls --json` para revisar rotas ativas.

## Comandos recomendados

```bash
odins doctor
odins doctor --json
odins ls --json
```

## Problemas comuns

- Homebrew ausente.
- `dnsmasq` parado.
- Resolver local ausente em `/etc/resolver/<tld>`.
- Backend configurado não instalado ou não iniciado.
- Certificado local não preparado.
- Nenhuma rota ativa no store.

## Regras para agentes

- Sempre começar por `odins doctor`.
- Não reinventar verificações que o `doctor` já cobre.
- Se o usuário pedir automação, preferir o JSON do `doctor`.
- Em falhas ligadas a `sudo`, explicar a ação exata que exigiu privilégio.
