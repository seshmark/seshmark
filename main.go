package main

import (
	"os"
	"path/filepath"

	"github.com/seshmark/seshmark/internal/cmd"
	"github.com/seshmark/seshmark/pkg/version"
)

func main() {
	base := filepath.Base(os.Args[0])
	switch base {
	case "agentblame", "git-agentblame", "aiblame", "git-aiblame":
		os.Args = append([]string{os.Args[0], "blame"}, os.Args[1:]...)
	}

	cmd.Execute(version.Version)
}
