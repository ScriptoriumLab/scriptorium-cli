package project

import (
	"fmt"
	"os"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/build/cmake"
)

type Brush struct {
	projectRootDir string
}

func (workspace *Workspace) Brush() *Brush {
	return &Brush{
		projectRootDir: workspace.root + `\scriptorium-brush`,
	}
}

func (brush *Brush) buildAndTest() (string, error) {
	fmt.Println("Building and testing Scriptorium Brush...")

	buildDir := brush.projectRootDir + `\build`

	fmt.Println("Cleaning existing Brush build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return "", fmt.Errorf("failed to clean Brush build directory: %w", err)
	}

	fmt.Println("Configuring Brush...")
	if err := cmake.Configure(brush.projectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Brush: %w", err)
	}

	fmt.Println("Building Brush...")
	if err := cmake.Build(brush.projectRootDir); err != nil {
		return "", fmt.Errorf("failed to build Brush: %w", err)
	}

	fmt.Println("Running Brush unit tests...")
	if err := cmake.RunTests(brush.projectRootDir, "brush-unit"); err != nil {
		return "", fmt.Errorf("scriptorium brush tests failed: %w", err)
	}

	artifact := brush.projectRootDir + `\build\ScriptoriumLabIME\scriptorium-brush.dll`
	return artifact, nil
}
