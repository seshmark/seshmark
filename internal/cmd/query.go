package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var queryCmd = &cobra.Command{
	Use:   "query [terms...]",
	Short: "Search AI-tagged commits",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		return runQuery(cmd, args)
	},
}

func init() {
	queryCmd.Flags().String("agent", "", "Filter by AI-Agent")
	queryCmd.Flags().String("model", "", "Filter by AI-Model")
	queryCmd.Flags().String("session", "", "Filter by AI-Session")
	queryCmd.Flags().String("path", "", "Filter by changed files")
	queryCmd.Flags().String("since", "", "Commits after date")
	queryCmd.Flags().String("until", "", "Commits before date")
	queryCmd.Flags().String("format", "text", "Output format: text, json")
}

type queryResult struct {
	Hash         string   `json:"hash"`
	ShortHash    string   `json:"short_hash"`
	Author       string   `json:"author"`
	Date         string   `json:"date"`
	Subject      string   `json:"subject"`
	AISession    string   `json:"ai_session"`
	AIAgent      string   `json:"ai_agent"`
	AIModel      string   `json:"ai_model"`
	FilesChanged []string `json:"files_changed"`
}

func runQuery(cmd *cobra.Command, terms []string) error {
	agent, _ := cmd.Flags().GetString("agent")
	model, _ := cmd.Flags().GetString("model")
	session, _ := cmd.Flags().GetString("session")
	pathFilter, _ := cmd.Flags().GetString("path")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	format, _ := cmd.Flags().GetString("format")

	// Get all commits with Seshmark-Version
	gitArgs := []string{"log", "--all", "--no-merges", "--grep=Seshmark-Version:", "--format=%H"}
	if since != "" {
		gitArgs = append(gitArgs, "--since="+since)
	}
	if until != "" {
		gitArgs = append(gitArgs, "--until="+until)
	}

	out, err := git.Exec(gitArgs...)
	if err != nil {
		return err
	}
	if out == "" {
		fmt.Println("No AI commits found.")
		return nil
	}

	var results []queryResult
	for _, hash := range strings.Fields(out) {
		// Parse commit body
		body, _ := git.Exec("log", "-1", "--format=%B", hash)
		subj, _ := git.Exec("log", "-1", "--format=%s", hash)
		date, _ := git.Exec("log", "-1", "--format=%ai", hash)
		authorName, _ := git.Exec("log", "-1", "--format=%an", hash)

		meta := parseTrailers(body)

		// Apply filters
		if agent != "" && !strings.Contains(meta["agent"], agent) {
			continue
		}
		if model != "" && !strings.Contains(meta["model"], model) {
			continue
		}
		if session != "" && !strings.Contains(meta["session"], session) {
			continue
		}

		// Path filter
		files := []string{}
		if pathFilter != "" {
			filesOut, _ := git.Exec("diff-tree", "--no-commit-id", "--name-only", "-r", hash)
			files = strings.Split(filesOut, "\n")
			found := false
			for _, f := range files {
				if strings.Contains(f, pathFilter) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Term search
		if len(terms) > 0 {
			found := false
			for _, t := range terms {
				if strings.Contains(strings.ToLower(subj), strings.ToLower(t)) ||
					strings.Contains(strings.ToLower(body), strings.ToLower(t)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		short, _ := git.Exec("log", "-1", "--format=%h", hash)
		results = append(results, queryResult{
			Hash:         hash,
			ShortHash:    short,
			Author:       strings.TrimSpace(authorName),
			Date:         date,
			Subject:      subj,
			AISession:    meta["session"],
			AIAgent:      meta["agent"],
			AIModel:      meta["model"],
			FilesChanged: files,
		})
	}

	if format == "json" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	for _, r := range results {
		fmt.Printf("%s | %s | %s | %s | %s\n", r.ShortHash, r.Date, r.AIAgent, r.AIModel, r.Subject)
	}
	return nil
}

