// Package cmake provides CMake and CTest operations for CMake-based projects.
package cmake

import (
	"fmt"
	"os"
	"os/exec"
)

func Configure(dir string) error {
	cmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure project: %w", err)
	}

	return nil
}

func Build(dir string) error {
	cmd := exec.Command(
		"cmake",
		"--build", "build",
	)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build project: %w", err)
	}

	return nil
}

func RunTests(dir string, label string) error {
	cmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^"+label+"$",
		"--schedule-random",
		"--output-on-failure",
	)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run tests failed: %w", err)
	}

	return nil
}
