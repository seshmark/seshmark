package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var blameCmd = &cobra.Command{
	Use:   "blame <file>",
	Short: "Show AI attribution per line",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		format, _ := cmd.Flags().GetString("format")
		return runBlame(args[0], format)
	},
}

func init() {
	blameCmd.Flags().String("format", "text", "Output format: text, json")
}

type blameEntry struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Commit     string `json:"commit"`
	Author     string `json:"author"`
	Date       string `json:"date"`
	AISession  string `json:"ai_session"`
	AIAgent    string `json:"ai_agent"`
	AIModel    string `json:"ai_model"`
	IsHuman    bool   `json:"is_human"`
}

func runBlame(file string, format string) error {
	out, err := git.Exec("blame", "--porcelain", file)
	if err != nil {
		return err
	}

	entries := parsePorcelainBlame(out)
	
	if format == "json" {
		data, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	for _, e := range entries {
		if e.IsHuman {
			fmt.Printf("\033[90m[human          ]\033[0m %s\n", e.Line)
		} else {
			prov := "---"
			sid := "---"
			if e.AISession != "" {
				parts := strings.SplitN(e.AISession, ":", 2)
				prov = parts[0]
				if len(parts) > 1 {
					sid = parts[1]
					if len(sid) > 12 {
						sid = sid[:12]
					}
				}
			}
			ag := e.AIAgent
			if len(ag) > 10 {
				ag = ag[:10]
			}
			fmt.Printf("\033[36m[%s|%s|%s]\033[0m %s\n", prov, sid, ag, e.Line)
		}
	}
	return nil
}

func parsePorcelainBlame(raw string) []blameEntry {
	var entries []blameEntry
	var current *blameEntry
	
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		
		// New commit block starts with hash
		if matched, _ := regexp.MatchString(`^[0-9a-f]{40}`, line); matched {
			if current != nil {
				entries = append(entries, *current)
			}
			parts := strings.Fields(line)
			current = &blameEntry{
				Commit: parts[0],
			}
			continue
		}
		
		if current == nil {
			continue
		}
		
		if strings.HasPrefix(line, "author ") {
			current.Author = strings.TrimPrefix(line, "author ")
		} else if strings.HasPrefix(line, "author-time ") {
			ts := strings.TrimPrefix(line, "author-time ")
			// Parse unix timestamp
			// current.Date = format time
			_ = ts
		} else if strings.HasPrefix(line, "\t") {
			current.Line = line[1:]
			current.LineNumber = len(entries) + 1
		}
	}
	
	if current != nil {
		entries = append(entries, *current)
	}
	
	// Now fetch AI metadata for each commit
	commitMeta := make(map[string]map[string]string)
	for _, e := range entries {
		if _, ok := commitMeta[e.Commit]; !ok {
			meta := fetchCommitMeta(e.Commit)
			commitMeta[e.Commit] = meta
		}
	}
	
	for i := range entries {
		meta := commitMeta[entries[i].Commit]
		entries[i].AISession = meta["session"]
		entries[i].AIAgent = meta["agent"]
		entries[i].AIModel = meta["model"]
		entries[i].IsHuman = meta["session"] == "" && meta["agent"] == ""
	}
	
	return entries
}

func fetchCommitMeta(commit string) map[string]string {
	result := map[string]string{"session": "", "agent": "", "model": ""}
	body, err := git.Exec("log", "-1", "--format=%B", commit)
	if err != nil {
		return result
	}
	
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "AI-Session: ") {
			result["session"] = strings.TrimPrefix(line, "AI-Session: ")
		} else if strings.HasPrefix(line, "AI-Agent: ") {
			result["agent"] = strings.TrimPrefix(line, "AI-Agent: ")
		} else if strings.HasPrefix(line, "AI-Model: ") {
			result["model"] = strings.TrimPrefix(line, "AI-Model: ")
		}
	}
	return result
}
