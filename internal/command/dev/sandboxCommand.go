package dev

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/sandbox"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

type sandboxCommand struct {
	sandbox *sandbox.Sandbox
	workspace *project.Workspace
}

func (sandboxCmd *sandboxCommand) setupScriptoriumEnv() error {
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
