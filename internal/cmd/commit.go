package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Commit with seshmark trailers",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		noMark, _ := cmd.Flags().GetBool("no-mark")
		message := args[0]
		
		if noMark {
			return exec.Command("git", "commit", "-m", message).Run()
		}
		
		// Build a temporary commit with trailers via the hook
		return exec.Command("git", "commit", "-m", message).Run()
	},
}

func init() {
	commitCmd.Flags().Bool("no-mark", false, "Commit without seshmark trailers")
}
