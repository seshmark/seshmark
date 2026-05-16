package resolver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/seshmark/seshmark/internal/config"
)

// Resolver represents an installed provider resolver
type Resolver struct {
	Provider string
	Path     string
	Exists   bool
}

// Find searches for a resolver for the given provider
func Find(provider string) *Resolver {
	// Normalize provider name
	provider = strings.ToLower(provider)

	// Search in XDG data dir
	resolverDir := filepath.Join(config.DataDir(), "resolvers")
	resolverPath := filepath.Join(resolverDir, provider)

	if _, err := os.Stat(resolverPath); err == nil {
		return &Resolver{
			Provider: provider,
			Path:     resolverPath,
			Exists:   true,
		}
	}

	return &Resolver{
		Provider: provider,
		Path:     resolverPath,
		Exists:   false,
	}
}

// List returns all installed resolvers
func List() ([]Resolver, error) {
	resolverDir := filepath.Join(config.DataDir(), "resolvers")

	entries, err := os.ReadDir(resolverDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var resolvers []Resolver
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		// Only include executable files
		if info.Mode()&0o111 != 0 {
			resolvers = append(resolvers, Resolver{
				Provider: entry.Name(),
				Path:     filepath.Join(resolverDir, entry.Name()),
				Exists:   true,
			})
		}
	}

	return resolvers, nil
}

// Exec runs the resolver for a provider with the given session ID
func (r *Resolver) Exec(sessionID string, timeout time.Duration) error {
	if !r.Exists {
		return fmt.Errorf("no resolver found for provider '%s' at %s", r.Provider, r.Path)
	}

	// Security: ensure resolver is a file, not a symlink to something dangerous
	info, err := os.Stat(r.Path)
	if err != nil {
		return fmt.Errorf("cannot access resolver: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		// Follow symlink but verify target
		resolved, err := filepath.EvalSymlinks(r.Path)
		if err != nil {
			return fmt.Errorf("cannot resolve symlink: %w", err)
		}
		info, err = os.Stat(resolved)
		if err != nil {
			return fmt.Errorf("cannot access resolved resolver: %w", err)
		}
	}

	// Ensure it's executable
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("resolver is not executable: %s", r.Path)
	}

	// Extract just the ID portion (after the colon)
	parts := strings.SplitN(sessionID, ":", 2)
	id := sessionID
	if len(parts) == 2 {
		id = parts[1]
	}

	cmd := exec.Command(r.Path, id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if timeout > 0 {
		done := make(chan error, 1)
		go func() {
			done <- cmd.Run()
		}()

		select {
		case err := <-done:
			return err
		case <-time.After(timeout):
			cmd.Process.Kill()
			return fmt.Errorf("resolver timed out after %v", timeout)
		}
	}

	return cmd.Run()
}

// Install copies a script to the resolver directory
func Install(provider, scriptPath string) error {
	provider = strings.ToLower(provider)

	resolverDir := filepath.Join(config.DataDir(), "resolvers")
	if err := os.MkdirAll(resolverDir, 0755); err != nil {
		return err
	}

	target := filepath.Join(resolverDir, provider)

	// Read source
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("cannot read script: %w", err)
	}

	// Write to resolver dir
	if err := os.WriteFile(target, data, 0755); err != nil {
		return fmt.Errorf("cannot install resolver: %w", err)
	}

	return nil
}

// Uninstall removes a resolver
func Uninstall(provider string) error {
	provider = strings.ToLower(provider)

	resolverDir := filepath.Join(config.DataDir(), "resolvers")
	target := filepath.Join(resolverDir, provider)

	if err := os.Remove(target); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("resolver '%s' is not installed", provider)
		}
		return err
	}

	return nil
}

// Template generates a starter resolver script
func Template(provider string) string {
	return fmt.Sprintf(`#!/bin/bash
# Seshmark resolver for %s
# Installed at: ~/.local/share/seshmark/resolvers/%s
#
# This script is called by: seshmark resume %s:<session-id>
# The session ID (without prefix) is passed as $1

SESSION_ID="$1"

# TODO: Implement your resolver logic here
# Examples:
#   - Open a native app: open "%s://$SESSION_ID"
#   - Open a URL: open "https://%s.com/session/$SESSION_ID"
#   - Print context: echo "Session: $SESSION_ID"
#   - Launch editor: cursor --open-session "$SESSION_ID"

echo "Resuming %s session: $SESSION_ID"
echo ""
echo "Edit this script at:"
echo "  ~/.local/share/seshmark/resolvers/%s"
echo ""
echo "Or create a new resolver:"
echo "  seshmark resolver create %s"
`, provider, provider, provider, provider, provider, provider, provider, provider)
}
