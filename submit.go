package acc

import (
	"io/ioutil"
	"time"
)

type Submitter interface {
	Submit(code, url string) (*Submission, error)
}

func Submit(filename, url string) (*Submission, error) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	submitter := InitSubmitter()
	return submitter.Submit(string(code), url)
}

type SubmissionStatus int

const (
	AC  SubmissionStatus = iota // Accepted
	WA                          // Wrong Answer
	TLE                         // Time Limit Exceed
	MLE                         // Memory Limit Exceed
	RE                          // Runtime Error
	CE                          // Compile Error
)

func (s SubmissionStatus) String() string {
	switch s {
	case AC:
		return "AC"
	case WA:
		return "WA"
	case TLE:
		return "TLE"
	case MLE:
		return "MLE"
	case RE:
		return "RE"
	case CE:
		return "CE"
	default:
		return ""
	}
}

type Submission struct {
	Status      SubmissionStatus
	CodeLength  int
	TimeScore   time.Duration
	MemoryScore int
	DetailUrl   string
}
