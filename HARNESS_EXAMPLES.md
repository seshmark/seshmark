# Seshmark Harness Integration Examples

Every example is **copy-paste ready**. Add the 3-line env var block before your harness calls `git commit`.

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
    # Find the most recent chat file
    chats = []
    for root, dirs, files in os.walk(cursor_dir):
        for f in files:
            if f.endswith('.json'):
                path = os.path.join(root, f)
                chats.append((path, os.path.getmtime(path)))
    if not chats:
        return None
    chats.sort(key=lambda x: x[1], reverse=True)
    # Extract chat ID from filename or path
    return os.path.basename(chats[0][0]).replace('.json', '')

def tag_cursor_commit():
    chat_id = get_latest_chat_id()
    if not chat_id:
        chat_id = f"manual-{datetime.now().strftime('%Y%m%d-%H%M%S')}"
    
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"cursor:{chat_id}"
    env["SESHMARK_AGENT"] = "cursor"
    env["SESHMARK_MODEL"] = os.environ.get("CURSOR_MODEL", "claude-sonnet-4-20250514")
    
    # Example: commit with the env vars set
    subprocess.run(["git", "commit", "-m", "feat: cursor changes"], env=env)

if __name__ == "__main__":
    tag_cursor_commit()
```

---

## Claude Code (Anthropic)

Claude Code stores threads in `~/.claude/`.

### Claude Code Built-in Integration

```bash
# Claude Code exposes the thread ID in env vars
# Add to your shell profile or Claude Code config

export SESHMARK_AGENT="claude-code"
export SESHMARK_MODEL="claude-sonnet-4-20250514"
```

### Wrapper Script

```bash
#!/bin/bash
# ~/.local/bin/claude-sesh

# Claude Code doesn't expose session ID directly,
# so we use a timestamp-based ID
SESSION_ID="$(date +%Y%m%d)-$(openssl rand -hex 4)"

export SESHMARK_SESSION_ID="claude:${SESSION_ID}"
export SESHMARK_AGENT="claude-code"
export SESHMARK_MODEL="${CLAUDE_MODEL:-claude-sonnet-4-20250514}"

claude "$@"
```

### Python Integration (If Building on Top of Claude Code)

```python
import os
import subprocess
import uuid
from datetime import datetime

def claude_commit(message: str, model: str = "claude-sonnet-4-20250514"):
    """Commit with Claude Code session metadata."""
    session_id = f"claude:thread-{uuid.uuid4().hex[:8]}"
    
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = session_id
    env["SESHMARK_AGENT"] = "claude-code"
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

# Codex doesn't persist sessions by default,
# so we generate one at invocation time
RUN_ID="run-$(date +%Y%m%d-%H%M%S)-$(openssl rand -hex 4)"

export SESHMARK_SESSION_ID="codex:${RUN_ID}"
export SESHMARK_AGENT="codex"
export SESHMARK_MODEL="${CODEX_MODEL:-gpt-4o}"

codex "$@"
```

### Post-Commit Hook for Codex

```bash
#!/bin/bash
# Run after codex makes changes

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

# Open VS Code
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
# Set in your shell before opening IDE
export SESHMARK_SESSION_ID="github-copilot:jetbrains-$(date +%Y%m%d)"
export SESHMARK_AGENT="github-copilot"
export SESHMARK_MODEL="copilot-chat"

# Open IDE (WebStorm, IntelliJ, etc.)
webstorm .
```

---

## Aider

Aider has built-in session awareness. It generates session IDs automatically.

### Aider Native Integration

```python
# In your Aider config or wrapper
import os

def get_aider_env(model: str = "gpt-4o"):
    """Get env vars for Aider commits."""
    session_id = os.environ.get("AIDER_SESSION_ID")
    if not session_id:
        import time
        session_id = f"aider-{int(time.time())}"
    
    return {
        "SESHMARK_SESSION_ID": f"aider:{session_id}",
        "SESHMARK_AGENT": "aider",
        "SESHMARK_MODEL": model,
    }

# Aider's commit callback
env = os.environ.copy()
env.update(get_aider_env(model="gpt-4o"))
```

### Aider Wrapper Script

```bash
#!/bin/bash
# ~/.local/bin/aider-sesh

SESSION_ID="${AIDER_SESSION_ID:-$(date +%Y%m%d-%H%M%S)}"
MODEL="${AIDER_MODEL:-gpt-4o}"

export SESHMARK_SESSION_ID="aider:${SESSION_ID}"
export SESHMARK_AGENT="aider"
export SESHMARK_MODEL="$MODEL"

aider "$@"
```

### Aider Config File (`~/.aider.conf.yml`)

```yaml
# Add to your Aider config
custom_commands:
  pre-commit: |
    export SESHMARK_SESSION_ID="aider:${AIDER_SESSION_ID}"
    export SESHMARK_AGENT="aider"
    export SESHMARK_MODEL="${AIDER_MODEL:-gpt-4o}"
```

---

## OpenCode

OpenCode is designed for extensibility.

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
    
    def pre_commit(self, message: str):
        # OpenCode calls this before git commit
        return {
            "SESHMARK_SESSION_ID": self.session_id,
            "SESHMARK_AGENT": "opencode",
            "SESHMARK_MODEL": os.environ.get("OPENCODE_MODEL", "qwen2.5-coder"),
        }

# Register with OpenCode
# In opencode config:
# plugins:
#   - path: opencode_plugin_seshmark.py
```

### OpenCode Wrapper

```bash
#!/bin/bash
# ~/.local/bin/opencode-sesh

SESSION_ID="opencode-$(date +%Y%m%d-%H%M%S)"

export SESHMARK_SESSION_ID="opencode:${SESSION_ID}"
export SESHMARK_AGENT="opencode"
export SESHMARK_MODEL="${OPENCODE_MODEL:-qwen2.5-coder}"

opencode "$@"
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
  ],
  "on_commit": {
    "env_vars": {
      "SESHMARK_AGENT": "continue",
      "SESHMARK_MODEL": "${config.model}"
    }
  }
}
```

### Continue Wrapper

```bash
#!/bin/bash
# ~/.local/bin/continue-sesh

SESSION_ID="continue-$(date +%Y%m%d-%H%M%S)"

export SESHMARK_SESSION_ID="continue:${SESSION_ID}"
export SESHMARK_AGENT="continue"
export SESHMARK_MODEL="${CONTINUE_MODEL:-claude-sonnet-4}"

# Continue runs as a VS Code extension, so we set env vars before opening VS Code
code "$@"
```

---

## Devin (Cognition AI)

Devin has a run-based model. Each task is a "run."

### Devin API Integration

```python
import os
import requests

def devin_commit(run_id: str, message: str):
    """Commit with Devin run metadata."""
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"devin:{run_id}"
    env["SESHMARK_AGENT"] = "devin"
    env["SESHMARK_MODEL"] = "devin-v1"
    
    subprocess.run(["git", "commit", "-m", message], env=env, check=True)

# Usage
# After Devin completes a run
devin_commit(run_id="run_abc123", message="feat: implement feature")
```

---

## SWE-Agent

SWE-Agent has explicit trajectory IDs.

### SWE-Agent Integration

```python
import os

def swe_agent_commit(trajectory_id: str, message: str):
    """Commit with SWE-Agent trajectory metadata."""
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = f"swe-agent:{trajectory_id}"
    env["SESHMARK_AGENT"] = "swe-agent"
    env["SESHMARK_MODEL"] = os.environ.get("SWE_AGENT_MODEL", "gpt-4o")
    
    subprocess.run(["git", "commit", "-m", message], env=env, check=True)

# Usage
swe_agent_commit(trajectory_id="traj_abc123", message="fix: resolve issue #42")
```

---

## Builder (Replit / Other Agents)

Generic wrapper for any agent tool.

### Universal Agent Wrapper

```bash
#!/bin/bash
# ~/.local/bin/agent-sesh
# Usage: agent-sesh <agent-name> <model> [args...]

AGENT="$1"
MODEL="$2"
shift 2

SESSION_ID="${AGENT}-$(date +%Y%m%d-%H%M%S)-$(openssl rand -hex 2)"

export SESHMARK_SESSION_ID="${AGENT}:${SESSION_ID}"
export SESHMARK_AGENT="$AGENT"
export SESHMARK_MODEL="$MODEL"

# Run the actual agent command
"$@"
```

### Python Universal Wrapper

```python
import os
import subprocess
import uuid
from datetime import datetime

def agent_commit(agent: str, model: str, message: str, tool_command: list[str] = None):
    """Universal commit wrapper for any AI agent.
    
    Args:
        agent: Agent name (cursor, claude, codex, etc.)
        model: Model identifier
        message: Git commit message
        tool_command: Optional command to run after setting env vars
    """
    session_id = f"{agent}:{datetime.now().strftime('%Y%m%d')}-{uuid.uuid4().hex[:6]}"
    
    env = os.environ.copy()
    env["SESHMARK_SESSION_ID"] = session_id
    env["SESHMARK_AGENT"] = agent
    env["SESHMARK_MODEL"] = model
    
    if tool_command:
        subprocess.run(tool_command, env=env, check=True)
    else:
        subprocess.run(["git", "commit", "-m", message], env=env, check=True)
    
    return session_id

# Examples
agent_commit("cursor", "claude-sonnet-4", "feat: auth middleware")
agent_commit("claude", "claude-sonnet-4", "fix: db race condition")
agent_commit("codex", "gpt-4o", "refactor: extract utils")
agent_commit("opencode", "qwen2.5-coder", "feat: add tests")
```

---

## Shell Alias Collection

Add these to `~/.bashrc` or `~/.zshrc`:

```bash
# Cursor
alias cursor-sesh='SESHMARK_SESSION_ID="cursor:$(basename $(pwd))-$(date +%Y%m%d)" SESHMARK_AGENT="cursor" cursor'

# Claude Code
alias claude-sesh='SESHMARK_SESSION_ID="claude:$(date +%Y%m%d)-$(openssl rand -hex 4)" SESHMARK_AGENT="claude-code" claude'

# Codex
alias codex-sesh='SESHMARK_SESSION_ID="codex:run-$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="codex" codex'

# Aider
alias aider-sesh='SESHMARK_SESSION_ID="aider:$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="aider" aider'

# OpenCode
alias opencode-sesh='SESHMARK_SESSION_ID="opencode:$(date +%Y%m%d-%H%M%S)" SESHMARK_AGENT="opencode" opencode'

# GitHub Copilot (VS Code)
alias code-copilot='SESHMARK_SESSION_ID="github-copilot:vscode-$(date +%Y%m%d)" SESHMARK_AGENT="github-copilot" code'
```

---

## Quick Reference

| Agent | Session ID Format | Model Env Var | Notes |
|-------|------------------|---------------|-------|
| Cursor | `cursor:chat-{id}` | `CURSOR_MODEL` | Chat ID from `~/.cursor/` |
| Claude Code | `claude:thread-{id}` | `CLAUDE_MODEL` | Thread ID from `~/.claude/` |
| Codex | `codex:run-{timestamp}` | `CODEX_MODEL` | Run-based, no persistence |
| GitHub Copilot | `github-copilot:{ide}-{date}` | `COPILOT_MODEL` | IDE + date fallback |
| Aider | `aider:{session-id}` | `AIDER_MODEL` | Built-in session tracking |
| OpenCode | `opencode:sess-{id}` | `OPENCODE_MODEL` | Extensible via plugin |
| Continue.dev | `continue:{timestamp}` | `CONTINUE_MODEL` | VS Code extension |
| Devin | `devin:{run-id}` | — | API-based run IDs |
| SWE-Agent | `swe-agent:{trajectory-id}` | `SWE_AGENT_MODEL` | Trajectory-based |

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
