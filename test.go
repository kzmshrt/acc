package acc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
	"golang.org/x/net/html"
)

func Test(filename, url string) error {
	testCases, err := getTestCases(url)
	if err != nil {
		return err
	}

	switch ext := filepath.Ext(filename); ext {
	case ".go":
		err = testGo(filename, testCases)
	default:
		return fmt.Errorf("unsupported file type: %v", ext)
	}
	if err != nil {
		return err
	}

	return nil
}

type testCase struct {
	Input  string
	Output string
}

func containsInArray(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func getTestCases(url string) ([]*testCase, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// prepare io.TeeReader
	var r io.Reader = resp.Body
	buf := bytes.NewBuffer(nil)
	r = io.TeeReader(r, buf)

	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	doc := goquery.NewDocumentFromNode(node)

	var inputs, outputs []string

	// input
	doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return containsInArray(s.Find("h3").Text(), []string{"入力例", "Sample Input"})
	}).Each(func(i int, s *goquery.Selection) {
		inputs = append(inputs, s.Find("pre").Text())
	})

	// output
	doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return containsInArray(s.Find("h3").Text(), []string{"出力例", "Sample Output"})
	}).Each(func(i int, s *goquery.Selection) {
		outputs = append(outputs, s.Find("pre").Text())
	})

	if ni, no := len(inputs), len(outputs); ni != no {
		return nil, fmt.Errorf("number of inputs and outputs does not match: (in, out) = (%d, %d)", ni, no)
	}

	testCases := make([]*testCase, len(inputs), len(inputs))

	for i := range testCases {
		testCases[i] = &testCase{
			Input:  inputs[i],
			Output: outputs[i],
		}
	}

	// write html content to file
	if err := saveFile(buf); err != nil {
		return nil, fmt.Errorf("error while saving HTML to file: %v", err)
	}

	return testCases, nil
}

// saveFile save buffer content to file.
// This function is used to check html content of AtCoder task page for debug.
func saveFile(buf *bytes.Buffer) error {
	filename := fmt.Sprintf("atcoder_task_doc_%s.html", time.Now().Format("2006-01-02-15-04-05"))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(buf.String())
	if err != nil {
		return err
	}
	return nil
}

func testGo(filename string, testCases []*testCase) error {
	outFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	buildGo, stderr := makeBuildGoCommand(filename, outFilename)
	if err := buildGo.Run(); err != nil {
		printErrorResult(buildError, stderr.Bytes())
		return err
	}
	defer removeFile(outFilename)

	for _, testCase := range testCases {
		runGo, stderr := makeRunGoCommand(outFilename, testCase)
		result, err := runGo.Output()
		if err != nil {
			printErrorResult(runtimeError, stderr.Bytes())
			return err
		}
		printResult(testCase, result)
	}

	return nil
}

type errorType int

const (
	buildError errorType = iota
	runtimeError
)

func printErrorResult(errorType errorType, result []byte) {
	errorTypeString := ""
	switch errorType {
	case buildError:
		errorTypeString = "COMPILE ERROR:"
	case runtimeError:
		errorTypeString = "RUNTIME ERROR:"
	}

	fmt.Println(aurora.Bold("=================================================="))
	fmt.Println(aurora.Blue(errorTypeString))
	fmt.Println(aurora.Red(string(result)))
	fmt.Println(aurora.Bold("=================================================="))
}

func printResult(testCase *testCase, result []byte) {
	status := aurora.Bold(aurora.Green("Consistent"))
	if string(result) != testCase.Output {
		status = aurora.Bold(aurora.Red("Inconsistent"))
	}

	fmt.Println(aurora.Bold("=================================================="))
	fmt.Println(aurora.Blue("INPUT:"))
	fmt.Println(testCase.Input)
	fmt.Println(aurora.Blue("EXPECTED OUTPUT:"))
	fmt.Println(testCase.Output)
	fmt.Println(aurora.Blue("ACTUAL OUTPUT:"))
	fmt.Println(string(result))
	fmt.Println(aurora.Blue("RESULT:"))
	fmt.Println(status)
	fmt.Println(aurora.Bold("=================================================="))
}

func makeBuildGoCommand(filename, outFilename string) (*exec.Cmd, *bytes.Buffer) {
	var stderrBuffer bytes.Buffer
	build := exec.Command("go", "build", "-o", outFilename, filename)
	build.Stderr = &stderrBuffer
	return build, &stderrBuffer
}

func makeRunGoCommand(outFilename string, testCase *testCase) (*exec.Cmd, *bytes.Buffer) {
	var stderrBuffer bytes.Buffer
	run := exec.Command(outFilename)
	run.Stdin = strings.NewReader(testCase.Input)
	run.Stderr = &stderrBuffer
	return run, &stderrBuffer
}

func removeFile(filename string) {
	if fileExists(filename) {
		if err := os.Remove(filename); err != nil {
			fmt.Printf("os.Remove(%s) failed: %v", filename, err)
			os.Exit(1)
		}
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
