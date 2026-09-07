package project

import (
	"fmt"
	"os"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/build/cmake"
)

type Felt struct {
	projectRootDir string
}

func (workspace *Workspace) Felt() *Felt {
	return &Felt{
		projectRootDir: workspace.root + `\scriptorium-felt`,
	}
}

func (felt *Felt) buildAndTest() error {
	fmt.Println("Building and testing Scriptorium Felt...")

	buildDir := felt.projectRootDir + `\build`

	fmt.Println("Cleaning existing Felt build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to clean Felt build directory: %w", err)
	}

	fmt.Println("Configuring Felt...")
	if err := cmake.Configure(felt.projectRootDir); err != nil {
		return fmt.Errorf("failed to configure Felt: %w", err)
	}

	fmt.Println("Building Felt...")
	if err := cmake.Build(felt.projectRootDir); err != nil {
		return fmt.Errorf("failed to build Felt: %w", err)
	}

	fmt.Println("Running Felt unit tests...")
	if err := cmake.RunTests(felt.projectRootDir, "felt-unit"); err != nil {
		return fmt.Errorf("scriptorium felt tests failed: %w", err)
	}

	return nil
}
