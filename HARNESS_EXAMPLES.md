# Seshmark Harness Integration Examples

> How to configure every major AI coding tool to set Seshmark environment variables for full metadata capture.

Every example is **copy-paste ready**. Add the 3-line env var block before your harness calls `git commit`.

---

## Quick Start

### For Agent-Based Tools (Pi, Claude Code, OpenCode)

Copy the relevant skill file to the tool's skill directory — the agent will then set env vars automatically before every commit:

```bash
# Pi
cp examples/harnesses/pi/SKILL.md ~/.cursor/skills/seshmark/SKILL.md

# Claude Code
cp examples/harnesses/claude-code/SKILL.md ~/.claude/skills/seshmark/SKILL.md

# OpenCode
cp examples/harnesses/opencode/SKILL.md ~/.opencode/skills/seshmark/SKILL.md
```

### For Wrapper-Based Tools (Aider, any CLI)

Install a wrapper script that intercepts the tool binary:

```bash
# Aider
cp examples/harnesses/aider/wrapper.sh ~/bin/aider
chmod +x ~/bin/aider

# Generic (any tool)
cp examples/harnesses/generic/wrapper.sh ~/bin/your-tool
# Edit: change TOOL_NAME to your tool's name
```

### For IDE-Based Tools (Cursor, Copilot)

Use `.cursorrules` or shell aliases (see tool-specific sections below).

---

## Verification

After setting up, make a commit and check:

```bash
seshmark who HEAD
# Should show:
#   AI-Session: <tool>:<uuid>
#   AI-Agent: <tool>
#   AI-Model: <model>
```

---

## Cursor

Cursor saves conversations in `~/.cursor/`. We reference the active conversation ID.

### Cursor IDE (Manual)

Cursor doesn't expose a session ID natively, but you can use the project directory name or a timestamp:

```bash
# In your terminal before using Cursor
export SESHMARK_SESSION_ID="cursor:$(basename $(pwd))-$(date +%Y%m%d)"
export SESHMARK_AGENT="cursor"
export SESHMARK_MODEL="claude-sonnet-4-20250514"

# Now open Cursor and work normally
cursor .
```

### Cursor via CLI (`cursor` command)

```bash
#!/bin/bash
# ~/.local/bin/cursor-sesh

PROJECT="$(basename $(pwd))"
DATE="$(date +%Y%m%d)"

export SESHMARK_SESSION_ID="cursor:${PROJECT}-${DATE}"
export SESHMARK_AGENT="cursor"
export SESHMARK_MODEL="${CURSOR_MODEL:-claude-sonnet-4-20250514}"

cursor "$@"
```

### Cursor with Explicit Chat ID

```python
# ~/.cursor/scripts/seshmark_hook.py
# Run this after each Cursor chat session

import os
import subprocess
import json
from datetime import datetime

def get_latest_chat_id():
    cursor_dir = os.path.expanduser("~/.cursor")
    chats = []
    for root, dirs, files in os.walk(cursor_dir):
        for f in files:
            if f.endswith('.json'):
                path = os.path.join(root, f)
                chats.append((path, os.path.getmtime(path)))
    if not chats:
        return None
    chats.sort(key=lambda x: x[1], reverse=True)
    return os.path.basename(chats[0][0]).replace('.json', '')

def tag_cursor_commit():
    chat_id = get_latest_chat_id()
    if not chat_id:
        chat_id = f"manual-{datetime.now().strftime('%Y%m%d-%H%M%S')}"

    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"cursor:{chat_id}"
    env["SESHMARK_AGENT"] = "cursor"
    env["SESHMARK_MODEL"] = os.environ.get("CURSOR_MODEL", "claude-sonnet-4-20250514")

    subprocess.run(["git", "commit", "-m", "feat: cursor changes"], env=env)

if __name__ == "__main__":
    tag_cursor_commit()
```

---

## Claude Code (Anthropic)

Claude Code stores threads in `~/.claude/`.

### Skill File (Recommended)

```bash
cp examples/harnesses/claude-code/SKILL.md ~/.claude/skills/seshmark/SKILL.md
```

### Wrapper Script

```bash
#!/bin/bash
# ~/.local/bin/claude-sesh

SESSION_ID="$(date +%Y%m%d)-$(openssl rand -hex 4)"

export SESHMARK_SESSION_ID="claude:${SESSION_ID}"
export SESHMARK_AGENT="claude"
export SESHMARK_MODEL="${CLAUDE_MODEL:-claude-sonnet-4-20250514}"

claude "$@"
```

### Python Integration (If Building on Top of Claude Code)

```python
import os
import subprocess
import uuid

def claude_commit(message: str, model: str = "claude-sonnet-4-20250514"):
    session_id = f"claude:thread-{uuid.uuid4().hex[:8]}"

    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = session_id
    env["SESHMARK_AGENT"] = "claude"
    env["SESHMARK_MODEL"] = model

    subprocess.run(["git", "commit", "-m", message], env=env, check=True)
    print(f"Tagged commit with session: {session_id}")

# Usage
claude_commit("feat: implement auth middleware", model="claude-sonnet-4-20250514")
```

---

## Codex CLI (OpenAI)

Codex is OpenAI's CLI tool. It exposes a run ID.

### Codex Native Integration

```bash
#!/bin/bash
# ~/.local/bin/codex-sesh

RUN_ID="run-$(date +%Y%m%d-%H%M%S)-$(openssl rand -hex 4)"

export SESHMARK_SESSION_ID="codex:${RUN_ID}"
export SESHMARK_AGENT="codex"
export SESHMARK_MODEL="${CODEX_MODEL:-gpt-4o}"

codex "$@"
```

### Post-Commit Hook for Codex

```bash
#!/bin/bash
CODEX_RUN_ID="${CODEX_RUN_ID:-$(date +%Y%m%d-%H%M%S)}"

export SESHMARK_SESSION_ID="codex:${CODEX_RUN_ID}"
export SESHMARK_AGENT="codex"
export SESHMARK_MODEL="gpt-4o"

git add -A
git commit -m "feat: codex changes"
```

---

## GitHub Copilot

Copilot doesn't have explicit sessions, but we can track by suggestion batch or file.

### Copilot in VS Code

```bash
# In your shell before opening VS Code
export SESHMARK_SESSION_ID="github-copilot:vscode-$(date +%Y%m%d)"
export SESHMARK_AGENT="github-copilot"
export SESHMARK_MODEL="copilot-chat"

code .
```

### Copilot CLI (`gh copilot`)

```bash
#!/bin/bash
# ~/.local/bin/gh-copilot-sesh

SESSION_ID="copilot-$(date +%Y%m%d-%H%M%S)"

export SESHMARK_SESSION_ID="github-copilot:${SESSION_ID}"
export SESHMARK_AGENT="github-copilot"
export SESHMARK_MODEL="copilot-chat"

gh copilot "$@"
```

### GitHub Copilot in JetBrains

```bash
export SESHMARK_SESSION_ID="github-copilot:jetbrains-$(date +%Y%m%d)"
export SESHMARK_AGENT="github-copilot"
export SESHMARK_MODEL="copilot-chat"

# Open IDE (WebStorm, IntelliJ, etc.)
webstorm .
```

---

## Aider

Aider has built-in session awareness. It generates session IDs automatically.

### Wrapper Script (Recommended)

```bash
cp examples/harnesses/aider/wrapper.sh ~/bin/aider
chmod +x ~/bin/aider
```

### Aider Config File (`~/.aider.conf.yml`)

```yaml
custom_commands:
  pre-commit: |
    export SESHMARK_SESSION_ID="aider:${AIDER_SESSION_ID}"
    export SESHMARK_AGENT="aider"
    export SESHMARK_MODEL="${AIDER_MODEL:-gpt-4o}"
```

---

## OpenCode

OpenCode is designed for extensibility.

### Skill File (Recommended)

```bash
cp examples/harnesses/opencode/SKILL.md ~/.opencode/skills/seshmark/SKILL.md
```

### OpenCode Native Plugin

```python
# opencode_plugin_seshmark.py
import os
import uuid

class SeshmarkPlugin:
    def __init__(self):
        self.session_id = None

    def on_session_start(self):
        self.session_id = f"opencode:sess-{uuid.uuid4().hex[:8]}"
        os.environ["SESHMARK_SESSION_ID"] = self.session_id
        os.environ["SESHMARK_AGENT"] = "opencode"
        os.environ["SESHMARK_MODEL"] = os.environ.get("OPENCODE_MODEL", "qwen2.5-coder")

    def pre_commit(self):
        return {
            "SESHMARK_SESSION_ID": self.session_id,
            "SESHMARK_AGENT": "opencode",
            "SESHMARK_MODEL": os.environ.get("OPENCODE_MODEL", "qwen2.5-coder"),
        }
```

---

## Pi

Pi uses a skill system — the agent reads behavioral instructions from skill files.

### Skill File (Recommended)

```bash
cp examples/harnesses/pi/SKILL.md ~/.cursor/skills/seshmark/SKILL.md
```

### Wrapper Script

```bash
#!/bin/bash
# ~/.local/bin/pi-sesh

SESSION_ID="pi:$(uuidgen 2>/dev/null || date +%s)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="pi"
export SESHMARK_MODEL="${PI_MODEL:-deepseek-v4}"

pi "$@"
```

---

## Continue.dev

Continue is a VS Code extension with explicit session tracking.

### Continue Config (`~/.continue/config.json`)

```json
{
  "custom_commands": [
    {
      "name": "seshmark-commit",
      "description": "Commit with seshmark metadata",
      "prompt": "Run a shell command that exports SESHMARK_SESSION_ID, SESHMARK_AGENT, SESHMARK_MODEL, then commits."
    }
  ]
}
```

---

## Devin (Cognition AI)

Devin has a run-based model. Each task is a "run."

```python
import os, subprocess

def devin_commit(run_id: str, message: str):
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"devin:{run_id}"
    env["SESHMARK_AGENT"] = "devin"
    env["SESHMARK_MODEL"] = "devin-v1"
    subprocess.run(["git", "commit", "-m", message], env=env, check=True)
```

---

## SWE-Agent

SWE-Agent has explicit trajectory IDs.

```python
import os, subprocess

def swe_agent_commit(trajectory_id: str, message: str):
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"swe-agent:{trajectory_id}"
    env["SESHMARK_AGENT"] = "swe-agent"
    env["SESHMARK_MODEL"] = os.environ.get("SWE_AGENT_MODEL", "gpt-4o")
    subprocess.run(["git", "commit", "-m", message], env=env, check=True)
```

---

## Shell Alias Collection

Add these to `~/.bashrc` or `~/.zshrc`:

```bash
# Cursor
alias cursor-sesh='SESHMARK_SESSION_ID="cursor:$(basename $(pwd))-$(date +%Y%m%d)" SESHMARK_AGENT="cursor" cursor'

# Claude Code
alias claude-sesh='SESHMARK_SESSION_ID="claude:$(date +%Y%m%d)" SESHMARK_AGENT="claude" claude'

# Codex
alias codex-sesh='SESHMARK_SESSION_ID="codex:run-$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="codex" codex'

# Aider
alias aider-sesh='SESHMARK_SESSION_ID="aider:$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="aider" aider'

# OpenCode
alias opencode-sesh='SESHMARK_SESSION_ID="opencode:$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="opencode" opencode'

# Pi
alias pi-sesh='SESHMARK_SESSION_ID="pi:$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="pi" pi'

# GitHub Copilot
alias code-copilot='SESHMARK_SESSION_ID="github-copilot:vscode-$(date +%Y%m%d)" SESHMARK_AGENT="github-copilot" code'
```

---

## Quick Reference

| Agent | Session ID Format | Skill File | Wrapper |
|-------|------------------|-----------|---------|
| Pi | `pi:<uuid>` | `examples/harnesses/pi/SKILL.md` | `pi-sesh` |
| Claude Code | `claude:<uuid>` | `examples/harnesses/claude-code/SKILL.md` | `claude-sesh` |
| OpenCode | `opencode:<uuid>` | `examples/harnesses/opencode/SKILL.md` | `opencode-sesh` |
| Cursor | `cursor:<project>-<date>` | `.cursorrules` | `cursor-sesh` |
| Aider | `aider:<session-id>` | — | `aider/wrapper.sh` |
| Codex | `codex:run-<timestamp>` | — | `codex-sesh` |
| GitHub Copilot | `github-copilot:<ide>-<date>` | — | `code-copilot` |
| Continue.dev | `continue:<timestamp>` | — | shell alias |
| Devin | `devin:<run-id>` | — | Python snippet |
| SWE-Agent | `swe-agent:<trajectory-id>` | — | Python snippet |
| Generic CLI | `<tool>:<uuid>` | — | `generic/wrapper.sh` |

---

## For Harness Builders

**The universal contract:**

```python
env["SESHMARK_SESSION_ID"] = f"{agent}:{session_id}"
env["SESHMARK_AGENT"] = agent
env["SESHMARK_MODEL"] = model
subprocess.run(["git", "commit", "-m", message], env=env)
```

That's all. No SDK. No API key. No dependency on seshmark being installed.

If seshmark is installed: commits are auto-tagged.
If seshmark is not installed: commit still works, trailers are just text in the message.
