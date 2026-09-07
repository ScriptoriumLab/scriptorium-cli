// Package pnpm provides functions to build Tauri-based projects using pnpm.
package pnpm

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildTauri(dir string) error {
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
