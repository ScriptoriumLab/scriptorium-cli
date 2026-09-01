package project

import (
	"fmt"
	"os"
	"os/exec"
)

func pnpmBuild(dir string) error {
	cmd := exec.Command(
		"pnpm",
		"tauri",
		"build",
	)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("project build failed: %w", err)
	}

	return nil
}
