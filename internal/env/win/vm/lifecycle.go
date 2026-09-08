package vm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (vm *VM) reset() error {
	fmt.Println("Resetting the development VM to baseline...")

	cmd := exec.Command(
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"revertToSnapshot",
		vm.config.VMImagePath,
		vm.config.VMSnapshotName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset development VM to baseline: %w", err)
	}

	return nil
}

func (vm *VM) start() error {
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
