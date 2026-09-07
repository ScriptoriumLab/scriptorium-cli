// Package sandbox provides functions to manage the development sandbox on Windows using VMware Workstation.
package sandbox

import (
	"fmt"
	"os/exec"
)

type Sandbox struct {}

func New() *Sandbox {
	return &Sandbox{}
}

func (sandbox *Sandbox) EnsureAvailable() error {
	fmt.Println("Ensuring Windows Sandbox is available...")

	if _, err := exec.LookPath("wsb"); err != nil {
		return fmt.Errorf("windows sandbox CLI 'wsb' was not found: %w", err)
	}

	return nil
}
