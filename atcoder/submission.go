package atcoder

import "time"

type SubmissionStatus int

const (
	SubmissionStatusAC  SubmissionStatus = iota // Accepted
	SubmissionStatusWA                          // Wrong Answer
	SubmissionStatusTLE                         // Time Limit Exceed
	SubmissionStatusMLE                         // Memory Limit Exceed
	SubmissionStatusRE                          // Runtime Error
	SubmissionStatusCE                          // Compile Error
)

var status2txt = map[SubmissionStatus]string{
	SubmissionStatusAC:  "AC",
	SubmissionStatusWA:  "WA",
	SubmissionStatusTLE: "TLE",
	SubmissionStatusMLE: "MLE",
	SubmissionStatusRE:  "RE",
	SubmissionStatusCE:  "CE",
}

var txt2status = map[string]SubmissionStatus{
	"AC":  SubmissionStatusAC,
	"WA":  SubmissionStatusWA,
	"TLE": SubmissionStatusTLE,
	"MLE": SubmissionStatusMLE,
	"RE":  SubmissionStatusRE,
	"CE":  SubmissionStatusCE,
}

func NewSubmissionStatusFromText(txt string) SubmissionStatus {
	return txt2status[txt]
}

func (s SubmissionStatus) String() string {
	return status2txt[s]
}

type Submission struct {
	Status      SubmissionStatus
	CodeLength  int
	TimeScore   time.Duration
	MemoryScore int
	DetailUrl   string
}
