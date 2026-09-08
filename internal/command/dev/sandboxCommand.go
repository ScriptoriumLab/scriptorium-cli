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
	sandboxCmd.sandbox.TempStagingDir = stagingDir

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

func (sandboxCmd *sandboxCommand) registerBrush() error {
	command := fmt.Sprintf(
		`regsvr32.exe /s "%s"`,
		sandboxCmd.product.Artifacts.BrushDLL,
	)

	if err := sandboxCmd.sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to register Brush in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) startInkstone() error {
	if err := sandboxCmd.sandbox.RunProgramDetached(
		sandboxCmd.product.Artifacts.InkstoneEXE,
	); err != nil {
		return fmt.Errorf(
			"failed to start Inkstone in Windows Sandbox: %w",
			err,
		)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) startInk() error {
	command := fmt.Sprintf(
		`powershell.exe -NoProfile -Command "$env:WEBVIEW2_BROWSER_EXECUTABLE_FOLDER='%s'; Start-Process -FilePath '%s'"`,
		sandboxCmd.sandbox.Config.WebView2BrowserExecutableFolder,
		sandboxCmd.product.Artifacts.InkEXE,
	)

	if err := sandboxCmd.sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to start Ink in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) startProduct() error {
	if err := sandboxCmd.registerBrush(); err != nil {
		return err
	}

	if err := sandboxCmd.startInkstone(); err != nil {
		return err
	}

	if err := sandboxCmd.startInk(); err != nil {
		return err
	}

	return nil
}

func (sandboxCmd *sandboxCommand) createTestTextFile() error {
	if err := sandboxCmd.sandbox.CreateFile(sandboxCmd.sandbox.Config.TestFilePath); err != nil {
		return fmt.Errorf(
			"failed to create manual test text file in Windows Sandbox: %w",
			err,
		)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) runNotepadPlusPlus() error {
	command := fmt.Sprintf(
		`powershell.exe -NoProfile -Command "Start-Process -FilePath '%s' -ArgumentList '%s'"`,
		sandboxCmd.sandbox.Config.NotepadPlusPlusPath,
		sandboxCmd.sandbox.Config.TestFilePath,
	)

	if err := sandboxCmd.sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to start Notepad++ in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandboxCmd *sandboxCommand) startManualTest() error {
	if err := sandboxCmd.createTestTextFile(); err != nil {
		return err
	}

	if err := sandboxCmd.runNotepadPlusPlus(); err != nil {
		return err
	}

	return nil
}

func newSandboxCommand() *sandboxCommand {
	return &sandboxCommand{}
}

func (sandboxCmd *sandboxCommand) execute() error {
	sandboxConfig, err := config.LoadSandbox()
	if err != nil {
		return err
	}
	sandboxCmd.sandbox = sandbox.New(sandboxConfig)
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

	if err := sandboxCmd.sandbox.SetupEnv(); err != nil {
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

	if err := sandboxCmd.startManualTest(); err != nil {
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
