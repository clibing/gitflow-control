package internal

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type SelectorModel struct {
	list   list.Model
	choice string
}

type SelectorItem struct {
	CommitType string
	Title      string
}

type SelectorDelegate struct{}

func (d SelectorDelegate) Height() int                             { return 1 }
func (d SelectorDelegate) Spacing() int                            { return 0 }
func (d SelectorDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SelectorDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectorItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. (%s)%s", index+1, i.CommitType, i.Title)
	if index == m.Index() {
		_, _ = fmt.Fprintf(w, str)
	} else {
		_, _ = fmt.Fprintf(w, Gray.Render(str))
	}

}

func (si SelectorItem) FilterValue() string { return si.Title }

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = m.list.SelectedItem().(SelectorItem).CommitType
			return m, func() tea.Msg { return Done{Next: ViewInputCommitMessage} }

		default:
			if !m.list.SettingFilter() && (keypress == "q" || keypress == "esc") {
				return m, tea.Quit
			}

			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

	default:
		return m, nil
	}
}

func (m SelectorModel) View() string {
	if m.choice != "" {
		return m.choice
	}
	return "\n" + m.list.View()
}

func NewSelectorModel() SelectorModel {
	items := []list.Item{
		SelectorItem{
			CommitType: Feat,
			Title:      CommitMessageType[Feat],
		},

		SelectorItem{
			CommitType: Fix,
			Title:      CommitMessageType[Fix],
		},

		SelectorItem{
			CommitType: Docs,
			Title:      CommitMessageType[Docs],
		},

		SelectorItem{
			CommitType: Style,
			Title:      CommitMessageType[Style],
		},

		SelectorItem{
			CommitType: Refactor,
			Title:      CommitMessageType[Refactor],
		},
		SelectorItem{
			CommitType: Test,
			Title:      CommitMessageType[Test],
		},

		SelectorItem{
			CommitType: Chore,
			Title:      CommitMessageType[Chore],
		},

		SelectorItem{
			CommitType: Hotfix,
			Title:      CommitMessageType[Hotfix],
		},
	}
	const defaultWidth = 20

	l := list.New(items, SelectorDelegate{}, defaultWidth, listHeight)
	l.Title = "Select Commit Type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SelectorModel{list: l}
}
