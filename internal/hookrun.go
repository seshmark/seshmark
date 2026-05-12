package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	return "", "", ""
}

func extractAgentFromBranch(branch string) string {
	re := regexp.MustCompile(`^([a-zA-Z0-9-]+)[/._-]`)
	matches := re.FindStringSubmatch(branch)
	if len(matches) < 2 {
		return ""
	}
	prefix := strings.ToLower(matches[1])
	for _, a := range knownAgents {
		if a == prefix {
			return a
		}
	}
	return ""
}
