// Package sandbox provides functions to manage the development sandbox on Windows using VMware Workstation.
package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"

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

func (sandbox *Sandbox) Prepare() error {
	fmt.Println("Preparing Windows Sandbox...")

	if err := sandbox.start(); err != nil {
		return err
	}

	if err := sandbox.connect(); err != nil {
		return err
	}

	if err := sandbox.setupEnv(); err != nil {
		return err
	}

	return nil
}

func (sandbox *Sandbox) Monitor() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	fmt.Println("Waiting for the development windows sandbox to stop...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt received. Stopping development widnows sandbox...")

			if err := sandbox.stop(); err != nil {
				return err
			}

			return nil
		default:
			running, err := sandbox.isRunning()
			if err != nil {
				return err
			}

			if !running {
				if err := sandbox.stop(); err != nil {
					return err
				}

				fmt.Println("Development windows sandbox stopped.")
				return nil
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func (sandbox *Sandbox) Cleanup() error {
	fmt.Println("Cleanup the development Windows Sandbox...")
	if err := os.RemoveAll(sandbox.TempStagingDir); err != nil {
		fmt.Printf("Failed to remove Sandbox staging directory %q: %v\n", sandbox.TempStagingDir, err)
	}

	return nil
}
