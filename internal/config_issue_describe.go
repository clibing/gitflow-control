package internal

import (
	"fmt"
	"time"
)

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

func GetIssueDescribe(project, branch string) (issue string) {
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
	issue = b.Number
	return
}

func IssueDescribe(project, branch string, time, all bool) {
	if len(project) == 0 || len(branch) == 0 {
		return
	}

	// project
	p, ok := config.Issue.Describes[project]
	if !ok {
		return
	}
	if all {
		for name, describe := range p {
			if time {
				fmt.Printf("[branch: %s, issue: %s]: \"%s\".(%s)\n", name, describe.Number, describe.Describe, describe.Time)
			} else {
				fmt.Printf("[branch: %s, issue: %s]: \"%s\"\n", name, describe.Number, describe.Describe)
			}
		}
		return
	}

	// 获取项目的分支描述map
	b, ok := p[branch]
	if !ok {
		return
	}
	if time {
		fmt.Printf("[%s]: \"%s\".(%s)\n", b.Number, b.Describe, b.Time)
	} else {
		fmt.Printf("[%s]: \"%s\"\n", b.Number, b.Describe)
	}
}

func IssueRecord(project, branch, issue, describe string) {
	if len(project) == 0 || len(branch) == 0 || len(issue) == 0 {
		fmt.Println("参数错误")
		return
	}

	// current time
	v := time.Now().Format("2006-01-02 15:04:05")
	// project
	p, ok := config.Issue.Describes[project]
	if !ok {
		b := make(map[string]*BugDescribe)
		b[branch] = &BugDescribe{
			Number:   issue,
			Describe: describe,
			Time:     v,
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
			Time:     v,
		}
		Rewrite()
		return
	}
	b.Number = issue
	b.Describe = describe
	b.Time = v
	Rewrite()
}

// 废弃
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
