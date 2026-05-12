package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
