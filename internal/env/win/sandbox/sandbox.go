// Package sandbox provides functions to manage the development sandbox on Windows using VMware Workstation.
package sandbox

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
	"github.com/google/uuid"
)

type Sandbox struct {
	Config           *config.SandboxConfig
	id               string
	remoteSessionPID int
	interactiveReady bool
	TempStagingDir   string
}

func New(config *config.SandboxConfig) *Sandbox {
	return &Sandbox{
		Config: config,
		id: uuid.NewString(),
	}
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
