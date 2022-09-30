package internal

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	message, _ := ExecGit("version")
	fmt.Println(message)

	rev_parse := []string{
		"--git-dir",       // 显示版本库 .git 目录所在的位置
		"--show-toplevel", // 显示工作区根目录
		"--show-prefix",   // 所在目录相对于工作区根目录的相对目录
		"--show-cdup",     // 显示从当前目录后退到工作区的根的深度
	}
	command := "rev-parse"
	for _, v := range rev_parse {
		message, _ = ExecGit(command, v)
		fmt.Printf("git %s %s, result: %s\n", command, v, message)
	}

	message, _ = Version()
	fmt.Println(message)

	message, _ = Branch(false, false)
	fmt.Println(message)
}
