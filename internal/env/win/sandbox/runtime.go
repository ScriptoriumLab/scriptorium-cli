package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func (sandbox *Sandbox) RunCommand(command string) error {
	if err := sandbox.waitForInteractiveSession(); err != nil {
		return err
	}

	fmt.Printf("Running in Windows Sandbox: %s\n", command)

	cmd := exec.Command(
		"wsb", "exec",
		"--id", sandbox.id,
		"--command", command,
		"--run-as", "ExistingLogin",
		"--raw",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"failed to execute command in Windows Sandbox: %q: %w",
			command,
			err,
		)
	}

	return nil
}

func (sandbox *Sandbox) CreateDir(path string) error {
	command := fmt.Sprintf(
		`powershell.exe -NoProfile -Command "New-Item -ItemType Directory -Force '%s' | Out-Null"`,
		path,
	)

	return sandbox.RunCommand(command)
}

func (sandbox *Sandbox) CopyFile(source, destination string) error {
	command := fmt.Sprintf(
		`powershell.exe -NoProfile -Command "Copy-Item -Force '%s' '%s'"`,
		source,
		destination,
	)

	if err := sandbox.RunCommand(command); err != nil {
		return fmt.Errorf(
			"failed to copy file in Windows Sandbox from %q to %q: %w",
			source,
			destination,
			err,
		)
	}

	return nil
}

func (sandbox *Sandbox) CreateFile(path string) error {
	command := fmt.Sprintf(
		`powershell.exe -NoProfile -Command "New-Item -ItemType File -Force '%s' | Out-Null"`,
		path,
	)

	return sandbox.RunCommand(command)
}

func (sandbox *Sandbox) ShareFolder(hostPath, sandboxPath string) error {
	cmd := exec.Command(
		"wsb", "share",
		"--id", sandbox.id,
		"--host-path", hostPath,
		"--sandbox-path", sandboxPath,
		"--raw",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to share folder with Windows Sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) waitForInteractiveSession() error {
	if sandbox.interactiveReady {
        return nil
    }

	fmt.Println("Waiting for Windows Sandbox interactive session...")

	for range 20 {
		cmd := exec.Command(
			"wsb", "exec",
			"--id", sandbox.id,
			"--command", "cmd.exe",
			"--run-as", "ExistingLogin",
		)

		if err := cmd.Run(); err == nil {
			sandbox.interactiveReady = true
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("windows sandbox interactive session did not become available")
}
