package internal

import "fmt"

func GetLatestIssue() string {

	name, err := GetProjectName()
	if err != nil {
		return ""
	}
	for _, item := range config.Issue.List {
		if item.ProjectName == name && item.Number != "" {
			return item.Number
		}
	}
	return ""
}

func IssueUpgrade() error {
	for project, issue := range config.Issue.List {
		fmt.Println(project, issue)
	}
	return nil
}

func IssueDescribe(project, branch string) {
	if len(project) == 0 || len(branch) == 0 {
		return
	}

	// project
	p, ok := config.Issue.Describes[project]
	if !ok {
		return
	}

	// 获取项目的分支描述map
	b, ok := p[branch]
	if !ok {
		return
	}
	fmt.Printf("[%s]: %s\n", b.Number, b.Describe)
}

func IssueRecord(project, branch, issue, describe string) {
	if len(project) == 0 || len(branch) == 0 || len(issue) == 0 {
		fmt.Println("参数错误")
		return
	}

	// project
	p, ok := config.Issue.Describes[project]
	if !ok {
		b := make(map[string]*BugDescribe)
		b[branch] = &BugDescribe{
			Number:   issue,
			Describe: describe,
		}
		config.Issue.Describes[project] = b
		Rewrite()
		return
	}

	// 获取项目的分支描述map
	b, ok := p[branch]
	if !ok {
		p[branch] = &BugDescribe{
			Number:   issue,
			Describe: describe,
		}
		Rewrite()
		return
	}
	b.Number = issue
	b.Describe = describe
	Rewrite()
}

func RecordIsuueHistory(project, issue string) {
	var skipAppend bool
	for _, item := range config.Issue.List {
		if item.ProjectName == project {
			item.Number = issue
			skipAppend = true
			break
		}
	}
	if !skipAppend {
		v := &RecordIssue{
			ProjectName: project,
			Number:      issue,
		}
		config.Issue.List = append(config.Issue.List, v)
	}
	Rewrite()
}
