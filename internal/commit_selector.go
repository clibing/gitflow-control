package internal

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

var (
	selectorTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#2E2E2E", Dark: "#DDDDDD"}).
		Background(lipgloss.AdaptiveColor{Light: "#19A04B", Dark: "#25A065"}).
		Bold(true).
		Padding(0, 1)

	selectorNormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"}).
		Padding(0, 0, 0, 2)

	selectorSelectedStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#9F72FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#9A4AFF", Dark: "#EE6FF8"}).
		Bold(true).
		Padding(0, 0, 0, 1)

	selectorPaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)

	selectorHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#6F6C6C", Dark: "#7A7A7A"})
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
		_, _ = fmt.Fprintf(w, selectorSelectedStyle.Render(str))
	} else {
		_, _ = fmt.Fprintf(w, selectorNormalStyle.Render(str))
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
	l := list.New(items, SelectorDelegate{}, 20, 12)
	l.Title = "Select Commit Type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = selectorTitleStyle
	l.Styles.PaginationStyle = selectorPaginationStyle
	h := help.New()
	h.Styles.ShortDesc = selectorHelpStyle
	h.Styles.ShortSeparator = selectorHelpStyle
	h.Styles.ShortKey = selectorHelpStyle
	l.Help = h

	return SelectorModel{list: l}
}
