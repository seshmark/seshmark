package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
)

var resumeCmd = &cobra.Command{
	Use:   "resume [commit|<provider>:<id>]",
	Short: "Resume a session via native resolver or context",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		
		target := "HEAD"
		if len(args) > 0 {
			target = args[0]
		}
		
		format, _ := cmd.Flags().GetString("format")
		return runResume(target, format)
	},
}

func init() {
	resumeCmd.Flags().String("format", "text", "Output format: text, json")
}

func runResume(target, format string) error {
	var sessionID string
	
	// Resolve target to session ID
	if strings.Contains(target, ":") {
		sessionID = target
	} else {
		body, err := git.Exec("log", "-1", "--format=%B", target)
		if err != nil {
			return err
		}
		meta := parseTrailers(body)
		sessionID = meta["session"]
		if sessionID == "" {
			return fmt.Errorf("no resumable session found for %s", target)
		}
	}
	
	parts := strings.SplitN(sessionID, ":", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid session ID: %s", sessionID)
	}
	provider := parts[0]
	id := parts[1]
	
	resolverPath := filepath.Join(config.DataDir(), "resolvers", provider)
	_, err := os.Stat(resolverPath)
	resolverAvailable := err == nil
	
	if format == "json" {
		// Just return metadata without executing
		fmt.Printf(`{"session_id": "%s", "provider": "%s", "resolver_available": %v}%s`,
			sessionID, provider, resolverAvailable, "\n")
		return nil
	}
	
	if resolverAvailable {
		fmt.Printf("Resuming %s via %s...\n", sessionID, provider)
		cmd := exec.Command(resolverPath, id)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	
	// Fallback: print context
	fmt.Printf("No resolver for '%s'.\n", provider)
	fmt.Printf("Session: %s\n\n", sessionID)
	
	// Try to find a commit with this session and show context
	commits, _ := git.Exec("log", "--all", "--grep=AI-Session: "+sessionID, "--format=%H", "--no-merges")
	if commits != "" {
		firstCommit := strings.Fields(commits)[0]
		return runContext(firstCommit, "markdown")
	}
	
	return nil
}
