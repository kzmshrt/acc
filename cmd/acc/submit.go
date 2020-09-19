package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kzmshrt/acc/atcoder"
	"github.com/urfave/cli/v2"
)

func Submit(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New(c.Command.UsageText)
	}

	filename := c.Args().Get(0)
	taskURL := c.Args().Get(1)

	// initialize client
	client, err := atcoder.NewRESTClient()
	if err != nil {
		return fmt.Errorf("client initialization failed: %v", err)
	}

	// login
	username := os.Getenv("ATCODER_USERNAME")
	password := os.Getenv("ATCODER_PASSWORD")
	if err := client.Authenticate(username, password); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	log.Printf("authentication completed.")

	// submit
	task, err := atcoder.ParseTaskURL(taskURL)
	if err != nil {
		return fmt.Errorf("failed parsing task URL: %s: %v", taskURL, err)
	}
	answer, err := atcoder.ParseAnswerFile(filename)
	if err != nil {
		return fmt.Errorf("failed reading file content: %s: %v", filename, err)
	}
	if err := client.Submit(task, answer); err != nil {
		return fmt.Errorf("failed submitting answer %v to task %v: %v", answer, task, err)
	}

	log.Printf("submission completed.")

	// wait judging
	submissions, err := client.ListSubmissions(task.ContestID)
	if err != nil {
		return fmt.Errorf("failed getting submissions: %v", err)
	}

	log.Printf("waiting judge...")

	for submissions[0].Judge == atcoder.JudgeWJ {
		time.Sleep(200 * time.Millisecond)

		submissions, err = client.ListSubmissions(task.ContestID)
		if err != nil {
			return fmt.Errorf("failed getting submissions while waiting judge: %v", err)
		}
	}

	log.Printf("judging completed.")

	printSubmission(submissions[0])

	return nil
}

func printSubmission(submission *atcoder.Submission) {
	fmt.Printf("================================================================================\n")
	fmt.Printf("Judge       : %s\n", submission.Judge)
	fmt.Printf("Time Score  : %d [ms]\n", submission.TimeScore)
	fmt.Printf("Memory Score: %d [KB]\n", submission.MemoryScore)
	fmt.Printf("Code Length : %d [Byte]\n", submission.CodeSize)
	fmt.Printf("Detail URL  : %s\n", submission.DetailURL)
	fmt.Printf("================================================================================\n")
}
