// Package config provides configuration loading and validation.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type VMConfig struct {
	VMRunPath            string
	VMImagePath          string
	VMEncryptionPassword string
	GuestUsername        string
	GuestPassword        string
}

func LoadVM() (*VMConfig, error) {
	fmt.Println("Loading VM configuration...")

	if err := godotenv.Load(".vm.env"); err != nil {
		return nil, fmt.Errorf("failed to load .vm.env file: %w", err)
	}

	config := &VMConfig{
		VMRunPath:            os.Getenv("ORIUM_VM_RUN_PATH"),
		VMImagePath:          os.Getenv("ORIUM_VM_IMAGE_PATH"),
		VMEncryptionPassword: os.Getenv("ORIUM_VM_ENCRYPTION_PASSWORD"),
		GuestUsername:        os.Getenv("ORIUM_GUEST_USERNAME"),
		GuestPassword:        os.Getenv("ORIUM_GUEST_PASSWORD"),
	}

	if config.VMEncryptionPassword == "" {
		return nil, fmt.Errorf("ORIUM_VM_ENCRYPTION_PASSWORD is not configured")
	}

	if config.GuestUsername == "" {
		return nil, fmt.Errorf("ORIUM_GUEST_USERNAME is not configured")
	}

	if config.GuestPassword == "" {
		return nil, fmt.Errorf("ORIUM_GUEST_PASSWORD is not configured")
	}

	return config, nil
}

type WorkspaceConfig struct {
	RootPath string
}

func LoadWorkspace() (*WorkspaceConfig, error) {
	fmt.Println("Loading workspace configuration...")

	if err := godotenv.Load(".scriptorium.src.env"); err != nil {
		return nil, fmt.Errorf("failed to load .scriptorium.env file: %w", err)
	}

	config := &WorkspaceConfig{
		RootPath: os.Getenv("SCRIPTORIUM_WORKSPACE_ROOT_PATH"),
	}

	if config.RootPath == "" {
		return nil, fmt.Errorf("SCRIPTORIUM_WORKSPACE_ROOT_PATH is not configured")
	}

	return config, nil
}
