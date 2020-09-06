package client

import (
	"testing"
)

func TestClient(t *testing.T) {
	cli, err := NewClient() // *Client
	if err != nil {
		t.Fatal(err)
	}

	problem := &Problem{}
	code := &Code{}

	submission, err := cli.Submit(problem, code) // *Submission
	if err != nil {
		t.Fatal(err)
	}

	submission.Status()      // SubmissionStatus
	submission.TimeScore()   // time.Duration
	submission.MemoryScore() // int
	submission.CodeLength()  // int
	submission.DetailURL()   // string

	testResults, err := cli.Test(problem, code) // []*TestResult
	if err != nil {
		t.Fatal(err)
	}

	testResults[0].Status() // TestResultStatus
	testResults[0].Input()  // string
	testResults[0].Output() // string
}
