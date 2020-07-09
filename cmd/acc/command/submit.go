package command

import (
	"log"

	"github.com/kzmshrt/acc"
	"github.com/urfave/cli/v2"
)

func Submit(c *cli.Context) error {
	if c.NArg() < 2 {
		log.Println(c.Command.UsageText)
		c.Done()
	}

	filename, url := c.Args().Get(0), c.Args().Get(1)

	_, err := acc.Submit(url, filename)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
