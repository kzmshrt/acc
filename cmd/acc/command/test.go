package command

import (
	"log"

	"github.com/kzmshrt/acc"
	"github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
	if c.NArg() < 2 {
		log.Fatal(c.Command.UsageText)
	}

	filename, url := c.Args().Get(0), c.Args().Get(1)

	err := acc.Test(url, filename)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
