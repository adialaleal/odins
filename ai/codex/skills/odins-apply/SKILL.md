# odins-apply

Use quando o usuário quiser aplicar o ODINS a um projeto, detectar a stack, gerar `.odins` ou organizar múltiplos serviços em um workspace.

## Workflow

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

---

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
