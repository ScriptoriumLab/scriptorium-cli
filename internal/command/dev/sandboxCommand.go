package dev

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	if err := sandboxCmd.sandbox.CreateDir(sandboxCmd.product.Config.ArtifactsPath); err != nil {
		return fmt.Errorf("failed to create artifacts directory in Windows Sandbox: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func (sandboxCmd *sandboxCommand) deployArtifacts(artifacts *project.ProjectArtifacts) error {
	fmt.Println("Deploying Scriptorium artifacts to development Windows Sandbox...")

	stagingDir, err := os.MkdirTemp("", "scriptorium-sandbox-*")
	if err != nil {
		return fmt.Errorf("failed to create Sandbox staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	if err := copyFile(
		artifacts.BrushDLL,
		filepath.Join(stagingDir, "scriptorium-brush.dll"),
	); err != nil {
		return err
	}

	if err := copyFile(
		artifacts.InkstoneEXE,
		filepath.Join(stagingDir, "scriptorium-inkstone.exe"),
	); err != nil {
		return err
	}

	if err := copyFile(
		artifacts.InkEXE,
		filepath.Join(stagingDir, "scriptorium-ink.exe"),
	); err != nil {
		return err
	}

	if err := copyFile(
		sandboxCmd.workspace.Dictionary().SourceFile(),
		filepath.Join(stagingDir, "pinyin_dictionary.txt"),
	); err != nil {
		return err
	}

	const sandboxStagingDir = `C:\ScriptoriumStaging`

	if err := sandboxCmd.sandbox.ShareFolder(
		stagingDir,
		sandboxStagingDir,
	); err != nil {
		return err
	}

	if err := sandboxCmd.sandbox.CopyFile(
		filepath.Join(sandboxStagingDir, "scriptorium-brush.dll"),
		sandboxCmd.product.Artifacts.BrushDLL,
	); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL to Windows Sandbox: %w", err)
	}

	if err := sandboxCmd.sandbox.CopyFile(
		filepath.Join(sandboxStagingDir, "scriptorium-inkstone.exe"),
		sandboxCmd.product.Artifacts.InkstoneEXE,
	); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable to Windows Sandbox: %w", err)
	}

	if err := sandboxCmd.sandbox.CopyFile(
		filepath.Join(sandboxStagingDir, "scriptorium-ink.exe"),
		sandboxCmd.product.Artifacts.InkEXE,
	); err != nil {
		return fmt.Errorf("failed to deploy Ink executable to Windows Sandbox: %w", err)
	}

	if err := sandboxCmd.sandbox.CopyFile(
		filepath.Join(sandboxStagingDir, "pinyin_dictionary.txt"),
		sandboxCmd.product.DictionaryPath,
	); err != nil {
		return fmt.Errorf("failed to deploy dictionary to Windows Sandbox: %w", err)
	}

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
