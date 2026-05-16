---
name: seshmark
description: Integrate with Seshmark to tag AI commits with session metadata. Use when committing code to ensure AI attribution is captured.
---

# Seshmark Integration for Pi

This skill tells Pi to automatically tag Git commits with AI metadata when Seshmark is installed.

## Behavior

Before running any `git commit` command, set these environment variables:

```bash
export SESHMARK_SESSION_ID="pi:$(uuidgen)"
export SESHMARK_AGENT="pi"
export SESHMARK_MODEL="<current-model>"
```

- **SESHMARK_SESSION_ID**: Generate a unique UUID per session. Format: `pi:<uuid>`
- **SESHMARK_AGENT**: Always set to `pi`
- **SESHMARK_MODEL**: Detect the current model from context (e.g., `deepseek-v4-flash`, `claude-sonnet-4-20250514`)

Run `git commit` with these variables set — either inline:

```bash
SESHMARK_SESSION_ID="pi:$(uuidgen)" SESHMARK_AGENT="pi" SESHMARK_MODEL="deepseek-v4" git commit -m "message"
```

Or exported before commit:

```bash
export SESHMARK_SESSION_ID="pi:$(uuidgen)"
export SESHMARK_AGENT="pi"
export SESHMARK_MODEL="deepseek-v4"
git commit -m "message"
```

This allows Seshmark to capture full metadata (agent + session + model) in every commit.
