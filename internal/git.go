package internal

import (
)

func Version() (string, error){
	msg, err := ExecGit("version")
	if err != nil {
		return "", err
	}
	return msg, err
}

func GitRepo() error {
	// 显示工作区根目录
	_, err := ExecGit("rev-parse", "--show-toplevel")
	// if err != nil {
	// 	fmt.Println("当前目录不是Git管理的项目，请检查.", err.Error())
	// 	os.Exit(1)
	// }
	return err
}
func Switch(name string) (string, error) {
	err := GitRepo()
	if err != nil {
		return "", err
	}
	return ExecGit("switch", "c", name)
}

func Branch(all, remote bool) (string, error) {
	err := GitRepo()
	if err != nil {
		return "", err
	}
	var cmd [3]string
	cmd[0] = "branch"
	cmd[1] = ""
	cmd[2] = ""

	if all {
		cmd[1] = "-a"
	}
	if remote {
		cmd[2] = "-r"
	}
	return ExecGit(cmd[0], cmd[1], cmd[2])
}
