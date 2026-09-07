// Package product provides the Product struct and related functionality for managing product artifacts and paths.
package product

import "github.com/ScriptoriumLab/scriptorium-cli/internal/config"

type ProductArtifacts struct {
	BrushDLL    string
	InkstoneEXE string
	InkEXE      string
}

type Product struct {
	Config         *config.ProductConfig
	LocalPath      string
	DictionaryPath string
	LogPath        string
	Artifacts      *ProductArtifacts
}

func NewProduct(config *config.ProductConfig) *Product {
	return &Product{
		Config:         config,
		LocalPath:      config.RootPath + `\Local`,
		DictionaryPath: config.RootPath + `\Local` + `\pinyin_dictionary.txt`,
		LogPath:        config.RootPath + `\Log`,
		Artifacts: &ProductArtifacts{
			BrushDLL:    config.ArtifactsPath + `\scriptorium-brush.dll`,
			InkstoneEXE: config.ArtifactsPath + `\scriptorium-inkstone.exe`,
			InkEXE:      config.ArtifactsPath + `\scriptorium-ink.exe`,
		},
	}
}
