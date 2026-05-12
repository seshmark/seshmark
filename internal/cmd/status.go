package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
	"github.com/seshmark/seshmark/internal/hook"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active session and installation state",
	RunE: func(cmd *cobra.Command, args []string) error {
		inRepo := true
		gitDir, err := git.Exec("rev-parse", "--git-dir")
		if err != nil {
			inRepo = false
		}

		// Active session
		var sessionID, agent, model string
		if inRepo {
			stateFile := filepath.Join(gitDir, "seshmark", "current")
			if data, err := os.ReadFile(stateFile); err == nil {
				var state map[string]interface{}
				if json.Unmarshal(data, &state) == nil {
					sessionID, _ = state["session_id"].(string)
					agent, _ = state["agent"].(string)
					model, _ = state["model"].(string)
				}
			}
		}
		if sessionID == "" {
			sessionID = os.Getenv("SESHMARK_SESSION_ID")
			agent = os.Getenv("SESHMARK_AGENT")
			model = os.Getenv("SESHMARK_MODEL")
		}

		// Branch
		var branch string
		if inRepo {
			branch, _ = git.Exec("branch", "--show-current")
		}

		// Hook installed?
		hookInstalled := false
		if inRepo {
			hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
			if _, err := os.Stat(hookPath); err == nil {
				data, _ := os.ReadFile(hookPath)
				hookInstalled = string(data) == hook.DelegatorScript
			}
		}

		fmt.Printf("Repository: %v\n", inRepo)
		if inRepo {
			fmt.Printf("Branch: %s\n", branch)
		}
		if sessionID != "" {
			fmt.Printf("Active session: %s\n", sessionID)
			if agent != "" {
				fmt.Printf("Agent: %s\n", agent)
			}
			if model != "" {
				fmt.Printf("Model: %s\n", model)
			}
		} else {
			fmt.Println("No active session.")
		}
		fmt.Printf("Hook installed: %v\n", hookInstalled)
		return nil
	},
}
