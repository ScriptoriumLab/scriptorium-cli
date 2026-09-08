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
	VMSnapshotName       string
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
		VMSnapshotName:       os.Getenv("ORIUM_VM_SNAPSHOT_NAME"),
		VMEncryptionPassword: os.Getenv("ORIUM_VM_ENCRYPTION_PASSWORD"),
		GuestUsername:        os.Getenv("ORIUM_GUEST_USERNAME"),
		GuestPassword:        os.Getenv("ORIUM_GUEST_PASSWORD"),
	}

	if config.VMRunPath == "" {
		return nil, fmt.Errorf("ORIUM_VM_RUN_PATH is not configured")
	}

	if config.VMImagePath == "" {
		return nil, fmt.Errorf("ORIUM_VM_IMAGE_PATH is not configured")
	}

	if config.VMSnapshotName == "" {
		return nil, fmt.Errorf("ORIUM_VM_SNAPSHOT_NAME is not configured")
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
		return nil, fmt.Errorf("failed to load .scriptorium.src.env file: %w", err)
	}

	config := &WorkspaceConfig{
		RootPath: os.Getenv("SCRIPTORIUM_WORKSPACE_ROOT_PATH"),
	}

	if config.RootPath == "" {
		return nil, fmt.Errorf("SCRIPTORIUM_WORKSPACE_ROOT_PATH is not configured")
	}

	return config, nil
}

type ProductConfig struct {
	RootPath string
	ArtifactsPath string
}

func LoadProduct() (*ProductConfig, error) {
	fmt.Println("Loading product configuration...")

	if err := godotenv.Load(".scriptorium.product.env"); err != nil {
		return nil, fmt.Errorf("failed to load .scriptorium.product.env file: %w", err)
	}

	config := &ProductConfig{
		RootPath:      os.Getenv("SCRIPTORIUM_PRODUCT_ROOT_PATH"),
		ArtifactsPath: os.Getenv("SCRIPTORIUM_PRODUCT_ARTIFACT_PATH"),
	}

	if config.RootPath == "" {
		return nil, fmt.Errorf("SCRIPTORIUM_PRODUCT_ROOT_PATH is not configured")
	}

	if config.ArtifactsPath == "" {
		return nil, fmt.Errorf("SCRIPTORIUM_PRODUCT_ARTIFACT_PATH is not configured")
	}

	return config, nil
}

type SandboxConfig struct {
    WebView2BrowserExecutableFolder string
    NotepadPlusPlusPath             string
    TestFilePath                    string
    HostDependenciesPath            string
	DependenciesPath                string
}


func LoadSandbox() (*SandboxConfig, error) {
	fmt.Println("Loading windows sandbox configuration...")

	if err := godotenv.Load(".sandbox.env"); err != nil {
		return nil, fmt.Errorf("failed to load .sandbox.env file: %w", err)
	}

	config := &SandboxConfig{
		WebView2BrowserExecutableFolder: os.Getenv("ORIUM_SANDBOX_WEBVIEW2_BROWSER_EXECUTABLE_FOLDER"),
		NotepadPlusPlusPath:             os.Getenv("ORIUM_SANDBOX_NOTEPAD_PLUS_PLUS_PATH"),
		TestFilePath:                    os.Getenv("ORIUM_SANDBOX_TEST_FILE_PATH"),
		HostDependenciesPath:            os.Getenv("ORIUM_SANDBOX_HOST_DEPENDENCIES_PATH"),
		DependenciesPath:                os.Getenv("ORIUM_SANDBOX_DEPENDENCIES_PATH"),
	}

	if config.WebView2BrowserExecutableFolder == "" {
		return nil, fmt.Errorf("ORIUM_SANDBOX_WEBVIEW2_BROWSER_EXECUTABLE_FOLDER is not configured")
	}

	if config.NotepadPlusPlusPath == "" {
		return nil, fmt.Errorf("ORIUM_SANDBOX_NOTEPAD_PLUS_PLUS_PATH is not configured")
	}

	if config.TestFilePath == "" {
		return nil, fmt.Errorf("ORIUM_SANDBOX_TEST_FILE_PATH is not configured")
	}

	if config.HostDependenciesPath == "" {
		return nil, fmt.Errorf("ORIUM_SANDBOX_HOST_DEPENDENCIES_PATH is not configured")
	}

	return config, nil
}
