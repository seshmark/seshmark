---
name: seshmark
description: Integrate with Seshmark to tag AI commits with session metadata. Use when committing code.
---

# Seshmark Integration for Claude Code

Before running any `git commit` command, set these environment variables:

```bash
export SESHMARK_SESSION_ID="claude:$(uuidgen)"
export SESHMARK_AGENT="claude"
export SESHMARK_MODEL="<current-model>"
```

- **SESHMARK_SESSION_ID**: Format `claude:<uuid>`, unique per session
- **SESHMARK_AGENT**: Always `claude`
- **SESHMARK_MODEL**: Current model name from context

Set them inline or export before the commit command.
