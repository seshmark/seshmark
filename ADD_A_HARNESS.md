# Adding a Custom AI Harness to Seshmark

> Three ways to make any AI coding tool work with Seshmark — no library, no API, no dependency.

Seshmark is designed to be **harness-agnostic**: you don't need to wait for us to add support for your tool. If your tool can set environment variables or follow a branch naming convention, it already works. If you want rich session resumption, you can write a tiny script.

---

## Table of Contents

- [How It Works (30 Seconds)](#how-it-works-30-seconds)
- [Method 1: Environment Variables (Easiest)](#method-1-environment-variables-easiest)
- [Method 2: Branch Naming (Zero Config)](#method-2-branch-naming-zero-config)
- [Method 3: Custom Resolver (Session Resumption)](#method-3-custom-resolver-session-resumption)
- [Testing Your Harness](#testing-your-harness)
- [Publishing Your Resolver](#publishing-your-resolver)
- [Example: Adding Support for a New Tool](#example-adding-support-for-a-new-tool)

---

## How It Works (30 Seconds)

Every time `git commit` runs, Seshmark's hook fires and detects the AI agent using this priority order:

| Priority | Method | What it captures | Set by |
|----------|--------|-----------------|--------|
| **1** | `seshmark track` session file | session + agent + model | You (manual) |
| **2** | `SESHMARK_*` env vars | session + agent + model | Your AI tool |
| **3** | `AI_*` env vars (legacy) | session + agent + model | Older tools |
| **4** | Branch name | agent only | `pi/feature` → `pi` |
| **5** | Parent process detection | agent only | Any CLI tool running `git` |
| **6** | Nothing found | → tagged as `[human]` | — |

The **first match wins**. So if a tool sets `SESHMARK_SESSION_ID`, that takes priority over branch naming or process detection.

Once detected, Seshmark appends Git trailers to the commit message automatically:

```
Seshmark-Version: 1.0.0
AI-Session: my-tool:sess-001
AI-Agent: my-tool
AI-Model: gpt-4o
```

That's it. No config file. No import. No central server.

---

## Method 1: Environment Variables (Easiest)

Set these three environment variables before `git commit`, and Seshmark handles the rest:

| Variable | Required | Example | Description |
|----------|----------|---------|-------------|
| `SESHMARK_SESSION_ID` | Yes | `my-tool:sess-001` | Unique session identifier. Convention: `<tool-name>:<id>` |
| `SESHMARK_AGENT` | Yes | `my-tool` | Name of the AI coding tool |
| `SESHMARK_MODEL` | No | `gpt-4o` | Model name (useful for stats) |

### Example: Wrapping a Tool

```bash
#!/bin/bash
# ~/bin/my-tool-wrapper.sh

# Capture the session start
SESSION_ID="my-tool:$(uuidgen)"
SESHMARK_SESSION_ID="$SESSION_ID" \
SESHMARK_AGENT="my-tool" \
SESHMARK_MODEL="gpt-4o" \
my-tool "$@"
```

Any commits made inside the tool are automatically tagged.

### Example: Programmatic (Python)

```python
import os, subprocess

env = os.environ.copy()
env["SESHMARK_SESSION_ID"] = "my-tool:sess-001"
env["SESHMARK_AGENT"] = "my-tool"
env["SESHMARK_MODEL"] = "gpt-4o"

subprocess.run(["git", "commit", "-m", "feat: implemented by my-tool"], env=env)
```

### Example: Programmatic (Node.js)

```javascript
const { execSync } = require('child_process');

execSync('git commit -m "feat: implemented by my-tool"', {
  env: {
    ...process.env,
    SESHMARK_SESSION_ID: 'my-tool:sess-001',
    SESHMARK_AGENT: 'my-tool',
    SESHMARK_MODEL: 'gpt-4o',
  },
});
```

---

## Method 2: Branch Naming (Zero Config)

No environment variables needed. Name your branch with the pattern `<agent-name>/<description>` and Seshmark infers the agent:

```bash
git checkout -b my-tool/implement-oauth
# commits on this branch get: AI-Agent: my-tool
```

This works automatically — no env vars, no config, no wrapper script. Seshmark watches branches pushed by any tool.

---

## Method 3: Custom Resolver (Session Resumption)

A resolver is a small script that tells Seshmark how to reopen a session in your tool. It's how `seshmark resume <session-id>` knows what to do.

### Resolver Contract

- **Location**: `~/.local/share/seshmark/resolvers/<agent-name>`
- **Input**: One argument — the session ID (e.g., `my-tool:sess-001`)
- **Output**: Open the session, or print instructions
- **Exit code**: 0 on success, non-zero on failure

### Minimal Resolver

```bash
#!/bin/bash
# ~/.local/share/seshmark/resolvers/my-tool
SESSION_ID="$1"
echo "Resuming session: $SESSION_ID"
echo "Open your tool and look for session: $SESSION_ID"
```

### Full Resolver with Tool Detection

```bash
#!/bin/bash
# ~/.local/share/seshmark/resolvers/my-tool
SESSION_ID="$1"

if command -v my-tool &>/dev/null; then
  # Tool is installed — open the session directly
  my-tool open --session "$SESSION_ID"
else
  echo "my-tool not found in PATH."
  echo "Install: npm install -g my-tool"
  echo ""
  echo "Session ID: $SESSION_ID"
  echo ""
  echo "Or continue manually: my-tool open --session $SESSION_ID"
  exit 1
fi
```

### Install Your Resolver

```bash
# Method A: Use seshmark's resolver command
seshmark resolver install my-tool ./my-resolver.sh

# Method B: Place it manually
cp my-resolver.sh ~/.local/share/seshmark/resolvers/my-tool
chmod +x ~/.local/share/seshmark/resolvers/my-tool

# Test it
seshmark resolver test my-tool
```

---

## Testing Your Harness

### Verify Environment Variables

```bash
# Set up a test session
export SESHMARK_SESSION_ID="my-tool:test-001"
export SESHMARK_AGENT="my-tool"
export SESHMARK_MODEL="test-model"

# Make a test commit
git commit --allow-empty -m "test: my-tool attribution"

# Verify the trailers were added
git log -1 --format="%B"
# Should show:
#   test: my-tool attribution
#   SESHMARK-Version: 1.0.0
#   AI-Session: my-tool:test-001
#   AI-Agent: my-tool
#   AI-Model: test-model
```

### Verify Query

```bash
seshmark query --agent my-tool
# Should show your test commit
```

### Verify Resolver (if you wrote one)

```bash
seshmark resolver test my-tool
# Should run your resolver script with a dummy session ID
```

---

## Publishing Your Resolver

Once your harness is working:

1. **Share the resolver** — commit it to your project or publish as a gist
2. **Submit to Seshmark** — open a PR adding your resolver to `examples/resolvers/`. See [`CONTRIBUTING.md`](CONTRIBUTING.md)
3. **Tell the community** — post about it on [GitHub Discussions](https://github.com/seshmark/seshmark/discussions)

Your resolver script becomes the reference for other users of your tool. A well-documented resolver is itself good documentation.

---

## Example: Adding Support for a New Tool

Let's walk through adding support for a hypothetical tool called `my-coder`.

### Step 1: Environment Variables

```bash
# In my-coder's source or wrapper
export SESHMARK_SESSION_ID="my-coder:$(uuidgen)"
export SESHMARK_AGENT="my-coder"
export SESHMARK_MODEL="claude-sonnet-4"
```

### Step 2: Resolver Script

```bash
#!/bin/bash
# ~/.local/share/seshmark/resolvers/my-coder
SESSION_ID="$1"

if command -v my-coder &>/dev/null; then
  # Strip prefix if present: "my-coder:abc123" -> "abc123"
  ID="${SESSION_ID#my-coder:}"
  my-coder resume --session "$ID"
else
  echo "my-coder not found." >&2
  echo "Install: pip install my-coder" >&2
  exit 1
fi
```

### Step 3: Test the Full Flow

```bash
# Tag a commit
export SESHMARK_SESSION_ID="my-coder:test-001"
export SESHMARK_AGENT="my-coder"
git commit --allow-empty -m "test: full flow"

# Blame it
git agentblame HEAD

# Resume it
seshmark resume HEAD
# -> Resolves to my-coder:test-001
# -> Opens my-coder resume --session test-001
```

### Step 4: Contribute

Open a PR adding your resolver to `examples/resolvers/my-coder`. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

---

## Key Concept: The Seshmark Trailers

When Seshmark tags a commit, it adds these trailers to the commit message:

```
Seshmark-Version: 1.0.0
AI-Session: <tool-name>:<session-id>
AI-Agent: <tool-name>
AI-Model: <model-name>
```

These are **standard Git trailers** — they parse with `git log --format` and work even if Seshmark isn't installed. You or your users can type them by hand and get full compatibility.

---

*Questions? Open a [discussion](https://github.com/seshmark/seshmark/discussions) or [issue](https://github.com/seshmark/seshmark/issues).*
