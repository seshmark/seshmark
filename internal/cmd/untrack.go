package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var untrackCmd = &cobra.Command{
	Use:   "untrack",
	Short: "Stop tracking the current session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		gitDir, err := git.Exec("rev-parse", "--git-dir")
		if err != nil {
			return err
		}
		stateFile := filepath.Join(gitDir, "seshmark", "current")
		if err := os.Remove(stateFile); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No active session.")
				return nil
			}
			return err
		}
		fmt.Println("Session untracked.")
		return nil
	},
}
