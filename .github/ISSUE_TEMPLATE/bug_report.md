---
name: Bug Report
about: Something isn't working as expected
title: "[bug] "
labels: bug
assignees: adialaleal
---

## Describe the bug

<!-- A clear and concise description of what the bug is. -->

## Steps to reproduce

1. Go to '...'
2. Run '...'
3. See error

## Expected behavior

<!-- What you expected to happen. -->

## Actual behavior

<!-- What actually happened. Include error output. -->

## Environment

- **macOS version:** <!-- e.g. Sonoma 14.4 -->
- **ODINS version:** <!-- run: odins --version -->
- **Proxy backend:** <!-- caddy / nginx / apache -->
- **TLD:** <!-- e.g. .odins -->

## Diagnostic output

```bash
# Please run these and paste the output:
odins --version
brew services info dnsmasq
brew services info caddy
curl -s http://localhost:2019/config/ | python3 -m json.tool
scutil --dns | grep -A3 odins
```

<details>
<summary>Output</summary>

```
paste here
```
</details>

## Additional context

<!-- Screenshots, logs, anything else that might help. -->
