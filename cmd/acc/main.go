package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	Name = "acc"
)

func main() {
	app := &cli.App{
		Name:  Name,
		Usage: "AtCoder Client",
		Commands: []*cli.Command{
			{
				Name:   "submit",
				Usage:  "submit answer",
				Action: Submit,
			},
			{
				Name:   "test",
				Usage:  "test answer",
				Action: Test,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
