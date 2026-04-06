# AI Packs

Este diretório contém adapters gerados a partir da fonte canônica em `docs/ai/`.

- `codex/skills/`: skills PT-BR para uso com Codex.
- `antigravity/PROJECT_RULES.md`: rules pack portátil para Antigravity.
- `../CLAUDE.md`: adapter de projeto para Claude Code.

Para regenerar:

```bash
go run ./tools/ai-pack-gen
```

Para validar sincronismo:

```bash
go run ./tools/ai-pack-gen --check
```
