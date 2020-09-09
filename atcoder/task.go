package atcoder

import (
	"net/url"
	"strings"
)

type Task struct {
	ContestID string
	TaskID    string
}

func ParseTaskURL(taskURL string) (*Task, error) {
	baseURL, err := url.Parse("https://atcoder.jp")
	if err != nil {
		return nil, err
	}
	u, err := baseURL.Parse(taskURL)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(u.Path, "/")
	return &Task{ContestID: parts[2], TaskID: parts[4]}, nil
}
