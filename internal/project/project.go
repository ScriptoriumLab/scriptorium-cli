// Package project contains functions to build and test the Scriptorium project components.
package project

import (
	"fmt"
	"os"
	"os/exec"
)

// TODO: Make the Scriptorium Project root directory configurable
const projectRootDir = `D:\Projects\Scriptorium`

const (
	feltProjectRootDir = projectRootDir + `\scriptorium-felt`

	brushProjectRootDir = projectRootDir + `\scriptorium-brush`

	inkstoneProjectRootDir = projectRootDir + `\scriptorium-inkstone`
	DictionarySourceFile = inkstoneProjectRootDir + `\data\pinyin_dictionary.txt`

	inkProjectRootDir = projectRootDir + `\scriptorium-ink`
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

	inkExe, err := buildInk()
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
	if err := cmakeConfig(feltProjectRootDir); err != nil {
		return fmt.Errorf("failed to configure Felt: %w", err)
	}

	fmt.Println("Building Felt...")
	buildCmd := exec.Command(
		"cmake",
		"--build", "build",
	)

	buildCmd.Dir = feltProjectRootDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Felt: %w", err)
	}

	fmt.Println("Running Felt unit tests...")
	testCmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^felt-unit$",
		"--schedule-random",
		"--output-on-failure",
	)

	testCmd.Dir = feltProjectRootDir
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	if err := testCmd.Run(); err != nil {
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
	if err := cmakeConfig(brushProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Brush: %w", err)
	}

	fmt.Println("Building Brush...")
	buildCmd := exec.Command(
		"cmake",
		"--build", "build",
	)

	buildCmd.Dir = brushProjectRootDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build Brush: %w", err)
	}

	fmt.Println("Running Brush unit tests...")
	testCmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^brush-unit$",
		"--schedule-random",
		"--output-on-failure",
	)

	testCmd.Dir = brushProjectRootDir
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	if err := testCmd.Run(); err != nil {
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
	if err := cmakeConfig(inkstoneProjectRootDir); err != nil {
		return "", fmt.Errorf("failed to configure Inkstone: %w", err)
	}

	fmt.Println("Building Inkstone...")
	buildCmd := exec.Command(
		"cmake",
		"--build", "build",
	)

	buildCmd.Dir = inkstoneProjectRootDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build Inkstone: %w", err)
	}

	fmt.Println("Running Inkstone unit tests...")
	unitTestCmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^inkstone-unit$",
		"--schedule-random",
		"--output-on-failure",
	)

	unitTestCmd.Dir = inkstoneProjectRootDir
	unitTestCmd.Stdout = os.Stdout
	unitTestCmd.Stderr = os.Stderr

	if err := unitTestCmd.Run(); err != nil {
		return "", fmt.Errorf("scriptorium inkstone unit tests failed: %w", err)
	}

	fmt.Println("Running Inkstone integration tests...")
	integrationTestCmd := exec.Command(
		"ctest",
		"--test-dir", "build",
		"-L", "^inkstone-integration$",
		"--schedule-random",
		"--output-on-failure",
	)

	integrationTestCmd.Dir = inkstoneProjectRootDir
	integrationTestCmd.Stdout = os.Stdout
	integrationTestCmd.Stderr = os.Stderr

	if err := integrationTestCmd.Run(); err != nil {
		return "", fmt.Errorf("scriptorium inkstone integration tests failed: %w", err)
	}

	exe := inkstoneProjectRootDir + `\build\ScriptoriumLabIME\scriptorium-inkstone.exe`
	return exe, nil
}

func buildInk() (string, error) {
	fmt.Println("Building Scriptorium Ink...")
	cmd := exec.Command(
		"pnpm",
		"tauri",
		"build",
	)

	cmd.Dir = inkProjectRootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("scriptorium ink build failed: %w", err)
	}

	exe := inkProjectRootDir + `\src-tauri\target\release\scriptorium-ink.exe`
	return exe, nil
}

func cmakeConfig(dir string) error {
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = dir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
		return fmt.Errorf("failed to configure project: %w", err)
	}

	return nil
}
