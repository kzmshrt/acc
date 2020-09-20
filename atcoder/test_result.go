package atcoder

type TestStatus string

const (
	TestStatusOK TestStatus = "OK"
	TestStatusWA TestStatus = "WA"
	TestStatusRE TestStatus = "RE"
	TestStatusCE TestStatus = "CE"
)

type TestResult struct {
	TestStatus     TestStatus
	Input          string
	ExpectedOutput string
	ActualOutput   string
}
