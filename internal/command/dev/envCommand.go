package dev

import "github.com/ScriptoriumLab/scriptorium-cli/internal/project"

type envCommand interface {
	ensureEnv() error
	prepareEnv() error

	setupProductPrerequisites() error
	deployArtifacts(artifacts *project.ProjectArtifacts, dictionarySourcePath string) error
	startProduct() error
	startManualTests() error

	monitorEnv() error
	cleanupEnv() error
}
