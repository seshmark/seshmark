package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/hook"
)

var hookCmd = &cobra.Command{
	Use:   "hook <install|uninstall>",
	Short: "Install or uninstall Git hooks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		global, _ := cmd.Flags().GetBool("global")
		force, _ := cmd.Flags().GetBool("force")

		switch args[0] {
		case "install":
			if !global {
				if err := requireGitRepo(); err != nil {
					return err
				}
			}
			return hook.Install(true, global, force)
		case "uninstall":
			return hook.Uninstall(global)
		default:
			return fmt.Errorf("unknown hook command: %s", args[0])
		}
	},
}

func init() {
	hookCmd.Flags().BoolP("global", "g", false, "Apply to global Git template")
	hookCmd.Flags().BoolP("force", "f", false, "Overwrite existing hook")
}
