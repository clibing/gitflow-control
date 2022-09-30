package internal

import (
	"errors"
	"fmt"
)

func Version() (string, error) {
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
	if all && remote {
		return ExecGit("branch", "-a", "-r")
	}
	if all {
		return ExecGit("branch", "-a")
	}
	if remote {
		return ExecGit("branch", "-r")
	}
	return ExecGit("branch")
}

func CurrentBranch() (string, error) {
	err := GitRepo()
	if err != nil {
		return "", err
	}
	return ExecGit("symbolic-ref", "--short", "HEAD")
}

func Author() (string, string, error) {
	name := ""
	email := ""

	if value, err := ExecGit("config", "user.name"); err == nil {
		name = value
	}

	if value, err := ExecGit("config", "user.email"); err == nil {
		email = value
	}

	return name, email, nil
}

func SignedOffBy() (string, error) {
	name, email, err := Author()

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Signed-off-by: %s <%s>", name, email), nil
}

func Push() (string, error) {
	err := GitRepo()
	if err != nil {
		return "", err
	}

	branch, err := CurrentBranch()
	if err != nil {
		return "", err
	}

	return ExecGit("push", "origin", branch)
}

func HasStagedFiles() error {
	msg, err := ExecGit("diff", "--cached", "--name-only")
	if err != nil {
		return err
	}
	if msg == "" {
		return errors.New("当前暂存区没有文件，执行`git add`增加文件后再次提交")
	}
	return nil
}