package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/seshmark/seshmark/internal/git"
)

type SessionState struct {
	Version   string `json:"version"`
	SessionID string `json:"session_id"`
	Agent     string `json:"agent"`
	Model     string `json:"model"`
}

var knownAgents = []string{
	"cursor", "claude", "codex", "opencode", "aider",
	"github-copilot", "devin", "swe-agent", "builder",
	"pi",
}

func HookRun(msgFile, source string) error {
	if source == "merge" || source == "squash" {
		return nil
	}

	// Deduplication: if already tagged, skip
	data, err := os.ReadFile(msgFile)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), "Seshmark-Version:") {
		return nil
	}

	sessionID, agent, model := resolveMetadata()
	if sessionID == "" && agent == "" {
		return nil
	}

	// Append trailers
	f, err := os.OpenFile(msgFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f)
	fmt.Fprintln(f, "Seshmark-Version: 1.0.0")
	if sessionID != "" {
		fmt.Fprintf(f, "AI-Session: %s\n", sessionID)
	}
	if agent != "" {
		fmt.Fprintf(f, "AI-Agent: %s\n", agent)
	}
	if model != "" {
		fmt.Fprintf(f, "AI-Model: %s\n", model)
	}
	return nil
}

func resolveMetadata() (sessionID, agent, model string) {
	// Priority 1: session file
	gitDir, err := git.Exec("rev-parse", "--git-dir")
	if err == nil {
		stateFile := filepath.Join(gitDir, "seshmark", "current")
		if data, err := os.ReadFile(stateFile); err == nil {
			var state SessionState
			if json.Unmarshal(data, &state) == nil {
				return state.SessionID, state.Agent, state.Model
			}
		}
	}

	// Priority 2: namespaced env vars
	if sid := os.Getenv("SESHMARK_SESSION_ID"); sid != "" {
		return sid, os.Getenv("SESHMARK_AGENT"), os.Getenv("SESHMARK_MODEL")
	}
	// Priority 2b: agent-only via env var (no session ID needed)
	if agent := os.Getenv("SESHMARK_AGENT"); agent != "" {
		return "", agent, os.Getenv("SESHMARK_MODEL")
	}

	// Priority 3: bare env vars (legacy)
	if sid := os.Getenv("AI_SESSION_ID"); sid != "" {
		return sid, os.Getenv("AI_AGENT"), os.Getenv("AI_MODEL")
	}

	// Priority 4: branch inference (agent only)
	branch, _ := git.Exec("branch", "--show-current")
	if branch != "" {
		if a := extractAgentFromBranch(branch); a != "" {
			return "", a, ""
		}
	}

	// Priority 5: process-based detection (agent only)
	// Check parent processes to detect the AI tool that's running
	if a := detectAgentFromProcess(); a != "" {
		return "", a, ""
	}

	return "", "", ""
}

func extractAgentFromBranch(branch string) string {
	re := regexp.MustCompile(`^([a-zA-Z0-9-]+)[/._-]`)
	matches := re.FindStringSubmatch(branch)
	if len(matches) < 2 {
		return ""
	}
	prefix := strings.ToLower(matches[1])
	// Skip common Git branch prefixes that aren't agents
	skipPrefixes := map[string]bool{
		"feature": true, "feat": true, "fix": true, "bugfix": true,
		"hotfix": true, "release": true, "chore": true, "docs": true,
		"refactor": true, "test": true, "main": true, "master": true,
		"develop": true, "dev": true,
	}
	if skipPrefixes[prefix] {
		return ""
	}
	return prefix
}

// detectAgentFromProcess walks the parent process tree looking for the
// AI coding tool that invoked this commit. It works for any tool by
// skipping known non-agent processes (git, shells, etc.) and returning
// the first unknown parent process name as the agent.
func detectAgentFromProcess() string {
	ppid := os.Getppid()
	if ppid <= 1 {
		return ""
	}

	// Skip these process names — they're never the AI tool
	skipProcesses := map[string]bool{
		// Git and VCS
		"git": true, "git2": true,
		// Shells
		"bash": true, "zsh": true, "sh": true, "dash": true,
		"fish": true, "ksh": true, "tcsh": true,
		// seshmark itself
		"seshmark": true, "hook-run": true,
		// Common languages that tools may be built with
		"python": true, "python3": true, "node": true, "nodejs": true,
		"deno": true, "bun": true,
		// Package managers / build tools
		"make": true, "npx": true, "npm": true, "yarn": true, "pnpm": true,
		// Terminal/session (never the tool itself)
		"login": true, "tmux": true, "screen": true,
		// Editors — we want the AI, not the editor
		"vim": true, "nvim": true, "emacs": true, "nano": true,
	}

	// Check if a process name might be an AI coding tool.
	// It must NOT be in the skip list and must NOT be a system process.
	isLikelyTool := func(name string) bool {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			return false
		}
		if skipProcesses[name] {
			return false
		}
		// Skip numbered exit codes from ps
		if _, err := fmt.Sscanf(name, "%d", new(int)); err == nil {
			return false
		}
		// Skip paths with common system dirs
		if strings.HasPrefix(name, "/usr/lib/") ||
			strings.HasPrefix(name, "/System/") ||
			strings.HasPrefix(name, "/Applications/") {
			return false
		}
		// Must have at least 2 chars to be meaningful
		if len(name) < 2 {
			return false
		}
		return true
	}

	// Walk up the process tree looking for the first non-skipped process
	// We use a set to detect cycles
	seen := map[int]bool{}
	pid := ppid
	for i := 0; i < 8; i++ {
		if seen[pid] {
			break
		}
		seen[pid] = true

		name, err := getProcessName(pid)
		if err != nil || name == "" {
			break
		}
		name = strings.TrimSpace(strings.TrimSuffix(name, "\n"))
		name = filepath.Base(name) // strip directory path

		if isLikelyTool(name) {
			return strings.ToLower(name)
		}

		// Walk up to parent
		pid, err = getParentPid(pid)
		if err != nil || pid <= 1 {
			break
		}
	}

	return ""
}

func getProcessName(pid int) (string, error) {
	if runtime.GOOS == "windows" {
		return "", nil
	}
	out, err := exec.Command("ps", "-o", "comm=", "-p", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getParentPid(pid int) (int, error) {
	if runtime.GOOS == "windows" {
		return 0, fmt.Errorf("not supported")
	}
	out, err := exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		return 0, err
	}
	ppid := 0
	fmt.Sscanf(string(out), "%d", &ppid)
	if ppid <= 0 {
		return 0, fmt.Errorf("no parent")
	}
	return ppid, nil
}
