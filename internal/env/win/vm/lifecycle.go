package vm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

// TODO: Make the development VM snapshot configurable.
const devVMSnapshot = "baseline"

func (vm *VM) Reset() error {
	fmt.Println("Resetting the development VM to baseline...")

	cmd := exec.Command(
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"revertToSnapshot",
		vm.config.VMImagePath,
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
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"start",
		vm.config.VMImagePath,
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
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"stop",
		vm.config.VMImagePath,
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
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"list",
	)

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list running VMs: %w", err)
	}

	return strings.Contains(string(output), vm.config.VMImagePath), nil
}
