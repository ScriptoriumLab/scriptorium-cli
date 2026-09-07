// Package sandbox provides functions to manage the development sandbox on Windows using VMware Workstation.
package sandbox

import (
	"fmt"
	"io"
	"os/exec"
)

type Sandbox struct {}

func New() *Sandbox {
	return &Sandbox{}
}

func (sandbox *Sandbox) EnsureAvailable() error {
	fmt.Println("Ensuring Windows Sandbox is available...")

	cmd := exec.Command("wsb", "--help")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("windows sandbox CLI 'wsb' was not found: %w", err)
	}

	return nil
}
