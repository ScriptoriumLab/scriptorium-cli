// Package env defines the interface for environment commands in the Scriptorium CLI tool. It provides methods for ensuring, preparing, deploying, starting, monitoring, and cleaning up the environment.
package env

import "github.com/ScriptoriumLab/scriptorium-cli/internal/project"

type Command interface {
	EnsureEnv() error
	PrepareEnv() error

	SetupProductPrerequisites() error
	DeployArtifacts(artifacts *project.ProjectArtifacts, dictionarySourcePath string) error
	StartProduct() error
	StartManualTests() error

	MonitorEnv() error
	CleanupEnv() error
}
