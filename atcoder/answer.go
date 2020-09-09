package atcoder

import (
	"io/ioutil"
	"path/filepath"
)

const (
	LanguageIDGo = "4026"
)

var ext2lang = map[string]string{
	".go": LanguageIDGo,
}

type Answer struct {
	SourceCode string
	LanguageID string
}

func ParseAnswerFile(filename string) (*Answer, error) {
	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Answer{SourceCode: string(sourceCode), LanguageID: ext2lang[filepath.Ext(filename)]}, nil
}
