package dev

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/sandbox"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/product"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

type sandboxCommand struct {
	sandbox   *sandbox.Sandbox
	workspace *project.Workspace
	product   *product.Product
}

func (sandboxCmd *sandboxCommand) setupScriptoriumEnv() error {
	if err := sandboxCmd.sandbox.CreateDir(sandboxCmd.product.LogPath); err != nil {
		return fmt.Errorf("failed to create log directory in Windows Sandbox: %w", err)
	}

	if err := sandboxCmd.sandbox.CreateDir(sandboxCmd.product.LocalPath); err != nil {
		return fmt.Errorf("failed to create local directory in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) deployArtifacts(artifacts *project.ProjectArtifacts) error {
	fmt.Println("Deploying Scriptorium artifacts to development windows sandbox...")
	return nil
}

func (sandboxCmd *sandboxCommand) startProduct() error {
	return nil
}

func newSandboxCommand() *sandboxCommand {
	return &sandboxCommand{}
}

func (sandboxCmd *sandboxCommand) execute() error {
	sandboxCmd.sandbox = sandbox.New()
	workspaceConfig, err := config.LoadWorkspace()
	if err != nil {
		return err
	}
	sandboxCmd.workspace = project.NewWorkspace(workspaceConfig)

	productConfig, err := config.LoadProduct()
	if err != nil {
		return err
	}
	sandboxCmd.product = product.NewProduct(productConfig)

	if err := sandboxCmd.sandbox.EnsureAvailable(); err != nil {
		return err
	}

	artifacts, err := sandboxCmd.workspace.BuildScriptoriumAndRunAllTests()
	if err != nil {
		return err
	}

	if err := sandboxCmd.sandbox.Start(); err != nil {
		return err
	}

	if err := sandboxCmd.sandbox.Connect(); err != nil {
		return err
	}

	if err := sandboxCmd.setupScriptoriumEnv(); err != nil {
		return err
	}

	if err := sandboxCmd.deployArtifacts(artifacts); err != nil {
		return err
	}

	if err := sandboxCmd.startProduct(); err != nil {
		return err
	}

	if err := sandboxCmd.sandbox.Monitor(); err != nil {
		return err
	}

	if err := sandboxCmd.sandbox.Cleanup(); err != nil {
		return err
	}

	return nil
}
