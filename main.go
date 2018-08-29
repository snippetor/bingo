package main

import (
	"os"
	"github.com/urfave/cli"
	"log"
	"sort"
	"github.com/snippetor/bingo/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "bingo"
	app.Usage = "A distributed, open source framework in Golang"
	app.Author = "snippetor@163.com"
	app.Version = "v1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init go project by bingo framework",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					command.Init(c.Args()[0])
				} else {
					command.Init("")
				}
				return nil
			},
		},
		{
			Name:      "run",
			Aliases:   []string{"r"},
			Usage:     "run the app by name and env",
			UsageText: "bingo run [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					command.Run(c.Args()[0], "")
				} else if c.NArg() == 2 {
					command.Run(c.Args()[0], c.Args()[1])
				}
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
