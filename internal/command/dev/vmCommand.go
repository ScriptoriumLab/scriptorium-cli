package dev

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/vm"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

// TODO: Make the Scriptorium product root directory configurable.
const productRootDir = `C:\Users\dev\Scriptorium`

const (
	productLocalDir      = productRootDir + `\Local`
	productDictionaryDir = productLocalDir + `\pinyin_dictionary.txt`
	productLogDir        = productRootDir + `\Log`
)

const productArtifactsDir = `C:\Users\dev\Desktop\ScriptoriumArtifacts`

const (
	productBrushDLL    = productArtifactsDir + `\scriptorium-brush.dll`
	productInkstoneEXE = productArtifactsDir + `\scriptorium-inkstone.exe`
	productInkEXE      = productArtifactsDir + `\scriptorium-ink.exe`
)

const devUseCaseTaskName = "Scriptorium Dev Use Case"

type vmCommand struct {
	machine *vm.VM
}

func (vmCmd *vmCommand) setupScriptoriumEnv() error {
	if err := vmCmd.machine.CreateDir(productLogDir); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	if err := vmCmd.machine.CreateDir(productLocalDir); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	if err := vmCmd.machine.CopyFile(project.DictionarySourceFile, productDictionaryDir); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) deployArtifacts(artifacts *project.ProjectArtifacts) error {
	fmt.Println("Deploying Scriptorium artifacts to development VM...")
	if err := vmCmd.machine.CreateDir(productArtifactsDir); err != nil {
		return fmt.Errorf("failed to create artifact directory in VM: %w", err)
	}

	fmt.Println("Deploying Brush DLL...")
	if err := vmCmd.machine.CopyFile(artifacts.BrushDLL, productBrushDLL); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL: %w", err)
	}

	fmt.Println("Deploying Inkstone executable...")
	if err := vmCmd.machine.CopyFile(artifacts.InkstoneEXE, productInkstoneEXE); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable: %w", err)
	}

	fmt.Println("Deploying Ink executable...")
	if err := vmCmd.machine.CopyFile(artifacts.InkEXE, productInkEXE); err != nil {
		return fmt.Errorf("failed to deploy Ink executable: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) registerBrush() error {
	fmt.Println("Registering Scriptorium Brush...")
	if err := vmCmd.machine.RunProgram(`C:\Windows\System32\regsvr32.exe`, "/s", productBrushDLL); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func (vmCmd *vmCommand) runUseCase() error {
	fmt.Println("Running Scriptorium development use case...")
	if err := vmCmd.machine.RunProgram(`C:\Windows\System32\schtasks.exe`, "/Run", "/TN", devUseCaseTaskName); err != nil {
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

func newVMCommand() *vmCommand {
	return &vmCommand{}
}

func (vmCmd *vmCommand) execute() error {
	config, err := config.Load()
	if err != nil {
		return err
	}

	vmCmd.machine = vm.New(config)

	if err := vmCmd.machine.EnsureAvailable(); err != nil {
		return err
	}

	artifacts, err := project.BuildScriptoriumAndRunAllTests()
	if err != nil {
		return err
	}

	if err := vmCmd.machine.Reset(); err != nil {
		return err
	}

	if err := vmCmd.machine.Start(); err != nil {
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

	if err := vmCmd.machine.Reset(); err != nil {
		return err
	}

	return nil
}
