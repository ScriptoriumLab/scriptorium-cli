package project

type Dictionary struct {
	sourceFile string
}

func (workspace *Workspace) Dictionary () *Dictionary {
	return &Dictionary{
		sourceFile: workspace.config.RootPath + `\scriptorium-inkstone\data\pinyin_dictionary.txt`,
	}
}

func (dictionary *Dictionary) SourceFile() string {
	return dictionary.sourceFile
}
