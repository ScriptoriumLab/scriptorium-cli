// Package dev provides the implementation of the `dev` subcommand.
package dev

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/spf13/cobra"
)

type devEnv string

const (
	devEnvVM devEnv = "vm"
)

var env string

// TODO: Make the Scriptorium product root directory configurable.
const productRootDir = `C:\Users\dev\Scriptorium`

const (
	productLocalDir = productRootDir + `\Local`
	productDictionaryDir = productLocalDir + `\pinyin_dictionary.txt`
	productLogDir = productRootDir + `\Log`
)

const productArtifactsDir = `C:\Users\dev\Desktop\ScriptoriumArtifacts`

const (
	productBrushDLL = productArtifactsDir + `\scriptorium-brush.dll`
	productInkstoneEXE = productArtifactsDir + `\scriptorium-inkstone.exe`
	productInkEXE = productArtifactsDir + `\scriptorium-ink.exe`
)

func setupScriptoriumEnv(config *config.Config) error {
	createLogCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productLogDir,
	)

	createLogCmd.Stdout = os.Stdout
	createLogCmd.Stderr = os.Stderr

	if err := createLogCmd.Run(); err != nil {
		return fmt.Errorf("failed to create log directory in VM: %w", err)
	}

	createLocalCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productLocalDir,
	)

	createLocalCmd.Stdout = os.Stdout
	createLocalCmd.Stderr = os.Stderr

	if err := createLocalCmd.Run(); err != nil {
		return fmt.Errorf("failed to create local directory in VM: %w", err)
	}

	copyDictCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		dictionarySourceFile,
		productDictionaryDir,
	)

	copyDictCmd.Stdout = os.Stdout
	copyDictCmd.Stderr = os.Stderr

	if err := copyDictCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy dictionary file to VM: %w", err)
	}

	return nil
}

func deployArtifacts(artifacts *ProjectArtifacts, config *config.Config) error {
	fmt.Println("Deploying Scriptorium artifacts to development VM...")

	createArtifactsDirCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"createDirectoryInGuest",
		devVMPath,
		productArtifactsDir,
	)

	createArtifactsDirCmd.Stdout = os.Stdout
	createArtifactsDirCmd.Stderr = os.Stderr

	if err := createArtifactsDirCmd.Run(); err != nil {
		return fmt.Errorf("failed to create artifact directory in VM: %w", err)
	}

	fmt.Println("Deploying Brush DLL...")
	copyBrushCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.BrushDLL,
		productBrushDLL,
	)

	copyBrushCmd.Stdout = os.Stdout
	copyBrushCmd.Stderr = os.Stderr

	if err := copyBrushCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Brush DLL: %w", err)
	}

	fmt.Println("Deploying Inkstone executable...")
	copyInkstoneCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.InkstoneEXE,
		productInkstoneEXE,
	)

	copyInkstoneCmd.Stdout = os.Stdout
	copyInkstoneCmd.Stderr = os.Stderr

	if err := copyInkstoneCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Inkstone executable: %w", err)
	}

	fmt.Println("Deploying Ink executable...")
	copyInkCmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"CopyFileFromHostToGuest",
		devVMPath,
		artifacts.InkEXE,
		productInkEXE,
	)

	copyInkCmd.Stdout = os.Stdout
	copyInkCmd.Stderr = os.Stderr

	if err := copyInkCmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy Ink executable: %w", err)
	}

	return nil
}

func registerBrush(config *config.Config) error {
	fmt.Println("Registering Scriptorium Brush...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"runProgramInGuest",
		devVMPath,
		`C:\Windows\System32\regsvr32.exe`,
		"/s",
		productBrushDLL,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to register Brush DLL: %w", err)
	}

	return nil
}

func startProduct(config *config.Config) error {
	if err := registerBrush(config); err != nil {
		return err
	}

	if err := RunUseCase(config); err != nil {
		return err
	}

	return nil
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Prepare and start a complete local development environment.",
	Long:  `The dev command sets up and starts a complete local development environment for Scriptorium.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch devEnv(env) {
		case devEnvVM:
			if err := EnsureVMAvailable(); err != nil {
				return err
			}

			config, err := config.LoadConfig()
			if err != nil {
				return err
			}

			artifacts, err := BuildScriptoriumAndRunAllTests()
			if err != nil {
				return err
			}

			if err := ResetVMEnv(config); err != nil {
				return err
			}

			if err := StartVM(config); err != nil {
				return err
			}

			if err := setupScriptoriumEnv(config); err != nil {
				return err
			}

			if err := deployArtifacts(artifacts, config); err != nil {
				return err
			}

			if err := startProduct(config); err != nil {
				return err
			}

			if err := WaitForVM(config); err != nil {
				return err
			}

			if err := ResetVMEnv(config); err != nil {
				return err
			}

			return nil
		default:
			return fmt.Errorf("unsupported development environment: %s", env)
		}
	},
}

func NewCommand() *cobra.Command {
	devCmd.Flags().StringVarP(&env, "env", "E", "vm", "development environment to use; currently only 'vm' is supported, with more environments planned")

	return devCmd
}
