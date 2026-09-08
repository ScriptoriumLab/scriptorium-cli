// Package dev provides the implementation of the `dev` subcommand.
package dev

import (
	"fmt"

	"github.com/spf13/cobra"
)

type devEnv string

const (
	devEnvVM devEnv = "vm"
	devEnvSandbox devEnv = "sandbox"
)

var env string

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Prepare and start a complete local development environment.",
	Long:  `The dev command sets up and starts a complete local development environment for Scriptorium.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch devEnv(env) {
		case devEnvVM:
			if err := newVMCommand().execute(); err != nil {
				return err
			}

			return nil
		case devEnvSandbox:
			if err := newSandboxCommand().execute(); err != nil {
				return err
			}

			return nil
		default:
			return fmt.Errorf("unsupported development environment: %s", env)
		}
	},
}

func NewCommand() *cobra.Command {
	devCmd.Flags().StringVarP(&env, "env", "E", "sandbox", "development environment to use")

	return devCmd
}
