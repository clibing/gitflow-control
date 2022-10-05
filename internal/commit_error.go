package internal

import tea "github.com/charmbracelet/bubbletea"

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
	// return layOutStyle.Render(errorStyle.Render(m.err.Error()))
	return m.Err.Error()
}
