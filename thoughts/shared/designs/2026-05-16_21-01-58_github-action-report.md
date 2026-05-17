# GitHub Action — AI Report on PRs

> Design for the Seshmark GitHub Action that posts AI attribution reports on Pull Requests.

---

## Overview

Monetize the CLI by building a GitHub Action that runs on every Pull Request, analyzes commits for Seshmark trailers, and posts an AI Attribution Report as a PR comment. Free for public repos; paid for private repos with policy enforcement.

Based on the monetization plan at `MONETIZATION_PLAN.md` and the backup spec at `seshmark-plans-backup/GITHUB_APP_SPEC.md`.

## Desired End State

- A GitHub Action at `seshmark/report-action` (separate repo) that users install via one YAML file
- On every PR, it posts a formatted AI report comment with stats, commit breakdown, and file-level analysis
- Published to GitHub Marketplace with free tier for public repos
- Paid tiers (Team, Enterprise) sold via GitHub Marketplace billing

## Architecture

### File Structure (separate repo: `seshmark/report-action`)

```
seshmark-report-action/
├── action.yml          # Action metadata (inputs, outputs, branding)
├── README.md           # Action docs
├── LICENSE             # MIT
├── src/
│   └── main.ts         # TypeScript entry point (~200 lines)
├── dist/
│   └── index.js        # Compiled bundle (checked in)
├── package.json
├── tsconfig.json
└── .github/
    └── workflows/
        └── ci.yml      # Test the action
```

### How It Works

1. User adds `.github/workflows/seshmark.yml` to their repo
2. On every PR, the Action: — bash: `curl -fsSL https://seshmark.dev/install | bash` — runs `seshmark stats --format json` — runs `seshmark query --format json` — formats JSON into a PR comment — posts/updates the comment
3. No server, no database, no external API call

### Key Design Decisions

- **Separate repo from CLI**: GitHub Actions need to be published as `{owner}/{action-name}` repos. Installable via `uses: seshmark/report-action@v1`
- **CLI is the engine**: The Action installs the real seshmark CLI and shells out to it. This keeps the Action as a thin ~200 line TypeScript wrapper, not a reimplementation
- **No cloud**: The Action runs entirely inside the GitHub Actions runner. No data leaves the repo

## File Map

| File | Purpose | Source |
|------|---------|--------|
| `action.yml` | Action metadata | New |
| `src/main.ts` | Main logic | ~200 lines from backup spec |
| `package.json` | Dependencies | New |
| `tsconfig.json` | TypeScript config | New |
| `README.md` | Action docs | New |
| `LICENSE` | MIT | Copy |
| `.github/workflows/ci.yml` | Test the action | New |

## Scope

### Building
- TypeScript Action that installs CLI, runs stats, formats report, posts comment
- README with usage and examples
- CI workflow that tests the action
- Published to GitHub Marketplace (free tier)

### Not Building (Phase 2)
- Policy enforcement (merge blocking, .github/seshmark.yml config) — Enterprise tier later
- Paid Marketplace tiers — after free tier gains traction
- Cross-repo dashboard — not needed until 100+ repos

## Architecture Decisions

1. **TypeScript over Go**: GitHub Actions ecosystem is Node/TypeScript-native. The Go CLI is already cross-compiled — the Action just shells out to it.
2. **ncc bundling**: Bundle TypeScript into a single `dist/index.js` with `@vercel/ncc`. No runtime dependencies for users.
3. **`exec` over `toolkit`**: Use `@actions/exec` to run seshmark CLI commands rather than reimplementing parsing logic.
4. **Comment update, not spam**: The Action finds any existing Seshmark comment on the PR and updates it, rather than posting a new one every time.
