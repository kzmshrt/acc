package atcoder

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type document struct {
	doc *goquery.Document
}

type LoginPageDocument document

func NewLoginPageDocument(root *html.Node) *LoginPageDocument {
	return &LoginPageDocument{
		goquery.NewDocumentFromNode(root),
	}
}

func (l *LoginPageDocument) GetCSRFToken() (token string) {
	l.doc.Find("form[action=''] > input").Each(func(_ int, s *goquery.Selection) {
		token, _ = s.Attr("value")
	})
	return token
}

type TaskPageDocument document

func NewTaskPageDocument(root *html.Node) *TaskPageDocument {
	return &TaskPageDocument{
		goquery.NewDocumentFromNode(root),
	}
}

func (t *TaskPageDocument) GetTestCases() ([]*TestCase, error) {
	containsInArray := func(s string, substrs []string) bool {
		for _, substr := range substrs {
			if strings.Contains(s, substr) {
				return true
			}
		}
		return false
	}

	var inputs, outputs []string

	// input
	t.doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return containsInArray(s.Find("h3").Text(), []string{"入力例", "Sample Input"})
	}).Each(func(i int, s *goquery.Selection) {
		inputs = append(inputs, s.Find("pre").Text())
	})

	// output
	t.doc.Find(".part").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return containsInArray(s.Find("h3").Text(), []string{"出力例", "Sample Output"})
	}).Each(func(i int, s *goquery.Selection) {
		outputs = append(outputs, s.Find("pre").Text())
	})

	if ni, no := len(inputs), len(outputs); ni != no {
		return nil, fmt.Errorf("number of inputs and outputs does not match: (in, out) = (%d, %d)", ni, no)
	}

	testCases := make([]*TestCase, len(inputs), len(inputs))

	for i := range testCases {
		testCases[i] = &TestCase{
			Input:  inputs[i],
			Output: outputs[i],
		}
	}

	return testCases, nil
}

func (t *TaskPageDocument) GetCSRFToken() (token string) {
	t.doc.Find(".form-code-submit > input").Each(func(_ int, s *goquery.Selection) {
		token, _ = s.Attr("value")
	})
	return token
}

type SubmissionsMePageDocument document

func NewSubmissionsMePageDocument(root *html.Node) *SubmissionsMePageDocument {
	return &SubmissionsMePageDocument{
		goquery.NewDocumentFromNode(root),
	}
}

func (t *SubmissionsMePageDocument) GetSubmissions() (submissions []*Submission, err error) {
	decideJudge := func(s string) Judge {
		switch s {
		case "WJ":
			return JudgeWJ
		case "AC":
			return JudgeAC
		case "WA":
			return JudgeWA
		case "TLE":
			return JudgeTLE
		case "MLE":
			return JudgeMLE
		case "RE":
			return JudgeRE
		case "CE":
			return JudgeCE
		default:
			return JudgeNA
		}
	}

	t.doc.Find(".panel-submission").Find("table").Find("tbody").Find("tr").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		submission := new(Submission)

		judgeString := strings.Split(s.Find("td:nth-child(7) > span").Text(), " ")[0]
		judge := decideJudge(judgeString)
		if t.doc.Find(".panel-submission").Find("table").Find("tbody").Find("tr:nth-child(1)").Find("td:nth-child(7)").HasClass("waiting-judge") {
			judge = JudgeWJ
		}
		submission.Judge = judge

		codeSize, err := strconv.Atoi(strings.Split(s.Find("td:nth-child(6)").Text(), " ")[0])
		if err != nil {
			return false
		}
		submission.CodeSize = codeSize

		switch judge {
		case JudgeAC, JudgeWA, JudgeTLE, JudgeMLE, JudgeRE:
			timeScore, err := strconv.Atoi(strings.Split(s.Find("td:nth-child(8)").Text(), " ")[0])
			if err != nil {
				return false
			}
			submission.TimeScore = timeScore

			memoryScore, err := strconv.Atoi(strings.Split(s.Find("td:nth-child(9)").Text(), " ")[0])
			if err != nil {
				return false
			}
			submission.MemoryScore = memoryScore

			detailURLPath, ok := s.Find("td:nth-child(10) > a").First().Attr("href")
			if !ok {
				err = fmt.Errorf("failed parsing detail URL")
				return false
			}
			baseURL, _ := url.Parse(defaultBaseURL)
			detailURL, err := baseURL.Parse(detailURLPath)
			if err != nil {
				return false
			}
			submission.DetailURL = detailURL.String()
		case JudgeWJ, JudgeCE:
			detailURLPath, ok := s.Find("td:nth-child(8) > a").First().Attr("href")
			if !ok {
				err = fmt.Errorf("failed parsing detail URL")
				return false
			}
			baseURL, _ := url.Parse(defaultBaseURL)
			detailURL, err := baseURL.Parse(detailURLPath)
			if err != nil {
				return false
			}
			submission.DetailURL = detailURL.String()
		default:
			err = fmt.Errorf("unknown status string: %s", judge)
			return false
		}

		submissions = append(submissions, submission)
		return true
	})

	if err != nil {
		return nil, err
	}
	return submissions, nil
}
