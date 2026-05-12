# Seshmark Skill System

## Overview

Skills are capabilities that **AI coding agents/harnesses** use to interact with Seshmark data. They live in the agent's context, not in Seshmark itself.

**Seshmark stays lightweight.** The skill system lets harness builders give their agents the ability to query, track, and resume sessions — without adding any code to Seshmark.

## How It Works

1. A harness (Cursor, Claude Code, etc.) loads the Seshmark skill into its context
2. The agent can now invoke Seshmark commands as tools
3. Seshmark CLI handles the actual Git operations
4. Results flow back to the agent for decision-making

## Skill Manifest

A skill is a YAML file that describes what the agent can do with Seshmark:

```yaml
name: seshmark
description: Query and manage AI session metadata in Git history
author: seshmark
version: "1.0.0"

tools:
  - name: blame
    description: Show AI attribution for each line in a file
    command: "seshmark blame {file} --format json"
    parameters:
      file: string
    
  - name: query
    description: Search AI-tagged commits
    command: "seshmark query {filters} --format json"
    parameters:
      agent: string?
      model: string?
      path: string?
      since: string?
      terms: [string]
    
  - name: who
    description: Show AI metadata for a commit
    command: "seshmark who {commit} --format json"
    parameters:
      commit: string  # default: HEAD
    
  - name: context
    description: Reconstruct session context for a commit
    command: "seshmark context {commit} --format prompt"
    parameters:
      commit: string
    
  - name: resume
    description: Get session info for resuming work
    command: "seshmark resume {commit} --format json"
    parameters:
      commit: string
    
  - name: track
    description: Start tracking a new session
    command: "seshmark track {session_id} --agent {agent} --model {model}"
    parameters:
      session_id: string
      agent: string
      model: string?
    
  - name: note
    description: Add a continuation note to current session
    command: "seshmark note {text}"
    parameters:
      text: string
    
  - name: stats
    description: Generate AI code usage report
    command: "seshmark stats --format json"
    parameters: {}
```

## How Harnesses Use Skills

### Example 1: Claude Code with Seshmark Skill

Claude Code loads the skill into its system prompt:

```
You have access to the Seshmark skill for tracking AI sessions in Git.

Available tools:
- seshmark_blame(file): Show which AI wrote each line
- seshmark_query(filters): Find commits by agent/model/date
- seshmark_who(commit): Show commit metadata
- seshmark_context(commit): Reconstruct session context
- seshmark_track(session_id, agent, model): Start tracking
- seshmark_note(text): Add a continuation note

Use these tools when:
- The user asks "who wrote this?"
- The user wants to continue previous work
- The user asks about AI usage stats
- Before committing, to check if a session is active
```

**User:** "Who wrote the auth middleware?"

**Claude:**
1. Calls `seshmark_blame("src/auth.ts")`
2. Gets JSON: `{line: 42, ai_agent: "cursor", ai_session: "cursor:chat-abc123"}`
3. Calls `seshmark_context("a3f9d2e")`
4. Gets continuation prompt
5. Responds: "The auth middleware was written by Cursor (session `cursor:chat-abc123`) on May 7. The session focused on OAuth2 implementation with token validation. Would you like me to resume that session?"

### Example 2: Cursor with Seshmark Skill

Cursor's agent sees the skill and can:

```javascript
// Before committing, Cursor checks if a session is active
const status = await runTool("seshmark_status");
if (!status.active_session) {
  await runTool("seshmark_track", {
    session_id: "cursor:" + currentChatId,
    agent: "cursor",
    model: currentModel
  });
}
```

### Example 3: OpenCode with Seshmark Skill

```python
# In the agent's reasoning loop
if user_input.contains("who wrote"):
    result = seshmark_blame(file=extracted_file)
    return format_blame_result(result)

if user_input.contains("continue"):
    commit = seshmark_query(terms=["auth", "oauth"], limit=1)
    context = seshmark_context(commit=commit.hash)
    return f"Resuming session {context.session_id}. Previous work: {context.continuation_prompt}"
```

## Skill Registration

### For Harness Builders

Add the Seshmark skill to your harness's skill registry:

```json
{
  "skills": [
    {
      "name": "seshmark",
      "source": "https://github.com/seshmark/seshmark/blob/main/skills/seshmark.yaml",
      "version": "1.0.0",
      "auto_enable": true
    }
  ]
}
```

### For Users

Users enable the skill in their harness config:

```yaml
# ~/.cursor/skills.yaml
skills:
  - name: seshmark
    enabled: true

# ~/.claude/skills.yaml  
skills:
  - name: seshmark
    enabled: true
    # Optional: customize which tools are available
    tools:
      - blame
      - query
      - context
```

## Built-in Skills

| Skill | Description | For Harness |
|-------|-------------|-------------|
| `seshmark` | Query Git history, blame, resume sessions | All |
| `seshmark-export` | Generate PR reports, markdown summaries | All |
| `seshmark-ci` | GitHub Actions integration | CI/CD |

## Custom Skills

Harness builders can extend the base skill:

```yaml
# ~/.config/my-agent/skills/seshmark-custom.yaml
extends: seshmark
name: seshmark-pro
description: Extended Seshmark with cost tracking

tools:
  - name: cost_report
    description: Generate cost analysis report
    command: "seshmark query --agent {agent} --format json | jq '{total_commits, avg_cost}'"
    parameters:
      agent: string
```

## The Seshmark Skill File

The canonical skill file lives at:
```
https://github.com/seshmark/seshmark/blob/main/skills/seshmark.yaml
```

Harness builders reference this URL. Updates are automatic when the file changes.

## Why This Design?

| Approach | Problem | Skill System |
|----------|---------|--------------|
| **Seshmark SDK** | Requires importing a library, version management | Agent uses CLI via shell calls |
| **Seshmark API** | Requires network, server, auth tokens | Runs locally via CLI |
| **GitHub App** | Requires webhooks, infrastructure, permissions | Agent decides when to invoke |
| **Hardcoded** | Every agent reimplements the logic | Single canonical skill definition |

The skill is a **declaration**, not code. The harness reads it and knows:
- What tools exist
- What parameters they take
- How to invoke them (shell command)
- When to use them (description)

Seshmark stays a CLI tool. The agent handles the orchestration.

---

## Example: Complete Agent Session with Skills

```
User: I need to fix a bug in auth. It was working last week.

Agent: I'll check what changed in the auth module.
[Calls seshmark_query with path="auth/", since="1 week ago"]

Agent: I found 5 commits. The most recent is by Claude:
- Commit: a3f9d2e
- Agent: claude
- Model: claude-sonnet-4
- Session: claude:thread-xyz789
- Subject: "feat: add OAuth2 middleware"

Agent: Let me get the full context of that session.
[Calls seshmark_context with commit="a3f9d2e"]

Agent: The session implemented OAuth2 with token validation. 
The continuation prompt mentions: "Still need to handle edge case: empty token"

User: Yes, that's the bug! The empty token case.

Agent: Let me check the current state of that file.
[Calls seshmark_blame with file="src/auth.ts"]

Agent: Lines 42-67 were written in that Claude session. 
Line 45 is where the empty token check should be.

User: Fix it.

Agent: I'll start a new tracked session to fix this.
[Calls seshmark_track with session_id="claude:fix-auth-20250514", agent="claude"]

[Agent writes the fix]

Agent: Fixed. I'll add a note about what we did.
[Calls seshmark_note with text="Fixed empty token edge case in validateToken()"]

[Git commit happens with SESHMARK_* env vars set]

Agent: Done. The fix is committed under session `claude:fix-auth-20250514`.
```

---

*Skills make Seshmark a first-class citizen in every agent's toolkit.*
