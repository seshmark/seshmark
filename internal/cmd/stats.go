package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/git"
)

var statsCmd = &cobra.Command{
	Use:   "stats [--since <date>] [--until <date>] [--format json|markdown]",
	Short: "Generate a shareable AI code report",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireGitRepo(); err != nil {
			return err
		}
		return runStats(cmd)
	},
}

func init() {
	statsCmd.Flags().String("since", "", "Start date (e.g., 'last week', '2025-01-01')")
	statsCmd.Flags().String("until", "", "End date")
	statsCmd.Flags().String("format", "text", "Output format: text, json, markdown")
	rootCmd.AddCommand(statsCmd)
}

type statsResult struct {
	TotalCommits  int                `json:"total_commits"`
	AICommits     int                `json:"ai_commits"`
	HumanCommits  int                `json:"human_commits"`
	AIPercent     float64            `json:"ai_percent"`
	ByAgent       map[string]int     `json:"by_agent"`
	ByModel       map[string]int     `json:"by_model"`
	ByFile        []fileStat         `json:"by_file"`
	FirstAICommit string             `json:"first_ai_commit,omitempty"`
	LastAICommit  string             `json:"last_ai_commit,omitempty"`
}

type fileStat struct {
	Path     string  `json:"path"`
	AICommits int     `json:"ai_commits"`
	TotalCommits int `json:"total_commits"`
	AIPercent float64 `json:"ai_percent"`
}

func runStats(cmd *cobra.Command) error {
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	format, _ := cmd.Flags().GetString("format")

	gitArgs := []string{"log", "--all", "--format=%H|%ai|%s"}
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
		fmt.Println("No commits found.")
		return nil
	}

	result := statsResult{
		ByAgent: make(map[string]int),
		ByModel: make(map[string]int),
	}

	fileStats := make(map[string]*fileStat)

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
		hash := parts[0]
		date := parts[1]
		// subject := parts[2]

		result.TotalCommits++

		body, _ := git.Exec("log", "-1", "--format=%B", hash)
		meta := parseTrailers(body)

		if meta["agent"] == "" && meta["session"] == "" {
			result.HumanCommits++
			continue
		}

		result.AICommits++
		if result.FirstAICommit == "" || date < result.FirstAICommit {
			result.FirstAICommit = date
		}
		if date > result.LastAICommit {
			result.LastAICommit = date
		}

		if meta["agent"] != "" {
			result.ByAgent[meta["agent"]]++
		}
		if meta["model"] != "" {
			result.ByModel[meta["model"]]++
		}

		// File stats
		files, _ := git.Exec("diff-tree", "--no-commit-id", "--name-only", "-r", hash)
		for _, f := range strings.Split(files, "\n") {
			if f == "" {
				continue
			}
			if fileStats[f] == nil {
				fileStats[f] = &fileStat{Path: f}
			}
			fileStats[f].TotalCommits++
			fileStats[f].AICommits++
		}
	}

	if result.TotalCommits > 0 {
		result.AIPercent = float64(result.AICommits) / float64(result.TotalCommits) * 100
	}

	// Build sorted file list
	for _, fs := range fileStats {
		fs.AIPercent = float64(fs.AICommits) / float64(fs.TotalCommits) * 100
		result.ByFile = append(result.ByFile, *fs)
	}
	sort.Slice(result.ByFile, func(i, j int) bool {
		return result.ByFile[i].AIPercent > result.ByFile[j].AIPercent
	})
	if len(result.ByFile) > 10 {
		result.ByFile = result.ByFile[:10]
	}

	if format == "json" {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if format == "markdown" {
		printMarkdownStats(result)
		return nil
	}

	printTextStats(result)
	return nil
}

func printTextStats(result statsResult) {
	fmt.Println()
	fmt.Println("AI Code Report")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("Total commits:          %d\n", result.TotalCommits)
	fmt.Printf("AI commits:             %d (%.1f%%)\n", result.AICommits, result.AIPercent)
	fmt.Printf("Human commits:          %d (%.1f%%)\n", result.HumanCommits, 100-result.AIPercent)
	fmt.Println()

	if len(result.ByAgent) > 0 {
		fmt.Println("By Agent:")
		printBarChart(result.ByAgent, result.AICommits)
		fmt.Println()
	}

	if len(result.ByModel) > 0 {
		fmt.Println("By Model:")
		printBarChart(result.ByModel, result.AICommits)
		fmt.Println()
	}

	if len(result.ByFile) > 0 {
		fmt.Println("By File (most AI):")
		for _, f := range result.ByFile {
			bar := renderBar(f.AIPercent, 20)
			fmt.Printf("  %-30s %s %.0f%%\n", f.Path, bar, f.AIPercent)
		}
		fmt.Println()
	}

	if result.FirstAICommit != "" {
		fmt.Printf("First AI commit:        %s\n", result.FirstAICommit)
		fmt.Printf("Latest AI commit:       %s\n", result.LastAICommit)
	}
	fmt.Println()
}

func printMarkdownStats(result statsResult) {
	fmt.Println("## AI Code Report")
	fmt.Println()
	fmt.Printf("| Metric | Value |\n")
	fmt.Printf("|--------|-------|\n")
	fmt.Printf("| Total commits | %d |\n", result.TotalCommits)
	fmt.Printf("| AI commits | %d (%.1f%%) |\n", result.AICommits, result.AIPercent)
	fmt.Printf("| Human commits | %d (%.1f%%) |\n", result.HumanCommits, 100-result.AIPercent)
	fmt.Println()

	if len(result.ByAgent) > 0 {
		fmt.Println("### By Agent")
		fmt.Println()
		fmt.Println("| Agent | Commits | % |")
		fmt.Println("|-------|---------|---|")
		for _, pair := range sortedPairs(result.ByAgent) {
			pct := float64(pair.value) / float64(result.AICommits) * 100
			fmt.Printf("| %s | %d | %.1f%% |\n", pair.key, pair.value, pct)
		}
		fmt.Println()
	}

	if len(result.ByFile) > 0 {
		fmt.Println("### By File (most AI)")
		fmt.Println()
		fmt.Println("| File | AI % |")
		fmt.Println("|------|------|")
		for _, f := range result.ByFile {
			fmt.Printf("| %s | %.0f%% |\n", f.Path, f.AIPercent)
		}
		fmt.Println()
	}

	fmt.Println("*Generated by [Seshmark](https://seshmark.dev)*")
}

func renderBar(percent float64, width int) string {
	filled := int(percent / 100.0 * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

func printBarChart(data map[string]int, total int) {
	pairs := sortedPairs(data)
	maxKeyLen := 0
	for _, p := range pairs {
		if len(p.key) > maxKeyLen {
			maxKeyLen = len(p.key)
		}
	}
	for _, p := range pairs {
		pct := float64(p.value) / float64(total) * 100
		bar := renderBar(pct, 20)
		fmt.Printf("  %-*s %s %d (%.0f%%)\n", maxKeyLen, p.key, bar, p.value, pct)
	}
}

type kvPair struct {
	key   string
	value int
}

func sortedPairs(m map[string]int) []kvPair {
	var pairs []kvPair
	for k, v := range m {
		pairs = append(pairs, kvPair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value > pairs[j].value
	})
	return pairs
}
