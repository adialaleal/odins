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
