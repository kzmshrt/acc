package atcoder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v response submitting login form", resp.StatusCode)
	}

	return nil
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

func (c *RESTClient) GetTestCases(task *Task) ([]*TestCase, error) {
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf(pathFormatTask, task.ContestID, task.TaskID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// prepare io.TeeReader
	var r io.Reader = resp.Body
	buf := bytes.NewBuffer(nil)
	r = io.TeeReader(r, buf)

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	// write html content to file
	if err := saveFile(buf); err != nil {
		return nil, fmt.Errorf("error while saving HTML to file: %v", err)
	}

	return NewTaskPageDocument(node).GetTestCases()
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return NewSubmissionsMePageDocument(node).GetSubmissions()
}
