package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	layOutStyle = lipgloss.NewStyle().
			Padding(1, 0, 1, 2)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Width(64).
			Foreground(lipgloss.AdaptiveColor{Light: "#E11C9C", Dark: "#FF62DA"}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#E11C9C", Dark: "#FF62DA"}).
			Padding(1, 3, 1, 3)
)

type ErrorModel struct {
	Err error
}

func NewErrorModel() ErrorModel {
	return ErrorModel{}
}

func (m ErrorModel) Init() tea.Cmd {
	return nil
}

func (m ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case error:
		m.Err = msg.(error)
	}
	return m, tea.Quit
}

func (m ErrorModel) View() string {
	if m.Err == nil {
		return ""
	}
	return layOutStyle.Render(errorStyle.Render(m.Err.Error()))
}
