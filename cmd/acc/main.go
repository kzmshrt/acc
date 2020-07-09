package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kzmshrt/acc/cmd/acc/command"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "acc"
)

func main() {
	acc := &cli.App{
		Name:  AppName,
		Usage: "AtCoder Client",
		Commands: []*cli.Command{
			{
				Name:      "test",
				Usage:     "test code with samples on question page",
				UsageText: fmt.Sprintf("test <filename> <url>"),
				Action:    command.Test,
			},
			{
				Name:      "submit",
				Usage:     "submit code",
				UsageText: fmt.Sprintf("submit <filename> <url>"),
				Action:    command.Submit,
			},
		},
	}
	err := acc.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
