# Seshmark Launch Copy — Copy-Paste Ready

## Twitter / X Post

**The hook (video required):**

```
I built `git agentblame` — like `git blame`, but it shows which AI (Cursor, Claude, Copilot) wrote each line of your code.

Install in 10 seconds:
```
curl -fsSL https://seshmark.dev/install | bash
```

Works with any AI tool. No vendor lock-in. No central database.

🔗 seshmark.dev

#buildinpublic #ai #devtools
```

**Follow-up tweet (thread, 2 hours later):**

```
How it works:

1. Install seshmark once
2. Use Cursor / Claude Code / whatever you want
3. Every commit gets tagged automatically
4. Run `git agentblame <file>` to see which AI wrote which line

Branch naming works too: `claude/feature-auth` → auto-tagged as `AI-Agent: claude`

No config needed.
```

**Third tweet (next day):**

```
Built a VS Code extension too.

See `[cursor]`, `[claude]` inline in your editor gutter. Hover for session ID and model.

Install: seshmark.dev/vscode
```

---

## Reddit Posts

### r/cursor

**Title:** `git agentblame` — a drop-in command that shows which Cursor/Claude/Copilot session wrote each line of your code

**Body:**
```
I built a tool that auto-tags every AI commit with metadata, so you can later trace any line back to the exact AI session that wrote it.

**What it does:**
- `git agentblame src/auth.ts` → shows `[cursor|chat-abc123|cursor]` per line
- `seshmark query --agent cursor --path auth` → find all Cursor commits touching auth
- `seshmark resume HEAD` → reopen the original session

**How it works:**
A Git hook reads env vars set by your AI tool (or infers from branch names like `claude/feature-auth`) and appends trailers to commit messages:

```
Seshmark-Version: 1.0.0
AI-Session: cursor:chat-abc123
AI-Agent: cursor
```

**Install:**
```bash
curl -fsSL https://seshmark.dev/install | bash
```

Zero config after install. Works with any AI tool. Open source.

GitHub: github.com/seshmark/seshmark
```

### r/LocalLLaMA

**Title:** Seshmark — trace AI-written code back to its original session, tool-agnostic

**Body:**
```
I built a lightweight convention for linking AI coding sessions to Git commits. Works with Cursor, Claude Code, OpenCode, Aider, Copilot — anything that can set env vars before `git commit`.

The core idea: commit message trailers (`AI-Session`, `AI-Agent`, `AI-Model`) that travel with the repo forever. No external database. No vendor API calls.

```bash
git agentblame src/auth.ts
# [cursor|abc123|cursor] import { oauth } from './oauth';
# [human| | ] const x = 1;
# [claude|xyz789|claude] function validateToken(token) {
```

Install: `curl -fsSL https://seshmark.dev/install | bash`

GitHub: github.com/seshmark/seshmark
```

### r/programming

**Title:** I built `git agentblame` — git blame for the AI coding era

**Body:**
```
Every day I use Cursor or Claude Code, commit the changes, and a week later I have no idea which session produced what.

So I built a convention + CLI tool that auto-tags every AI commit, then lets you query and blame later.

```bash
$ git agentblame src/auth.ts
[cursor|chat-abc|cursor] import { oauth } from './oauth';
[human          |       ] const x = 1;
[claude|thread-x|claude] function validateToken(token) {
```

**Features:**
- Zero config — install once, hook auto-tags commits
- Tool agnostic — works with any AI tool via env vars
- Branch inference — `claude/feature` auto-tagged without explicit tracking
- Orchestrator-ready — every command outputs `--format json`

Open source, MIT. Would love feedback.

GitHub: github.com/seshmark/seshmark
```

### r/github

**Title:** Seshmark — a Git-native convention for tracking AI code attribution

**Body:**
```
A convention for embedding AI session metadata directly into Git commit messages via trailers (`AI-Session`, `AI-Agent`, `AI-Model`).

The metadata travels with the repo, is visible in GitHub's commit UI, and is queryable via standard `git log --grep`.

No external dependencies. No API keys. No database.

Install: `curl -fsSL https://seshmark.dev/install | bash`

GitHub: github.com/seshmark/seshmark
```

---

## LinkedIn Post

**Tone:** Professional, solves a real problem, appeals to engineering managers

```
I built a tool to solve a problem every team using AI coding assistants faces: we have no idea which AI session wrote what code.

Six months from now, when a bug appears in a file that was "written by Cursor or Claude or maybe Copilot," how do you find the original session to fix it?

Enter Seshmark.

It is a lightweight Git convention that auto-tags every AI commit with metadata — which agent, which model, which session — embedded directly in the commit message as standard trailers.

Later, you can run `git agentblame` (yes, like git blame) and see exactly which AI wrote each line.

Install: curl -fsSL https://seshmark.dev/install | bash

Open source. Works with Cursor, Claude Code, Copilot, and any tool that can set an env var before `git commit`.

Would love feedback from teams actively using AI assistants in production.

#ai #devtools #softwareengineering #github #cursor #claude
```

**Follow-up comment (if post gets traction):**

```
Update: I also built a VS Code extension that shows agent labels inline in the editor gutter — [cursor], [claude], etc. Hover for model and session details.

Install: seshmark.dev/vscode
```

---

## GitHub Release Notes (v0.1.0)

```
## Seshmark v0.1.0

**The problem:** AI writes your code. You commit it. Six months later, you have no idea which session produced it.

**The solution:** A Git-native convention that auto-tags AI commits with session metadata, then lets you query and blame later.

### Install
```bash
curl -fsSL https://seshmark.dev/install | bash
```

### Try it
```bash
git agentblame src/auth.ts
seshmark query --agent cursor
seshmark who HEAD
```

### Features
- `git agentblame` — show AI attribution per line
- `seshmark query` — search commits by agent, model, path, date
- `seshmark context` — reconstruct session context for orchestrators
- `seshmark resume` — jump back into the original AI session
- Branch inference — `claude/feature` auto-tagged without config
- VS Code extension — inline gutter labels with hover details
- Tool agnostic — works with Cursor, Claude Code, Copilot, any AI tool
- Zero infrastructure — no APIs, no database, no vendor lock-in

### The Convention
```
Seshmark-Version: 1.0.0
AI-Session: cursor:chat-abc123
AI-Agent: cursor
AI-Model: claude-sonnet-4-20250514
```

Read the full spec: seshmark.dev
```

---

## HN Show HN (If You Get Someone to Post It)

**Title:** Show HN: git agentblame — like git blame, but for AI sessions

**Body:**
```
I built a Git-native convention for tracking which AI wrote each line of code.

Install:
curl -fsSL https://seshmark.dev/install | bash

Then:
git agentblame src/auth.ts
# [cursor|abc123|cursor] import { oauth } from './oauth';
# [human|       |       ] const x = 1;

It auto-tags commits by reading env vars your AI tool sets (or infers from branch names like `claude/feature`). No config needed after install.

Open source, MIT. Would love feedback.
```

---

## Email to Engineering Managers / CTOs

**Subject:** Do you know which AI wrote your production code?

```
Hi [Name],

Quick question: if a critical bug appears in a file that was written by Cursor or Claude 6 months ago, can your team find the exact AI session to fix it?

Most teams can't. They commit AI-generated code with zero attribution.

I built Seshmark to solve this. It is a lightweight Git convention that auto-tags every AI commit with agent, model, and session ID — embedded directly in commit messages as standard trailers.

Later, you can:
- Run `git agentblame` to see which AI wrote each line
- Query your entire codebase by agent (`seshmark query --agent cursor`)
- Reconstruct session context for orchestrators and compliance

Zero infrastructure. No API keys. No vendor lock-in. Open source.

Install: https://seshmark.dev
GitHub: https://github.com/seshmark/seshmark

Would you be open to a 10-minute call to see if this fits your team's workflow?

Best,
[Your name]
```

---

## All Links (After You Set Up)

| Destination | URL |
|-------------|-----|
| Landing page | `https://seshmark.dev` or `https://YOURNAME.github.io/seshmark` |
| GitHub repo | `https://github.com/seshmark/seshmark` or `https://github.com/YOURNAME/seshmark` |
| VS Code extension | `https://seshmark.dev/vscode` (redirect to marketplace listing) |
| Install script | `https://seshmark.dev/install.sh` |
| Docs | `https://github.com/seshmark/seshmark/blob/main/README.md` |

---

*Copy-paste and ship.*
