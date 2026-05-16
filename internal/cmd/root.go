package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var rootCmd = &cobra.Command{
	Use:   "seshmark",
	Short: "Seshmark — Agent-aware Git history",
	Long: `Seshmark is a zero-infrastructure, tool-agnostic convention
that links AI coding agent sessions to Git commits.

Also known as: git agentblame, agentblame, aiblame`,
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(blameCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(whoCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(trackCmd)
	rootCmd.AddCommand(untrackCmd)
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(upgradeCmd)
}

func requireGitRepo() error {
	_, err := git.Exec("rev-parse", "--git-dir")
	return err
}
