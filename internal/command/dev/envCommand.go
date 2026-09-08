package dev

import "github.com/ScriptoriumLab/scriptorium-cli/internal/project"

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
