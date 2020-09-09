package atcoder

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RESTClient struct {
	url      *url.URL
	client   *http.Client
	username string
	password string
}

func NewRESTClient() (*RESTClient, error) {
	u, err := url.Parse("https://atcoder.jp/")
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	client.Jar = jar

	username := os.Getenv("ATCODER_USERNAME")
	password := os.Getenv("ATCODER_PASSWORD")

	return &RESTClient{
		url:      u,
		client:   client,
		username: username,
		password: password,
	}, nil
}

func (c *RESTClient) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := c.url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *RESTClient) getCsrfToken(resBody io.Reader, formSelectorStr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(resBody)
	if err != nil {
		return "", err
	}
	attributes := doc.Find(formSelectorStr).Find("input[name='csrf_token']").Nodes[0].Attr
	for _, attr := range attributes {
		if attr.Key == "value" {
			return attr.Val, nil
		}
	}
	return "", errors.New("cannot find csrf_token")
}

func (c *RESTClient) buildLoginForm(csrfToken string) *url.Values {
	form := make(url.Values)
	form.Set("username", c.username)
	form.Set("password", c.password)
	form.Set("csrf_token", csrfToken)
	return &form
}

func (c *RESTClient) Authenticate() (*http.Response, error) {
	// get token
	req, err := c.newRequest(http.MethodGet, "/login", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	csrfToken, err := c.getCsrfToken(res.Body, "form[action='']")
	if err != nil {
		return nil, err
	}

	// post login
	loginForm := c.buildLoginForm(csrfToken)
	req, err = c.newRequest(http.MethodPost, "/login", strings.NewReader(loginForm.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RESTClient) buildSubmitForm(answer *Answer, problem *Problem, csrfToken string) *url.Values {
	form := make(url.Values)
	form.Set("data.TaskScreenName", problem.TaskID)
	form.Set("data.LanguageId", answer.Lang.AtCoderLangID())
	form.Set("sourceCode", answer.Code)
	form.Set("csrf_token", csrfToken)
	return &form
}

func (c *RESTClient) SubmitFile(filename, problemURL string) (*Submission, error) {
	answer, err := NewAnswerFromFile(filename)
	if err != nil {
		return nil, err
	}
	problem, err := NewProblemFromURL(problemURL)
	if err != nil {
		return nil, err
	}
	return c.Submit(answer, problem)
}

func (c *RESTClient) getIsJudgingLatestSubmission(resBody io.Reader) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(resBody)
	if err != nil {
		return false, err
	}
	return doc.
			Find(".panel-submission").
			Find("table").
			Find("tbody").
			Find("tr:nth-child(1)").
			Find("td:nth-child(7)").
			HasClass("waiting-judge"),
		nil
}

func (c *RESTClient) getIsJudging(problem *Problem) (bool, error) {
	req, err := c.newRequest(http.MethodGet, problem.SubmissionURL(), nil)
	if err != nil {
		return false, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	return c.getIsJudgingLatestSubmission(res.Body)
}

func (c *RESTClient) extractFirstSubmission(resBody io.Reader) (*Submission, error) {
	doc, err := goquery.NewDocumentFromReader(resBody)
	if err != nil {
		return nil, err
	}

	firstSubmissionRow := doc.
		Find(".panel-submission").
		Find("table").
		Find("tbody").
		Find("tr:nth-child(1)")

	codeLengthCell := firstSubmissionRow.Find("td:nth-child(6)")
	codeLength, err := strconv.Atoi(strings.Split(codeLengthCell.Text(), " ")[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse code length: %v", err)
	}

	statusCell := firstSubmissionRow.Find("td:nth-child(7)")
	status := NewSubmissionStatusFromText(statusCell.Find("span").Text())

	switch status {
	case SubmissionStatusAC, SubmissionStatusWA:
		timeScoreCell := firstSubmissionRow.Find("td:nth-child(8)")
		timeScore, err := strconv.Atoi(strings.Split(timeScoreCell.Text(), " ")[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time score: %v", err)
		}

		memoryScoreCell := firstSubmissionRow.Find("td:nth-child(9)")
		memoryScore, err := strconv.Atoi(strings.Split(memoryScoreCell.Text(), " ")[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse memory score: %v", err)
		}

		detailURLCell := firstSubmissionRow.Find("td:nth-child(10)")
		detailURLPath, _ := detailURLCell.Find("a").First().Attr("href")
		detailURL, _ := c.url.Parse(detailURLPath)

		return &Submission{Status: status, CodeLength: codeLength, TimeScore: timeScore, MemoryScore: memoryScore, DetailUrl: detailURL.String()}, nil
	default:
		detailURLCell := firstSubmissionRow.Find("td:nth-child(8)")
		detailURLPath, _ := detailURLCell.Find("a").First().Attr("href")
		detailURL, _ := c.url.Parse(detailURLPath)

		return &Submission{CodeLength: codeLength, Status: status, DetailUrl: detailURL.String()}, nil
	}
}

func (c *RESTClient) getFirstSubmission(problem *Problem) (*Submission, error) {
	req, err := c.newRequest(http.MethodGet, problem.SubmissionURL(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return c.extractFirstSubmission(res.Body)
}

func (c *RESTClient) Submit(answer *Answer, problem *Problem) (*Submission, error) {
	// get token
	req, err := c.newRequest(http.MethodGet, problem.URL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	csrfToken, err := c.getCsrfToken(res.Body, ".form-code-submit")
	if err != nil {
		return nil, err
	}

	// post answer
	submitForm := c.buildSubmitForm(answer, problem, csrfToken)
	req, err = c.newRequest(http.MethodPost, problem.ActionPath(), strings.NewReader(submitForm.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// wait judge
	judging, err := c.getIsJudging(problem)
	if err != nil {
		return nil, err
	}
	for judging {
		time.Sleep(200 * time.Millisecond)
		judging, err = c.getIsJudging(problem)
		if err != nil {
			return nil, err
		}
	}

	// check submission
	return c.getFirstSubmission(problem)
}
