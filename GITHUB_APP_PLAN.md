# Seshmark GitHub App — Architecture Plan

## Overview

A GitHub App that reads `Seshmark-Version` trailers on every Pull Request and adds an automated "AI Attribution Report" comment. Free for open-source repos. Paid for private repos via GitHub Marketplace.

This is the **primary monetization path**: build once, run on GitHub's infrastructure, collect recurring revenue without ongoing sales or support.

---

## What the App Does

### On Every Pull Request Opened/Updated

1. Scans all commits in the PR for `Seshmark-Version:` trailers
2. Generates a report comment on the PR:

   ```markdown
   ## 🤖 Seshmark AI Report

   | Metric | Value |
   |--------|-------|
   | **AI Commits** | 12 of 15 (80%) |
   | **Agents Used** | cursor (8), claude (3), codex (1) |
   | **Models** | claude-sonnet-4, gpt-4o, qwen2.5-coder |
   | **Largest AI Change** | `src/auth.ts` (+245 lines, Cursor) |
   | **Human Review Required?** | ⚠️ Yes — `auth/` module has 90% AI code |

   ### Commits
   | Commit | Agent | Model | Session |
   |--------|-------|-------|---------|
   | `a3f9d2e` | cursor | claude-sonnet | `cursor:chat-abc123` |
   | `b8e1c4f` | claude | claude-sonnet | `claude:thread-xyz` |

   ---
   *Powered by [Seshmark](https://seshmark.dev)*
   ```

3. Optionally applies labels: `ai-cursor`, `ai-claude`, `ai-heavy`, `ai-mixed`, `human-only`

### Policy Enforcement (Paid Tier)

```yaml
# .github/seshmark.yml in the repo
policies:
  require_review_for_ai:
    enabled: true
    threshold: 0.5  # 50%+ AI code requires human review
  allowed_agents:
    - cursor
    - claude
    - codex
  blocked_agents:
    - unknown
    - local
  require_signed_ai_commits: true
```

If a policy is violated, the app:
- Fails the PR check (red X on merge button)
- Adds a comment explaining the violation
- Prevents merge until resolved

---

## Architecture

### Serverless (Recommended)

Use **GitHub Actions** as the runtime. No server to maintain.

```
GitHub Webhook
    │ (PR opened/updated)
    ▼
GitHub Actions workflow
    │ (triggered by App)
    ▼
Seshmark GitHub App Action
    │ (reads commits, generates report)
    ▼
PR Comment / Check / Label
```

**Why Actions:**
- Free compute for public repos
- Runs inside the user's repo (no external infra)
- GitHub handles all authentication
- You publish an Action, not a server

### Marketplace Billing

| Plan | Price | What you get |
|------|-------|--------------|
| **Free** | $0 | Public repos only. Basic report. |
| **Team** | $0.50/repo/month | Private repos. Full report + labels. |
| **Enterprise** | $2/repo/month | Policy enforcement + audit API + SSO. |

GitHub Marketplace handles all billing, invoicing, and seat management. You just get a monthly payout.

---

## Implementation

The GitHub App itself is a thin wrapper:

1. **GitHub App registration** (one-time setup)
2. **GitHub Action** (`seshmark/report-action@v1`) — the actual logic
3. **Marketplace listing** — description, pricing, screenshots

The Action code:
- Clones the PR branch
- Runs `git log --grep="Seshmark-Version"` 
- Parses trailers
- Posts comment via `github-script` action
- Optionally fails the check based on `.github/seshmark.yml`

---

## Revenue Model

### Conservative Projection

| Stage | Repos | Revenue | Timeline |
|-------|-------|---------|----------|
| Free adoption | 1,000 public repos | $0 | Month 0–6 |
| First paid users | 50 private repos @ $0.50 | $25/mo | Month 6–12 |
| Growth | 500 private repos @ $0.50 | $250/mo | Month 12–18 |
| Enterprise kick-in | 100 enterprise @ $2 + 500 team @ $0.50 | $450/mo | Month 18–24 |

**Realistic year-2 revenue: $200–$500/month.** Not life-changing, but genuinely passive. Covers the domain and a nice dinner.

### The Real Win

The GitHub App drives **CLI adoption** (developers see the report and install the tool locally). The CLI drives **brand recognition**. The brand drives **consulting and enterprise deals** if you want to actively pursue them.

The app is not the main revenue stream — it is the **funnel** that makes the tool famous.

---

## Maintenance Burden

| Task | Frequency | Effort |
|------|-----------|--------|
| Update GitHub Action version | Quarterly | 30 min |
| Respond to GitHub issues | Weekly | 1 hour |
| Update Marketplace listing | When features ship | 1 hour |
| Fix GitHub API changes | Annually | 2–4 hours |

**Total: ~2 hours/month** after initial build.

---

## Build vs. Skip

| Approach | Build Effort | Revenue Potential | Maintenance |
|----------|-------------|-------------------|-------------|
| GitHub App (Actions) | Medium | Low but real | Very low |
| SaaS Dashboard | High | Medium-High | High |
| IDE Extension Pro | Medium | Low | Medium |
| Enterprise On-Prem | Very High | Very High | Very High |
| GitHub Sponsors only | Zero | Near-zero | Zero |

**Recommendation:** Build the GitHub App Action. It is the only monetization path that is both low-effort and actually generates recurring revenue. Skip the SaaS dashboard until you have 1,000+ GitHub stars and users asking for it.

---

## Next Steps

1. Register a GitHub App at `github.com/settings/apps`
2. Create `.github/workflows/seshmark.yml` in your own repo as proof of concept
3. Publish the Action to GitHub Marketplace
4. Set pricing: free for public, $0.50/repo for private
5. Add a call-to-action in the CLI: `seshmark status` says "Enable PR reports: install the GitHub App"

---

*This is the "set it and forget it" path. Build it in a weekend. Ship it. Let GitHub handle the rest.*
