package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
)

var noteCmd = &cobra.Command{
	Use:   "note <text...>",
	Short: "Append a note to the current session",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		return runNote(strings.Join(args, " "))
	},
}

func runNote(text string) error {
	// Find current session
	gitDir, err := git.Exec("rev-parse", "--git-dir")
	if err != nil {
		return err
	}

	stateFile := filepath.Join(gitDir, "seshmark", "current")
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("no active session")
	}

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	sessionID, ok := state["session_id"].(string)
	if !ok {
		return fmt.Errorf("no active session")
	}

	parts := strings.SplitN(sessionID, ":", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid session ID")
	}

	sessionDir := filepath.Join(config.DataDir(), "sessions", parts[1])
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return err
	}

	notesFile := filepath.Join(sessionDir, "notes.md")
	f, err := os.OpenFile(notesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "- [%s] %s\n", time.Now().UTC().Format(time.RFC3339), text)
	fmt.Println("Note saved.")
	return nil
}
