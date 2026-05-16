# Contributing to Seshmark

> *git agentblame — Know which AI wrote every line.*

First off, thanks for being here. Seshmark is an open-source, community-driven project, and every contribution — whether it's a PR, an issue, a resolver script, or a typo fix — makes this better for everyone.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [What You Can Contribute](#what-you-can-contribute)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Adding a Resolver for a New AI Tool](#adding-a-resolver-for-a-new-ai-tool)
- [Style Guide](#style-guide)
- [Getting Help](#getting-help)

---

## Code of Conduct

This project follows a **no-drama policy**. Be respectful, assume good intent, and remember that everyone was a beginner once. Harassment, gatekeeping, and personal attacks are not tolerated.

## What You Can Contribute

| Area | Good For | Effort |
|------|----------|--------|
| **Bug reports** | Everyone | 5 min |
| **Documentation fixes** | Everyone | 10 min |
| **Resolver scripts** | Intermediate | 30 min |
| **New commands** | Experienced Go devs | 2-4 hr |
| **Feature requests** | Everyone | 10 min |
| **VS Code extension** | TS devs | 1-3 hr |
| **Website / landing page** | CSS/HTML folks | 30 min |

## Development Setup

### Prerequisites

- **Go 1.21+** ([install](https://go.dev/dl/))
- **Git** (obviously)
- A terminal you're comfortable in

### Clone & Build

```bash
git clone https://github.com/seshmark/seshmark.git
cd seshmark
go build -o seshmark .
./seshmark --help
```

### Install Locally (for testing)

```bash
make install
# or
cp seshmark ~/.local/bin/
```

The binary supports symlinks for its aliases:

```bash
ln -sf ~/.local/bin/seshmark ~/.local/bin/agentblame
ln -sf ~/.local/bin/seshmark ~/.local/bin/git-agentblame
```

## Project Structure

```
seshmark/
├── main.go                  # Entry point, alias detection
├── go.mod / go.sum          # Dependencies
├── Makefile                 # Build, test, release targets
├── install.sh / install.ps1 # One-liner install scripts
├── index.html               # Landing page (GitHub Pages)
├── img_assets/              # Brand images and logos
│
├── internal/
│   ├── cmd/                 # CLI commands (cobra)
│   │   ├── root.go          # Root command + registration
│   │   ├── blame.go         # git agentblame
│   │   ├── query.go         # seshmark query
│   │   ├── resume.go        # seshmark resume
│   │   ├── resolver.go      # Resolver management
│   │   ├── hook.go          # Hook install/uninstall
│   │   ├── hookrun.go       # Internal hook execution
│   │   ├── track.go         # Manual session tracking
│   │   ├── untrack.go       # Stop session tracking
│   │   ├── note.go          # Continuation notes
│   │   ├── commit.go        # Explicit commit with trailers
│   │   ├── log.go           # Session commit log
│   │   ├── context.go       # Session context reconstruction
│   │   ├── who.go           # Commit metadata
│   │   ├── stats.go         # AI usage stats
│   │   ├── status.go        # Active session state
│   │   └── doctor.go        # Installation diagnostics
│   │
│   ├── config/              # XDG config paths
│   ├── git/                 # Git execution helpers
│   ├── hook/                # Hook scripts installer
│   └── resolver/            # Resolver find/exec system
│
├── pkg/
│   └── version/             # Build version injection
│
├── examples/resolvers/      # Reference resolvers
│   ├── cursor
│   ├── claude
│   ├── opencode
│   └── pi
│
└── vscode-extension/        # VS Code extension (TypeScript)
```

## Making Changes

### Branch Naming

We use a lightweight convention:

- `fix/<short-description>` — Bug fixes
- `feat/<short-description>` — New features
- `docs/<short-description>` — Documentation
- `refactor/<short-description>` — Code improvements

Please keep branches focused on a single concern.

### Commit Messages

We follow [conventional commits](https://www.conventionalcommits.org/):

```
feat: add --json flag to blame command

Explain what this does and why. Reference issues if applicable.
```

Seshmark uses Git trailers extensively — your commits should too when relevant:

```
feat: add opencode resolver

AI-Agent: claude
```

### Go Conventions

- Run `go fmt` before committing
- Run `go vet ./...` to catch issues
- Avoid external dependencies where possible (Seshmark's dependency tree is intentionally tiny)
- Prefer the standard library over frameworks
- Error messages should be lowercase (Go convention) and actionable
- Use `internal/` for packages you don't want to expose publicly

## Testing

```bash
# Run unit tests
go test ./...

# Run with verbose output
go test -v ./...

# Run a specific test
go test -v ./internal/cmd -run TestBlame
```

### Manual Testing

After building, test the installation flow:

```bash
# Create a test repo
cd /tmp
mkdir test-seshmark && cd test-seshmark
git init

# Install the hook
seshmark hook install

# Make a tracked commit
SESHMARK_SESSION_ID="test:manual-001" \
SESHMARK_AGENT="test-agent" \
SESHMARK_MODEL="test-model" \
git commit --allow-empty -m "test: manual attribution"

# Verify
seshmark who HEAD
git agentblame HEAD
```

## Submitting a Pull Request

1. **Fork the repo** on GitHub
2. **Create a branch** for your change
3. **Make your changes** with clear commits
4. **Run tests** — make sure `go test ./...` passes
5. **Push to your fork** and open a PR
6. **Describe your change** — what it does, why, and how to test it
7. **Wait for review** — we try to respond within a few days

### PR Checklist

- [ ] Code compiles (`go build .`)
- [ ] Tests pass (`go test ./...`)
- [ ] `go fmt` has been run
- [ ] `go vet` shows no issues
- [ ] New functionality includes tests (if applicable)
- [ ] Documentation updated (if applicable)

## Adding a Resolver for a New AI Tool

This is the most common and most welcome contribution. A resolver is a script that tells Seshmark how to reopen a session in a specific AI tool.

### Quick Start

```bash
# Generate a starter script
seshmark resolver create my-tool

# Edit it
vim ~/.local/share/seshmark/resolvers/my-tool

# Make it executable
chmod +x ~/.local/share/seshmark/resolvers/my-tool

# Test it
seshmark resolver test my-tool
```

### Resolver Script Contract

```bash
#!/bin/bash
# ~/.local/share/seshmark/resolvers/my-tool
SESSION_ID="$1"
# Your logic here — open the session in your tool
my-tool open --session "$SESSION_ID"
```

The resolver receives exactly one argument: the session ID (e.g., `cursor:chat-abc123`).

### What Makes a Great Resolver

- **Handles missing tool gracefully**: If the tool isn't installed, print a helpful message with install instructions
- **Session format flexibility**: Handle both prefixed (`cursor:chat-abc123`) and bare (`chat-abc123`) session IDs
- **Clear output**: If you can't open the session directly, at least print the session ID so the user can find it manually

### Submit Your Resolver

Once your resolver is working:

1. Add it to `examples/resolvers/<your-tool>`
2. Submit a PR — we'd love to include it!

## Style Guide

### Go

- Follow `gofmt` output exactly (no exceptions)
- Use `CamelCase` for exported names, `camelCase` for unexported
- Error values start with lowercase (e.g., `fmt.Errorf("not a git repository")`)
- Comment your public types and functions
- Use table-driven tests

### Documentation

- Use `code blocks` for commands, file paths, and env vars
- Use `--format` for CLI flags (double-dash)
- Prefer examples over explanations
- Keep sentences short (non-native English speakers should be able to read it)
- Oxford comma always

### Markdown

- Sentence case for headings
- One blank line before and after code blocks
- Use `[text](link)` for references
- Relative links for repo-internal docs

## Getting Help

- **Open an issue** on [GitHub](https://github.com/seshmark/seshmark/issues) for bugs and feature requests
- **Start a discussion** on [GitHub Discussions](https://github.com/seshmark/seshmark/discussions) for questions and ideas
- **Check existing issues** before opening a new one
- **Tag your issue** appropriately (bug, enhancement, documentation, question)

---

*Every PR, issue, and resolver script makes this better. We're glad you're here.*
