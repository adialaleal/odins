# odins-setup

Use quando o usuário quiser instalar, inicializar ou validar o ambiente local do ODINS no macOS.

## Workflow

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
