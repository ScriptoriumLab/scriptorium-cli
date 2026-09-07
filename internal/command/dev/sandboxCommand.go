package dev

import "github.com/ScriptoriumLab/scriptorium-cli/internal/env/win/sandbox"

type sandboxCommand struct {
	sandbox *sandbox.Sandbox
}

func newSandboxCommand() *sandboxCommand {
	return &sandboxCommand{}
}

func (sandboxCmd *sandboxCommand) execute() error {
	sandboxCmd.sandbox = sandbox.New()

	if err := sandboxCmd.sandbox.EnsureAvailable(); err != nil {
		return err
	}

	return nil
}
