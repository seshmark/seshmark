package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
	"github.com/seshmark/seshmark/internal/hook"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose seshmark installation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor()
	},
}

func runDoctor() error {
	ok := true

	// Check seshmark in PATH
	if _, err := exec.LookPath("seshmark"); err != nil {
		fmt.Println("❌ seshmark not found in PATH")
		ok = false
	} else {
		fmt.Println("✅ seshmark in PATH")
	}

	// Check git
	if _, err := git.Exec("--version"); err != nil {
		fmt.Println("❌ Git not found")
		ok = false
	} else {
		ver, _ := git.Exec("--version")
		fmt.Printf("✅ %s\n", ver)
	}

	// Check config dir
	cfgDir := config.ConfigDir()
	if _, err := os.Stat(cfgDir); err == nil {
		fmt.Printf("✅ Config dir: %s\n", cfgDir)
	} else {
		fmt.Printf("⚠️  Config dir missing: %s\n", cfgDir)
	}

	// Check data dir
	dataDir := config.DataDir()
	if _, err := os.Stat(dataDir); err == nil {
		fmt.Printf("✅ Data dir: %s\n", dataDir)
	} else {
		fmt.Printf("⚠️  Data dir missing: %s\n", dataDir)
	}

	// Check resolvers
	resolverDir := filepath.Join(dataDir, "resolvers")
	if entries, err := os.ReadDir(resolverDir); err == nil && len(entries) > 0 {
		fmt.Printf("✅ Resolvers found (%d):\n", len(entries))
		for _, e := range entries {
			info, _ := e.Info()
			fmt.Printf("   - %s (mode: %o)\n", e.Name(), info.Mode().Perm())
		}
	} else {
		fmt.Println("ℹ️  No resolvers installed")
	}

	// Check current repo hook
	gitDir, err := git.Exec("rev-parse", "--git-dir")
	if err == nil {
		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
		if data, err := os.ReadFile(hookPath); err == nil {
			if string(data) == hook.DelegatorScript {
				fmt.Println("✅ Hook installed in current repo")
			} else {
				fmt.Println("⚠️  Different hook found in current repo")
			}
		} else {
			fmt.Println("⚠️  No hook in current repo")
		}
	} else {
		fmt.Println("ℹ️  Not in a Git repository")
	}

	if !ok {
		return fmt.Errorf("some checks failed")
	}
	return nil
}
