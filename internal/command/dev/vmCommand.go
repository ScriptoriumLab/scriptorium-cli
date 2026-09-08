package dev

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/vm"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/product"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

const devUseCaseTaskName = "Scriptorium Dev Use Case"

type vmCommand struct {
	machine *vm.VM
	workspace *project.Workspace
	product *product.Product
}

func (vmCmd *vmCommand) setupScriptoriumEnv() error {
	if err := vmCmd.machine.CreateDir(vmCmd.product.LogPath); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	if err := vmCmd.machine.CreateDir(vmCmd.product.LocalPath); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	if err := vmCmd.machine.CopyFile(vmCmd.workspace.Dictionary().SourceFile(), vmCmd.product.DictionaryPath); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) deployArtifacts(artifacts *project.ProjectArtifacts) error {
	fmt.Println("Deploying Scriptorium artifacts to development VM...")
	if err := vmCmd.machine.CreateDir(vmCmd.product.Config.ArtifactsPath); err != nil {
		return fmt.Errorf("failed to create artifact directory in VM: %w", err)
	}

	fmt.Println("Deploying Brush DLL...")
	if err := vmCmd.machine.CopyFile(artifacts.BrushDLL, vmCmd.product.Artifacts.BrushDLL); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL: %w", err)
	}

	fmt.Println("Deploying Inkstone executable...")
	if err := vmCmd.machine.CopyFile(artifacts.InkstoneEXE, vmCmd.product.Artifacts.InkstoneEXE); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable: %w", err)
	}

	fmt.Println("Deploying Ink executable...")
	if err := vmCmd.machine.CopyFile(artifacts.InkEXE, vmCmd.product.Artifacts.InkEXE); err != nil {
		return fmt.Errorf("failed to deploy Ink executable: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) registerBrush() error {
	fmt.Println("Registering Scriptorium Brush...")
	if err := vmCmd.machine.RunProgramDetached(`C:\Windows\System32\regsvr32.exe`, "/s", vmCmd.product.Artifacts.BrushDLL); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) runUseCase() error {
	fmt.Println("Running Scriptorium development use case...")
	if err := vmCmd.machine.RunProgramDetached(`C:\Windows\System32\schtasks.exe`, "/Run", "/TN", devUseCaseTaskName); err != nil {
		return fmt.Errorf("failed to run Scriptorium development use case: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) startProduct() error {
	if err := vmCmd.registerBrush(); err != nil {
		return err
	}

	if err := vmCmd.runUseCase(); err != nil {
		return err
	}

	return nil
}

func newVMCommand(devCommand *devCommand) *vmCommand {
	return &vmCommand{
		workspace: devCommand.workspace,
		product: devCommand.product,
	}
}

func (vmCmd *vmCommand) execute() error {
	vmConfig, err := config.LoadVM()
	if err != nil {
		return err
	}
	vmCmd.machine = vm.New(vmConfig)

	if err := vmCmd.machine.EnsureAvailable(); err != nil {
		return err
	}

	artifacts, err := vmCmd.workspace.BuildScriptoriumAndRunAllTests()
	if err != nil {
		return err
	}

	if err := vmCmd.machine.Prepare(); err != nil {
		return err
	}

	if err := vmCmd.setupScriptoriumEnv(); err != nil {
		return err
	}

	if err := vmCmd.deployArtifacts(artifacts); err != nil {
		return err
	}

	if err := vmCmd.startProduct(); err != nil {
		return err
	}

	if err := vmCmd.machine.Monitor(); err != nil {
		return err
	}

	if err := vmCmd.machine.Cleanup(); err != nil {
		return err
	}

	return nil
}
