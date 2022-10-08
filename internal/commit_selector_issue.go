package internal

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectorIssueModel struct {
	list   list.Model
	choice string
}

type SelectorIssueItem struct {
	Title string
}

func (si SelectorIssueItem) FilterValue() string { return si.Title }

type SelectorIssueDelegate struct{}

func (d SelectorIssueDelegate) Height() int                             { return 1 }
func (d SelectorIssueDelegate) Spacing() int                            { return 0 }
func (d SelectorIssueDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SelectorIssueDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectorItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. (%s)%s", index+1, i.CommitType, i.Title)
	if index == m.Index() {
		_, _ = fmt.Fprintf(w, selectorSelectedStyle.Render(str))
	} else {
		_, _ = fmt.Fprintf(w, selectorNormalStyle.Render(str))
	}

}

func (m SelectorIssueModel) Init() tea.Cmd {
	return nil
}

func (m SelectorIssueModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, func() tea.Msg { return Next{Next: ViewInputCommitMessage} }

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

func (m SelectorIssueModel) View() string {
	if m.choice != "" {
		return m.choice
	}
	return "\n" + m.list.View()
}

func NewSelectorIssueModel() SelectorIssueModel {
	items := []list.Item{
		SelectorIssueItem{
			Title: "backend-001",
		},
		SelectorIssueItem{
			Title: "backend-002",
		},
	}
	l := list.New(items, SelectorDelegate{}, 20, 12)
	l.Title = "Select Issue"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = selectorTitleStyle
	l.Styles.PaginationStyle = selectorPaginationStyle
	h := help.New()
	h.Styles.ShortDesc = selectorHelpStyle
	h.Styles.ShortSeparator = selectorHelpStyle
	h.Styles.ShortKey = selectorHelpStyle
	l.Help = h

	return SelectorIssueModel{list: l}
}
