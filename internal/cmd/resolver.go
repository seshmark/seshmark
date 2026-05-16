package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/resolver"
)

var resolverCmd = &cobra.Command{
	Use:   "resolver",
	Short: "Manage session resolvers for AI harnesses",
	Long: `Resolvers are scripts that open or resume native AI sessions.

Install a resolver to make 'seshmark resume' open your AI tool directly.

Examples:
  seshmark resolver list                    # List installed resolvers
  seshmark resolver create cursor           # Create a starter resolver
  seshmark resolver install cursor ./open-cursor.sh
  seshmark resolver uninstall cursor
  seshmark resolver test cursor             # Test a resolver`,
}

var resolverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed resolvers",
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvers, err := resolver.List()
		if err != nil {
			return err
		}
		if len(resolvers) == 0 {
			fmt.Println("No resolvers installed.")
			fmt.Println("")
			fmt.Println("Create one:")
			fmt.Println("  seshmark resolver create <provider>")
			return nil
		}
		fmt.Println("Installed resolvers:")
		for _, r := range resolvers {
			fmt.Printf("  %-15s %s\n", r.Provider, r.Path)
		}
		return nil
	},
}

var resolverCreateCmd = &cobra.Command{
	Use:   "create <provider>",
	Short: "Create a starter resolver script",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		r := resolver.Find(provider)
		if r.Exists {
			return fmt.Errorf("resolver '%s' already exists at %s", provider, r.Path)
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(r.Path), 0755); err != nil {
			return fmt.Errorf("cannot create resolver directory: %w", err)
		}

		// Write template
		tmpl := resolver.Template(provider)
		if err := os.WriteFile(r.Path, []byte(tmpl), 0755); err != nil {
			return fmt.Errorf("cannot create resolver: %w", err)
		}

		fmt.Printf("Created resolver: %s\n", r.Path)
		fmt.Println("")
		fmt.Println("Edit this script to implement your resolver logic.")
		fmt.Println("Then test it:")
		fmt.Printf("  seshmark resolver test %s\n", provider)
		return nil
	},
}

var resolverInstallCmd = &cobra.Command{
	Use:   "install <provider> <script-path>",
	Short: "Install a resolver from a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		scriptPath := args[1]
		if err := resolver.Install(provider, scriptPath); err != nil {
			return err
		}
		fmt.Printf("Installed resolver: %s\n", provider)
		return nil
	},
}

var resolverUninstallCmd = &cobra.Command{
	Use:   "uninstall <provider>",
	Short: "Remove a resolver",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := resolver.Uninstall(args[0]); err != nil {
			return err
		}
		fmt.Printf("Uninstalled resolver: %s\n", args[0])
		return nil
	},
}

var resolverTestCmd = &cobra.Command{
	Use:   "test <provider>",
	Short: "Test a resolver with a dummy session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		r := resolver.Find(provider)
		if !r.Exists {
			return fmt.Errorf("resolver '%s' not found. Create it: seshmark resolver create %s", provider, provider)
		}

		fmt.Printf("Testing resolver for '%s'...\n", provider)
		fmt.Printf("Path: %s\n", r.Path)
		fmt.Println("")

		err := r.Exec(provider+":test-session", 5*time.Second)
		if err != nil {
			fmt.Printf("Resolver exited with error: %v\n", err)
			return nil
		}

		fmt.Println("Resolver executed successfully.")
		return nil
	},
}

func init() {
	resolverCmd.AddCommand(resolverListCmd)
	resolverCmd.AddCommand(resolverCreateCmd)
	resolverCmd.AddCommand(resolverInstallCmd)
	resolverCmd.AddCommand(resolverUninstallCmd)
	resolverCmd.AddCommand(resolverTestCmd)
	rootCmd.AddCommand(resolverCmd)
}
