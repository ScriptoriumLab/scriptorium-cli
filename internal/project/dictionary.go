package project

type Dictionary struct {
	sourceFile string
}

func NewDictionary (workspaceRoot string) *Dictionary {
	return &Dictionary{
		sourceFile: workspaceRoot + `\scriptorium-inkstone\data\pinyin_dictionary.txt`,
	}
}

func (dictionary *Dictionary) SourceFile() string {
	return dictionary.sourceFile
}
