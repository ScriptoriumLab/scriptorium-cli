package project

import (
	"fmt"
	"os"
	"os/exec"
)

func cmakeConfig(dir string) error {
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

func cmakeBuild(dir string) error {
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

func ctestRun(testTag string, dir string) error {
	cmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^"+testTag+"$",
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
