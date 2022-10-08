package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	committingStyle = lipgloss.NewStyle().
			Padding(1, 1, 1, 2)

	committingIsuueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF0", Dark: "#FFFDF0"}).
				Background(lipgloss.AdaptiveColor{Light: "#5B43FF", Dark: "#7652FF"})

	committingTypeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}).
				Background(lipgloss.AdaptiveColor{Light: "#5B44FF", Dark: "#7653FF"})

	committingScopeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}).
				Background(lipgloss.AdaptiveColor{Light: "#1FD314", Dark: "#2AD67F"})

	committingSubjectStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}).
				Background(lipgloss.AdaptiveColor{Light: "#E11C9C", Dark: "#EE6FF8"})

	committingBodyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#2AD67F"})

	committingFooterStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#2AD67F"})

	committingSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#2AD67F"})

	committingFailedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#D63B3A", Dark: "#D63B3A"})
)

/**
 * 提交
 */

type SubmitModel struct {
	Err     error
	Next    bool
	Msg     CommitMessage
	Spinner spinner.Model
}

func (m SubmitModel) Init() tea.Cmd {
	return nil
}
func (m SubmitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case CommitMessage:
		m.Msg = msg
		return m, func() tea.Msg {
			time.Sleep(time.Second)
			return Commit(msg, GetConfig())
		}
	case error:
		m.Next = true
		m.Err = msg
		return m, tea.Quit
	case nil:
		m.Next = true
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
}

func (m SubmitModel) View() string {
	issue := ""
	if RequiredFooter() {
		issue = committingIsuueStyle.Render(fmt.Sprintf("%s%s%s\n", GetConfig().Issue.LeftMarker, m.Msg.Footer, GetConfig().Issue.RightMarker))
	}
	header := committingTypeStyle.Render(m.Msg.Type) + committingScopeStyle.Render("("+m.Msg.Scope+")") + committingSubjectStyle.Render(": "+m.Msg.Subject) + "\n"
	body := committingBodyStyle.Render(m.Msg.Body)
	footer := committingFooterStyle.Render(m.Msg.Footer+"\n"+m.Msg.SOB) + "\n"

	msg := m.Spinner.View()
	if m.Next {
		if m.Err != nil {
			msg = committingFailedStyle.Render("( ●●● ) Commit Failed: \n" + m.Err.Error())
		} else {
			msg = committingSuccessStyle.Render("◉◉◉◉ Always code as if the guy who ends up maintaining your \n◉◉◉◉ code will be a violent psychopath who knows where you live...")
		}
	}
	return committingStyle.Render(lipgloss.JoinVertical(lipgloss.Left, issue, header, body, footer, msg))
}

func NewSubmitModel() SubmitModel {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{
			"(●    ) C",
			"( ●   ) Co",
			"(  ●  ) Com",
			"(   ● ) Comm",
			"(    ●) Commi",
			"(    ●) Commit",
			"(   ● ) Committ",
			"(  ●  ) Committi",
			"( ●   ) Committin",
			"(●    ) Committing",
			"( ●   ) Committing.",
			"(  ●  ) Committing..",
			"(   ● ) Committing...",
			"(    ●) Committing...",
			"(   ● ) Committing...",
			"(  ●  ) Committing...",
			"( ●   ) Committing...",
			"(●    ) Committing...",
		},
		FPS: time.Second / 15,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#19F896"}).Bold(true)
	return SubmitModel{Spinner: s}
	//return SubmitModel{}
}
