// Package project contains functions to build and test the Scriptorium project components.
package project

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/config"
)

type Artifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

type Workspace struct {
	config *config.WorkspaceConfig
}

func NewWorkspace(config *config.WorkspaceConfig) *Workspace {
	return &Workspace{
		config: config,
	}
}

func (workspace *Workspace) BuildScriptoriumAndRunAllTests() (*Artifacts, error) {
	fmt.Println("Building Scriptorium and running all tests...")
	if err := workspace.Felt().buildAndTest(); err != nil {
		return nil, err
	}

	brushDll, err := workspace.Brush().buildAndTest()
	if err != nil {
		return nil, err
	}

	inkstoneExe, err := workspace.Inkstone().buildAndTest()
	if err != nil {
		return nil, err
	}

	inkExe, err := workspace.Ink().build()
	if err != nil {
		return nil, err
	}

	return &Artifacts{
		BrushDLL:    brushDll,
		InkstoneEXE: inkstoneExe,
		InkEXE:      inkExe,
	}, nil
}
