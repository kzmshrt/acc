package atcoder

import (
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
	"golang.org/x/net/publicsuffix"
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

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
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
	token, exists := doc.
		Find(formSelectorStr).
		Find(fmt.Sprintf("input[name=\"%s\"]", "csrf_token")).
		Attr("value")
	if !exists {
		return "", fmt.Errorf("cannot find input[name=csrf_token]")
	}
	return token, nil
}

func (c *RESTClient) buildLoginForm(csrfToken string) *url.Values {
	form := new(url.Values)
	form.Set("username", c.username)
	form.Set("password", c.password)
	form.Set("csrf_token", csrfToken)
	return form
}

func (c *RESTClient) Authenticate() error {
	req, err := c.newRequest(http.MethodGet, "/login", nil)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	csrfToken, err := c.getCsrfToken(res.Body, "form[action='']")
	if err != nil {
		return err
	}
	loginForm := c.buildLoginForm(csrfToken)
	req, err = c.newRequest(http.MethodPost, "/login", strings.NewReader(loginForm.Encode()))
	if err != nil {
		return err
	}
	_, err = c.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTClient) buildSubmitForm(answer *Answer, problem *Problem, csrfToken string) *url.Values {
	form := new(url.Values)
	form.Set("data.TaskScreenName", problem.TaskID)
	form.Set("data.LanguageId", answer.Lang.AtCoderLangID())
	form.Set("sourceCode", answer.Code)
	form.Set("csrf_token", csrfToken)
	return form
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
			Find("panel-submission").
			Find("table").
			Find("tbody").
			Find("tr:nth-child(1)").
			Find("td:nth-child(7)").
			Find("span").
			HasClass("label-default"),
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

	firstSubmissionRow := doc.Find("panel-submission").Find("table").Find("tbody").Find("tr:nth-child(1)")

	codeLengthCell := firstSubmissionRow.Find("td:nth-child(6)")
	statusCell := firstSubmissionRow.Find("td:nth-child(7)")

	codeLength, err := strconv.Atoi(strings.Split(codeLengthCell.Text(), " ")[0])
	if err != nil {
		return nil, err
	}
	status := NewSubmissionStatusFromText(statusCell.Find("span").Text())

	if status != SubmissionStatusAC {
		return &Submission{
			CodeLength: codeLength,
			Status:     status,
		}, nil
	}

	timeScoreCell := firstSubmissionRow.Find("td:nth-child(8)")
	memoryScoreCell := firstSubmissionRow.Find("td:nth-child(9)")
	detailURLCell := firstSubmissionRow.Find("td:nth-child(10)")

	timeScore, err := strconv.Atoi(strings.Split(timeScoreCell.Text(), " ")[0])
	if err != nil {
		return nil, err
	}
	memoryScore, err := strconv.Atoi(strings.Split(memoryScoreCell.Text(), " ")[0])
	if err != nil {
		return nil, err
	}
	detailURL := strings.Split(detailURLCell.Text(), " ")[0]

	return &Submission{
		Status:      status,
		CodeLength:  codeLength,
		TimeScore:   time.Duration(timeScore) * time.Millisecond,
		MemoryScore: memoryScore,
		DetailUrl:   detailURL,
	}, nil
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
	// post answer
	req, err := c.newRequest(http.MethodGet, problem.URL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	csrfToken, err := c.getCsrfToken(res.Body, "form .form-code-submit")
	if err != nil {
		return nil, err
	}
	submitForm := c.buildSubmitForm(answer, problem, csrfToken)
	req, err = c.newRequest(http.MethodPost, problem.URL, strings.NewReader(submitForm.Encode()))
	if err != nil {
		return nil, err
	}
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
		time.Sleep(300 * time.Millisecond)
		judging, err = c.getIsJudging(problem)
		if err != nil {
			return nil, err
		}
	}

	// check submission
	return c.getFirstSubmission(problem)
}
