// Package config provides configuration loading and validation.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VMEncryptionPassword string
	GuestUsername        string
	GuestPassword        string
}

func LoadConfig() (*Config, error) {
	fmt.Println("Loading configuration files and environment variables...")

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	config := &Config{
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
