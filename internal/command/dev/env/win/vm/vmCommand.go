// Package vm provides the implementation of the EnvCommand interface for managing a development environment on a Windows virtual machine.
package vm

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/command/dev/env"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	vmenv "github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/vm"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/product"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

const devUseCaseTaskName = "Scriptorium Dev Use Case"

type vmCommand struct {
	machine *vmenv.VM
	product *product.Product
}

// Ensure *vmCommand implements env.EnvCommand.
var _ env.Command = (*vmCommand)(nil)

func NewCommand(product *product.Product) (*vmCommand, error) {
	vmConfig, err := config.LoadVM()
	if err != nil {
		return nil, err
	}

	return &vmCommand{
		machine: vmenv.New(vmConfig),
		product: product,
	}, nil
}

func (vmCmd *vmCommand) EnsureEnv() error {
	return vmCmd.machine.EnsureAvailable()
}

func (vmCmd *vmCommand) PrepareEnv() error {
	return vmCmd.machine.Prepare()
}

func (vmCmd *vmCommand) SetupProductPrerequisites() error {
	if err := vmCmd.machine.CreateDir(vmCmd.product.LogPath); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	if err := vmCmd.machine.CreateDir(vmCmd.product.LocalPath); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) DeployArtifacts(artifacts *project.ProjectArtifacts, dictionarySourcePath string) error {
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

	if err := vmCmd.machine.CopyFile(dictionarySourcePath, vmCmd.product.DictionaryPath); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) StartProduct() error {
	if err := vmCmd.registerBrush(); err != nil {
		return err
	}

	return nil
}

func (vmCmd *vmCommand) StartManualTests() error {
	if err := vmCmd.runUseCase(); err != nil {
		return err
	}

	return nil
}

func (vmCmd *vmCommand) MonitorEnv() error {
	return vmCmd.machine.Monitor()
}

func (vmCmd *vmCommand) CleanupEnv() error {
	return vmCmd.machine.Cleanup()
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

