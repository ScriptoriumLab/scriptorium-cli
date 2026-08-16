package main

import (
	"fmt"
	"os"

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

func main() {
	cli.AddCommand(greetingCmd)

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
