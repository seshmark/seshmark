# Seshmark Plugin System — Custom Trailers

## Philosophy

Seshmark ships with **4 core trailers** and nothing else. The binary stays lightweight.

Users who want custom metadata (token usage, cost, latency, confidence scores) can define them via **declarative config** — no code changes, no plugins to install, no runtime bloat.

## How It Works

### 1. User defines custom trailers in config

```yaml
# ~/.config/seshmark/trailers.yaml
custom:
  - name: Token-Usage
    env: SESHMARK_TOKEN_USAGE
    description: "Total tokens consumed in the session"
  
  - name: Cost
    env: SESHMARK_COST
    format: "${value} USD"
    description: "Estimated API cost for this session"
  
  - name: Latency
    env: SESHMARK_LATENCY_MS
    description: "Average response latency in milliseconds"
  
  - name: Confidence
    env: SESHMARK_CONFIDENCE
    description: "AI confidence score (0-1)"
```

### 2. Harness sets env vars before commit

```python
env["SESHMARK_TOKEN_USAGE"] = "15000"
env["SESHMARK_COST"] = "0.023"
env["SESHMARK_LATENCY_MS"] = "450"
env["SESHMARK_CONFIDENCE"] = "0.92"
```

### 3. Hook appends them automatically

The commit message becomes:

```
feat: implement OAuth2 middleware

Seshmark-Version: 1.0.0
AI-Session: cursor:chat-abc123
AI-Agent: cursor
AI-Model: claude-sonnet-4-20250514
Token-Usage: 15000
Cost: 0.023 USD
Latency: 450
Confidence: 0.92
```

### 4. Query and blame support them

```bash
seshmark query --custom Cost --since "last week"
```

### 5. Built-in protections

- Custom trailers are **appended after core trailers** (so core spec stability is guaranteed)
- Names must match `[A-Za-z0-9-]+` (no colons, no spaces)
- Max 10 custom trailers per commit (prevents spam)
- Custom trailers are **never required** for spec compliance
- Unknown custom trailers are silently ignored by parsers that don't support them

## Why This Design

| Approach | Problem | Our Solution |
|----------|---------|--------------|
| **Built-in trailers for everything** | Binary bloat, opinionated, can't satisfy everyone | Core stays minimal, user defines extras |
| **Plugin scripts** | Security nightmare, version hell, dependency management | Declarative YAML, zero code execution |
| **API / registry** | Requires network, centralization, vendor lock-in | Env vars, local config, fully offline |
| **Git notes** | Invisible, hard to query, not cloned by default | Trailers are in commit messages, always visible |

## Config Locations (XDG)

```
~/.config/seshmark/trailers.yaml     # User-defined custom trailers
~/.config/seshmark/agents.json       # Agent name mappings
```

## For Harness Builders

Harnesses expose custom metrics by setting env vars. No schema registration, no API calls.

```python
# Cursor could expose this natively
env["SESHMARK_TOKEN_USAGE"] = str(session.total_tokens)
env["SESHMARK_COST"] = str(session.estimated_cost)
```

If the user has `Token-Usage` in their `trailers.yaml`, it gets appended. If not, the env var is ignored. Zero coupling.

## For Tooling

```bash
# List all custom trailers the user has configured
seshmark config list-custom

# Show example commit with custom trailers
seshmark config preview

# Validate trailers.yaml
seshmark config validate
```

## Future: Shared Schemas

If a community converges on standard custom trailers, they can be promoted to "extended core" without breaking anything:

```yaml
# ~/.config/seshmark/trailers.yaml
extends: "github.com/seshmark/trailers/v1/cost-schema.yaml"
custom:
  - name: My-Custom-Field
    env: MY_VAR
```

This is Phase 3. For now, local YAML is sufficient.
