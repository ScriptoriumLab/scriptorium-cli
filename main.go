package main

import (
	"fmt"
	"os"
	"os/exec"

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

const devVMSnapshot = "baseline"

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

func buildScriptorium() {
	fmt.Println("Building Scriptorium IME...")
}

func runAllTests() {
	fmt.Println("Running all tests under the Scriptorium IME.")
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

func setupVM(config *Config) error {
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

func cleanUp() {
	fmt.Println("Cleaning up temporary files and resources...")
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

			buildScriptorium()
			runAllTests()

			if err := setupVM(config); err != nil {
				return err
			}

			if err := startVM(config); err != nil {
				return err
			}

			cleanUp()
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
