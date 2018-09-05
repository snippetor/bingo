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
			Name:      "pack",
			Aliases:   []string{"b"},
			Usage:     "pack bingo app, if app name is * or app name and env is empty, pack all app in one package.",
			UsageText: "bingo publish [app name] [env]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					command.Pack("*", "")
				} else if c.NArg() == 1 {
					command.Pack(c.Args()[0], "")
				} else if c.NArg() == 2 {
					command.Pack(c.Args()[0], c.Args()[1])
				} else {
					cli.ShowCommandHelp(c, "pack")
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
		{
			Name:      "pb",
			Usage:     "gen protobuf go file using gogofaster(https://github.com/gogo/protobuf)",
			UsageText: "bingo pb [.proto dir]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "d",
					Usage: "output directory",
				},
				cli.BoolFlag{
					Name:  "s",
					Usage: "output file into single directory",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					command.Pb(c.Args()[0], c.String("d"), c.Bool("s"))
				} else {
					cli.ShowCommandHelp(c, "pb")
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
