package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cli = &cobra.Command{
	Use:   "scriptorium-cli",
	Short: "A CLI too for enhancing the Scriptorium IME Developer Experience",
	Long:  "A CLI tool for enhancing the Scriptorium IME Developer Experience",
}

var greetingCmd = &cobra.Command{
	Use:   "greet",
	Short: "Prints a greeting message",
	Long:  "Prints a greeting message to the console",
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
