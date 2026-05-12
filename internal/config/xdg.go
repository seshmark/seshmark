package config

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "seshmark")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "seshmark")
}

func DataDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "seshmark")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "seshmark")
}

func CacheDir() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "seshmark")
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", "seshmark")
}

// LegacyDir returns the old ~/.seshmark path for backward compatibility
func LegacyDir() string {
	return filepath.Join(os.Getenv("HOME"), ".seshmark")
}
