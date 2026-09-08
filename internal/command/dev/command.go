// Package dev provides the implementation of the `dev` subcommand.
package dev

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/product"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
	"github.com/spf13/cobra"
)

type devCommand struct {
	workspace *project.Workspace
	product   *product.Product
}

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
		devCommand, err := newDevCommand()
		if err != nil {
			return err
		}

		switch devEnv(env) {
		case devEnvVM:
			if err := newVMCommand(devCommand).execute(); err != nil {
				return err
			}

			return nil
		case devEnvSandbox:
			if err := newSandboxCommand(devCommand).execute(); err != nil {
				return err
			}

			return nil
		default:
			return fmt.Errorf("unsupported development environment: %s", env)
		}
	},
}

func newDevCommand() (*devCommand, error) {
	devCmd := &devCommand{}
	workspaceConfig, err := config.LoadWorkspace()
	if err != nil {
		return nil, err
	}
	devCmd.workspace = project.NewWorkspace(workspaceConfig)

	productConfig, err := config.LoadProduct()
	if err != nil {
		return nil, err
	}
	devCmd.product = product.NewProduct(productConfig)

	return devCmd, nil
}

func NewCommand() *cobra.Command {
	devCmd.Flags().StringVarP(&env, "env", "E", "sandbox", "development environment to use")

	return devCmd
}
