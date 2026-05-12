package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var whoCmd = &cobra.Command{
	Use:   "who [commit]",
	Short: "Show AI metadata for a commit",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		commit := "HEAD"
		if len(args) > 0 {
			commit = args[0]
		}
		format, _ := cmd.Flags().GetString("format")
		return runWho(commit, format)
	},
}

func init() {
	whoCmd.Flags().String("format", "text", "Output format: text, json")
}

func runWho(commit, format string) error {
	body, err := git.Exec("log", "-1", "--format=%B", commit)
	if err != nil {
		return err
	}

	meta := make(map[string]string)
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "Seshmark-Version: ") {
			meta["seshmark_version"] = strings.TrimPrefix(line, "Seshmark-Version: ")
		} else if strings.HasPrefix(line, "AI-Session: ") {
			meta["ai_session"] = strings.TrimPrefix(line, "AI-Session: ")
		} else if strings.HasPrefix(line, "AI-Agent: ") {
			meta["ai_agent"] = strings.TrimPrefix(line, "AI-Agent: ")
		} else if strings.HasPrefix(line, "AI-Model: ") {
			meta["ai_model"] = strings.TrimPrefix(line, "AI-Model: ")
		}
	}

	if format == "json" {
		out := map[string]interface{}{
			"hash":           commit,
			"has_seshmark":   meta["seshmark_version"] != "",
			"seshmark_version": meta["seshmark_version"],
			"ai_session":     meta["ai_session"],
			"ai_agent":       meta["ai_agent"],
			"ai_model":       meta["ai_model"],
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if meta["seshmark_version"] == "" {
		fmt.Println("No AI metadata.")
		return nil
	}

	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Seshmark-Version: %s\n", meta["seshmark_version"])
	if meta["ai_session"] != "" {
		fmt.Printf("AI-Session: %s\n", meta["ai_session"])
	}
	if meta["ai_agent"] != "" {
		fmt.Printf("AI-Agent: %s\n", meta["ai_agent"])
	}
	if meta["ai_model"] != "" {
		fmt.Printf("AI-Model: %s\n", meta["ai_model"])
	}
	return nil
}
