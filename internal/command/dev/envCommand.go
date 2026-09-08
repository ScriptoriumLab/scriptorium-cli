package dev

import (
	"github.com/ScriptoriumLab/scriptorium-cli/internal/command/dev/env/win/sandbox"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/command/dev/env/win/vm"
	"github.com/ScriptoriumLab/scriptorium-cli/internal/project"
)

type envCommand interface {
	EnsureEnv() error
	PrepareEnv() error

	SetupProductPrerequisites() error
	DeployArtifacts(artifacts *project.Artifacts, dictionarySourcePath string) error
	StartProduct() error
	StartManualTests() error

	MonitorEnv() error
	CleanupEnv() error
}

// Ensure specific environment command implements dev.envCommand.
var (
    _ envCommand = (*vm.Command)(nil)
    _ envCommand = (*sandbox.Command)(nil)
)
