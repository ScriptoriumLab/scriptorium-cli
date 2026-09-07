// Package project contains functions to build and test the Scriptorium project components.
package project

import (
	"fmt"
)

// TODO: Make the Scriptorium Project root directory configurable
const workspaceRoot = `D:\Projects\Scriptorium`

const DictionarySourceFile = workspaceRoot + `\scriptorium-inkstone\data\pinyin_dictionary.txt`

type ProjectArtifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

func BuildScriptoriumAndRunAllTests() (*ProjectArtifacts, error) {
	fmt.Println("Building Scriptorium and running all tests...")
	if err := NewFelt(workspaceRoot).buildAndTest(); err != nil {
		return nil, err
	}

	brushDll, err := NewBrush(workspaceRoot).buildAndTest()
	if err != nil {
		return nil, err
	}

	inkstoneExe, err := NewInkstone(workspaceRoot).buildAndTest()
	if err != nil {
		return nil, err
	}

	inkExe, err := NewInk(workspaceRoot).build()
	if err != nil {
		return nil, err
	}

	return &ProjectArtifacts{
		BrushDLL:    brushDll,
		InkstoneEXE: inkstoneExe,
		InkEXE:      inkExe,
	}, nil
}
