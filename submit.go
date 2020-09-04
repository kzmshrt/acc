package acc

import (
	"io/ioutil"
)

type Submitter interface {
	Submit(code, url string) (*Submission, error)
}

func Submit(url, filename string) (*Submission, error) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	submitter := InitSubmitter()
	return submitter.Submit(string(code), url)
}
