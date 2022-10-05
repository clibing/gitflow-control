package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"time"

	"github.com/clibing/gitflow-control/internal"
	"github.com/urfave/cli/v2"
)

type Control struct {
	Version     string
	BuildDate   string
	BuildCommit string
}

var subApps = make([]*cli.App, 10)

func (c *Control) Init() {
	subApps[0] = c.newBranchCliApp(internal.Feat)
	subApps[1] = c.newBranchCliApp(internal.Fix)
	subApps[2] = c.newBranchCliApp(internal.Docs)
	subApps[3] = c.newBranchCliApp(internal.Style)
	subApps[4] = c.newBranchCliApp(internal.Refactor)
	subApps[5] = c.newBranchCliApp(internal.Test)
	subApps[6] = c.newBranchCliApp(internal.Chore)
	subApps[7] = c.newBranchCliApp("hotfix")
	subApps[8] = c.CommitApp()
	subApps[9] = checkMessageApp()
}

func (m *Control) DefaultCliApp() *cli.App {
	return &cli.App{
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

func (m *Control) Install() *cli.Command {
	return &cli.Command{
		Name:    "install",
		Aliases: []string{"git-install"},
		Usage:   "Install gitflow-control",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return internal.Install(c.String("path"))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Install path",
				Value:   "/usr/local/bin",
			},
		},
	}
}

func (m *Control) Uninstall() *cli.Command {
	return &cli.Command{
		Name:    "uninstall",
		Aliases: []string{"git-uninstall"},
		Usage:   "Uninstall gitflow-control",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return internal.UnInstall(c.String("path"))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "UnInstall path",
				Value:   "/usr/local/bin",
			},
		},
	}
}

func (c *Control) newBranchCliApp(ct string) *cli.App {
	return &cli.App{
		Name:                 "git-" + string(ct),
		Usage:                fmt.Sprintf("Create %s branch", ct),
		UsageText:            fmt.Sprintf("git %s BRANCH", ct),
		Version:              fmt.Sprintf("%s %s %s", c.Version, c.BuildDate, c.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			branhName := fmt.Sprintf("%s/%s", ct, c.Args().First())
			_, err := internal.Switch(branhName)
			return err
		},
	}
}

func (c *Control) GetBranchCliApp(binName string) *cli.App {
	for _, app := range subApps {
		if app != nil && binName == app.Name {
			return app
		}
	}
	return nil
}

func (c *Control) CommitApp() *cli.App {
	return &cli.App{
		Name:                 "git-ci",
		Usage:                "Interactive commit",
		UsageText:            "git ci",
		Version:              fmt.Sprintf("%s %s %s", c.Version, c.BuildDate, c.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}

			m := internal.CommitModel{
				Views: []tea.Model{
					internal.NewSelectorModel(),
					internal.NewInputsModel(),
					internal.NewSubmitModel(),
					internal.NewErrorModel(),
				},
			}

			return tea.NewProgram(&m).Start()
		},
	}

}

func checkMessageApp() *cli.App {
	return nil
}
