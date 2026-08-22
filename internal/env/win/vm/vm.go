// Package vm provides functions to manage the development virtual machine (VM) on Windows using VMware Workstation.
package vm

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

// VmrunPath TODO: Discover vmrun dynamically instead of relying on the default
// VMware Workstation installation path.
const VmrunPath = `C:\Program Files\VMware\VMware Workstation\vmrun.exe`

// DevVMPath TODO: Make the development VM path configurable.
const DevVMPath = `D:\Projects\Scriptorium\dev-env\scriptorium-dev\scriptorium-dev.vmx`

// TODO: Make the development VM snapshot configurable.
const devVMSnapshot = "baseline"

type VM struct {
	config *config.Config
}

func New(config *config.Config) *VM {
	return &VM{
		config: config,
	}
}

func (vm *VM) CreateDir(path string) error {
	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"-gu", vm.config.GuestUsername,
		"-gp", vm.config.GuestPassword,
		"createDirectoryInGuest",
		DevVMPath,
		path,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create directory in VM: %w", err)
	}

	return nil
}

func (vm *VM) CopyFile(src string, target string) error {
	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"-gu", vm.config.GuestUsername,
		"-gp", vm.config.GuestPassword,
		"CopyFileFromHostToGuest",
		DevVMPath,
		src,
		target,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy file to VM: %w", err)
	}

	return nil
}

func (vm *VM) EnsureAvailable() error {
	fmt.Println("Ensuring VMware is available...")
	if _, err := os.Stat(VmrunPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VMware CLI 'vmrun' was not found at %s", VmrunPath)
		}

		return fmt.Errorf("failed to check VMware CLI: %w", err)
	}

	return nil
}

func (vm *VM) Reset() error {
	fmt.Println("Resetting the development VM to baseline...")

	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"revertToSnapshot",
		DevVMPath,
		devVMSnapshot,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset development VM to baseline: %w", err)
	}

	return nil
}

func (vm *VM) Start() error {
	fmt.Println("Starting the development VM...")

	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"start",
		DevVMPath,
		"gui",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start development VM: %w", err)
	}

	return nil
}

func (vm *VM) Monitor() error {
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

			if err := vm.stopVM(); err != nil {
				return err
			}

			return nil
		default:
			running, err := vm.isRunning()
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

func (vm *VM) stopVM() error {
	fmt.Println("Stopping the development VM...")

	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"stop",
		DevVMPath,
		"hard",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop development VM: %w", err)
	}

	return nil
}

func (vm *VM) isRunning() (bool, error) {
	cmd := exec.Command(
		VmrunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"list",
	)

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list running VMs: %w", err)
	}

	return strings.Contains(string(output), DevVMPath), nil
}
