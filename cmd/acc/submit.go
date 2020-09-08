package main

import (
	"errors"
	"fmt"

	"github.com/kzmshrt/acc/atcoder"
	"github.com/urfave/cli/v2"
)

func Submit(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New(c.Command.UsageText)
	}

	filename := c.Args().Get(0)
	url := c.Args().Get(1)

	// client
	client, err := atcoder.NewRESTClient()
	if err != nil {
		return fmt.Errorf("client initialization failed: %v", err)
	}
	// login
	if err := client.Authenticate(); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	// submit
	submission, err := client.SubmitFile(filename, url)
	if err != nil {
		return fmt.Errorf("submission failed: %v", err)
	}

	// submission
	printSubmission(submission)

	return nil
}

func printSubmission(submission *atcoder.Submission) {
	fmt.Println("========================================")
	fmt.Printf("Status:       %s\n", submission.Status)
	fmt.Printf("Time Score:   %d [ms]\n", submission.TimeScore)
	fmt.Printf("Memory Score: %d [KB]\n", submission.MemoryScore)
	fmt.Printf("Code Length:  %d [Byte]\n", submission.CodeLength)
	fmt.Printf("Detail URL:   %s\n", submission.DetailUrl)
	fmt.Println("========================================")
}
