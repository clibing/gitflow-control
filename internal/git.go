package internal

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
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

const commitMessageCheckPatternV1 = `^(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*`
const commitMessageCheckPatternV2 = `^\%s[a-zA-Z]+\-[0-9]+\%s[\n\r]+(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*`

const commitMessageCheckFailedMsgV1 = `
╭────────────────────────────────────────────────────────────────────────────────────────╮
│ ✗ The commit message is not standardized.                                              │
│ ✗ It must match the regular expression:                                                │
│                                                                                        │
│ ^(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*     │
╰────────────────────────────────────────────────────────────────────────────────────────╯`

const commitMessageCheckFailedMsgV2 = `
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ ✗ The commit message is not standardized.                                                                       │
│ ✗ It must match the regular expression:                                                                         │
│                                                                                                                 │
│ ^\%s[a-zA-Z]+\-[0-9]+\%s[\n\r]+(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*  │
│                                                                                                                 │
│ example:                                                                                                        │
│ [BACKEND-001]                                                                                                   │
│                                                                                                                 │
│ chore(pom): add pom dep version                                                                                 │
│                                                                                                                 │
│ add pom dep version                                                                                             │
│                                                                                                                 │
│ Signed-off-by: clibing <wmsjhappy@gmail.com>                                                                    │
│                                                                                                                 │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
注意：issue的格式为[英文字母+引文短接线+数字]`

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

func CheckRepo() error {
	// 显示工作区根目录
	_, err := ExecGit("rev-parse", "--show-toplevel")
	// if err != nil {
	// 	fmt.Println("当前目录不是Git管理的项目，请检查.", err.Error())
	// 	os.Exit(1)
	// }
	return err
}
func Switch(name string) (string, error) {
	err := CheckRepo()
	if err != nil {
		return "", err
	}
	return ExecGit("switch", "-c", name)
}

func Branch(all, remote bool) (string, error) {
	err := CheckRepo()
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
	err := CheckRepo()
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
	err := CheckRepo()
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

func GetOriginUrl() (string, error) {
	err := CheckRepo()
	if err != nil {
		return "", err
	}

	msg, err := ExecGit("remote", "get-url", "--push", "origin")
	if err != nil {
		return "", err
	}
	return msg, nil
}

type CommitMessage struct {
	Type    string // 已经初始化
	Scope   string // 本地提交影响的范围，例如:数据层、控制层、视图层等等，视项目不同而不同。
	Subject string // 不超过50个字符
	Body    string // 具体的描述信息 建议72个字符
	Footer  string // 允许支持关闭issue
	SOB     string // name, email
}

type CommentMessageType struct {
	Name string
	Note string
}

func Commit(msg CommitMessage, config *Config) error {
	if err := HasStagedFiles(); err != nil {
		return err
	}

	f, err := os.CreateTemp("", "gitflow-commit")
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

	if RequiredFooter() {
		_, err = fmt.Fprintf(f, "%s%s%s\n\n%s(%s): %s\n\n%s\n\n%s\n\n%s\n", config.Issue.LeftMarker, msg.Footer, config.Issue.RightMarker, msg.Type, msg.Scope, msg.Subject, msg.Body, msg.Footer, msg.SOB)
	} else {
		_, err = fmt.Fprintf(f, "%s(%s): %s\n\n%s\n\n%s\n\n%s\n", msg.Type, msg.Scope, msg.Subject, msg.Body, msg.Footer, msg.SOB)
	}
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	_, err = ExecGit("commit", "-F", f.Name())
	return err
}

func CheckCommitMessage(message string, config *Config) error {
	rg := commitMessageCheckPatternV1
	if RequiredFooter() {
		rg = fmt.Sprintf(commitMessageCheckPatternV2, config.Issue.LeftMarker, config.Issue.RightMarker)
	}
	// 增加 commit-msg hook时使用
	reg := regexp.MustCompile(rg)
	bs, err := os.ReadFile(message)
	if err != nil {
		return err
	}

	msgs := reg.FindStringSubmatch(string(bs))
	if RequiredFooter() {
		if len(msgs) != 4 {
			return fmt.Errorf(commitMessageCheckFailedMsgV2, config.Issue.LeftMarker, config.Issue.RightMarker)
		}
	} else {
		if len(msgs) != 4 {
			return fmt.Errorf(commitMessageCheckFailedMsgV1)
		}
	}

	return nil
}

func GetProjectName() (string, error) {
	url, err := GetOriginUrl()
	if err != nil {
		return "", err
	}
	var reg = regexp.MustCompile(`(?m)\/([a-zA-Z_\-0-9]+)\.git`)

	msgs := reg.FindStringSubmatch(url)

	if len(msgs) != 2 {
		return "", fmt.Errorf("current git push url: %s, not found name", url)
	}
	return msgs[1], nil
}

func NameAndEmail(name, email string, global bool) {
	_, err := CurrentBranch()
	if err != nil {
		fmt.Printf("设置错误：%s\n", err.Error())
		return
	}
	// 读模式
	if len(name) == 0 || len(email) == 0 {
		if global {
			n, e := ExecGit("config", "--global", "user.name")
			if e == nil {
				fmt.Println("name:  ", n)
			}
			p, e := ExecGit("config", "--global", "user.email")
			if e == nil {
				fmt.Println("email: ", p)
			}
		} else {
			n, e := ExecGit("config", "user.name")
			if e == nil {
				fmt.Println("name:  ", n)
			}
			p, e := ExecGit("config", "user.email")
			if e == nil {
				fmt.Println("email: ", p)
			}
		}
		return
	}

	if global {
		ExecGit("config", "--global", "user.name", name)
		ExecGit("config", "--global", "user.email", email)
		return
	} else {
		ExecGit("config", "user.name", name)
		ExecGit("config", "user.email", email)
		return
	}
}
