package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var logCmd = &cobra.Command{
	Use:   "log <session-id>",
	Short: "List commits in a session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		return runLog(args[0])
	},
}

func runLog(sessionID string) error {
	out, err := git.Exec("log", "--all", "--grep=AI-Session: "+sessionID, "--format=%h %ai %s", "--no-merges")
	if err != nil {
		return err
	}
	if out == "" {
		fmt.Println("No commits found for session.")
		return nil
	}
	fmt.Println(out)
	return nil
}
