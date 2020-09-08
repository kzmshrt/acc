package atcoder

import (
	"io/ioutil"
)

type Lang int

const (
	LangGo Lang = iota
)

var ext2lang = map[string]Lang{
	".go": LangGo,
}

var lang2id = map[Lang]string{
	LangGo: "4026",
}

func (l Lang) AtCoderLangID() string {
	return lang2id[l]
}

type Answer struct {
	Lang Lang
	Code string
}

func NewAnswerFromFile(filename string) (*Answer, error) {
	codeBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Answer{
		Lang: ext2lang[filename],
		Code: string(codeBytes),
	}, nil
}
