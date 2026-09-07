// Package project contains functions to build and test the Scriptorium project components.
package project

import (
	"fmt"
)

type ProjectArtifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

type Workspace struct {
	root string
}

func NewWorkspace(root string) *Workspace {
	return &Workspace{
		root: root,
	}
}

func (workspace *Workspace) BuildScriptoriumAndRunAllTests() (*ProjectArtifacts, error) {
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

	return &ProjectArtifacts{
		BrushDLL:    brushDll,
		InkstoneEXE: inkstoneExe,
		InkEXE:      inkExe,
	}, nil
}
