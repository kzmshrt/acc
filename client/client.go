package client

import (
	"time"
)

type Problem struct {
}

type Code struct {
}

type Submission struct {
	submissionStatus SubmissionStatus
	timeScore        time.Duration
	memoryScore      int
	codeLength       int
	detailURL        string
}

type SubmissionStatus int

const (
	SubmissionStatusAC SubmissionStatus = iota
)

func (s *Submission) Status() SubmissionStatus {
	return s.submissionStatus
}

func (s *Submission) TimeScore() time.Duration {
	return s.timeScore
}

func (s *Submission) MemoryScore() int {
	return s.memoryScore
}

func (s *Submission) CodeLength() int {
	return s.codeLength
}

func (s *Submission) DetailURL() string {
	return s.detailURL
}

type TestStatus int

const (
	TestStatusOK TestStatus = iota
)

type TestResult struct {
	testStatus TestStatus
	input      string
	output     string
}

func (t *TestResult) Status() TestStatus {
	return t.testStatus
}

func (t *TestResult) Input() string {
	return t.input
}

func (t *TestResult) Output() string {
	return t.output
}

type Client struct {
}

func NewClient() (*Client, error) {
	return nil, nil
}

func (c *Client) Submit(problem *Problem, code *Code) (*Submission, error) {
	return &Submission{}, nil
}

func (c *Client) Test(problem *Problem, code *Code) ([]*TestResult, error) {
	return []*TestResult{&TestResult{}}, nil
}
