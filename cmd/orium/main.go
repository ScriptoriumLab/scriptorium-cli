package main

import (
	"fmt"
	"os"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/command/dev"

	"github.com/spf13/cobra"
)

var cli = &cobra.Command{
	Use:   "orium",
	Short: "Developer toolchain for Scriptorium",
	Long: `Orium is the developer toolchain for Scriptorium.

It provides a unified command-line interface for building, running,
testing, diagnosing, and maintaining local development workflows.`,
}

func init() {
	cli.AddCommand(dev.NewCommand())
}

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
