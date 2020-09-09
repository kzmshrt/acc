package atcoder

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
)

const (
	defaultBaseURL = "https://atcoder.jp/"
)

const (
	pathLogin = "/login"

	pathFormatTask         = "/contests/%s/tasks/%s"
	pathFormatSubmit       = "/contests/%s/submit"
	pathFormatSubmissionMe = "/contests/%s/submissions/me"
)

type RESTClient struct {
	BaseURL *url.URL
	Client  *http.Client
}

func NewRESTClient() (*RESTClient, error) {
	u, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	client.Jar = jar

	return &RESTClient{BaseURL: u, Client: client}, nil
}

func (c *RESTClient) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *RESTClient) Authenticate(username, password string) error {
	buildLoginForm := func(username, password, csrfToken string) url.Values {
		form := make(url.Values)
		form.Set("username", username)
		form.Set("password", password)
		form.Set("csrf_token", csrfToken)
		return form
	}

	req, err := c.NewRequest(http.MethodGet, pathLogin, nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	csrfToken := NewLoginPageDocument(node).GetCSRFToken()

	u, _ := c.BaseURL.Parse(pathLogin)
	resp, err = c.Client.PostForm(u.String(), buildLoginForm(username, password, csrfToken))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v response submitting login form", resp.StatusCode)
	}

	return nil
}

func (c *RESTClient) SubmitFile(taskURL, filename string) error {
	task, err := ParseTaskURL(taskURL)
	if err != nil {
		return err
	}
	answer, err := ParseAnswerFile(filename)
	if err != nil {
		return err
	}
	return c.Submit(task, answer)
}

func (c *RESTClient) Submit(task *Task, answer *Answer) error {
	buildSubmitForm := func(task *Task, answer *Answer, csrfToken string) url.Values {
		form := make(url.Values)
		form.Set("data.TaskScreenName", task.TaskID)
		form.Set("data.LanguageId", answer.LanguageID)
		form.Set("sourceCode", answer.SourceCode)
		form.Set("csrf_token", csrfToken)
		return form
	}

	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf(pathFormatTask, task.ContestID, task.TaskID), nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	csrfToken := NewTaskPageDocument(node).GetCSRFToken()

	u, _ := c.BaseURL.Parse(fmt.Sprintf(pathFormatSubmit, task.ContestID))
	resp, err = c.Client.PostForm(u.String(), buildSubmitForm(task, answer, csrfToken))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v response submitting submit form", resp.StatusCode)
	}

	return nil
}

func (c *RESTClient) ListSubmissions(contestID string) ([]*Submission, error) {
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf(pathFormatSubmissionMe, contestID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	submissions, err := NewSubmissionsMePageDocument(node).GetSubmissions()
	if err != nil {
		return nil, err
	}

	return submissions, nil
}
