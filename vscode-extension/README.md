# Seshmark for VS Code

See AI session attribution inline in your code editor.

## Features

- **Inline gutter labels** showing which AI agent wrote each line
- **Hover tooltips** with session ID, model, and commit info
- **Status bar** toggle to enable/disable decorations
- **Command palette** integration for blame, query, and who

## Usage

1. Install the Seshmark CLI: `curl -fsSL https://seshmark.github.io/seshmark/install.sh | bash`
2. Install this extension
3. Open any file in a repo with seshmark-tagged commits
4. See `[cursor]`, `[claude]`, etc. in the gutter

## Commands

| Command | Description |
|---------|-------------|
| `Seshmark: Agent Blame` | Run `git agentblame` in terminal |
| `Seshmark: Query AI Commits` | Search commits by agent |
| `Seshmark: Show Commit Info` | Show info for current line |
| `Seshmark: Toggle Inline Decorations` | Show/hide gutter labels |

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `seshmark.enabled` | `true` | Show AI attribution inline |
| `seshmark.showAgentInGutter` | `true` | Show agent name in gutter |
| `seshmark.colors` | `{...}` | Colors per agent |

## Requirements

- Seshmark CLI installed and in PATH
- Git repository with seshmark-tagged commits
