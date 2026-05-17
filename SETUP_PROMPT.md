# Seshmark — One-Shot Setup for AI Agents

Give this prompt to any AI coding agent (Pi, Claude Code, Cursor, OpenCode, etc.) and it will set up Seshmark automatically.

## The Prompt

Copy and paste this into your AI agent:

```
Please set up Seshmark (git agentblame) for me:

1. Install the CLI:
   curl -fsSL https://seshmark.github.io/seshmark/install.sh | bash

2. Load the Seshmark skill so all future commits get tagged:
   - If you're Pi: use /load-skill <path> after downloading
   - If you're Claude Code: copy examples/harnesses/claude-code/SKILL.md to ~/.claude/skills/seshmark/
   - If you're Cursor: copy examples/harnesses/cursor/.cursorrules to the project root
   - If you're OpenCode: copy examples/harnesses/opencode/SKILL.md to ~/.opencode/skills/seshmark/
   - If you're Aider: copy examples/harnesses/aider/wrapper.sh to ~/bin/aider

3. Before every git commit, always set these three environment variables:
   export SESHMARK_SESSION_ID="<tool-name>:$(uuidgen)"
   export SESHMARK_AGENT="<tool-name>"
   export SESHMARK_MODEL="<current-model>"

4. Make a test commit to verify it works:
   echo "seshmark setup" > test.txt && git add test.txt
   SESHMARK_SESSION_ID="<tool-name>:test-$(date +%s)" \
   SESHMARK_AGENT="<tool-name>" \
   SESHMARK_MODEL="<current-model>" \
   git commit -m "test: seshmark setup"
   seshmark who HEAD
```

## What This Does

| Step | What happens |
|------|-------------|
| 1 | Installs the seshmark CLI and Git hook |
| 2 | Loads the behavioral skill so the agent remembers to set env vars |
| 3 | Tells the agent to always tag commits with session + agent + model |
| 4 | Verifies everything works end-to-end |

## For Power Users

If you want full control, see:
- [ADD_A_HARNESS.md](ADD_A_HARNESS.md) — three ways to add any tool
- [HARNESS_EXAMPLES.md](HARNESS_EXAMPLES.md) — per-tool configuration guides
- [examples/harnesses/](examples/harnesses/) — skill files and wrapper scripts
