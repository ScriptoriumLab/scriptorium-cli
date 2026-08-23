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

func setupScriptoriumEnv(machine *vm.VM) error {
	if err := machine.CreateDir(productLogDir); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	if err := machine.CreateDir(productLocalDir); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	if err := machine.CopyFile(project.DictionarySourceFile, productDictionaryDir); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func deployArtifacts(artifacts *project.ProjectArtifacts, machine *vm.VM) error {
	fmt.Println("Deploying Scriptorium artifacts to development VM...")
	if err := machine.CreateDir(productArtifactsDir); err != nil {
		return fmt.Errorf("failed to create artifact directory in VM: %w", err)
	}

	fmt.Println("Deploying Brush DLL...")
	if err := machine.CopyFile(artifacts.BrushDLL, productBrushDLL); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL: %w", err)
	}

	fmt.Println("Deploying Inkstone executable...")
	if err := machine.CopyFile(artifacts.InkstoneEXE, productInkstoneEXE); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable: %w", err)
	}

	fmt.Println("Deploying Ink executable...")
	if err := machine.CopyFile(artifacts.InkEXE, productInkEXE); err != nil {
		return fmt.Errorf("failed to deploy Ink executable: %w", err)
	}

	return nil
}

func registerBrush(machine *vm.VM) error {
	fmt.Println("Registering Scriptorium Brush...")
	if err := machine.RunProgram(`C:\Windows\System32\regsvr32.exe`, "/s", productBrushDLL); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func runUseCase(machine *vm.VM) error {
	fmt.Println("Running Scriptorium development use case...")
	if err := machine.RunProgram(`C:\Windows\System32\schtasks.exe`, "/Run", "/TN", devUseCaseTaskName); err != nil {
		return fmt.Errorf("failed to run Scriptorium development use case: %w", err)
	}

	return nil
}

func startProduct(machine *vm.VM) error {
	if err := registerBrush(machine); err != nil {
		return err
	}

	if err := runUseCase(machine); err != nil {
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

	machine := vm.New(config)

	if err := machine.EnsureAvailable(); err != nil {
		return err
	}

	artifacts, err := project.BuildScriptoriumAndRunAllTests()
	if err != nil {
		return err
	}

	if err := machine.Reset(); err != nil {
		return err
	}

	if err := machine.Start(); err != nil {
		return err
	}

	if err := setupScriptoriumEnv(machine); err != nil {
		return err
	}

	if err := deployArtifacts(artifacts, machine); err != nil {
		return err
	}

	if err := startProduct(machine); err != nil {
		return err
	}

	if err := machine.Monitor(); err != nil {
		return err
	}

	if err := machine.Reset(); err != nil {
		return err
	}

	return nil
}
