// Package project contains functions to build and test the Scriptorium project components.
package project

import (
	"fmt"
	"os"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/build/cmake"
)

// TODO: Make the Scriptorium Project root directory configurable
const projectRootDir = `D:\Projects\Scriptorium`

const (
	feltProjectRootDir = projectRootDir + `\scriptorium-felt`

	brushProjectRootDir = projectRootDir + `\scriptorium-brush`

	inkstoneProjectRootDir = projectRootDir + `\scriptorium-inkstone`
	DictionarySourceFile = inkstoneProjectRootDir + `\data\pinyin_dictionary.txt`
)

type ProjectArtifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

func BuildScriptoriumAndRunAllTests() (*ProjectArtifacts, error) {
	fmt.Println("Building Scriptorium and running all tests...")
	if err := buildAndTestFelt(); err != nil {
		return nil, err
	}

	brushDll, err := buildAndTestBrush()
	if err != nil {
		return nil, err
	}

	inkstoneExe, err := buildAndTestInkstone()
	if err != nil {
		return nil, err
	}

	inkExe, err := NewInk(projectRootDir).build()
	if err != nil {
		return nil, err
	}

	return &ProjectArtifacts{
		BrushDLL:    brushDll,
		InkstoneEXE: inkstoneExe,
		InkEXE:      inkExe,
	}, nil
}

func buildAndTestFelt() error {
	fmt.Println("Building and testing Scriptorium Felt...")

	buildDir := feltProjectRootDir + `\build`

	fmt.Println("Cleaning existing Felt build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to clean Felt build directory: %w", err)
	}

	fmt.Println("Configuring Felt...")
	if err := cmake.Configure(feltProjectRootDir); err != nil {
		return fmt.Errorf("failed to configure Felt: %w", err)
	}

	fmt.Println("Building Felt...")
	if err := cmake.Build(feltProjectRootDir); err != nil {
		return fmt.Errorf("failed to build Felt: %w", err)
	}

	fmt.Println("Running Felt unit tests...")
	if err := cmake.RunTests(feltProjectRootDir, "felt-unit"); err != nil {
		return fmt.Errorf("scriptorium felt tests failed: %w", err)
	}

	return nil
}

func buildAndTestBrush() (string, error) {
	fmt.Println("Building and testing Scriptorium Brush...")

	buildDir := brushProjectRootDir + `\build`

	fmt.Println("Cleaning existing Brush build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return "", fmt.Errorf("failed to clean Brush build directory: %w", err)
	}

	fmt.Println("Configuring Brush...")
	if err := cmake.Configure(brushProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Brush: %w", err)
	}

	fmt.Println("Building Brush...")
	if err := cmake.Build(brushProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to build Brush: %w", err)
	}

	fmt.Println("Running Brush unit tests...")
	if err := cmake.RunTests(brushProjectRootDir, "brush-unit"); err != nil {
		return "", fmt.Errorf("scriptorium brush tests failed: %w", err)
	}

	artifact := brushProjectRootDir + `\build\ScriptoriumLabIME\scriptorium-brush.dll`
	return artifact, nil
}

func buildAndTestInkstone() (string, error) {
	fmt.Println("Building and testing Scriptorium Inkstone...")

	buildDir := inkstoneProjectRootDir + `\build`

	fmt.Println("Cleaning existing Inkstone build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return "", fmt.Errorf("failed to clean Inkstone build directory: %w", err)
	}

	fmt.Println("Configuring Inkstone...")
	if err := cmake.Configure(inkstoneProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Inkstone: %w", err)
	}

	fmt.Println("Building Inkstone...")
	if err := cmake.Build(inkstoneProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to build Inkstone: %w", err)
	}

	fmt.Println("Running Inkstone unit tests...")
	if err := cmake.RunTests(inkstoneProjectRootDir, "inkstone-unit"); err != nil {
		return "", fmt.Errorf("scriptorium inkstone unit tests failed: %w", err)
	}

	fmt.Println("Running Inkstone integration tests...")
	if err := cmake.RunTests(inkstoneProjectRootDir, "inkstone-integration"); err != nil {
		return "", fmt.Errorf("scriptorium inkstone integration tests failed: %w", err)
	}

	exe := inkstoneProjectRootDir + `\build\ScriptoriumLabIME\scriptorium-inkstone.exe`
	return exe, nil
}
