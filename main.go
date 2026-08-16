package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

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

func ensureSandboxAvailable() error {
	fmt.Println("Checking Windows Sandbox availability...")

	cmd := exec.Command(
		"dism.exe",
		"/Online",
		"/English",
		"/Get-FeatureInfo",
		"/FeatureName:Containers-DisposableClientVM",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to query Windows Sandbox feature: %w\n%s",
			err,
			string(output),
		)
	}

	result := string(output)

	if strings.Contains(result, "State : Enabled") {
		fmt.Println("Windows Sandbox is available.")
		return nil
	}

	if strings.Contains(result, "State : Disabled") {
		return fmt.Errorf(`sandbox is not enabled.

Enable it from an elevated PowerShell:

  dism.exe /Online /Enable-Feature /FeatureName:Containers-DisposableClientVM /All

A system restart may be required.

After restarting, run:

  orium dev`)
	}

	return fmt.Errorf(
		"unable to determine Windows Sandbox state:\n%s",
		result,
	)
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
		if err := ensureSandboxAvailable(); err != nil {
			return err
		}

		buildScriptorium()
		runAllTests()
		loadConfig()
		startSandbox()
		cleanUp()

		return nil
	},
}

func main() {
	cli.AddCommand(greetingCmd)
	cli.AddCommand(devCmd)

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
