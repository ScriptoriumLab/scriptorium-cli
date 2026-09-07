package project

import (
	"fmt"

	"github.com/ScriptoriumLab/scriptorium-cli/internal/build/pnpm"
)

type Ink struct{
	projectRootDir string
}

func (workspace *Workspace) Ink() *Ink {
	return &Ink{
		projectRootDir: workspace.config.RootPath + `\scriptorium-ink`,
	}
}

func (ink *Ink) build() (string, error) {
	fmt.Println("Building Scriptorium Ink...")

	if err := pnpm.BuildTauri(ink.projectRootDir); err != nil {
		return "", fmt.Errorf("scriptorium ink build failed: %w", err)
	}

	exe := ink.projectRootDir + `\src-tauri\target\release\scriptorium-ink.exe`
	return exe, nil
}
