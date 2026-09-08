// Package vm provides functions to manage the development virtual machine (VM) on Windows using VMware Workstation.
package vm

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
)

type VM struct {
	config *config.VMConfig
}

func New(config *config.VMConfig) *VM {
	return &VM{
		config: config,
	}
}

func (vm *VM) EnsureAvailable() error {
	fmt.Println("Ensuring VMware is available...")
	if _, err := os.Stat(vm.config.VMRunPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VMware CLI 'vmrun' was not found at %s", vm.config.VMRunPath)
		}

		return fmt.Errorf("failed to check VMware CLI: %w", err)
	}

	return nil
}

func (vm *VM) Prepare() error {
	fmt.Println("Preparing VM...")

	if err := vm.reset(); err != nil {
		return err
	}

	if err := vm.start(); err != nil {
		return err
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

func (vm *VM) Cleanup() error {
	return vm.reset()
}
