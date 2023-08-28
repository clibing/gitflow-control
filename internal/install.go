package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

/**
 * 安装目录
 * 是否 配置 hook
 * 主要的动作:
 * 1. 清理symlinks, bin
 * 2. unset commit-msg
 * 3. create gitflow control home
 * 4. install current bin
 * 5. create symlinks
 * 6. set commit hooks
 */
func Install(path string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// bin home dir
	controlHome := filepath.Join(home, HomeDir)
	// bin
	controlBin := filepath.Join(path, BinName)
	// hook global config
	controlHooks := filepath.Join(controlHome, "hooks")

	// 1. remove bin home ...
	err = os.RemoveAll(controlHome)
	if err != nil {
		return err
	}

	for _, symlink := range gitCommandSymlinks(path) {
		// 获取 sym link链接
		if _, err := os.Lstat(symlink); err == nil {
			err = os.RemoveAll(symlink)
			if err != nil {
				return fmt.Errorf("删除符号链接失败, sym link: %s:%s", symlink, err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("获取当前符号链接失败, sym link: %s:%s", symlink, err)
		}
	}

	// 2. unset commit-msg
	cv, _ := ExecGit("config", "--global", "--get", "core.hooksPath")
	if len(cv) > 0 {
		_, err = ExecGit("config", "--global", "--unset", "core.hooksPath")
		if err != nil {
			return fmt.Errorf("取消Git的commit-msg配置失败: %s", err)
		}
	}

	// 3
	err = os.MkdirAll(controlHome, 0755)
	if err != nil {
		return fmt.Errorf("创建工作目录[%s]失败: %s", controlHome, err)
	}
	RecoverConfigFile()

	// 4 install bin
	binPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return fmt.Errorf("获取当前可执行文件[%s],失败: %s", os.Args[0], err)
	}

	currentFile, err := os.Open(binPath)
	if err != nil {
		return fmt.Errorf("读取当前可执行文件[%s]文件信息，失败: %s", os.Args[0], err)
	}
	defer func() { _ = currentFile.Close() }()

	installFile, err := os.OpenFile(filepath.Join(path, BinName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("创建可执行文件[%s]失败: %s", filepath.Join(controlHome, BinName), err)
	}
	defer func() { _ = installFile.Close() }()

	_, err = io.Copy(installFile, currentFile)
	if err != nil {
		return fmt.Errorf("拷贝文件[%s]失败: %s", filepath.Join(controlHome, BinName), err)
	}

	// 5 create sym link
	for _, symlink := range gitCommandSymlinks(path) {
		err = os.Symlink(controlBin, symlink)
		if err != nil {
			return fmt.Errorf("创建Symlink[%s]失败: %s", symlink, err)
		}
	}

	// set commit hook
	err = os.MkdirAll(controlHooks, 0755)
	if err != nil {
		return fmt.Errorf("创建目录[%s]失败: %s", controlHooks, err)
	}
	err = os.Symlink(controlBin, filepath.Join(controlHooks, "commit-msg"))
	if err != nil {
		return fmt.Errorf("创建Symlink[%s]失败: %s", filepath.Join(controlHooks, "commit-msg"), err)
	}
	_, _ = ExecGit("config", "--global", "core.hooksPath", controlHooks)

	fmt.Println("安装成功...")
	return nil
}

func UnInstall(path string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// bin home dir
	controlHome := filepath.Join(home, ".gitflow-control")
	// bin
	controlBin := filepath.Join(path, BinName)

	// 1. rmeove work home and remove syslink
	err = os.RemoveAll(controlHome)
	if err != nil {
		return err
	}
	for _, symlink := range gitCommandSymlinks(path) {
		// 获取 sym link链接
		if _, err := os.Lstat(symlink); err == nil {
			err = os.RemoveAll(symlink)
			if err != nil {
				return fmt.Errorf("删除符号链接失败, sym link: %s:%s", symlink, err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("获取当前符号链接失败, sym link: %s:%s", symlink, err)
		}
	}

	// 2. unset commit-msg
	_, err = ExecGit("config", "--global", "--unset", "core.hooksPath")
	if err != nil {
		return fmt.Errorf("取消Git的commit-msg配置失败: %s", err)
	}

	// 3 remove bin file
	err = os.Remove(controlBin)
	if err != nil {
		return fmt.Errorf("删除二进制异常[%s]失败: %s", controlBin, err)
	}

	fmt.Println("卸载成功...")
	return nil
}

// 扩展的git命令列表
func gitCommandSymlinks(path string) []string {
	return []string{
		filepath.Join(path, "git-ci"),         // 自定义提交
		filepath.Join(path, "git-feat"),       // 创建feat分支
		filepath.Join(path, "git-fix"),        // 创建fit分支
		filepath.Join(path, "git-docs"),       // 创建docs类分支
		filepath.Join(path, "git-style"),      // 创建sytle的分支
		filepath.Join(path, "git-refactor"),   // 创建refactory的分支
		filepath.Join(path, "git-test"),       // 创建test分支
		filepath.Join(path, "git-chore"),      // 创建chore分支
		filepath.Join(path, "git-hotfix"),     // 创建hotfix分支
		filepath.Join(path, "git-issue"),      // 记录最近一次的issue号
		filepath.Join(path, "git-record"),     // 记录当前分支描述信息， 主要用于描述当前分支业务类型
		filepath.Join(path, "git-name-email"), // 设置当前提交的账号和邮箱，支持全局设置
	}
}
