package main

import (
	"errors"

	"github.com/kzmshrt/acc"
	"github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New(c.Command.UsageText)
	}

	filename := c.Args().Get(0)
	url := c.Args().Get(1)

	err := acc.Test(filename, url)
	if err != nil {
		return err
	}

	return nil
}
