package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/config"
	"github.com/seshmark/seshmark/internal/git"
)

// agentColors maps agent names to ANSI color codes for blame output
var agentColors = map[string]string{
	"cursor":         "33",       // yellow/orange
	"claude":         "38;5;173", // terracotta (256-color)
	"pi":             "32",       // green
	"opencode":       "34",       // blue
	"aider":          "35",       // magenta/purple
	"copilot":        "90",       // gray
	"github-copilot": "90",       // gray
	"codex":          "32",       // green
}

var blameCmd = &cobra.Command{
	Use:   "blame <file>",
	Short: "Show AI attribution per line",
	Long: `Show which AI agent wrote each line, with session and model info.

The output format is customizable. By default it shows [agent|session|model].
Use --fields to pick which fields to show and in what order.

Examples:
  git agentblame file.ts
  git agentblame file.ts --fields agent,model
  git agentblame file.ts --template "[{agent} @ {model}]"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		format, _ := cmd.Flags().GetString("format")
		fieldsStr, _ := cmd.Flags().GetString("fields")
		return runBlame(args[0], format, fieldsStr)
	},
}

func init() {
	blameCmd.Flags().String("format", "text", "Output format: text, json")
	blameCmd.Flags().String("fields", "", "Comma-separated fields to show: agent,session,model,commit,author,date (default from .seshmark.yml)")
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

func runBlame(file string, format string, fieldsStr string) error {
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

	// Resolve blame fields
	fields := resolveBlameFields(fieldsStr)

	for _, e := range entries {
		if e.IsHuman {
			fmt.Printf("\033[90m[human          ]\033[0m %s\n", e.Line)
		} else {
			label := formatBlameLabel(e, fields)
			ag := e.AIAgent
			if ag == "" {
				// Extract agent from session ID if available
				if e.AISession != "" {
					parts := strings.SplitN(e.AISession, ":", 2)
					ag = parts[0]
				} else {
					ag = "ai"
				}
			}
			color, ok := agentColors[strings.ToLower(ag)]
			if !ok {
				color = "36" // cyan fallback
			}
			fmt.Printf("\033[%sm%s\033[0m %s\n", color, label, e.Line)
		}
	}
	return nil
}

// resolveBlazeFields returns the fields to display, checking in order:
// 1. --fields CLI flag
// 2. .seshmark.yml config file
// 3. Default: [agent, session, model]
func resolveBlameFields(fieldsStr string) []string {
	// 1. CLI flag
	if fieldsStr != "" {
		parts := strings.Split(fieldsStr, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	// 2. Config file
	cfg, err := config.Load()
	if err == nil && len(cfg.Blame.Fields) > 0 {
		return cfg.Blame.Fields
	}

	// 3. Default
	return config.DefaultBlameFields
}

// formatBlameLabel builds the blame tag like [agent|session|model]
func formatBlameLabel(e blameEntry, fields []string) string {
	vals := make([]string, len(fields))
	for i, f := range fields {
		switch f {
		case "agent":
			ag := e.AIAgent
			if ag == "" && e.AISession != "" {
				parts := strings.SplitN(e.AISession, ":", 2)
				ag = parts[0]
			}
			vals[i] = truncate(ag, 10)
			if vals[i] == "" {
				vals[i] = "--"
			}
		case "session":
			sid := ""
			if e.AISession != "" {
				parts := strings.SplitN(e.AISession, ":", 2)
				if len(parts) > 1 {
					sid = parts[1]
				} else {
					sid = parts[0]
				}
			}
			vals[i] = truncate(sid, 12)
			if vals[i] == "" {
				vals[i] = "--"
			}
		case "model":
			vals[i] = truncate(e.AIModel, 14)
			if vals[i] == "" {
				vals[i] = "--"
			}
		case "commit":
			vals[i] = truncate(e.Commit, 8)
			if vals[i] == "" {
				vals[i] = "--"
			}
		case "author":
			vals[i] = truncate(e.Author, 10)
			if vals[i] == "" {
				vals[i] = "--"
			}
		case "date":
			vals[i] = truncate(e.Date, 10)
			if vals[i] == "" {
				vals[i] = "--"
			}
		default:
			vals[i] = "?"
		}
	}
	return "[" + strings.Join(vals, "│") + "]"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
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
