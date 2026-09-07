package project

type Dictionary struct {
	sourceFile string
}

func (workspace *Workspace) NewDictionary () *Dictionary {
	return &Dictionary{
		sourceFile: workspace.root + `\scriptorium-inkstone\data\pinyin_dictionary.txt`,
	}
}

func (dictionary *Dictionary) SourceFile() string {
	return dictionary.sourceFile
}
