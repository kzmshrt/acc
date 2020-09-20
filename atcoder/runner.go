package atcoder

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Runner struct {
}

func NewRunner() (*Runner, error) {
	return &Runner{}, nil
}

func makeBuildCmdGo(filename, outFilename string) (*exec.Cmd, *bytes.Buffer) {
	var stderrBuffer bytes.Buffer
	build := exec.Command("go", "build", "-o", outFilename, filename)
	build.Stderr = &stderrBuffer
	return build, &stderrBuffer
}

func makeRunCmdGo(outFilename string, testCase *TestCase) (*exec.Cmd, *bytes.Buffer) {
	var stderrBuffer bytes.Buffer
	run := exec.Command(outFilename)
	run.Stdin = strings.NewReader(testCase.Input)
	run.Stderr = &stderrBuffer
	return run, &stderrBuffer
}

func testGo(filename string, testCases []*TestCase) ([]*TestResult, error) {
	removeFile := func(filename string) {
		fileExists := func(filename string) bool {
			_, err := os.Stat(filename)
			return err == nil
		}
		if fileExists(filename) {
			if err := os.Remove(filename); err != nil {
				fmt.Printf("os.Remove(%s) failed: %v", filename, err)
				os.Exit(1)
			}
		}
	}

	var testResults []*TestResult

	outFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	buildGo, stderr := makeBuildCmdGo(filename, outFilename)
	if err := buildGo.Run(); err != nil {
		testResults = append(testResults, &TestResult{TestStatus: TestStatusCE, ActualOutput: stderr.String()})
		return testResults, nil
	}

	defer removeFile(outFilename)

	for _, testCase := range testCases {
		runGo, stderr := makeRunCmdGo(outFilename, testCase)
		result, err := runGo.Output()

		if err != nil {
			testResults = append(testResults, &TestResult{
				TestStatus:     TestStatusRE,
				Input:          testCase.Input,
				ExpectedOutput: testCase.Output,
				ActualOutput:   stderr.String(),
			})
			continue
		}

		status := TestStatusOK
		if testCase.Output != string(result) {
			status = TestStatusWA
		}

		testResults = append(testResults, &TestResult{
			TestStatus:     status,
			Input:          testCase.Input,
			ExpectedOutput: testCase.Output,
			ActualOutput:   string(result),
		})
	}

	return testResults, nil
}

func (r *Runner) Run(filename string, testCases []*TestCase) ([]*TestResult, error) {
	switch ext := filepath.Ext(filename); ext {
	case ".go":
		return testGo(filename, testCases)
	default:
		return nil, fmt.Errorf("unsupported file type: %v", ext)
	}
}
