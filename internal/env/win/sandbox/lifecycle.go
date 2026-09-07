package sandbox

import (
	"fmt"
	"os"
	"os/exec"
)

func (sandbox *Sandbox) Start() error {
	fmt.Println("Starting the development Windows Sandbox...")
	cmd := exec.Command("wsb", "start",
		"--id", sandbox.id,
		"--raw",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start development windows sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) Connect() error {
	fmt.Println("Connecting to the development Windows Sandbox...")
	cmd := exec.Command("wsb", "connect",
		"--id", sandbox.id,
		"--raw",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to connect to the development windows sandbox: %w", err)
	}

	return nil
}
