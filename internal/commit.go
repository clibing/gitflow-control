package internal

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

const (
	ViewSelectCommitType   = iota // 选择 commit的type类型
	ViewInputCommitMessage        // 输入提交所需信息
	ViewCommitMessage             // 调动git commit 提交信息
	ViewError                     // 错误
)

type CommitModel struct {
	Views  []tea.Model // 提交时所需要的 tea model
	Index int         // 当前展示的view
	Error error       // 错误信息
}

type Next struct {
	Error error // 错误信息
	Next  int   // CommitModel.View 的下标
}

func (m CommitModel) Init() tea.Cmd {
	return func() tea.Msg {
		// 检查当前目录是否为git工作目录
		err := CheckRepo()
		if err != nil {
			return Next{Next: ViewError, Error: err}
		}

		// 检查是否存在 暂存的文件
		err = HasStagedFiles()
		if err != nil {
			return Next{Next: ViewError, Error: err}
		}
		return nil
	}
}

func (m CommitModel) View() string {
	return m.Views[m.Index].View()
}

func (m CommitModel) ShowErrorMessage() tea.Msg {
	return m.Error
}

func (m CommitModel) Inputs() tea.Msg {
	return m.Views[ViewSelectCommitType].(SelectorModel).choice
}

func (m CommitModel) Commit() tea.Msg {
	sob, err := SignedOffBy()
	if err != nil {
		return Next{Error: err}
	}

	msg := CommitMessage{
		Type:    m.Views[ViewSelectCommitType].(SelectorModel).choice,
		Scope:   m.Views[ViewInputCommitMessage].(InputsModel).Inputs[0].Input.Value(),
		Subject: m.Views[ViewInputCommitMessage].(InputsModel).Inputs[1].Input.Value(),
		Body:    strings.Replace(m.Views[ViewInputCommitMessage].(InputsModel).Inputs[2].Input.Value(), `\\n `, "\n", -1),
		Footer:  m.Views[ViewInputCommitMessage].(InputsModel).Inputs[3].Input.Value(),
		SOB:     sob,
	}

	if msg.Body == "" {
		msg.Body = msg.Subject
	}

	return msg
}

func (m CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case Next: // If the view returns a done message, it means that the stage has been processed
		m.Error = msg.(Next).Error
		m.Index = msg.(Next).Next

		// some special views need to determine the state of the data to update
		switch m.Index {
		case ViewInputCommitMessage:
			// textinput.Blink: 光标
			// spinner.Tick 后台任务正在运行的 提示
			// 当前 model的方法 ，主要获取上一步selector的commit type类型
			return m, tea.Batch(textinput.Blink, spinner.Tick, m.Inputs)
		case ViewCommitMessage:
			return m, tea.Batch(spinner.Tick, m.Commit)
		case ViewError:
			return m, m.ShowErrorMessage
		default:
			return m, tea.Quit
		}
	default: // By default, the cmd returned by the view needs to be processed by itself
		var cmd tea.Cmd
		m.Views[m.Index], cmd = m.Views[m.Index].Update(msg)
		return m, cmd
	}
}
