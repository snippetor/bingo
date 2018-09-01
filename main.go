package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/snippetor/bingo/command"
	"sort"
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
			Usage:   "init go project with bingo framework",
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
			Name:      "debug",
			Aliases:   []string{"r"},
			Usage:     "debug the app with name and env",
			UsageText: "bingo debug [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					command.Debug(c.Args()[0], "")
				} else if c.NArg() == 2 {
					command.Debug(c.Args()[0], c.Args()[1])
				} else {
					cli.ShowCommandHelp(c, "debug")
				}
				return nil
			},
		},
		{
			Name:      "new",
			Aliases:   []string{"n"},
			Usage:     "new a app config",
			UsageText: "bingo new [app package] [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 2 {
					command.New(c.Args()[0], c.Args()[1], "")
				} else if c.NArg() == 3 {
					command.New(c.Args()[0], c.Args()[1], c.Args()[2])
				} else {
					cli.ShowCommandHelp(c, "new")
				}
				return nil
			},
		},
		{
			Name:      "build",
			Aliases:   []string{"b"},
			Usage:     "build bingo app, if app name is * or app name and env is empty, build all app.",
			UsageText: "bingo build [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					command.Build("*", "")
				} else if c.NArg() == 1 {
					command.Build(c.Args()[0], "")
				} else if c.NArg() == 2 {
					command.Build(c.Args()[0], c.Args()[1])
				} else {
					cli.ShowCommandHelp(c, "build")
				}
				return nil
			},
		},
		{
			Name:      "publish",
			Aliases:   []string{"b"},
			Usage:     "publish bingo app, if app name is * or app name and env is empty, publish all app in one package.",
			UsageText: "bingo publish [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					command.Publish("*", "")
				} else if c.NArg() == 1 {
					command.Publish(c.Args()[0], "")
				} else if c.NArg() == 2 {
					command.Publish(c.Args()[0], c.Args()[1])
				} else {
					cli.ShowCommandHelp(c, "publish")
				}
				return nil
			},
		},
		{
			Name:      "run",
			Aliases:   []string{"r"},
			Usage:     "run the app's executable file with name and env",
			UsageText: "bingo run [app name] [env] [options]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "b",
					Usage: "running the app in background",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					command.Run(c.Args()[0], "", c.Bool("b"))
				} else if c.NArg() == 2 {
					command.Run(c.Args()[0], c.Args()[1], c.Bool("b"))
				} else {
					cli.ShowCommandHelp(c, "run")
				}
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
