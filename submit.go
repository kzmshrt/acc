package acc

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sclevine/agouti"
)

type resultStatus int

const (
	ac  resultStatus = iota // Accepted
	wa                      // Wrong Answer
	tle                     // Time Limit Exceed
	mle                     // Memory Limit Exceed
	re                      // Runtime Error
	ce                      // Compile Error
)

func (s resultStatus) String() string {
	switch s {
	case ac:
		return "AC"
	case wa:
		return "WA"
	case tle:
		return "TLE"
	case mle:
		return "MLE"
	case re:
		return "RE"
	case ce:
		return "CE"
	default:
		return ""
	}
}

type submitCodeLength int

func (cl submitCodeLength) String() string {
	return fmt.Sprintf("%d Byte", cl)
}

type submitTimeScore int

func (ts submitTimeScore) String() string {
	return fmt.Sprintf("%d ms", ts)
}

type submitMemoryScore int

func (ms submitMemoryScore) String() string {
	return fmt.Sprintf("%d KB", ms)
}

type submitResult struct {
	Status      resultStatus
	CodeLength  submitCodeLength
	TimeScore   time.Duration
	MemoryScore submitMemoryScore
	DetailUrl   string
}

type seleniumSubmitter struct{}

func newSeleniumSubmitter() *seleniumSubmitter {
	return &seleniumSubmitter{}
}

func (*seleniumSubmitter) Submit(taskUrl string, sourceCode []byte) (*submitResult, error) {
	driver := agouti.ChromeDriver()
	err := driver.Start()
	if err != nil {

		return nil, err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}

	// login
	err = page.Navigate("https://atcoder.jp/login")
	if err != nil {
		return nil, err
	}
	page.FindByID("username").SendKeys(os.Getenv("ATCODER_USERNAME"))
	page.FindByID("password").SendKeys(os.Getenv("ATCODER_PASSWORD"))
	page.FindByID("submit").Click()

	// task
	err = page.Navigate(taskUrl)
	if err != nil {
		return nil, err
	}
	page.FindByName("data.LanguageId").FindByXPath("//option[contains(text(), 'Go')]").Click()
	page.RunScript("$('.editor').data('editor').setValue(sourceCode)", map[string]interface{}{"sourceCode": string(sourceCode)}, nil)
	page.FindByID("submit").Click()

	time.Sleep(time.Second * 20)

	return nil, nil
}

func Submit(url, filename string) (*submitResult, error) {
	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return newSeleniumSubmitter().Submit(url, sourceCode)
}
