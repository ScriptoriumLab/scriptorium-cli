package project

type Dictionary struct {
	SourceFile string
}

func NewDictionary (workspaceRoot string) *Dictionary {
	return &Dictionary{
		SourceFile: workspaceRoot + `\scriptorium-inkstone\data\pinyin_dictionary.txt`,
	}
}
