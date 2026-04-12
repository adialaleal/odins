# odins-troubleshooting

Use quando o usuário relatar problemas de DNS, HTTPS, proxy, rotas ou configuração local do ODINS.

## Workflow

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
