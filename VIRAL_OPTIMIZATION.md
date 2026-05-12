# Seshmark — Viral Optimization Playbook

## The One Rule

> People do not share tools. People share screenshots of tools making them feel smart.

Your entire strategy: make the output so visually striking that developers screenshot it and post it unprompted.

---

## Tier 1: Do These Before Launch (Highest Impact)

### 1. The Hero GIF (Non-Negotiable)

**What:** A 6–10 second screen recording showing the full loop.

**Script:**
1. Cursor writes code (2s)
2. `git commit -m "feat: auth"` (1s)
3. `git agentblame src/auth.ts` (2s)
4. Output appears with colors — `[cursor]`, `[claude]`, `[human]` (3s)
5. Cursor hovers over `[cursor]` label, tooltip shows session ID (2s)

**Why this works:** The entire value proposition is communicated in 6 seconds. No reading required.

**Where it goes:**
- Top of README (above all text)
- Top of `index.html` (full-width, above the fold)
- Twitter post (native video, not a link)
- Reddit post (linked in the body)

**Tool to make it:**
- macOS: `⌘+Shift+5` → screen recording, then compress with `ffmpeg -i input.mov -vf "fps=30,scale=1200:-1:flags=lanczos" -c:v libx264 -crf 28 -preset fast -an -movflags +faststart output.mp4`
- Or use CleanShot X / Screen Studio for polished recordings

**Rules:**
- Dark terminal (GitHub Dark or Dracula theme)
- Large font (18px+)
- No typing mistakes
- No pauses longer than 1 second
- No audio needed

---

### 2. The Demo Repo (Instant Gratification)

**What:** A separate repo `seshmark/demo` with pre-tagged commits that anyone can clone and run `git agentblame` on immediately.

**Why this works:** 90% of people who see your post will not install anything. But they will clone a demo repo and run one command. If the output is beautiful, they screenshot it and share.

**Contents:**
```bash
git clone https://github.com/seshmark/demo
cd demo
git agentblame src/api.ts
```

The repo should contain:
- A fake `src/api.ts` with 30 lines
- 5 pre-tagged commits mixing human, Cursor, Claude, and Copilot
- A README that says: "Run `git agentblame src/api.ts` to see the magic"

**Critical:** The output must look stunning. The blame output is your advertisement.

---

### 3. Beautiful Blame Output (The Screenshot Engine)

**Current output:**
```
[cursor|chat-abc123|cursor] import { oauth } from './oauth';
```

**Viral output:**
```
  1  [human]           const express = require('express');
  2  [cursor]          import { oauth } from './oauth';
  3  [claude]          function validateToken(token: string) {
  4  [human]             if (!token) return null;
  5  [copilot]           const hash = crypto.sha256(token);
```

**Changes:**
- Drop the session ID from the default text view (too noisy for screenshots)
- Use emojis or color blocks instead of bracket text where possible
- Pad everything so columns align perfectly
- Add a header row

**Implementation:** Modify `internal/cmd/blame.go` to produce a cleaner default output. Keep `--verbose` or `--full` for the session ID.

---

### 4. Add `seshmark stats` (The Shareable Chart)

**What:** A command that generates a shareable summary of AI usage in the repo.

```bash
$ seshmark stats

AI Code Report — my-project
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total commits:          1,247
AI commits:               843 (68%)
Human commits:            404 (32%)

By Agent:
  cursor        ████████████████████  482 (57%)
  claude        ████████████          241 (29%)
  copilot       ███                    97 (11%)
  codex         █                      23 (3%)

By Model:
  claude-sonnet-4     ████████████████  340
  gpt-4o              ██████████       210
  qwen2.5-coder       ██████           120

By File (most AI):
  src/auth.ts         ████████████████████  94%
  src/payments.ts     ██████████████        78%
  src/db.ts           ██                    12%

First AI commit:        2024-11-03
Latest AI commit:       2025-05-12

Run `seshmark stats --format json` for machine-readable output.
```

**Why this works:**
- Every engineering manager who sees this will screenshot it and post it in Slack
- It answers the question: "How much of our code is AI-written?" (the #1 question right now)
- It makes teams look data-driven for sharing it
- `--format markdown` lets them paste it directly into PR descriptions

**Implementation:** New file `internal/cmd/stats.go`. Runs `git log --grep="Seshmark-Version:"` and aggregates data. Use simple ASCII bars (20 chars wide, proportional fill).

---

### 5. Shields.io Badges on README (Trust Signals)

Add these to the top of README.md, immediately under the title:

```markdown
![Version](https://img.shields.io/github/v/release/seshmark/seshmark?color=3fb950&label=version)
![License](https://img.shields.io/badge/license-MIT-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen)
```

**Why:** Developers subconsciously skip repos without badges. Badges signal "this is maintained, tested, and real."

---

### 6. README Reconstruction

**Current README structure:** Paragraphs first.

**Viral README structure:**

```markdown
# Seshmark — Agent-Aware Git History

> `git agentblame` — like `git blame`, but shows which AI (Cursor, Claude, Copilot) wrote each line.

[HERO GIF HERE — full width, above everything]

```bash
curl -fsSL https://seshmark.dev/install | bash
```

## 10-Second Demo

```bash
git clone https://github.com/seshmark/demo
cd demo
git agentblame src/api.ts
```

[SCREENSHOT OF OUTPUT]

## What It Does
...
```

**Rules:**
- The first 300 pixels of the README must contain the GIF + install command
- No paragraphs above the GIF
- The install command must be copy-pasteable with one click
- A screenshot of the output must appear within 2 scrolls

---

## Tier 2: Do These Within 2 Weeks

### 7. VS Code Extension: Colored Backgrounds in Gutter

**Current extension:** Text labels `[cursor]` after lines.

**Viral extension:** Colored background pills in the gutter, like GitHub issue labels.

```
1 | const express = require('express');       |          |
2 | import { oauth } from './oauth';          | [cursor] |  ← small colored pill
3 | function validateToken(token) {           | [claude] |  ← different color
4 |   if (!token) return null;                |          |
```

**Visual spec:**
- Background color per agent (Cursor = orange, Claude = terracotta, etc.)
- White text, rounded corners, 8px font
- Appears only in the gutter (right side, not inline with code)
- Hover = full tooltip with session, model, commit

**Why this works:** It turns every file in your codebase into a heatmap. Developers will open random files just to see the colors. That is addictive.

**Implementation:** In `vscode-extension/src/extension.ts`, use `before` decoration instead of `after`, with `backgroundColor` and `borderRadius`.

---

### 8. Add `seshmark export` (PR Report Generator)

```bash
seshmark export --since "last week" --format markdown
```

Generates:
```markdown
## AI Code Report (May 6–12)

| File | AI % | Agents |
|------|------|--------|
| src/auth.ts | 94% | cursor, claude |
| src/payments.ts | 78% | cursor |
| src/db.ts | 12% | copilot |

**Sessions to review:**
- `cursor:chat-abc123` — 12 commits, touched auth, payments
- `claude:thread-xyz` — 4 commits, touched auth

*Generated by [Seshmark](https://seshmark.dev)*
```

**Why:** Engineering leads paste this into weekly standup docs. Every paste is free marketing.

---

### 9. Landing Page: Replace Fake Terminal with Real Video

**Current `index.html`:** Fake CSS terminal.

**Viral `index.html`:**
- Full-width hero video (autoplay, muted, loop)
- The video is the GIF from #1, but embedded as MP4
- Below the video: one install command
- Below that: 3 feature cards
- No scrolling required to understand the product

**Why:** Videos autoplaying in the browser get 3x the engagement of static images.

---

### 10. The "Trending on GitHub" Boost

GitHub's trending algorithm considers:
- Stars per day (velocity matters more than total)
- Forks
- Traffic from external sources

**To game it:**
1. Post on Twitter + Reddit on a **Tuesday at 9am ET** (optimal for developer attention)
2. Ask 5 friends to star it in the first 2 hours (initial velocity signals GitHub's algo)
3. Pin a tweet the day before saying "Launching tomorrow" (builds anticipation)
4. Cross-post to Dev.to the same day (Dev.to articles rank well on Google)

---

## Tier 3: Advanced (Do These If Traction Hits)

### 11. The Hacker News Surrogate Strategy

Since you cannot post Show HN:
1. Find someone with HN karma (friend, Twitter mutual, colleague)
2. Send them the exact title and body from `LAUNCH_COPY.md`
3. Ask them to post at **Tuesday 9am ET**
4. Be the first comment with technical details

**If nobody you know has karma:**
- Comment on existing AI-related HN threads: "We built something related — seshmark.dev — auto-tags AI commits so you can blame them later"
- Do not spam. One relevant comment per week.

### 12. GitHub Topics (SEO)

Add these topics to your repo (Settings → Topics):
```
ai, git, cli, developer-tools, cursor, claude, copilot, code-review, attribution, machine-learning, openai, vscode-extension
```

**Why:** GitHub search and Google both index these.

### 13. The "One Weird Trick" Tweet Thread

A 5-tweet thread that teaches something, with Seshmark as the punchline:

```
Tweet 1/5: Most teams using AI coding tools have a secret problem.
They have no idea which AI wrote what code.

Tweet 2/5: In 6 months, when a bug appears in a file that was "written by Cursor or Claude," nobody can find the original session.

Tweet 3/5: Git commit messages are the answer. We just need to agree on a convention.

Tweet 4/5: I built Seshmark — it auto-tags every AI commit with the agent, model, and session ID.

Tweet 5/5: Then you run `git agentblame` and see exactly which AI wrote each line.

Install: curl -fsSL https://seshmark.dev/install | bash
```

**Why:** Education-first threads get 5x the engagement of product announcements.

### 14. Influencer Seeding (Low Effort, High Variance)

DM or reply to these account types on Twitter:
- AI tool account managers (`@cursor`, `@anthropicAI`, `@github`)
- Dev tool reviewers (`@sama` sometimes, `@paulg` occasionally)
- Newsletter writers (Console.dev, TLDR, Pointer.io)

**The message:**
```
Hi [Name], I built a lightweight open-source convention for tracking 
which AI session wrote each line of code in Git.

Thought it might be relevant to your audience. Zero pressure.

https://seshmark.dev
```

No follow-up. If they are interested, they will post. If not, move on.

---

## The Checklist

Print this and check items off:

### Before Launch
- [ ] Hero GIF recorded and compressed (< 5MB)
- [ ] Demo repo created with pre-tagged commits
- [ ] `seshmark stats` command implemented
- [ ] Blame output visually optimized
- [ ] README rebuilt with GIF on top
- [ ] Shields.io badges added
- [ ] VS Code extension with colored gutter pills
- [ ] Landing page uses real video, not CSS fake terminal

### Launch Day
- [ ] GitHub repo public
- [ ] GitHub Pages enabled
- [ ] Twitter post with video (9am ET Tuesday)
- [ ] Reddit posts on r/cursor, r/LocalLLaMA, r/programming
- [ ] LinkedIn post
- [ ] Dev.to article
- [ ] Ask 5 friends to star within first 2 hours

### Week After
- [ ] Reply to every GitHub issue within 6 hours
- [ ] Quote-retweet anyone who mentions it
- [ ] Post follow-up: "I added VS Code support"
- [ ] Post follow-up: "The convention spec explained"
- [ ] Email 3 newsletters with one-line pitch

---

## The Brutal Truth

You can optimize everything above and still get 50 stars. Or you can do half of it and get 5,000. Virality is not deterministic.

But the baseline for success is this:

**If someone can clone a repo, run one command, and screenshot the output in under 30 seconds — you have a chance.**

**If they need to read docs, install dependencies, and configure paths — you do not.**

The demo repo + beautiful blame output + hero GIF are the only things that truly matter. Everything else is optimization.

---

*End of playbook.*
