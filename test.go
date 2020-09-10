package acc

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func getTestCases(url string) ([]*testCase, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := goquery.NewDocumentFromNode(node)

	var inputs, outputs []string

	doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.Find("h3").Text(), "入力例")
	}).Each(func(i int, s *goquery.Selection) {
		inputs = append(inputs, s.Find("pre").Text())
	})
	doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.Find("h3").Text(), "出力例")
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

	return testCases, nil
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
