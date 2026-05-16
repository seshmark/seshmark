# Harness Integration — Env Var Configuration Research

> How each major AI coding harness lets users set environment variables before `git commit`, enabling Seshmark's full metadata capture.

---

## Overview

Seshmark detects AI agents via 5 methods. The richest metadata (session + agent + model) comes from `SESHMARK_*` env vars (Priority 2). This document researches how users can inject these env vars for each popular harness without modifying the tool's source code.

---

## Tool-by-Tool Research

### 1. Pi

| Aspect | Detail |
|--------|--------|
| CLI binary | `pi` |
| Config file | `~/.pi/config` or project `.pi/config.yml` |
| Skill system | `~/.cursor/skills/<skill>/SKILL.md` (also cloned into `.claude/skills`, `.opencode/skills`) |
| Commit behavior | Pi runs `git commit` directly via shell |
| Env var injection method | **Skill instruction** — tell Pi to set SESHMARK_* before git commit |
| Best approach | Create a `seshmark` skill at `~/.cursor/skills/seshmark/SKILL.md` that instructs Pi: "Before running any git commit command, set SESHMARK_SESSION_ID, SESHMARK_AGENT, and SESHMARK_MODEL environment variables." |

**Skill file content** (what Pi should follow):

```
When committing code, before running git commit, always set:
  SESHMARK_SESSION_ID=pi:<unique-session-id>
  SESHMARK_AGENT=pi
  SESHMARK_MODEL=<current-model>

Use a UUID for the session ID. Detect the current model from context.
```

### 2. Claude Code

| Aspect | Detail |
|--------|--------|
| CLI binary | `claude` |
| Config file | `~/.claude/settings.json`, project-level `CLAUDE.md` |
| Skill system | `~/.claude/skills/<skill>/SKILL.md` |
| Commit behavior | Claude Code runs `git commit` via shell |
| Env var injection method | **CLAUDE.md instruction** or **Skill file** |
| Best approach | Create a skill at `~/.claude/skills/seshmark/SKILL.md` with instructions to set env vars before commit. Alternatively, create a wrapper script. |

### 3. Cursor

| Aspect | Detail |
|--------|--------|
| CLI binary | `cursor` (also GUI app) |
| Config file | `~/.cursor/settings.json`, project `.cursorrules` |
| Skill/system prompt | Cursor follows `.cursorrules` and system prompt |
| Commit behavior | Cursor's AI agent runs `git commit` via integrated terminal |
| Env var injection method | **Wrapper script** — create a `git-agentblame-wrapper` that sets env vars and aliases `git commit` |
| Best approach | Users configure Cursor to use a custom commit command, or create a global git hook (which seshmark already installs). For env vars specifically, users set them in their shell profile before launching Cursor, or use a `.cursorrules` instruction. |

### 4. OpenCode

| Aspect | Detail |
|--------|--------|
| CLI binary | `opencode` |
| Config file | `~/.opencode/config.json`, project-level config |
| Skill system | `~/.opencode/skills/<skill>/SKILL.md` |
| Commit behavior | OpenCode runs `git commit` via shell |
| Env var injection method | **Skill instruction** or **wrapper script** |
| Best approach | Create a skill at `~/.opencode/skills/seshmark/SKILL.md` that instructs OpenCode to set env vars before committing. |

### 5. Aider

| Aspect | Detail |
|--------|--------|
| CLI binary | `aider` |
| Config file | `~/.aider.conf.yml`, project `.aider.conf.yml` |
| Env var support | Aider has `--env` flag and reads env vars |
| Commit behavior | Aider runs `git commit` via shell |
| Env var injection method | **Shell wrapper** or **aider config** |
| Best approach | Users create a wrapper script `~/bin/aider` that sets SESHMARK_* before calling the real `aider`. Or use aider's `--env` flag if available. |

### 6. GitHub Copilot (VS Code)

| Aspect | Detail |
|--------|--------|
| Integration | VS Code extension |
| Config | VS Code settings |
| Commit behavior | VS Code's built-in git UI calls `git commit` |
| Env var injection method | **Global git config hook** or **VS Code task** |
| Best approach | Use VS Code's `git.enableCommitSigning` or configure a VS Code task that runs before commit. The seshmark global hook template (already installed) handles this — but Copilot doesn't set env vars natively. Users rely on branch inference or process detection instead. |

### 7. Other CLI tools (generic)

| Aspect | Detail |
|--------|--------|
| Method | **Wrapper script** pattern |
| Pattern | Users create a script that wraps the tool binary, sets SESHMARK_* env vars, then execs the real tool. |
| Installation | Drop the wrapper in `~/bin/` earlier in PATH than the real tool |

---

## Common Patterns

### Pattern A: Skill/Agent Instruction

For tools that read agent instructions (Pi, Claude Code, OpenCode, Cursor), create a skill file at the tool's skill directory instructing it to set env vars before every git commit.

**Universal instruction block:**

```
## Seshmark Integration

Before running any `git commit` command, you MUST set these environment variables:

SESHMARK_SESSION_ID=<your-tool-name>:<unique-session-id>
SESHMARK_AGENT=<your-tool-name>
SESHMARK_MODEL=<current-model>

Generate a unique session ID (UUID) for each session.
The agent name should match the tool name used in Seshmark's known agents list.
```

### Pattern B: Shell Wrapper

For any CLI tool, create a wrapper script that intercepts the tool binary:

```bash
#!/bin/bash
# ~/bin/pi — wrapper for pi that adds Seshmark metadata
REAL_PI=$(which -a pi | grep -v "$HOME/bin" | head -1)

# Generate session ID
SESSION_ID="pi:$(uuidgen)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="pi"
export SESHMARK_MODEL="$(detect_model_hint)"

exec "$REAL_PI" "$@"
```

### Pattern C: Git Hook (already exists)

Seshmark's own `prepare-commit-mssg` hook intercepts commits. Users don't need to do anything else for basic agent detection — process detection and branch inference work out of the box. Env vars are needed only for **full metadata** (session + model).

---

## Files to Create

| File | Purpose | Content |
|------|---------|---------|
| `examples/harnesses/pi/SKILL.md` | Pi skill to set env vars | Instruction block + env var template |
| `examples/harnesses/claude-code/SKILL.md` | Claude Code skill | Instruction block |
| `examples/harnesses/opencode/SKILL.md` | OpenCode skill | Instruction block |
| `examples/harnesses/cursor/.cursorrules` | Cursor rules snippet | Instruction block |
| `examples/harnesses/aider/wrapper.sh` | Aider wrapper script | Shell script template |
| `examples/harnesses/generic/wrapper.sh` | Generic wrapper script | Shell script template for any tool |
| `HARNESS_EXAMPLES.md` | Updated integration guide | Tool-by-tool examples |
