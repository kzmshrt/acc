package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kzmshrt/acc/atcoder"
	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

func printTestResults(results []*atcoder.TestResult) {
	nok, nwa, nre, nce := 0, 0, 0, 0

	for _, result := range results {
		switch result.TestStatus {
		case atcoder.TestStatusOK:
			nok++
			fmt.Println(aurora.Bold("=================================================="))
			fmt.Println(aurora.Blue("Test Status:"))
			fmt.Println(aurora.Green(result.TestStatus))
			fmt.Println(aurora.Blue("Input:"))
			fmt.Println(result.Input)
			fmt.Println(aurora.Blue("Expected Output:"))
			fmt.Println(result.ExpectedOutput)
			fmt.Println(aurora.Blue("Actual Output:"))
			fmt.Println(result.ActualOutput)
			fmt.Println(aurora.Bold("=================================================="))
		case atcoder.TestStatusWA:
			nwa++
			fmt.Println(aurora.Bold("=================================================="))
			fmt.Println(aurora.Blue("Test Status:"))
			fmt.Println(aurora.Red(result.TestStatus))
			fmt.Println(aurora.Blue("Input:"))
			fmt.Println(result.Input)
			fmt.Println(aurora.Blue("Expected Output:"))
			fmt.Println(result.ExpectedOutput)
			fmt.Println(aurora.Blue("Actual Output:"))
			fmt.Println(result.ActualOutput)
			fmt.Println(aurora.Bold("=================================================="))
		case atcoder.TestStatusRE:
			nre++
			fmt.Println(aurora.Bold("=================================================="))
			fmt.Println(aurora.Blue("Test Status:"))
			fmt.Println(aurora.BrightYellow(result.TestStatus))
			fmt.Println(aurora.Blue("Input:"))
			fmt.Println(result.Input)
			fmt.Println(aurora.Blue("Expected Output:"))
			fmt.Println(result.ExpectedOutput)
			fmt.Println(aurora.Blue("Actual Output:"))
			fmt.Println(result.ActualOutput)
			fmt.Println(aurora.Bold("=================================================="))
		case atcoder.TestStatusCE:
			nce++
			fmt.Println(aurora.Bold("=================================================="))
			fmt.Println(aurora.Blue("Test Status:"))
			fmt.Println(aurora.BrightYellow(result.TestStatus))
			fmt.Println(aurora.Blue("Error:"))
			fmt.Println(aurora.Red(result.ActualOutput))
			fmt.Println(aurora.Bold("=================================================="))
		}
	}

	fmt.Println(aurora.Bold("=================================================="))
	fmt.Println("OK:", nok)
	fmt.Println("WA:", nwa)
	fmt.Println("RE:", nre)
	fmt.Println("CE:", nce)
	fmt.Println(aurora.Bold("=================================================="))
}

func Test(c *cli.Context) error {
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

	// get tasks
	task, err := atcoder.ParseTaskURL(taskURL)
	if err != nil {
		return fmt.Errorf("failed parsing task URL: %s: %v", taskURL, err)
	}
	testCases, err := client.GetTestCases(task)
	if err != nil {
		return fmt.Errorf("failed getting test cases from task page: %s: %v", taskURL, err)
	}

	log.Printf("fetching test cases completed.")
	log.Printf("running test...")

	// run test
	runner, _ := atcoder.NewRunner()
	testResults, err := runner.Run(filename, testCases)
	if err != nil {
		return fmt.Errorf("error occurred while running test: %v", err)
	}

	printTestResults(testResults)

	return nil
}
