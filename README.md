# Seshmark — Agent-Aware Git History

> Also known as: `git agentblame`, `agentblame`, `aiblame`

Seshmark is a zero-infrastructure, tool-agnostic convention that links AI coding agent sessions to Git commits. It enables you to:

- **Blame** any line to see which AI session wrote it: `git agentblame src/auth.ts`
- **Query** your AI commit history: `seshmark query --agent claude --path auth`
- **Resume** sessions: `seshmark resume HEAD` opens the original Cursor/Claude chat
- **Context** for orchestrators: `seshmark context HEAD --format json` reconstructs the full session

Works with **any** AI tool — Cursor, Claude Code, OpenCode, Codex, Copilot, Aider, and more.

## Install

```bash
curl -fsSL https://seshmark.dev/install | bash
```

Then try it:

```bash
git agentblame src/auth.ts
seshmark who HEAD
seshmark query --agent cursor
```

## How It Works

1. **Install** seshmark once. It adds a Git hook that auto-tags commits.
2. **Use your AI tool normally.** Commits get stamped with `AI-Session`, `AI-Agent`, `AI-Model` trailers.
3. **Or use branch naming.** A branch like `claude/auth-refactor` automatically tags commits with `AI-Agent: claude`.
4. **Query later.** Find who wrote what, when, and with which session.

## The Convention

AI-assisted commits include trailers:

```
feat: implement OAuth2 middleware

Seshmark-Version: 1.0.0
AI-Session: cursor:chat-abc123
AI-Agent: cursor
AI-Model: claude-sonnet-4-20250514
```

This convention works even without installing the seshmark CLI — just type the trailers manually.

## Commands

| Command | Description |
|---------|-------------|
| `git agentblame <file>` | Show AI attribution per line |
| `seshmark query [filters]` | Search AI-tagged commits |
| `seshmark who [commit]` | Show metadata for a commit |
| `seshmark context <commit>` | Reconstruct session context |
| `seshmark resume [commit]` | Resume session via native resolver |
| `seshmark log <session>` | List commits in a session |
| `seshmark status` | Show active session |
| `seshmark doctor` | Diagnose installation |
| `seshmark track <id>` | Manual session tracking |
| `seshmark untrack` | Stop tracking |

## For Harness Builders

Add three env vars before `git commit`:

```python
env["SESHMARK_SESSION_ID"] = "opencode:sess-001"
env["SESHMARK_AGENT"] = "opencode"
env["SESHMARK_MODEL"] = "qwen2.5-coder"
```

That's it. No library import. No API call. No dependency.

## For Orchestrators

Every command supports `--format json`:

```bash
seshmark blame src/auth.ts --format json
seshmark context HEAD --format json
```

This enables autonomous agents to discover sessions and reconstruct context programmatically.

## License

MIT
