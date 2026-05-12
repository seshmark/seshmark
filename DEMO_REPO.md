# Seshmark Demo Repo — SEO & Viral Optimization

## Goal

A clone-and-run demo that produces beautiful `git agentblame` output instantly. Optimized for search, sharing, and screenshots.

## Repo Setup

**Name:** `seshmark/demo` (under the org)
**Visibility:** Public
**Description:** "Clone this repo and run `git agentblame` to see Seshmark in action — AI attribution for Git commits"

### README.md (SEO-Optimized)

```markdown
# Seshmark Demo — See AI Attribution in Action

> Clone this repo and run one command to see which AI (Cursor, Claude, Copilot) wrote each line.

## Try It (10 Seconds)

```bash
git clone https://github.com/seshmark/demo.git
cd demo
git agentblame src/auth.ts
```

## What You'll See

```
  1  [human]           import express from 'express';
  2  [cursor]          import { OAuth2Client } from 'google-auth-library';
  3  [claude]          async function validateToken(token: string) {
  4  [human]             if (!token) throw new Error('Missing token');
  5  [copilot]           const decoded = jwt.verify(token, SECRET);
  6  [cursor]            return decoded.userId;
  7  [human]           }
```

Each `[agent]` label shows which AI session wrote that line. Hover in VS Code for details.

## How It Works

This repo's commit history uses the [Seshmark convention](https://seshmark.dev) — 
lightweight trailers embedded in Git commit messages:

```
Seshmark-Version: 1.0.0
AI-Session: cursor:chat-abc123
AI-Agent: cursor
AI-Model: claude-sonnet-4-20250514
```

## Install Seshmark on Your Own Repos

```bash
curl -fsSL https://seshmark.dev/install | bash
```

Works with any AI tool — Cursor, Claude Code, OpenCode, Copilot, Aider, and more.

## License

MIT
```

### GitHub Topics

```
ai, git, blame, developer-tools, cursor, claude, copilot, code-review, attribution, ai-coding
```

### File Structure (Realistic Project)

```
demo/
├── .github/
│   └── workflows/
│       └── ci.yml              # Runs git agentblame in CI (shows green checkmark)
├── src/
│   ├── auth.ts                 # 30 lines, mixed human + AI commits
│   ├── middleware.ts           # 25 lines, mostly Cursor
│   ├── db.ts                   # 40 lines, mostly human with Copilot
│   └── payment.ts              # 35 lines, Claude + Cursor
├── tests/
│   └── auth.test.ts            # Human-written
├── package.json                # Real-looking project
├── tsconfig.json
└── README.md
```

### Commit History (Pre-Tagged)

| Commit | Agent | Model | Files | Description |
|--------|-------|-------|-------|-------------|
| `a1b2c3d` | human | — | `package.json` | Initial setup |
| `e4f5g6h` | cursor | claude-sonnet | `src/auth.ts` | OAuth2 implementation |
| `i7j8k9l` | claude | claude-sonnet | `src/middleware.ts` | Auth middleware |
| `m0n1o2p` | copilot | gpt-4o | `src/db.ts` | DB connection |
| `q3r4s5t` | cursor | claude-sonnet | `src/payment.ts` | Stripe integration |
| `u6v7w8x` | human | — | `tests/auth.test.ts` | Add tests |
| `y9z0a1b` | claude | claude-sonnet | `src/auth.ts` | Fix token edge case |

### CI Workflow (`.github/workflows/ci.yml`)

```yaml
name: Demo
on: [push, pull_request]
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install seshmark
        run: curl -fsSL https://seshmark.dev/install | bash
      - name: Show AI attribution
        run: git agentblame src/auth.ts
```

**Why:** Green checkmark on every commit = trust signal. Also proves the install script works.

### SEO Optimizations

1. **Repo name includes keywords:** `demo` is generic, but the description includes "AI attribution" and "git blame"
2. **README H1:** "Seshmark Demo — See AI Attribution in Action" (contains brand + keywords)
3. **First code block:** The install command (copy-pasteable)
4. **Second code block:** The actual output (screenshot-worthy)
5. **Links back to main repo:** Passes link equity
6. **License file:** Required for GitHub search indexing
7. **Issues enabled:** So people can ask questions (activity signals)
8. **Pinned to org profile:** First repo visitors see

### Viral Optimizations

1. **Commit messages are realistic:** Not "test1", "fix" — real descriptions like "Implement OAuth2 middleware with token validation"
2. **Mixed human + AI lines:** Shows the tool handles both, makes output more interesting
3. **Multiple agents:** Cursor, Claude, Copilot all represented — broader appeal
4. **Different models:** Shows `AI-Model` trailer in action
5. **Complex enough:** 4 files, 7 commits, realistic structure — not a toy example
6. **Test file included:** Shows human-only commits too

### Screenshot Guide for Users

Add a "Share" section to README:

```markdown
## Share Your Screenshot

Run `git agentblame` on any file and screenshot the output. Tag us on Twitter @seshmarkdev.

The most beautiful screenshots get featured on our homepage.
```

---

## Checklist

- [ ] Create `seshmark/demo` repo
- [ ] Add realistic files with mixed human/AI content
- [ ] Create 7 pre-tagged commits with proper trailers
- [ ] Write SEO-optimized README
- [ ] Add GitHub topics
- [ ] Add CI workflow showing green checkmark
- [ ] Enable issues
- [ ] Pin to org profile
- [ ] Link from main `seshmark/seshmark` README
- [ ] Test clone-and-run flow end-to-end
