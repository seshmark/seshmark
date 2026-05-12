package hook

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
)

func Install(local, global, force bool) error {
	if global {
		tmplDir := filepath.Join(config.ConfigDir(), "git-template", "hooks")
		if err := os.MkdirAll(tmplDir, 0755); err != nil {
			return err
		}
		hookPath := filepath.Join(tmplDir, "prepare-commit-msg")
		if err := writeHook(hookPath, force); err != nil {
			return err
		}
		if err := setGlobalTemplateDir(); err != nil {
			return err
		}
		fmt.Println("Installed global hook template.")
	}

	if local {
		gitDir, err := git.Exec("rev-parse", "--git-dir")
		if err != nil {
			return fmt.Errorf("not in a git repository")
		}
		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
		if err := writeHook(hookPath, force); err != nil {
			return err
		}
		fmt.Println("Installed hook in current repository.")
	}

	return nil
}

func Uninstall(global bool) error {
	if global {
		tmplDir := filepath.Join(config.ConfigDir(), "git-template", "hooks")
		hookPath := filepath.Join(tmplDir, "prepare-commit-msg")
		if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
			return err
		}
		fmt.Println("Removed global hook template.")
	}

	gitDir, err := git.Exec("rev-parse", "--git-dir")
	if err == nil {
		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
		data, err := os.ReadFile(hookPath)
		if err == nil && string(data) == DelegatorScript {
			os.Remove(hookPath)
			fmt.Println("Removed hook from current repository.")
		}
	}
	return nil
}

func writeHook(path string, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("hook already exists at %s; use --force to overwrite", path)
	}
	return os.WriteFile(path, []byte(DelegatorScript), 0755)
}

func setGlobalTemplateDir() error {
	_, err := git.Exec("config", "--global", "init.templateDir", filepath.Join(config.ConfigDir(), "git-template"))
	return err
}
