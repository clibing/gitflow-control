package cmd

import (
	"fmt"
	"time"

	"github.com/clibing/gitflow-control/internal"
	"github.com/urfave/cli/v2"
)

type Control struct {
	Version string
	BuildDate string
	BuildCommit string
}

var mainApp *cli.App


func (m *Control) Init(){
	mainApp = &cli.App{
		Name:                 "gitflow-control",
		Usage:                "Git Flow Control",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			m.Install(),
			m.Uninstall(),
		},
	}
}

func (m *Control) Install() *cli.Command{
	return &cli.Command{
		Name:    "install",
		Aliases: []string{"git-install"},
		Usage:   "Install gitflow-control",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return internal.Install(c.String("dir"), c.Bool("hook"))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "Install dir",
				Value:   "/usr/local/bin",
			},
			&cli.BoolFlag{
				Name:  "hook",
				Usage: "Install Commit Message hook",
				Value: false,
			},
		},
	}
}

func (m *Control) Uninstall() *cli.Command {
	return nil
}