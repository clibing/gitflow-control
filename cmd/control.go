package cmd

import (
	"fmt"
	"runtime"
	"strings"
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

var subApps = make([]*cli.App, 12)

func (m *Control) Init() {
	subApps[0] = m.NewBranchCliApp(internal.Feat)
	subApps[1] = m.NewBranchCliApp(internal.Fix)
	subApps[2] = m.NewBranchCliApp(internal.Docs)
	subApps[3] = m.NewBranchCliApp(internal.Style)
	subApps[4] = m.NewBranchCliApp(internal.Refactor)
	subApps[5] = m.NewBranchCliApp(internal.Test)
	subApps[6] = m.NewBranchCliApp(internal.Chore)
	subApps[7] = m.NewBranchCliApp("hotfix")
	subApps[8] = m.CommitApp()
	subApps[9] = m.CheckMessageApp()
	subApps[10] = m.IssueApp()
	subApps[11] = m.BranchRecord()
	m.Config = internal.GetConfig()
}

func (m *Control) DefaultCliApp() *cli.App {
	return &cli.App{
		Name: "gitflow-control",
		UsageText: `Git Flow Control

git ci: 自定义提交
git feat: 创建feat分支
git fix: 创建fit分支
git docs: 创建docs类分支
git style: 创建sytle的分支
git refactor: 创建refactory的分支
git test: 创建test分支
git chore: 创建chore分支
git hotfix: 创建hotfix分支
git issue: 记录最近一次的issue号
git record: 记录当前分支描述信息， 主要用于描述当前分支业务类型`,
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

func (m *Control) NewBranchCliApp(ct string) *cli.App {
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

func (m *Control) GetSubCliApp(binName string) *cli.App {
	n := binName
	if runtime.GOOS == "windows" {
		n = strings.ReplaceAll(binName, ".exe", "")
	}
	for _, app := range subApps {
		if app != nil && n == app.Name {
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

func (m *Control) IssueApp() *cli.App {
	return &cli.App{
		Name:                 "git-issue",
		Usage:                "Git Issue",
		UsageText:            "git issue --bug issue-number",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			issue := c.String("bug")
			project := c.String("project")
			if len(project) == 0 {
				project, _ = internal.GetProjectName()
			}
			internal.RecordIsuueHistory(project, issue)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bug",
				Aliases: []string{"b"},
				Usage:   "issue number",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "project name",
				Value:   "",
			},
		},
	}
}

// 分支  记录器
func (m *Control) BranchRecord() *cli.App {
	return &cli.App{
		Name:                 "git-record",
		Usage:                "Git Record",
		UsageText:            "git record --title \"当前分支描述信息\"",
		Version:              fmt.Sprintf("%s %s %s", m.Version, m.BuildDate, m.BuildCommit),
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			project, _ := internal.GetProjectName()
			branch, _ := internal.CurrentBranch()
			title := c.String("title")
			if len(title) == 0 {
				title = internal.GetBranchRecord(project, branch)
				if len(title) > 0 {
					fmt.Printf("%s\n", title)
				}
				return nil
			}

			internal.BranchRecord(project, branch, title)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "title description",
				Value:   "",
			},
		},
	}
}
