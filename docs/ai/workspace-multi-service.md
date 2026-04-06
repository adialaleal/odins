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
