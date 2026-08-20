package main

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

type Config struct {
	VMEncryptionPassword string
	GuestUsername        string
	GuestPassword        string
}

var cli = &cobra.Command{
	Use:   "orium",
	Short: "Developer toolchain for Scriptorium",
	Long: `Orium is the developer toolchain for Scriptorium.

It provides a unified command-line interface for building, running,
testing, diagnosing, and maintaining local development workflows.`,
}

var greetingCmd = &cobra.Command{
	Use:   "greet",
	Short: "Print a greeting message",
	Long:  "Print a simple greeting message from Orium.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, Scriptorium!")
	},
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
		"-DCMAKE_BUILD_TYPE=Debug",
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

func buildAndTestBrush() error {
	fmt.Println("Building and testing Scriptorium Brush...")

	buildDir := brushProjectRootDir + `\build`

	fmt.Println("Cleaning existing Brush build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to clean Brush build directory: %w", err)
	}

	fmt.Println("Configuring Brush...")
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Debug",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = brushProjectRootDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Brush: %w", err)
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
		return fmt.Errorf("failed to build Brush: %w", err)
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
		return fmt.Errorf("scriptorium brush tests failed: %w", err)
	}

	return nil
}

func buildAndTestInkstone() error {
	fmt.Println("Building and testing Scriptorium Inkstone...")

	buildDir := inkstoneProjectRootDir + `\build`

	fmt.Println("Cleaning existing Inkstone build directory...")
	if err := os.RemoveAll(buildDir); err != nil {
		return fmt.Errorf("failed to clean Inkstone build directory: %w", err)
	}

	fmt.Println("Configuring Inkstone...")
	configureCmd := exec.Command(
		"cmake",
		"-S", ".",
		"-B", "build",
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Debug",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	)

	configureCmd.Dir = inkstoneProjectRootDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr

	if err := configureCmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Inkstone: %w", err)
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
		return fmt.Errorf("failed to build Inkstone: %w", err)
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
		return fmt.Errorf("scriptorium inkstone unit tests failed: %w", err)
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
		return fmt.Errorf("scriptorium inkstone integration tests failed: %w", err)
	}

	return nil
}

func buildInk() error {
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
		return fmt.Errorf("scriptorium ink build failed: %w", err)
	}

	return nil
}

func buildScriptoriumAndRunAllTests() error {
	fmt.Println("Building Scriptorium and running all tests...")
	if err := buildAndTestFelt(); err != nil {
		return err
	}

	if err := buildAndTestBrush(); err != nil {
		return err
	}

	if err := buildAndTestInkstone(); err != nil {
		return err
	}

	if err := buildInk(); err != nil {
		return err
	}

	return nil
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

			if err := buildScriptoriumAndRunAllTests(); err != nil {
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

func init() {
	cli.AddCommand(greetingCmd)
	devCmd.Flags().StringVarP(&env, "env", "E", "vm", "development environment to use; currently only 'vm' is supported, with more environments planned")
	cli.AddCommand(devCmd)
}

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
