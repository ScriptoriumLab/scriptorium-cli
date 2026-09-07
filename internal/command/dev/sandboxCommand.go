package dev

type sandboxCommand struct {}

func newSandboxCommand() *sandboxCommand {
	return &sandboxCommand{}
}

func (sandboxCmd *sandboxCommand) execute() error {
	return nil
}
