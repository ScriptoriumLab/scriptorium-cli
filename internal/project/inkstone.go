package project

import (
	"fmt"
	"os"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/build/cmake"
)

type Inkstone struct {
	projectRootDir string
}

func (workspace *Workspace) Inkstone() *Inkstone {
	return &Inkstone{
		projectRootDir: workspace.root + `\scriptorium-inkstone`,
	}
}

func (inkstone *Inkstone)buildAndTest() (string, error) {
	fmt.Println("Building and testing Scriptorium Inkstone...")

	buildDir := inkstone.projectRootDir + `\build`

	fmt.Println("Cleaning existing Inkstone build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return "", fmt.Errorf("failed to clean Inkstone build directory: %w", err)
	}

	fmt.Println("Configuring Inkstone...")
	if err := cmake.Configure(inkstone.projectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Inkstone: %w", err)
	}

	fmt.Println("Building Inkstone...")
	if err := cmake.Build(inkstone.projectRootDir); err != nil {
		return "", fmt.Errorf("failed to build Inkstone: %w", err)
	}

	fmt.Println("Running Inkstone unit tests...")
	if err := cmake.RunTests(inkstone.projectRootDir, "inkstone-unit"); err != nil {
		return "", fmt.Errorf("scriptorium inkstone unit tests failed: %w", err)
	}

	fmt.Println("Running Inkstone integration tests...")
	if err := cmake.RunTests(inkstone.projectRootDir, "inkstone-integration"); err != nil {
		return "", fmt.Errorf("scriptorium inkstone integration tests failed: %w", err)
	}

	exe := inkstone.projectRootDir + `\build\ScriptoriumLabIME\scriptorium-inkstone.exe`
	return exe, nil
}
