package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/clibing/gitflow-control/cmd"
	"github.com/mattn/go-runewidth"
)

var (
	version     string
	buildDate   string
	buildCommit string
)

func main0() {

	control := &cmd.Control{
		Version:     version,
		BuildDate:   buildDate,
		BuildCommit: buildCommit,
	}
	control.Init()

	err := control.CheckMessageApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main2() {

	control := &cmd.Control{
		Version:     version,
		BuildDate:   buildDate,
		BuildCommit: buildCommit,
	}
	control.Init()

	err := control.CommitApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main3() {
	control := &cmd.Control{
		Version:     version,
		BuildDate:   buildDate,
		BuildCommit: buildCommit,
	}
	control.Init()

	err := control.NewBranchCliApp("hotfix").Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	control := &cmd.Control{
		Version:     version,
		BuildDate:   buildDate,
		BuildCommit: buildCommit,
	}
	control.Init()

	app := control.DefaultCliApp()
	// 可执行文件的路径
	bin, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 获取当前命令的名字
	binName := filepath.Base(bin)
	// 检查是否存在扩展的命令
	result := control.GetSubCliApp(binName)
	if result != nil {
		app = result
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// See also: https://github.com/charmbracelet/lipgloss/issues/40#issuecomment-891167509
func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}
