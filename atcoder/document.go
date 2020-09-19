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

func (t *TaskPageDocument) GetTaskCases() ([]*TestCase, error) {
	return nil, nil
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
	t.doc.Find(".panel-submission").Find("table").Find("tbody").Find("tr").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		submission := new(Submission)

		judge := strings.Split(s.Find("td:nth-child(7) > span").Text(), " ")[0]
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
