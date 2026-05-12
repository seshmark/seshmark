package cmd

import (
	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal"
)

var hookRunCmd = &cobra.Command{
	Use:    "hook-run <msg-file> [<source>]",
	Short:  "Internal: run by Git hook (do not use directly)",
	Hidden: true,
	Args:   cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		msgFile := args[0]
		source := ""
		if len(args) > 1 {
			source = args[1]
		}
		return internal.HookRun(msgFile, source)
	},
}

func init() {
	rootCmd.AddCommand(hookRunCmd)
}
