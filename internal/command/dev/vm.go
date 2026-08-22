package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
)

// TODO: Discover vmrun dynamically instead of relying on the default
// VMware Workstation installation path.
const vmrunPath = `C:\Program Files\VMware\VMware Workstation\vmrun.exe`

// TODO: Make the development VM path configurable.
const devVMPath = `D:\Projects\Scriptorium\dev-env\scriptorium-dev\scriptorium-dev.vmx`

// TODO: Make the development VM snapshot configurable.
const devVMSnapshot = "baseline"

const devUseCaseTaskName = "Scriptorium Dev Use Case"

func EnsureVMAvailable() error {
	fmt.Println("Ensuring VMware is available...")
	if _, err := os.Stat(vmrunPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VMware CLI 'vmrun' was not found at %s", vmrunPath)
		}

		return fmt.Errorf("failed to check VMware CLI: %w", err)
	}

	return nil
}

func ResetVMEnv(config *config.Config) error {
	fmt.Println("Resetting the development VM to baseline...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"revertToSnapshot",
		devVMPath,
		devVMSnapshot,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset development VM to baseline: %w", err)
	}

	return nil
}

func StartVM(config *config.Config) error {
	fmt.Println("Starting the development VM...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"start",
		devVMPath,
		"gui",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start development VM: %w", err)
	}

	return nil
}

func WaitForVM(config *config.Config) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	fmt.Println("Waiting for the development VM to stop...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received. Stopping development VM...")

			if err := StopVM(config); err != nil {
				return err
			}

			return nil
		default:
			running, err := isVMRunning(config)
			if err != nil {
				return err
			}

			if !running {
				fmt.Println("Development VM stopped.")
				return nil
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func StopVM(config *config.Config) error {
	fmt.Println("Stopping the development VM...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"stop",
		devVMPath,
		"hard",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop development VM: %w", err)
	}

	return nil
}

func RunUseCase(config *config.Config) error {
	fmt.Println("Running Scriptorium development use case...")

	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"-gu", config.GuestUsername,
		"-gp", config.GuestPassword,
		"runProgramInGuest",
		devVMPath,
		`C:\Windows\System32\schtasks.exe`,
		"/Run",
		"/TN",
		devUseCaseTaskName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Scriptorium development use case: %w", err)
	}

	return nil
}

func isVMRunning(config *config.Config) (bool, error) {
	cmd := exec.Command(
		vmrunPath,
		"-T", "ws",
		"-vp", config.VMEncryptionPassword,
		"list",
	)

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list running VMs: %w", err)
	}

	return strings.Contains(string(output), devVMPath), nil
}

