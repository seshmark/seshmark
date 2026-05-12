package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) (string, string) {
	tmpDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	// Build seshmark binary into tmpDir so hook can find it
	binPath := filepath.Join(tmpDir, "seshmark")
	buildCmd := exec.Command("go", "build", "-o", binPath, "github.com/seshmark/seshmark")
	buildCmd.Env = os.Environ()
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	return tmpDir, tmpDir
}

func TestHookInstallAndCommit(t *testing.T) {
	tmpDir, binDir := setupTestRepo(t)
	hookPath := filepath.Join(tmpDir, ".git", "hooks", "prepare-commit-msg")

	bin := filepath.Join(binDir, "seshmark")

	// Install hook
	installCmd := exec.Command(bin, "hook", "install")
	installCmd.Dir = tmpDir
	if out, err := installCmd.CombinedOutput(); err != nil {
		t.Fatalf("hook install failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("hook file not found: %v", err)
	}
	if !strings.Contains(string(data), "seshmark hook-run") {
		t.Errorf("hook does not delegate to seshmark")
	}

	// Track session
	trackCmd := exec.Command(bin, "track", "cursor:test-session", "--agent", "cursor", "--model", "claude-sonnet")
	trackCmd.Dir = tmpDir
	if out, err := trackCmd.CombinedOutput(); err != nil {
		t.Fatalf("track failed: %v\n%s", err, out)
	}

	// Create a file and commit
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello\n"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()

	// Commit with seshmark in PATH
	commitCmd := exec.Command("git", "-C", tmpDir, "commit", "-m", "feat: add test")
	commitCmd.Env = append(os.Environ(), "PATH="+binDir+":"+os.Getenv("PATH"))
	if out, err := commitCmd.CombinedOutput(); err != nil {
		t.Fatalf("commit failed: %v\n%s", err, out)
	}

	// Check trailers
	logCmd := exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%B")
	out, _ := logCmd.CombinedOutput()
	body := string(out)

	if !strings.Contains(body, "Seshmark-Version:") {
		t.Errorf("missing Seshmark-Version trailer, got:\n%s", body)
	}
	if !strings.Contains(body, "AI-Session: cursor:test-session") {
		t.Errorf("missing AI-Session trailer, got:\n%s", body)
	}
	if !strings.Contains(body, "AI-Agent: cursor") {
		t.Errorf("missing AI-Agent trailer, got:\n%s", body)
	}
}

func TestBranchInference(t *testing.T) {
	tmpDir, binDir := setupTestRepo(t)
	bin := filepath.Join(binDir, "seshmark")

	// Install hook
	hookInstall := exec.Command(bin, "hook", "install")
	hookInstall.Dir = tmpDir
	hookInstall.Run()

	// Create branch with agent prefix
	exec.Command("git", "-C", tmpDir, "checkout", "-b", "claude/feature").Run()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello\n"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()

	commitCmd := exec.Command("git", "-C", tmpDir, "commit", "-m", "feat: add test")
	commitCmd.Env = append(os.Environ(), "PATH="+binDir+":"+os.Getenv("PATH"))
	commitCmd.Run()

	logCmd := exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%B")
	out, _ := logCmd.CombinedOutput()
	body := string(out)

	if !strings.Contains(body, "AI-Agent: claude") {
		t.Errorf("branch inference failed, got:\n%s", body)
	}
	if strings.Contains(body, "AI-Session:") {
		t.Errorf("branch-only should not include AI-Session, got:\n%s", body)
	}
}

func TestWhoAndBlame(t *testing.T) {
	tmpDir, binDir := setupTestRepo(t)
	bin := filepath.Join(binDir, "seshmark")

	// Install hook
	hookInstall := exec.Command(bin, "hook", "install")
	hookInstall.Dir = tmpDir
	hookInstall.Run()

	// Track and commit
	trackCmd := exec.Command(bin, "track", "cursor:test", "--agent", "cursor")
	trackCmd.Dir = tmpDir
	trackCmd.Run()

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("line1\nline2\n"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "test.txt").Run()

	commitCmd := exec.Command("git", "-C", tmpDir, "commit", "-m", "feat: add lines")
	commitCmd.Env = append(os.Environ(), "PATH="+binDir+":"+os.Getenv("PATH"))
	commitCmd.Run()

	// Test who
	whoCmd := exec.Command(bin, "who", "HEAD", "--format", "json")
	whoCmd.Dir = tmpDir
	out, _ := whoCmd.CombinedOutput()
	if !strings.Contains(string(out), `"ai_agent": "cursor"`) {
		t.Errorf("who did not find agent, got: %s", out)
	}

	// Test blame
	blameCmd := exec.Command(bin, "blame", "test.txt", "--format", "json")
	blameCmd.Dir = tmpDir
	out, _ = blameCmd.CombinedOutput()
	if !strings.Contains(string(out), `"ai_agent": "cursor"`) {
		t.Errorf("blame did not find agent, got: %s", out)
	}
}
