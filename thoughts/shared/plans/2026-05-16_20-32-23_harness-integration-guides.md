---
date: 2026-05-16T20:32:23+0530
author: Sahil Choudhary
commit: 6ef1b14
branch: main
repository: seshmark
topic: "Harness Integration — Env Var Configuration Guides"
tags: [plan, harness-integration, documentation, examples]
status: ready
parent: "thoughts/shared/solutions/2026-05-16_harness-env-var-config.md"
last_updated: 2026-05-16T20:32:23+0530
last_updated_by: Sahil Choudhary
---

# Harness Integration — Env Var Configuration Guides

## Overview

Create per-harness configuration guides (skills, wrappers, rules files) that tell each AI coding tool how to set `SESHMARK_SESSION_ID`, `SESHMARK_AGENT`, and `SESHMARK_MODEL` environment variables before every `git commit`. This gives users full metadata capture (session + agent + model) for every harness, beyond what process detection alone provides.

Based on research in `thoughts/shared/solutions/2026-05-16_harness-env-var-config.md`.

## Desired End State

- A user of any major harness (Pi, Claude Code, OpenCode, Cursor, Aider) can follow a 1-minute setup to get full Seshmark metadata
- Each harness has a documented example in `examples/harnesses/<tool>/`
- `HARNESS_EXAMPLES.md` links to all harness examples and explains the three setup patterns (skill, wrapper, git hook)
- CI tests verify that each harness's config method works end-to-end

## What We're NOT Doing

- NOT modifying any harness's source code (these are user-side configurations only)
- NOT building a new detection method (env vars already work — this is about helping users set them)
- NOT creating a universal auto-installer (users choose their harness and apply the relevant config)
- NOT covering every edge case (focused on the 5 most popular harnesses)

## Phase 1: Pi Skill

### Overview
Create a Pi skill file at `examples/harnesses/pi/SKILL.md` that instructs Pi to set `SESHMARK_*` env vars before every `git commit`. Pi reads skills from `~/.cursor/skills/<name>/SKILL.md` (or symlinked in `.claude/skills`, `.opencode/skills`).

### Changes Required:

#### 1. Create Pi skill
**File**: `examples/harnesses/pi/SKILL.md`
**Changes**: New file — Pi skill that instructs the agent to set env vars before committing

```markdown
---
name: seshmark
description: Integrate with Seshmark to tag AI commits with session metadata. Use when committing code to ensure AI attribution is captured.
---

# Seshmark Integration for Pi

This skill tells Pi to automatically tag Git commits with AI metadata when Seshmark is installed.

## Behavior

Before running any `git commit` command, you MUST set these environment variables:

```bash
export SESHMARK_SESSION_ID="pi:$(uuidgen)"
export SESHMARK_AGENT="pi"
export SESHMARK_MODEL="<current-model>"
```

- **SESHMARK_SESSION_ID**: Generate a unique UUID per session. Format: `pi:<uuid>`
- **SESHMARK_AGENT**: Always set to `pi`
- **SESHMARK_MODEL**: Detect the current model from context (e.g., `deepseek-v4-flash`, `claude-sonnet-4-20250514`)

Run `git commit` WITH these variables set — either inline:

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
```

### Success Criteria:

#### Automated Verification:
- [ ] File exists at `examples/harnesses/pi/SKILL.md`

#### Manual Verification:
- [ ] Copy skill to `~/.cursor/skills/seshmark/` and verify Pi follows the instructions on next commit
- [ ] After commit, `seshmark who HEAD` shows `AI-Agent: pi`, `AI-Session: pi:<uuid>`, and `AI-Model: <model>`

---

## Phase 2: Claude Code Skill

### Overview
Create a Claude Code skill at `examples/harnesses/claude-code/SKILL.md`. Claude Code reads skills from `~/.claude/skills/<name>/SKILL.md` and also respects `CLAUDE.md` in the project root.

### Changes Required:

#### 1. Create Claude Code skill
**File**: `examples/harnesses/claude-code/SKILL.md`
**Changes**: New file — skill telling Claude Code to set env vars before commit

```markdown
---
name: seshmark
description: Integrate with Seshmark to tag AI commits with session metadata. Use when committing code.
---

# Seshmark Integration for Claude Code

Before running any `git commit` command, you MUST set these environment variables:

```bash
export SESHMARK_SESSION_ID="claude:$(uuidgen)"
export SESHMARK_AGENT="claude"
export SESHMARK_MODEL="<current-model>"
```

- **SESHMARK_SESSION_ID**: Format `claude:<uuid>`, unique per session
- **SESHMARK_AGENT**: Always `claude`
- **SESHMARK_MODEL**: Current model name from context

Set them inline or export before the commit command.
```

### Success Criteria:

#### Automated Verification:
- [ ] File exists at `examples/harnesses/claude-code/SKILL.md`
- [ ] File content includes `SESHMARK_AGENT="claude"`

#### Manual Verification:
- [ ] Copy skill to `~/.claude/skills/seshmark/` and verify Claude Code sets env vars on next commit

---

## Phase 3: OpenCode Skill

### Overview
Create an OpenCode skill at `examples/harnesses/opencode/SKILL.md`. OpenCode reads skills from `~/.opencode/skills/<name>/SKILL.md`.

### Changes Required:

#### 1. Create OpenCode skill
**File**: `examples/harnesses/opencode/SKILL.md`
**Changes**: New file

```markdown
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
```

### Success Criteria:

#### Automated Verification:
- [ ] File exists at `examples/harnesses/opencode/SKILL.md`

#### Manual Verification:
- [ ] Copy to `~/.opencode/skills/seshmark/` and verify the skill is loaded

---

## Phase 4: Cursor Rules + Aider Wrapper

### Overview
Cursor reads `.cursorrules` for agent instructions. Aider uses a shell wrapper pattern (no skill system). Create both in parallel.

### Changes Required:

#### 1. Create Cursor rules snippet
**File**: `examples/harnesses/cursor/.cursorrules`
**Changes**: New file

```markdown
# Seshmark Integration

Before committing code via git commit, set these environment variables:

SESHMARK_SESSION_ID=cursor:<unique-uuid>
SESHMARK_AGENT=cursor
SESHMARK_MODEL=<current-model>

Run git commit with these variables exported so Seshmark can capture full metadata.
```

#### 2. Create Aider wrapper script
**File**: `examples/harnesses/aider/wrapper.sh`
**Changes**: New file

```bash
#!/bin/bash
# Seshmark wrapper for Aider
#
# Install:
#   cp examples/harnesses/aider/wrapper.sh ~/bin/aider
#   chmod +x ~/bin/aider
#
# Make sure ~/bin comes before the real aider in PATH.

REAL_AIDER=$(which -a aider | grep -v "$HOME/bin" | head -1)

if [ -z "$REAL_AIDER" ]; then
  echo "Error: Real aider not found in PATH"
  exit 1
fi

# Generate session ID
SESSION_ID="aider:$(uuidgen 2>/dev/null || python3 -c 'import uuid; print(uuid.uuid4())' 2>/dev/null || date +%s)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="aider"

# Set model from aider's output if available, or let user configure
export SESHMARK_MODEL="${SESHMARK_MODEL:-}"

exec "$REAL_AIDER" "$@"
```

### Success Criteria:

#### Automated Verification:
- [ ] File exists at `examples/harnesses/cursor/.cursorrules`
- [ ] File exists at `examples/harnesses/aider/wrapper.sh`
- [ ] `wrapper.sh` is executable-checkable (contains `#!/bin/bash`)

#### Manual Verification:
- [ ] Copy wrapper to `~/bin/aider`, ensure `~/bin` is first in PATH, run `aider --version` and verify env vars are set
- [ ] Place `.cursorrules` in project root and verify Cursor loads it

---

## Phase 5: Generic Wrapper + HARNESS_EXAMPLES.md

### Overview
Create a generic wrapper script for any tool, and update `HARNESS_EXAMPLES.md` with all per-harness guides in one place.

### Changes Required:

#### 1. Create generic wrapper
**File**: `examples/harnesses/generic/wrapper.sh`
**Changes**: New file

```bash
#!/bin/bash
# Seshmark wrapper for any CLI tool
#
# Usage:
#   1. Copy this script and rename it to match your tool
#   2. Replace TOOL_NAME with your tool's name
#   3. Make sure ~/bin comes first in PATH
#   4. Run your tool normally

TOOL_NAME="${TOOL_NAME:-my-tool}"
REAL_BIN=$(which -a "$TOOL_NAME" | grep -v "$HOME/bin" | head -1)

if [ -z "$REAL_BIN" ]; then
  echo "Error: Real $TOOL_NAME not found in PATH"
  exit 1
fi

SESSION_ID="${TOOL_NAME}:$(uuidgen 2>/dev/null || python3 -c 'import uuid; print(uuid.uuid4())' 2>/dev/null || date +%s)"
export SESHMARK_SESSION_ID="$SESSION_ID"
export SESHMARK_AGENT="$TOOL_NAME"
export SESHMARK_MODEL="${SESHMARK_MODEL:-}"

exec "$REAL_BIN" "$@"
```

#### 2. Update HARNESS_EXAMPLES.md
**File**: `HARNESS_EXAMPLES.md`
**Changes**: Add per-harness sections with install instructions

```markdown
# Harness Integration Examples

> How to configure each AI coding tool to set Seshmark environment variables for full metadata capture.

## Quick Reference

| Harness | Method | Setup Time | Full Metadata? |
|---------|--------|-----------|---------------|
| Pi | Skill file | 30s | ✅ session + agent + model |
| Claude Code | Skill file | 30s | ✅ session + agent + model |
| OpenCode | Skill file | 30s | ✅ session + agent + model |
| Cursor | .cursorrules | 30s | ✅ session + agent + model |
| Aider | Wrapper script | 1 min | ✅ session + agent |
| Any CLI tool | Wrapper script | 1 min | ✅ session + agent |

## Pi

Copy the skill to your Pi skills directory:

```bash
cp examples/harnesses/pi/SKILL.md ~/.cursor/skills/seshmark/SKILL.md
```

Pi will now set `SESHMARK_SESSION_ID`, `SESHMARK_AGENT`, and `SESHMARK_MODEL` before every `git commit`.

## Claude Code

```bash
cp examples/harnesses/claude-code/SKILL.md ~/.claude/skills/seshmark/SKILL.md
```

## OpenCode

```bash
cp examples/harnesses/opencode/SKILL.md ~/.opencode/skills/seshmark/SKILL.md
```

## Cursor

Copy `.cursorrules` to your project root:

```bash
cp examples/harnesses/cursor/.cursorrules .cursorrules
```

## Aider

Install the wrapper:

```bash
cp examples/harnesses/aider/wrapper.sh ~/bin/aider
chmod +x ~/bin/aider
# Ensure ~/bin is first in PATH
export PATH="$HOME/bin:$PATH"
```

## Generic (Any CLI Tool)

```bash
cp examples/harnesses/generic/wrapper.sh ~/bin/your-tool
# Edit the wrapper: change TOOL_NAME to your tool's name
chmod +x ~/bin/your-tool
```

## Verification

After setting up, make a commit and check:

```bash
seshmark who HEAD
# Should show:
#   AI-Session: <tool>:<uuid>
#   AI-Agent: <tool>
#   AI-Model: <model>
```
```

### Success Criteria:

#### Automated Verification:
- [ ] File exists at `examples/harnesses/generic/wrapper.sh`
- [ ] `HARNESS_EXAMPLES.md` mentions all 5 harnesses
- [ ] All skill files referenced in `HARNESS_EXAMPLES.md` actually exist

#### Manual Verification:
- [ ] Follow each harness guide end-to-end and verify metadata appears in commits

---

## Phase 6: CI Test for Harness Configs

### Overview
Add a GitHub Actions job that tests each harness's skill/wrapper pattern by simulating what the harness does.

### Changes Required:

#### 1. Add harness config tests to CI
**File**: `.github/workflows/test-harnesses.yml` (append to existing)
**Changes**: Add a new job that tests the skill/wrapper files exist, are parseable, and produce correct env vars

```yaml
  test-harness-configs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Verify all harness example files exist
        run: |
          for file in \
            examples/harnesses/pi/SKILL.md \
            examples/harnesses/claude-code/SKILL.md \
            examples/harnesses/opencode/SKILL.md \
            examples/harnesses/cursor/.cursorrules \
            examples/harnesses/aider/wrapper.sh \
            examples/harnesses/generic/wrapper.sh; do
            if [ ! -f "$file" ]; then
              echo "❌ Missing: $file"
              exit 1
            fi
            echo "✅ $file"
          done

      - name: Verify skill files have correct agent names
        run: |
          grep -q 'SESHMARK_AGENT="pi"' examples/harnesses/pi/SKILL.md || exit 1
          grep -q 'SESHMARK_AGENT="claude"' examples/harnesses/claude-code/SKILL.md || exit 1
          grep -q 'SESHMARK_AGENT="opencode"' examples/harnesses/opencode/SKILL.md || exit 1
          grep -q 'SESHMARK_AGENT="cursor"' examples/harnesses/cursor/.cursorrules || exit 1
          echo "✅ All skill files reference correct agent names"

      - name: Verify wrapper scripts are valid shell
        run: |
          bash -n examples/harnesses/aider/wrapper.sh || exit 1
          bash -n examples/harnesses/generic/wrapper.sh || exit 1
          echo "✅ Wrapper scripts have valid shell syntax"

      - name: Verify HARNESS_EXAMPLES.md links to all harnesses
        run: |
          for harness in pi claude-code opencode cursor aider; do
            grep -q "examples/harnesses/$harness" HARNESS_EXAMPLES.md || exit 1
          done
          echo "✅ HARNESS_EXAMPLES.md links to all harnesses"
```

### Success Criteria:

#### Automated Verification:
- [ ] CI job passes: all harness files exist, agent names are correct, shell syntax is valid
- [ ] `HARNESS_EXAMPLES.md` links to all 5 harness directories

---

## Testing Strategy

### Automated:
- `go test ./...` — all existing tests still pass
- CI workflow test-harness-configs job — verifies file existence, shell syntax, agent name correctness

### Manual Testing Steps:
1. Install Pi skill: `cp examples/harnesses/pi/SKILL.md ~/.cursor/skills/seshmark/SKILL.md`, make a commit via Pi, verify `seshmark who HEAD` shows full metadata
2. Install Claude Code skill, repeat verification
3. Install OpenCode skill, repeat verification
4. Set up Aider wrapper, run `aider --version`, check env vars are set
5. Place `.cursorrules` in a test project, verify Cursor loads it

## Performance Considerations

None — these are static configuration files, no runtime impact.

## Migration Notes

- The old `examples/resolvers/` directory contains session resumption scripts only. Keep those as-is.
- The new `examples/harnesses/` directory is for pre-commit env var configuration only. They are complementary.
- No breaking changes to existing functionality.

## References

- Research: `thoughts/shared/solutions/2026-05-16_harness-env-var-config.md`
- Existing resolvers: `examples/resolvers/`
- Existing docs: `ADD_A_HARNESS.md`
