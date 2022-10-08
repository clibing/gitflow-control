package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/clibing/gitflow-control/internal"
	"github.com/urfave/cli/v2"
)

type Control struct {
	Version     string
	BuildDate   string
	BuildCommit string
	Config      *internal.Config // 配置信息
}

var subApps = make([]*cli.App, 10)

func (m *Control) Init() {
	subApps[0] = m.newBranchCliApp(internal.Feat)
	subApps[1] = m.newBranchCliApp(internal.Fix)
	subApps[2] = m.newBranchCliApp(internal.Docs)
	subApps[3] = m.newBranchCliApp(internal.Style)
	subApps[4] = m.newBranchCliApp(internal.Refactor)
	subApps[5] = m.newBranchCliApp(internal.Test)
	subApps[6] = m.newBranchCliApp(internal.Chore)
	subApps[7] = m.newBranchCliApp("hotfix")
	subApps[8] = m.CommitApp()
	subApps[9] = m.CheckMessageApp()
	m.Config = internal.GetConfig()
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

func (m *Control) newBranchCliApp(ct string) *cli.App {
	return &cli.App{
		Name:                 "git-" + ct,
		Usage:                fmt.Sprintf("Create %s branch", ct),
		UsageText:            fmt.Sprintf("git %s BRANCH", ct),
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			branchName := fmt.Sprintf("%s/%s", ct, c.Args().First())
			_, err := internal.Switch(branchName)
			return err
		},
	}
}

func (m *Control) GetBranchCliApp(binName string) *cli.App {
	for _, app := range subApps {
		if app != nil && binName == app.Name {
			return app
		}
	}
	return nil
}

func (m *Control) CommitApp() *cli.App {
	return &cli.App{
		Name:                 "git-ci",
		Usage:                "Interactive commit",
		UsageText:            "git ci",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}

			model := internal.CommitModel{
				Views: []tea.Model{
					internal.NewSelectorModel(),
					internal.NewInputsModel(),
					internal.NewSubmitModel(),
					internal.NewErrorModel(),
				},
			}

			return tea.NewProgram(&model).Start()
		},
	}

}

func (m *Control) CheckMessageApp() *cli.App {
	return &cli.App{
		Name:                 "commit-msg",
		Usage:                "Commit message hook",
		UsageText:            "commit-msg FILE",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			return internal.CheckCommitMessage(c.Args().First(), m.Config)
		},
	}
}

// TODO issue的管理
func (m *Control) IssueApp() *cli.App {
	return &cli.App{
		Name:                 "git-issue",
		Usage:                "git issue",
		UsageText:            "git issue command #ISSUE",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return internal.CheckCommitMessage(c.Args().First(), m.Config)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "-c append|remove|clean",
				Value:   "append",
			},
		},
	}
}
