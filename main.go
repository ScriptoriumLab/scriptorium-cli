package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type devEnv string

const (
	devEnvVM devEnv = "vm"
)

var env string

const vmrunPath = `C:\Program Files\VMware\VMware Workstation\vmrun.exe`

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
	fmt.Println("Ensuring Sandbox is available...")
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

func loadConfig() {
	fmt.Println("Loading configuration files and environment variables...")
}

func startSandbox() {
	fmt.Println("Starting the sandbox environment...")
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

			buildScriptorium()
			runAllTests()
			loadConfig()
			startSandbox()
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
