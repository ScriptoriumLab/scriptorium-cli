package vm

import (
	"fmt"
	"os"
	"os/exec"
)

func (vm *VM) CreateDir(path string) error {
	cmd := exec.Command(
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"-gu", vm.config.GuestUsername,
		"-gp", vm.config.GuestPassword,
		"createDirectoryInGuest",
		vm.config.VMImagePath,
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
		vm.config.VMRunPath,
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"-gu", vm.config.GuestUsername,
		"-gp", vm.config.GuestPassword,
		"CopyFileFromHostToGuest",
		vm.config.VMImagePath,
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

func (vm *VM) RunProgramDetached(program string, args ...string) error {
	vmArgs := []string{
		"-T", "ws",
		"-vp", vm.config.VMEncryptionPassword,
		"-gu", vm.config.GuestUsername,
		"-gp", vm.config.GuestPassword,
		"runProgramInGuest",
		vm.config.VMImagePath,
		program,
	}

	vmArgs = append(vmArgs, args...)

	cmd := exec.Command(vm.config.VMRunPath, vmArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run program in guest: %w", err)
	}

	return nil
}
