package internal

import "fmt"

func BranchRecord(project, branch, title string) (value string) {
	var exist bool
	for _, item := range config.Record {
		if item.Project == project {
			if v, ok := item.Describe[branch]; ok {
				value = v
			}
			// 相同 直接返回
			if item.Describe[branch] == title {
				return
			}
			exist = true
			item.Describe[branch] = title
			break
		}
	}

	if !exist {
		config.Record = append(config.Record, &Record{
			Project: project,
			Describe: map[string]string{
				branch: title,
			},
		})
	}
	Rewrite()
	return
}

func GetBranchRecord(project, branch string) (value string) {
	for _, item := range config.Record {
		if item.Project == project {
			if v, ok := item.Describe[branch]; ok {
				value = v
				return
			}
		}
	}
	value = ""
	return
}

func ShowBranchRecord(project, title string) (err error) {
	fmt.Printf("%s\n\n", project)
	for _, item := range config.Record {
		if item.Project == project {
			for k, v := range item.Describe {
				fmt.Printf("%10.24s: \"%s\"\n", k, v)
			}
			return nil
		}
	}
	return nil
}
