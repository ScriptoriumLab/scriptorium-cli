// Package dev provides the implementation of the `dev` subcommand.
package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type devEnv string

const (
	devEnvVM devEnv = "vm"
)

var env string

const devUseCaseTaskName = "Scriptorium Dev Use Case"

// TODO: Discover vmrun dynamically instead of relying on the default
// VMware Workstation installation path.
const vmrunPath = `C:\Program Files\VMware\VMware Workstation\vmrun.exe`

// TODO: Make the development VM path configurable.
const devVMPath = `D:\Projects\Scriptorium\dev-env\scriptorium-dev\scriptorium-dev.vmx`

// TODO: Make the development VM snapshot configurable.
const devVMSnapshot = "baseline"

// TODO: Make the Scriptorium Project root directory configurable
const projectRootDir = `D:\Projects\Scriptorium`

const (
	feltProjectRootDir = projectRootDir + `\scriptorium-felt`

	brushProjectRootDir = projectRootDir + `\scriptorium-brush`

	inkstoneProjectRootDir = projectRootDir + `\scriptorium-inkstone`
	dictionarySourceFile = inkstoneProjectRootDir + `\data\pinyin_dictionary.txt`

	inkProjectRootDir = projectRootDir + `\scriptorium-ink`
)

// TODO: Make the Scriptorium product root directory configurable.
const productRootDir = `C:\Users\dev\Scriptorium`

const (
	productLocalDir = productRootDir + `\Local`
	productDictionaryDir = productLocalDir + `\pinyin_dictionary.txt`
	productLogDir = productRootDir + `\Log`
)

const productArtifactsDir = `C:\Users\dev\Desktop\ScriptoriumArtifacts`

const (
	productBrushDLL = productArtifactsDir + `\scriptorium-brush.dll`
	productInkstoneEXE = productArtifactsDir + `\scriptorium-inkstone.exe`
	productInkEXE = productArtifactsDir + `\scriptorium-ink.exe`
)

type ProjectArtifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

type Config struct {
	VMEncryptionPassword string
	GuestUsername        string
	GuestPassword        string
}

func ensureVMAvailable() error {
	fmt.Println("Ensuring VMware is available...")
	if _, err := os.Stat(vmrunPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VMware CLI 'vmrun' was not found at %s", vmrunPath)
		}

		return fmt.Errorf("failed to check VMware CLI: %w", err)
	}

	return nil
}

func buildAndTestFelt() error {
	fmt.Println("Building and testing Scriptorium Felt...")

	buildDir := feltProjectRootDir + `\build`

	fmt.Println("Cleaning existing Felt build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to clean Felt build directory: %w", err)
	}

	fmt.Println("Configuring Felt...")
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = feltProjectRootDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
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
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = brushProjectRootDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
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
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = inkstoneProjectRootDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
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

func buildScriptoriumAndRunAllTests() (*ProjectArtifacts, error) {
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
		BrushDLL: brushDll,
		InkstoneEXE: inkstoneExe,
		InkEXE: inkExe,
	}, nil
}

func loadConfig() (*Config, error) {
	fmt.Println("Loading configuration files and environment variables...")

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	config := &Config{
		VMEncryptionPassword: os.Getenv("ORIUM_VM_ENCRYPTION_PASSWORD"),
		GuestUsername:        os.Getenv("ORIUM_GUEST_USERNAME"),
		GuestPassword:        os.Getenv("ORIUM_GUEST_PASSWORD"),
	}

	if config.VMEncryptionPassword == "" {
		return nil, fmt.Errorf("ORIUM_VM_ENCRYPTION_PASSWORD is not configured")
	}

	if config.GuestUsername == "" {
		return nil, fmt.Errorf("ORIUM_GUEST_USERNAME is not configured")
	}

	if config.GuestPassword == "" {
		return nil, fmt.Errorf("ORIUM_GUEST_PASSWORD is not configured")
	}

	return config, nil
}

func resetVMEnv(config *Config) error {
	fmt.Println("Resetting the development VM to baseline...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"revertToSnapshot",
		devVMPath,
		devVMSnapshot,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset development VM to baseline: %w", err)
	}

	return nil
}

func startVM(config *Config) error {
	fmt.Println("Starting the development VM...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"start",
		devVMPath,
		"gui",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start development VM: %w", err)
	}

	return nil
}

func isVMRunning(config *Config) (bool, error) {
	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"list",
	)

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list running VMs: %w", err)
	}

	return strings.Contains(string(output), devVMPath), nil
}

func setupScriptoriumEnv(config *Config) error {
	createLogCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productLogDir,
	)

	createLogCmd.Stdout = os.Stdout
	createLogCmd.Stderr = os.Stderr

	if err := createLogCmd.Run(); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	createLocalCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productLocalDir,
	)

	createLocalCmd.Stdout = os.Stdout
	createLocalCmd.Stderr = os.Stderr

	if err := createLocalCmd.Run(); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	copyDictCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		dictionarySourceFile,
		productDictionaryDir,
	)

	copyDictCmd.Stdout = os.Stdout
	copyDictCmd.Stderr = os.Stderr

	if err := copyDictCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func deployArtifacts(artifacts *ProjectArtifacts, config *Config) error {
	fmt.Println("Deploying Scriptorium artifacts to development VM...")

	createArtifactsDirCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productArtifactsDir,
	)

	createArtifactsDirCmd.Stdout = os.Stdout
	createArtifactsDirCmd.Stderr = os.Stderr

	if err := createArtifactsDirCmd.Run(); err != nil {
		return fmt.Errorf("failed to create artifact directory in VM: %w", err)
	}

	fmt.Println("Deploying Brush DLL...")
	copyBrushCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.BrushDLL,
		productBrushDLL,
	)

	copyBrushCmd.Stdout = os.Stdout
	copyBrushCmd.Stderr = os.Stderr

	if err := copyBrushCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL: %w", err)
	}

	fmt.Println("Deploying Inkstone executable...")
	copyInkstoneCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.InkstoneEXE,
		productInkstoneEXE,
	)

	copyInkstoneCmd.Stdout = os.Stdout
	copyInkstoneCmd.Stderr = os.Stderr

	if err := copyInkstoneCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable: %w", err)
	}

	fmt.Println("Deploying Ink executable...")
	copyInkCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.InkEXE,
		productInkEXE,
	)

	copyInkCmd.Stdout = os.Stdout
	copyInkCmd.Stderr = os.Stderr

	if err := copyInkCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Ink executable: %w", err)
	}

	return nil
}

func registerBrush(config *Config) error {
	fmt.Println("Registering Scriptorium Brush...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"runProgramInGuest",
		devVMPath,
		`C:\Windows\System32\regsvr32.exe`,
		"/s",
		productBrushDLL,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func runUseCase(config *Config) error {
	fmt.Println("Running Scriptorium development use case...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"runProgramInGuest",
		devVMPath,
		`C:\Windows\System32\schtasks.exe`,
		"/Run",
		"/TN",
		devUseCaseTaskName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Scriptorium development use case: %w", err)
	}

	return nil
}

func startProduct(config *Config) error {
	if err := registerBrush(config); err != nil {
		return err
	}

	if err := runUseCase(config); err != nil {
		return err
	}

	return nil
}

func stopVM(config *Config) error {
	fmt.Println("Stopping the development VM...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"stop",
		devVMPath,
		"hard",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop development VM: %w", err)
	}

	return nil
}

func waitForVM(config *Config) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	fmt.Println("Waiting for the development VM to stop...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received. Stopping development VM...")

			if err := stopVM(config); err != nil {
				return err
			}

			return nil
		default:
			running, err := isVMRunning(config)
			if err != nil {
				return err
			}

			if !running {
				fmt.Println("Development VM stopped.")
				return nil
			}

			time.Sleep(1 * time.Second)
		}
	}
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Prepare and start a complete local development environment.",
	Long:  `The dev command sets up and starts a complete local development environment for Scriptorium.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch devEnv(env) {
		case devEnvVM:
			if err := ensureVMAvailable(); err != nil {
				return err
			}

			config, err := loadConfig()
			if err != nil {
				return err
			}

			artifacts, err := buildScriptoriumAndRunAllTests()
			if err != nil {
				return err
			}

			if err := resetVMEnv(config); err != nil {
				return err
			}

			if err := startVM(config); err != nil {
				return err
			}

			if err := setupScriptoriumEnv(config); err != nil {
				return err
			}

			if err := deployArtifacts(artifacts, config); err != nil {
				return err
			}

			if err := startProduct(config); err != nil {
				return err
			}

			if err := waitForVM(config); err != nil {
				return err
			}

			if err := resetVMEnv(config); err != nil {
				return err
			}

			return nil
		default:
			return fmt.Errorf("unsupported development environment: %s", env)
		}
	},
}

func NewCommand() *cobra.Command {
	devCmd.Flags().StringVarP(&env, "env", "E", "vm", "development environment to use; currently only 'vm' is supported, with more environments planned")

	return devCmd
}
