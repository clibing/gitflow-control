package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	Feat     string = "feat"
	Fix      string = "fix"
	Docs     string = "docs"
	Style    string = "style"
	Refactor string = "refactor"
	Test     string = "test"
	Chore    string = "chore"
	Hotfix   string = "hotfix"
)

const commitMessageCheckPattern = `^(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*`

var CommitMessageType = map[string]string{
	Feat:     "新功能（feature）",
	Fix:      "修补bug",
	Docs:     "文档（documentation）",
	Style:    "格式（不影响代码运行的变动）",
	Refactor: "重构（即不是新增功能，也不是修改bug的代码变动）",
	Test:     "增加测试",
	Chore:    "构建过程或辅助工具的变动",
	Hotfix:   "紧急修复线上bug",
}

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

type CommitMessage struct {
	Issue   string              // 问题
	Type    string              // 已经初始化
	Scope   string              // 本地提交影响的范围，例如:数据层、控制层、视图层等等，视项目不同而不同。
	Subject string              // 不超过50个字符
	Body    string              // 具体的描述信息 建议72个字符
	Footer  CommitMessageFooter // 允许支持关闭issue
	SOB     string              // name, email
}

type CommitMessageFooter struct {
	Value  string   // 如果当前代码与上一个版本不兼容，则 Footer 部分以BREAKING CHANGE开头，后面是对变动的描述、以及变动理由和迁移方法。
	Closes []string // 将要关闭的issue 如果存在值，则将Value中存在的Closes去掉，在Value后面追加 Closes #Issue[0],#Issue[1]....
}

type CommentMessageType struct {
	Name string
	Note string
}

func Commit(msg CommitMessage) error {
	if err := HasStagedFiles(); err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "gitflow-commit")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	/**
		 * 参考: https://www.ruanyifeng.com/blog/2016/01/commit_message_change_log.html
		 *
		 * 这个是标准规范
		 * <type>(<scope>): <subject>
	     * // 空一行
	     * <body>
	     * // 空一行
	     * <footer>
		 *
		 * 支持YX讲issue放在头尾
	     * [issue]<type>(<scope>): <subject>
	     * // 空一行
	     * <body>
	     * // 空一行
	     * <footer>
		 *
	*/
	footer := msg.Footer

	footer_message := ""
	if len(footer.Value) > 0 {
		footer_message = fmt.Sprintf("BREAKING CHANGE: %s", footer.Value)
	}

	if len(footer.Closes) > 0 {
		issues := make([]string, len(footer.Closes))
		for _, issue := range footer.Closes {
			if strings.HasPrefix(issue, "#") {
				issues = append(issues, issue)
			} else {
				issues = append(issues, fmt.Sprintf("#%s", issue))
			}
		}
		footer_message = fmt.Sprintf("%s\n\nClosesCloses: %s", footer_message, strings.Join(issues, ","))
	}

	_, err = fmt.Fprintf(f, "[%s]%s(%s): %s\n\n%s\n\n%s\n\n%s\n", msg.Issue, msg.Type, msg.Scope, msg.Subject, msg.Body, footer_message, msg.SOB)
	if err != nil {
		return err
	}

	_, err = ExecGit("commit", "-F", f.Name())
	return err
}

func CheckCommitMessage(message string) error {
	// 增加 commit-msg hook时使用
	return nil
}
