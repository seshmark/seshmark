package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
)

var contextCmd = &cobra.Command{
	Use:   "context <commit>",
	Short: "Reconstruct session context for a commit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		format, _ := cmd.Flags().GetString("format")
		return runContext(args[0], format)
	},
}

func init() {
	contextCmd.Flags().String("format", "markdown", "Output format: json, markdown, prompt")
}

func parseTrailers(body string) map[string]string {
	result := map[string]string{"session": "", "agent": "", "model": ""}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "AI-Session: ") {
			result["session"] = strings.TrimPrefix(line, "AI-Session: ")
		} else if strings.HasPrefix(line, "AI-Agent: ") {
			result["agent"] = strings.TrimPrefix(line, "AI-Agent: ")
		} else if strings.HasPrefix(line, "AI-Model: ") {
			result["model"] = strings.TrimPrefix(line, "AI-Model: ")
		}
	}
	return result
}

func runContext(commit, format string) error {
	body, err := git.Exec("log", "-1", "--format=%B", commit)
	if err != nil {
		return err
	}

	meta := parseTrailers(body)
	if meta["session"] == "" && meta["agent"] == "" {
		return fmt.Errorf("no AI metadata found for %s", commit)
	}

	// Get diff
	diff, _ := git.Exec("show", "--stat", commit)
	rawDiff, _ := git.Exec("show", commit)

	// Get all commits in session
	var commitsInSession []map[string]string
	if meta["session"] != "" {
		logOut, _ := git.Exec("log", "--all", "--grep=AI-Session: "+meta["session"], "--format=%H|%s|%ai", "--no-merges")
		for _, line := range strings.Split(logOut, "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "|", 3)
			if len(parts) == 3 {
				commitsInSession = append(commitsInSession, map[string]string{
					"hash": parts[0],
					"subject": parts[1],
					"date": parts[2],
				})
			}
		}
	}

	// Check for resolver
	resolverPath := ""
	resolverAvailable := false
	if meta["session"] != "" {
		parts := strings.SplitN(meta["session"], ":", 2)
		if len(parts) > 0 {
			resolverPath = filepath.Join(config.DataDir(), "resolvers", parts[0])
			if _, err := os.Stat(resolverPath); err == nil {
				resolverAvailable = true
			}
		}
	}

	// Notes
	var notes []string
	if meta["session"] != "" {
		parts := strings.SplitN(meta["session"], ":", 2)
		if len(parts) > 1 {
			notesFile := filepath.Join(config.DataDir(), "sessions", parts[1], "notes.md")
			if data, err := os.ReadFile(notesFile); err == nil {
				notes = strings.Split(strings.TrimSpace(string(data)), "\n")
			}
		}
	}

	// Build context package
	ctx := map[string]interface{}{
		"session_id":        meta["session"],
		"agent":             meta["agent"],
		"model":             meta["model"],
		"commit":            commit,
		"diff_stat":         diff,
		"raw_diff":          rawDiff,
		"commits_in_session": commitsInSession,
		"notes":             notes,
		"resolver": map[string]interface{}{
			"available": resolverAvailable,
			"path":      resolverPath,
		},
	}

	if format == "json" {
		data, _ := json.MarshalIndent(ctx, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if format == "prompt" {
		fmt.Printf("You previously worked on: %s\n", meta["session"])
		fmt.Printf("Agent: %s\n", meta["agent"])
		if meta["model"] != "" {
			fmt.Printf("Model: %s\n", meta["model"])
		}
		fmt.Printf("\nCommit history:\n")
		for _, c := range commitsInSession {
			fmt.Printf("- %s: %s\n", c["hash"], c["subject"])
		}
		fmt.Printf("\nDiff stat:\n%s\n", diff)
		if len(notes) > 0 {
			fmt.Printf("\nNotes:\n")
			for _, n := range notes {
				fmt.Printf("- %s\n", n)
			}
		}
		fmt.Printf("\nContinue the work from here.\n")
		return nil
	}

	// Markdown
	fmt.Printf("# Session Context: %s\n\n", meta["session"])
	fmt.Printf("- **Agent:** %s\n", meta["agent"])
	if meta["model"] != "" {
		fmt.Printf("- **Model:** %s\n", meta["model"])
	}
	fmt.Printf("- **Commit:** %s\n\n", commit)
	fmt.Printf("## Diff\n\n```\n%s\n```\n\n", diff)
	if len(notes) > 0 {
		fmt.Printf("## Notes\n\n")
		for _, n := range notes {
			fmt.Printf("- %s\n", n)
		}
		fmt.Println()
	}
	if resolverAvailable {
		fmt.Printf("**Resolver available:** %s\n", resolverPath)
	}
	return nil
}
