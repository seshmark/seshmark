---
date: 2026-05-16T21:01:58+0530
author: Sahil Choudhary
commit: 16051bb
branch: main
repository: seshmark
topic: "GitHub Action — AI Report on PRs (seshmark/agentblame)"
tags: [plan, github-action, marketplace, monetization]
status: ready
parent: "thoughts/shared/designs/2026-05-16_21-01-58_github-action-report.md"
last_updated: 2026-05-16T21:01:58+0530
last_updated_by: Sahil Choudhary
---

# GitHub Action — AI Report on PRs (seshmark/agentblame)

## Overview

Build and publish a GitHub Action at `seshmark/agentblame` that runs on every Pull Request, analyzes commits for Seshmark AI trailers, and posts a formatted AI attribution report as a PR comment. Free for public repos. Paid tiers for private repos with policy enforcement.

The Action is a thin TypeScript wrapper (~200 lines) around the seshmark CLI — it installs the CLI, runs `seshmark stats` and `seshmark query`, formats the JSON output into a markdown report, and posts it via the GitHub API.

## Desired End State

- Repository `github.com/seshmark/agentblame` exists with the Action code
- `action.yml` defines inputs (token, show-models, show-sessions, fail-on-unknown-agent)
- `src/main.ts` implements: install CLI → run stats → format report → post/update PR comment
- `dist/index.js` is bundled and checked in
- `README.md` documents usage with copy-paste YAML example
- CI workflow tests the action on every push
- Published to GitHub Marketplace (free tier)
- CLI repo's README references the Action in a "GitHub App" section

## What We're NOT Doing

- NOT building policy enforcement (merge blocking, `.github/seshmark.yml` config) — that's Phase 2 (Enterprise tier)
- NOT setting up paid Marketplace billing — that's after free tier gains traction
- NOT building a cross-repo dashboard — not needed until 100+ repos use the action
- NOT duplicating seshmark CLI logic — the Action shells out to the real CLI

## Phase 1: Repository + Scaffolding

### Overview
Create the `seshmark/agentblame` repo with action.yml, package.json, tsconfig.json, and initial file structure.

### Changes Required:

#### 1. Create repo and local setup
**Repo**: `github.com/seshmark/agentblame` (create via GitHub UI)
**Local**: Clone, add initial files

```bash
# Create on GitHub first, then:
git clone git@github.com:seshmark/agentblame.git
cd agentblame
```

#### 2. Create `action.yml`
**File**: `action.yml`

```yaml
name: 'Seshmark Agent Blame'
description: 'Auto-generated AI attribution report for every Pull Request'
author: 'seshmark'
branding:
  icon: 'git-commit'
  color: 'green'

inputs:
  token:
    description: 'GitHub token'
    required: true
    default: ${{ github.token }}
  show-models:
    description: 'Include model names in report'
    required: false
    default: 'true'
  show-sessions:
    description: 'Include session IDs in report'
    required: false
    default: 'false'
  fail-on-unknown-agent:
    description: 'Fail check if unknown AI agent detected'
    required: false
    default: 'false'

outputs:
  ai-percent:
    description: 'Percentage of AI commits in PR'
  agents-found:
    description: 'Comma-separated list of agents found'
```

#### 3. Create `package.json`
**File**: `package.json`

```json
{
  "name": "seshmark-agentblame",
  "version": "1.0.0",
  "description": "AI attribution report for GitHub PRs",
  "main": "dist/index.js",
  "scripts": {
    "build": "ncc build src/main.ts -o dist",
    "package": "npm run build && git add dist"
  },
  "dependencies": {
    "@actions/core": "^1.10.0",
    "@actions/exec": "^1.1.1",
    "@actions/github": "^6.0.0"
  },
  "devDependencies": {
    "@vercel/ncc": "^0.38.0",
    "typescript": "^5.3.0"
  }
}
```

#### 4. Create `tsconfig.json`
**File**: `tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "es2020",
    "module": "commonjs",
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": false
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

#### 5. Create `src/` directory
```bash
mkdir src
```

### Success Criteria:

#### Automated Verification:
- [ ] `npm install` succeeds
- [ ] `npx tsc --noEmit` passes (TypeScript compiles)
- [ ] `npm run build` produces `dist/index.js`

#### Manual Verification:
- [ ] `action.yml` is valid (GitHub validates on push)
- [ ] Repo `seshmark/agentblame` exists on GitHub

---

## Phase 2: Action Logic (src/main.ts)

### Overview
Write the core TypeScript logic: install seshmark CLI, run stats + query, format as markdown, post/update PR comment.

### Changes Required:

#### 1. Create `src/main.ts`
**File**: `src/main.ts`

```typescript
import * as core from '@actions/core';
import * as github from '@actions/github';
import { exec } from '@actions/exec';
import * as fs from 'fs';

interface ReportData {
  total_commits: number;
  ai_commits: number;
  human_commits: number;
  ai_percent: number;
  by_agent: Record<string, number>;
  by_model: Record<string, number>;
  commits: Array<{
    hash: string;
    author: string;
    ai_agent: string | null;
    ai_model: string | null;
    ai_session: string | null;
    subject: string;
  }>;
  files: Array<{ path: string; ai_percent: number }>;
}

async function execAndCapture(cmd: string, args: string[]): Promise<string> {
  let stdout = '';
  await exec(cmd, args, {
    listeners: { stdout: (data) => { stdout += data.toString(); } },
  });
  return stdout;
}

async function installSeshmark(): Promise<void> {
  await exec('bash', ['-c', 'curl -fsSL https://seshmark.github.io/seshmark/install.sh | bash']);
  core.addPath(`${process.env.HOME}/.local/bin`);
}

async function getReportData(baseSha: string, headSha: string): Promise<ReportData> {
  const statsOutput = await execAndCapture('seshmark', [
    'stats', '--since', baseSha, '--until', headSha, '--format', 'json'
  ]);
  const stats = JSON.parse(statsOutput);

  const queryOutput = await execAndCapture('seshmark', [
    'query', '--since', baseSha, '--until', headSha, '--format', 'json'
  ]);
  const commits = JSON.parse(queryOutput);

  return {
    total_commits: stats.total_commits || 0,
    ai_commits: stats.ai_commits || 0,
    human_commits: stats.human_commits || 0,
    ai_percent: stats.ai_percent || 0,
    by_agent: stats.by_agent || {},
    by_model: stats.by_model || {},
    commits: commits || [],
    files: stats.by_file || [],
  };
}

function formatReport(data: ReportData): string {
  const showModels = core.getInput('show-models') === 'true';
  const showSessions = core.getInput('show-sessions') === 'true';

  let md = `## \u{1F916} Seshmark AI Report\n\n`;
  md += `| Metric | Value |\n|--------|-------|\n`;
  md += `| **Total Commits** | ${data.total_commits} |\n`;
  md += `| **AI Commits** | ${data.ai_commits} (${data.ai_percent.toFixed(1)}%) |\n`;
  md += `| **Human Commits** | ${data.human_commits} |\n`;

  if (Object.keys(data.by_agent).length > 0) {
    const agents = Object.entries(data.by_agent).map(([k, v]) => `${k} (${v})`).join(', ');
    md += `| **Agents Used** | ${agents} |\n`;
  }
  if (showModels && Object.keys(data.by_model).length > 0) {
    md += `| **Models** | ${Object.keys(data.by_model).join(', ')} |\n`;
  }
  md += '\n';

  if (data.commits.length > 0) {
    md += `### Commits in this PR\n\n| Commit | Author | Agent |`;
    if (showModels) md += ` Model |`;
    if (showSessions) md += ` Session |`;
    md += `\n|--------|--------|-------|`;
    if (showModels) md += `-------|`;
    if (showSessions) md += `---------|`;
    md += '\n';

    for (const c of data.commits) {
      md += `| ${c.hash.substring(0, 7)} | @${c.author} | ${c.ai_agent || 'human'} |`;
      if (showModels) md += ` ${c.ai_model || '\u2014'} |`;
      if (showSessions) md += ` ${c.ai_session ? c.ai_session.substring(0, 20) + '...' : '\u2014'} |`;
      md += '\n';
    }
    md += '\n';
  }

  const aiFiles = data.files.filter((f: any) => f.ai_percent > 0).slice(0, 10);
  if (aiFiles.length > 0) {
    md += `### Files with AI Code\n\n| File | AI % |\n|------|------|\n`;
    for (const f of aiFiles) {
      md += `| ${f.path} | ${f.ai_percent.toFixed(0)}% |\n`;
    }
    md += '\n';
  }

  md += `---\n*Report generated by [Seshmark](https://seshmark.dev)*`;
  return md;
}

async function postComment(body: string): Promise<void> {
  const token = core.getInput('token');
  const octokit = github.getOctokit(token);
  const context = github.context;
  const prNumber = context.issue.number;

  // Find existing Seshrark comment to update
  const { data: comments } = await octokit.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: prNumber,
  });

  const existing = comments.find(c =>
    c.user?.login === 'github-actions[bot]' &&
    c.body?.includes('Seshmark AI Report')
  );

  if (existing) {
    await octokit.rest.issues.updateComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      comment_id: existing.id,
      body,
    });
  } else {
    await octokit.rest.issues.createComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: prNumber,
      body,
    });
  }
}

async function run(): Promise<void> {
  try {
    const pr = github.context.payload.pull_request;
    const baseSha = pr?.base?.sha;
    const headSha = pr?.head?.sha;

    if (!baseSha || !headSha) {
      core.setFailed('Could not determine PR commits');
      return;
    }

    await installSeshmark();
    const data = await getReportData(baseSha, headSha);
    const report = formatReport(data);
    await postComment(report);

    core.setOutput('ai-percent', data.ai_percent.toFixed(1));
    core.setOutput('agents-found', Object.keys(data.by_agent).join(', '));

    if (core.getInput('fail-on-unknown-agent') === 'true') {
      if (data.by_agent['unknown'] || data.by_agent['local']) {
        core.setFailed('PR contains commits from unapproved AI agents');
      }
    }
  } catch (error) {
    core.setFailed(error instanceof Error ? error.message : 'Unknown error');
  }
}

run();
```

### Success Criteria:

#### Automated Verification:
- [ ] TypeScript compiles: `npx tsc --noEmit`
- [ ] Build succeeds: `npm run build`
- [ ] `dist/index.js` is non-empty

#### Manual Verification:
- [ ] Run action locally with `act` or review the TypeScript logic for correctness

---

## Phase 3: README + License

### Overview
Write the Action's README with copy-paste YAML usage, input reference, and examples.

### Changes Required:

#### 1. Create `README.md`
**File**: `README.md`

```markdown
# Seshmark Agent Blame — AI Report for PRs

> Auto-generated AI attribution report on every Pull Request. See which AI agents wrote what, which models were used, and which files have the most AI code.

## Usage

Add this workflow to any repo:

```yaml
# .github/workflows/seshmark.yml
name: Seshmark AI Report
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  report:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
      contents: read
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: seshmark/agentblame@v1
```

That's it. Every PR now gets an AI attribution report.

## Example Output

![AI Report Screenshot](docs/screenshot.png)

| Metric | Value |
|--------|-------|
| AI Commits | 10 of 12 (83%) |
| Agents Used | cursor (7), claude (2), copilot (1) |
| Models | claude-sonnet-4, gpt-4o |
| Files with AI Code | src/auth.ts (100%), src/api.ts (78%) |

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `token` | Yes | `${{ github.token }}` | GitHub token for posting comments |
| `show-models` | No | `true` | Include model names in the report |
| `show-sessions` | No | `false` | Include session IDs in the report |
| `fail-on-unknown-agent` | No | `false` | Fail the check if an unapproved AI agent is detected |

## Outputs

| Output | Description |
|--------|-------------|
| `ai-percent` | Percentage of AI commits in this PR |
| `agents-found` | Comma-separated list of AI agents detected |

## How It Works

1. Installs the seshmark CLI
2. Runs `seshmark stats` and `seshmark query` on the PR's commits
3. Formats the data into a markdown table
4. Posts (or updates) a comment on the PR

Works with any AI coding tool that uses Seshmark's commit trailers: Cursor, Claude Code, Copilot, Pi, Aider, OpenCode, and more.

## License

MIT
```

#### 2. Create `LICENSE`
**File**: `LICENSE`

```
MIT License

Copyright (c) 2026 Seshmark Contributors

Permission is hereby granted...
```

### Success Criteria:

#### Automated Verification:
- [ ] `README.md` exists and contains `uses: seshmark/agentblame@v1`
- [ ] `LICENSE` exists

#### Manual Verification:
- [ ] README passes the 5-second test: what is this, why should I care, how do I start

---

## Phase 4: CI Workflow

### Overview
Add a CI workflow that builds and validates the action on every push.

### Changes Required:

#### 1. Create `.github/workflows/ci.yml`
**File**: `.github/workflows/ci.yml`

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci
      - run: npm run build
      - run: node -e "require('./dist/index.js')" 2>&1 | head -5 || true

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - run: npm ci
      - run: npm run build
      - name: Test the action in a dummy PR
        env:
          GITHUB_TOKEN: ${{ github.token }}
        run: |
          # Create a test commit with AI trailers
          git config user.email "test@test.com"
          git config user.name "Test"
          echo "test" > f.txt && git add f.txt
          SESHMARK_SESSION_ID="cursor:test" SESHMARK_AGENT="cursor" git commit -m "test: ai commit"
          # Run the action locally (dry-run)
          node -e "
            const main = require('./dist/index.js');
            console.log('Action loaded successfully');
          "
```

### Success Criteria:

#### Automated Verification:
- [ ] CI workflow runs on push and PR
- [ ] Build step succeeds
- [ ] Action loads without runtime errors

---

## Phase 5: Publish to GitHub Marketplace

### Overview
Publish the action to the GitHub Marketplace as a free public action.

### Changes Required:

#### 1. Create a release on GitHub
```bash
# From the agentblame repo
git tag v1.0.0
git push origin v1.0.0

# Create a release — the "Publish to Marketplace" checkbox appears
gh release create v1.0.0 \
  --title "v1.0.0 — AI attribution reports for every PR" \
  --notes "See which AI agents wrote code in every Pull Request."
```

#### 2. Add reference in CLI repo's README
**File**: `README.md` (in seshmark/seshmark)

```markdown
## GitHub Action

Add AI attribution reports to every Pull Request:

```yaml
- uses: seshmark/agentblame@v1
```

See the [action repo](https://github.com/seshmark/agentblame) for full docs.
```

### Success Criteria:

#### Manual Verification:
- [ ] Marketplace listing is live at `github.com/marketplace/actions/seshmark-agent-blame`
- [ ] Action is installable: `uses: seshmark/agentblame@v1`
- [ ] CLI repo README links to the action

---

## Phase 6: Paid Tiers (Future)

### Overview
Add Team ($4/seat/month) and Enterprise ($10/seat/month) tiers with policy enforcement, merge blocking, and audit export. Only build after free tier has traction (50+ repos using it).

### Changes Required:

#### 1. Add policy config support to `src/main.ts`
```typescript
// Read .github/seshmark.yml for policies
// Apply allowed_agents, blocked_agents, require_review_for_ai thresholds
```

#### 2. Separate free vs paid logic
```typescript
// Check GitHub App installation plan metadata
// Free: basic report
// Team: + labels, model breakdown
// Enterprise: + policy enforcement, merge blocking, audit export
```

### Success Criteria:

#### Manual Verification:
- [ ] `.github/seshmark.yml` policies are enforced
- [ ] Pull Request fails check when policy is violated
- [ ] Enterprise audit export produces valid JSON/CSV

---

## Testing Strategy

### Automated:
- `npm run build` — TypeScript compiles and bundles
- CI workflow — builds action and validates it loads
- `go test ./...` in CLI repo — ensures CLI output format hasn't changed

### Manual Testing Steps:
1. Create a test repo with `.github/workflows/seshmark.yml`
2. Create a branch, add AI-tagged commits
3. Open a PR — verify the comment appears with correct data
4. Push more commits — verify the comment updates
5. Test with `show-sessions: true` and `show-models: false`

## Performance Considerations

- Action installs seshmark CLI on every run (~5 seconds). Cache the binary if runs are frequent.
- For PRs with 100+ commits, stats query may take 2-3 seconds. Acceptable for CI.
- Comment updates use the existing comment to avoid spamming the PR.

## Migration Notes

- The action is a new repo separate from the CLI. No migration needed.
- Users who already have seshmark CLI installed don't need to change anything — the Action is optional.
- The Action always installs the latest seshmark version. No version conflicts.

## References

- Design: `thoughts/shared/designs/2026-05-16_21-01-58_github-action-report.md`
- Monetization plan: `MONETIZATION_PLAN.md`
- Full spec (backup): `seshmark-plans-backup/GITHUB_APP_SPEC.md`
