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
	workspace  *project.Workspace
	product    *product.Product
	envCommand envCommand
}

type devEnv string

const (
	devEnvVM      devEnv = "vm"
	devEnvSandbox devEnv = "sandbox"
)

var env string

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Prepare and start a complete local development environment.",
	Long:  `The dev command sets up and starts a complete local development environment for Scriptorium.`,
	RunE: func(_ *cobra.Command, args []string) error {
		cmd, err := newDevCommand()
		if err != nil {
			return err
		}

		switch devEnv(env) {
		case devEnvVM:
			vmCmd, err := newVMCommand(cmd.product)
			if err != nil {
				return err
			}
			cmd.envCommand = vmCmd

		case devEnvSandbox:
			sandboxCmd, err := newSandboxCommand(cmd.product)
			if err != nil {
				return err
			}
			cmd.envCommand = sandboxCmd

		default:
			return fmt.Errorf("unsupported development environment: %s", env)
		}

		if err := cmd.execute(); err != nil {
			return err
		}

		return nil
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

func (cmd *devCommand) execute() error {
	if err := cmd.envCommand.ensureEnv(); err != nil {
		return err
	}

	artifacts, err := cmd.workspace.BuildScriptoriumAndRunAllTests()
	if err != nil {
		return err
	}

	if err := cmd.envCommand.prepareEnv(); err != nil {
		return err
	}

	if err := cmd.envCommand.setupProductPrerequisites(); err != nil {
		return err
	}

	if err := cmd.envCommand.deployArtifacts(artifacts, cmd.workspace.Dictionary().SourceFile()); err != nil {
		return err
	}

	if err := cmd.envCommand.startProduct(); err != nil {
		return err
	}

	if err := cmd.envCommand.startManualTests(); err != nil {
		return err
	}

	if err := cmd.envCommand.monitorEnv(); err != nil {
		return err
	}

	if err := cmd.envCommand.cleanupEnv(); err != nil {
		return err
	}

	return nil
}

func NewCommand() *cobra.Command {
	devCmd.Flags().StringVarP(&env, "env", "E", "sandbox", "development environment to use")

	return devCmd
}
