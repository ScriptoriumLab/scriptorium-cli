// Package vm provides the implementation of the EnvCommand interface for managing a development environment on a Windows virtual machine.
package vm

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	vmenv "github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/vm"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/product"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

const devUseCaseTaskName = "Scriptorium Dev Use Case"

type Command struct {
	machine *vmenv.VM
	product *product.Product
}

func NewCommand(product *product.Product) (*Command, error) {
	vmConfig, err := config.LoadVM()
	if err != nil {
		return nil, err
	}

	return &Command{
		machine: vmenv.New(vmConfig),
		product: product,
	}, nil
}

func (vmCmd *Command) EnsureEnv() error {
	return vmCmd.machine.EnsureAvailable()
}

func (vmCmd *Command) PrepareEnv() error {
	return vmCmd.machine.Prepare()
}

func (vmCmd *Command) SetupProductPrerequisites() error {
	if err := vmCmd.machine.CreateDir(vmCmd.product.LogPath); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	if err := vmCmd.machine.CreateDir(vmCmd.product.LocalPath); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	return nil
}

func (vmCmd *Command) DeployArtifacts(artifacts *project.Artifacts, dictionarySourcePath string) error {
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

func (vmCmd *Command) StartProduct() error {
	if err := vmCmd.registerBrush(); err != nil {
		return err
	}

	return nil
}

func (vmCmd *Command) StartManualTests() error {
	if err := vmCmd.runUseCase(); err != nil {
		return err
	}

	return nil
}

func (vmCmd *Command) MonitorEnv() error {
	return vmCmd.machine.Monitor()
}

func (vmCmd *Command) CleanupEnv() error {
	return vmCmd.machine.Cleanup()
}

func (vmCmd *Command) registerBrush() error {
	fmt.Println("Registering Scriptorium Brush...")
	if err := vmCmd.machine.RunProgramDetached(`C:\Windows\System32\regsvr32.exe`, "/s", vmCmd.product.Artifacts.BrushDLL); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func (vmCmd *Command) runUseCase() error {
	fmt.Println("Running Scriptorium development use case...")
	if err := vmCmd.machine.RunProgramDetached(`C:\Windows\System32\schtasks.exe`, "/Run", "/TN", devUseCaseTaskName); err != nil {
		return fmt.Errorf("failed to run Scriptorium development use case: %w", err)
	}

	return nil
}

