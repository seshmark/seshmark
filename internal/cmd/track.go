package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var trackCmd = &cobra.Command{
	Use:   "track [provider:]<id>",
	Short: "Start tracking a session manually",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		
		var sessionID string
		if len(args) > 0 {
			sessionID = args[0]
		} else {
			fromBranch, _ := cmd.Flags().GetBool("from-branch")
			if fromBranch {
				branch, err := git.Exec("branch", "--show-current")
				if err != nil || branch == "" {
					return fmt.Errorf("could not determine current branch")
				}
				sessionID = "local:" + branch + "-" + time.Now().Format("20060102") + "-manual"
			} else {
				sessionID = "local:" + time.Now().Format("20060102-150405")
			}
		}
		
		agent, _ := cmd.Flags().GetString("agent")
		model, _ := cmd.Flags().GetString("model")
		force, _ := cmd.Flags().GetBool("force")
		
		return runTrack(sessionID, agent, model, force)
	},
}

func init() {
	trackCmd.Flags().Bool("from-branch", false, "Infer agent from branch name")
	trackCmd.Flags().String("agent", "", "Agent name")
	trackCmd.Flags().String("model", "", "Model name")
	trackCmd.Flags().Bool("force", false, "Overwrite existing session")
}

func runTrack(sessionID, agent, model string, force bool) error {
	gitDir, err := git.Exec("rev-parse", "--git-dir")
	if err != nil {
		return err
	}
	
	stateFile := filepath.Join(gitDir, "seshmark", "current")
	if _, err := os.Stat(stateFile); err == nil && !force {
		return fmt.Errorf("session already active; use --force to overwrite")
	}
	
	if err := os.MkdirAll(filepath.Dir(stateFile), 0755); err != nil {
		return err
	}
	
	// Infer agent from session ID prefix if not provided
	if agent == "" {
		parts := strings.SplitN(sessionID, ":", 2)
		if len(parts) > 1 {
			agent = parts[0]
		}
	}
	
	state := map[string]string{
		"version":    "1",
		"session_id": sessionID,
		"agent":      agent,
		"model":      model,
		"started_at": time.Now().UTC().Format(time.RFC3339),
	}
	
	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return err
	}
	
	fmt.Printf("Active session: %s\n", sessionID)
	return nil
}
