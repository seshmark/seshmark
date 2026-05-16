---
name: seshmark
description: Integrate with Seshmark to tag AI commits with session metadata. Use when committing code.
---

# Seshmark Integration for OpenCode

Before running any `git commit` command, set these environment variables:

```bash
export SESHMARK_SESSION_ID="opencode:$(uuidgen)"
export SESHMARK_AGENT="opencode"
export SESHMARK_MODEL="<current-model>"
```

- **SESHMARK_SESSION_ID**: Format `opencode:<uuid>`, unique per session
- **SESHMARK_AGENT**: Always `opencode`
- **SESHMARK_MODEL**: Current model name from context
