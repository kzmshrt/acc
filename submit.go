package acc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	// driver
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	err := driver.Start()
	if err != nil {

		return nil, err
	}
	defer driver.Stop()

	// page
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

	// submissions
	time.Sleep(2 * time.Second)
	newSubmission := page.Find("table").Find("tbody").First("tr")
	cols := newSubmission.All("td")
	ncol, err := cols.Count()
	if err != nil {
		return nil, err
	}
	lastElem := cols.At(ncol - 1)
	detailUrl, err := lastElem.Find("a").Attribute("href")
	if err != nil {
		return nil, err
	}

	// submission detail
	err = page.Navigate(detailUrl)
	if err != nil {
		return nil, err
	}
	for isLabelDefault(page) {
		time.Sleep(500 * time.Millisecond)
	}

	tables := page.AllByClass("table-bordered")

	codeLength := tables.At(0).Find("tbody").All("tr").At(5)
	codeLengthTitle, _ := codeLength.Find("th").Text()
	codeLengthTitle = strings.TrimSpace(codeLengthTitle)
	codeLengthValue, _ := codeLength.Find("td").Text()
	codeLengthValue = strings.TrimSpace(codeLengthValue)

	resultStatus := tables.At(0).Find("tbody").All("tr").At(6)
	resultStatusTitle, _ := resultStatus.Find("th").Text()
	resultStatusTitle = strings.TrimSpace(resultStatusTitle)
	resultStatusValue, _ := resultStatus.Find("td").Text()
	resultStatusValue = strings.TrimSpace(resultStatusValue)

	timeScore := tables.At(0).Find("tbody").All("tr").At(7)
	timeScoreTitle, _ := timeScore.Find("th").Text()
	timeScoreTitle = strings.TrimSpace(timeScoreTitle)
	timeScoreValue, _ := timeScore.Find("td").Text()
	timeScoreValue = strings.TrimSpace(timeScoreValue)

	memoryScore := tables.At(0).Find("tbody").All("tr").At(8)
	memoryScoreTitle, _ := memoryScore.Find("th").Text()
	memoryScoreTitle = strings.TrimSpace(memoryScoreTitle)
	memoryScoreValue, _ := memoryScore.Find("td").Text()
	memoryScoreValue = strings.TrimSpace(memoryScoreValue)

	fmt.Println(resultStatusTitle, ":", resultStatusValue)
	fmt.Println(timeScoreTitle, ":", timeScoreValue)
	fmt.Println(memoryScoreTitle, ":", memoryScoreValue)
	fmt.Println(codeLengthTitle, ":", codeLengthValue)

	return nil, nil
}

func isLabelDefault(page *agouti.Page) bool {
	classAttribute, _ := page.FindByID("judge-status").Find("span").Attribute("class")
	return strings.Contains(classAttribute, "label-default")
}

func Submit(url, filename string) (*submitResult, error) {
	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return newSeleniumSubmitter().Submit(url, sourceCode)
}
