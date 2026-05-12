# Seshmark

> `git agentblame` — like `git blame`, but shows which AI (Cursor, Claude, Copilot) wrote each line.

<!-- Replace with your actual GIF: -->
<!-- ![Demo](docs/demo.gif) -->

[![Version](https://img.shields.io/github/v/release/seshmark/seshmark?color=3fb950&label=version)](https://github.com/seshmark/seshmark/releases)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](go.mod)

```bash
curl -fsSL https://seshmark.dev/install | bash
```

**Try it without installing:**
```bash
git clone https://github.com/seshmark/demo
cd demo
git agentblame src/api.ts
```

## What It Does

Seshmark auto-tags every AI commit with metadata — agent, model, session ID — embedded directly in Git commit messages. Later, you can query, blame, and resume sessions.

**Zero configuration.** Install once. Use any AI tool. The Git hook handles everything.

```bash
$ git agentblame src/auth.ts

  1  [human]           const express = require('express');
  2  [cursor]          import { oauth } from './oauth';
  3  [claude]          function validateToken(token: string) {
  4  [human]             if (!token) return null;
  5  [copilot]           const hash = crypto.sha256(token);
```

## Install

```bash
curl -fsSL https://seshmark.dev/install | bash
```

Works on macOS, Linux, and WSL. One command. No dependencies.

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
| `seshmark stats` | Shareable report: AI % by agent, model, file |
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

## VS Code Extension

See `[cursor]`, `[claude]` inline in your editor gutter. Hover for session details.

[Install from VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=seshmark.seshmark)

## License

MIT
