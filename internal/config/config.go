package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// BlameConfig controls the blame output format
type BlameConfig struct {
	Fields   []string `yaml:"fields"`
	Template string   `yaml:"template"`
}

// Config is the root seshmark configuration
type Config struct {
	Blame BlameConfig `yaml:"blame"`
}

// DefaultBlameFields is used when no config is provided
var DefaultBlameFields = []string{"agent", "session", "model"}

// Load reads the config from the repo or global location.
// It checks in order: repo .seshmark.yml → ~/.config/seshmark/config.yml
// Returns a default config if neither exists.
func Load() (*Config, error) {
	cfg := &Config{
		Blame: BlameConfig{
			Fields: DefaultBlameFields,
		},
	}

	// Try repo-level config first
	repoCfg := findRepoConfig()
	if repoCfg != "" {
		if err := loadFile(repoCfg, cfg); err == nil {
			return cfg, nil
		}
	}

	// Fall back to global config
	globalCfg := filepath.Join(ConfigDir(), "config.yml")
	if data, err := os.ReadFile(globalCfg); err == nil {
		if err := yaml.Unmarshal(data, cfg); err == nil {
			return cfg, nil
		}
	}

	return cfg, nil
}

// ValidateFields checks that all requested fields are valid
func ValidateFields(fields []string) error {
	valid := map[string]bool{
		"agent":   true,
		"session": true,
		"model":   true,
		"commit":  true,
		"author":  true,
		"date":    true,
	}
	for _, f := range fields {
		if !valid[f] {
			return fmt.Errorf("unknown blame field: %q (valid: agent, session, model, commit, author, date)", f)
		}
	}
	if len(fields) == 0 {
		return fmt.Errorf("at least one field is required")
	}
	return nil
}

func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func findRepoConfig() string {
	// Look for .seshmark.yml in current directory and parents
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		candidate := filepath.Join(dir, ".seshmark.yml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
