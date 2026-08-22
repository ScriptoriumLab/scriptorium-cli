// Package vm provides functions to manage the development virtual machine (VM) on Windows using VMware Workstation.
package vm

import (
	"fmt"
	"os"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
)

// vmrunPath TODO: Discover vmrun dynamically instead of relying on the default
// VMware Workstation installation path.
const vmrunPath = `C:\Program Files\VMware\VMware Workstation\vmrun.exe`

// devVMPath TODO: Make the development VM path configurable.
const devVMPath = `D:\Projects\Scriptorium\dev-env\scriptorium-dev\scriptorium-dev.vmx`

type VM struct {
	config *config.Config
}

func New(config *config.Config) *VM {
	return &VM{
		config: config,
	}
}

func (vm *VM) EnsureAvailable() error {
	fmt.Println("Ensuring VMware is available...")
	if _, err := os.Stat(vmrunPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("VMware CLI 'vmrun' was not found at %s", vmrunPath)
		}

		return fmt.Errorf("failed to check VMware CLI: %w", err)
	}

	return nil
}

