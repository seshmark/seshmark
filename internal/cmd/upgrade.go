package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/seshmark/seshmark/internal/hook"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Check for updates and upgrade seshmark",
	Long: `Check the latest release on GitHub and upgrade to it if newer.

Downloads the latest binary, replaces the current installation,
and re-installs Git hooks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpgrade()
	},
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func runUpgrade() error {
	// 1. Get current version
	currentVer := rootCmd.Version
	if currentVer == "" || currentVer == "dev" {
		fmt.Println("⚠️  Development build — skipping version check")
	}

	fmt.Print("🔍 Checking for updates... ")

	// 2. Fetch latest release from GitHub
	latest, err := fetchLatestRelease()
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("could not check for updates: %w", err)
	}

	latestVer := strings.TrimPrefix(latest.TagName, "v")
	currentVer = strings.TrimPrefix(currentVer, "v")

	// Strip git describe suffix for comparison
	if idx := strings.Index(currentVer, "-"); idx > 0 {
		currentVer = currentVer[:idx]
	}

	if currentVer == latestVer {
		fmt.Println("already up-to-date!")
		fmt.Printf("  Current: v%s\n", currentVer)
		return nil
	}

	// If we can't do a proper semver compare, just check if tags differ
	if currentVer != "dev" && currentVer != "" && latestVer != "" {
		fmt.Printf("found v%s\n", latestVer)
		fmt.Printf("  Current: v%s → Latest: v%s\n", currentVer, latestVer)
	} else {
		fmt.Printf("found %s\n", latest.TagName)
	}

	// 3. Determine current binary path
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine binary path: %w", err)
	}

	// Resolve symlinks to get the real path
	realPath, err := filepath.EvalSymlinks(binaryPath)
	if err == nil {
		binaryPath = realPath
	}

	// 4. Determine OS/arch for download
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goarch {
	case "x86_64":
		goarch = "amd64"
	case "aarch64":
		goarch = "arm64"
	}

	fileName := fmt.Sprintf("seshmark-%s-%s", goos, goarch)
	if goos == "windows" {
		fileName += ".exe"
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/seshmark/seshmark/releases/download/%s/%s",
		latest.TagName, fileName,
	)

	// 5. Download to temp file
	fmt.Println("\n⬇️  Downloading...")

	tmpFile := filepath.Join(os.TempDir(), "seshmark-upgrade")
	if goos == "windows" {
		tmpFile += ".exe"
	}

	if err := downloadFile(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if err := os.Chmod(tmpFile, 0755); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("could not set permissions: %w", err)
	}

	// 6. Verify the downloaded binary works
	output, err := exec.Command(tmpFile, "--version").Output()
	if err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("downloaded binary failed verification: %w", err)
	}

	downloadedVer := strings.TrimSpace(string(output))
	fmt.Printf("  Downloaded: %s\n", downloadedVer)

	// 7. Replace current binary
	backupPath := binaryPath + ".bak"
	os.Remove(backupPath) // remove any previous backup

	if err := os.Rename(binaryPath, backupPath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("could not back up current binary: %w", err)
	}

	if err := os.Rename(tmpFile, binaryPath); err != nil {
		// Restore backup
		os.Rename(backupPath, binaryPath)
		return fmt.Errorf("could not replace binary: %w", err)
	}

	fmt.Printf("  Installed: %s\n", binaryPath)

	// 8. Re-install hooks (preserving --global if applicable)
	fmt.Print("🔧 Re-installing hooks... ")
	if err := hook.Install(false, false, true); err != nil {
		fmt.Println("warning: could not re-install hook")
		fmt.Printf("  You can re-run: %s hook install\n", binaryPath)
	} else {
		fmt.Println("done")
	}

	// 9. Clean up backup
	os.Remove(backupPath)

	fmt.Println("\n✅ Upgrade complete!")
	fmt.Printf("   Run '%s --version' to verify.\n", filepath.Base(binaryPath))

	return nil
}

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/seshmark/seshmark/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "seshmark")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	if release.TagName == "" {
		return nil, fmt.Errorf("no releases found")
	}

	return &release, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned %d for %s", resp.StatusCode, url)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
