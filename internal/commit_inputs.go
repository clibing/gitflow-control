package internal

import (
	"errors"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

var (
	inputsTitleLayout = lipgloss.NewStyle().
		Padding(1, 0, 1, 2)

	inputsTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#333333", Dark: "#DDDDDD"}).
		Background(lipgloss.AdaptiveColor{Light: "#7653FF", Dark: "#7653FF"}).
		Bold(true).
		Padding(0, 1, 0, 1)

	inputsBlockLayout = lipgloss.NewStyle().
		Padding(0, 0, 1, 0)

	inputsCursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#25A065"})

	inputsPromptFocusStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#9F72FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#9A4AFF", Dark: "#EE6FF8"}).
		Bold(true).
		Padding(0, 0, 0, 1)

	inputsPromptNormalStyle = lipgloss.NewStyle().
		Padding(0, 0, 0, 2)

	inputsTextFocusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FFFDF5"}).
		Bold(true)

	inputsTextNormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FFFDF5"})

	inputsButtonLayout = lipgloss.NewStyle().
		Padding(2, 0, 1, 2)

	inputsButtonFocusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}).
		Background(lipgloss.AdaptiveColor{Light: "#25A065", Dark: "#25A065"}).
		Padding(0, 1, 0, 1).
		Bold(true)

	inputsButtonNormalStyle = inputsButtonFocusStyle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#626262", Dark: "#DDDDDD"}).
		Background(lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#626262"})

	inputsErrLayout = lipgloss.NewStyle().
		Padding(0, 0, 0, 1)

	inputsErrStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF6037")).
		Padding(0, 1, 0, 1).
		Bold(true)

	spinnerMetaFrame1 = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("❯")
	spinnerMetaFrame2 = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("❯")
	spinnerMetaFrame3 = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("❯")
)
type InputWithCheck struct {
	Input   textinput.Model
	Checker func(s string) error
}


type InputsModel struct {
	FocusIndex int
	Title      string
	Inputs     []InputWithCheck
	Err        error
	ErrSpinner spinner.Model
	EditMode   bool
}

func (m InputsModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, spinner.Tick)
}

func (m InputsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.EditMode {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		var renderCursor bool
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.FocusIndex == len(m.Inputs) {
				for _, iwc := range m.Inputs {
					if iwc.Checker != nil {
						m.Err = iwc.Checker(iwc.Input.Value())
						if m.Err != nil {
							return m, spinner.Tick
						}
					}
				}
				return m, func() tea.Msg { return Next{Next: ViewCommitMessage} }
			}
			fallthrough
		case "tab", "down":
			m.FocusIndex++
			if m.FocusIndex > len(m.Inputs) {
				m.FocusIndex = 0
			}
			renderCursor = true
		case "shift+tab", "up":
			m.FocusIndex--
			if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Inputs)
			}
			renderCursor = true
		}

		if renderCursor {
			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Input.Focus()
					m.Inputs[i].Input.PromptStyle = inputsPromptFocusStyle
					m.Inputs[i].Input.TextStyle = inputsTextFocusStyle
					continue
				}
				// Remove focused state
				m.Inputs[i].Input.Blur()
				m.Inputs[i].Input.PromptStyle = inputsPromptNormalStyle
				m.Inputs[i].Input.TextStyle = inputsTextNormalStyle
			}

			return m, tea.Batch(cmds...)
		}

	case string:
		m.Title = "✔ Commit Type: " + msg
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.ErrSpinner, cmd = m.ErrSpinner.Update(msg)
		return m, cmd
	}

	// Handle character input and blinking
	return m, m.updateInputs(msg)
}

func (m *InputsModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.Inputs)+1)

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i].Input, cmds[i] = m.Inputs[i].Input.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputsModel) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].Input.View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := inputsButtonNormalStyle.Render("➜ Submit")
	if m.FocusIndex == len(m.Inputs) {
		button = inputsButtonFocusStyle.Render("➜ Submit")
	}

	// check input value
	for _, iwc := range m.Inputs {
		if iwc.Checker != nil {
			m.Err = iwc.Checker(iwc.Input.Value())
			if m.Err != nil {
				button += inputsErrLayout.Render(m.ErrSpinner.View() + " " + inputsErrStyle.Render(m.Err.Error()))
				break
			}
		}
	}

	b.WriteString(inputsButtonLayout.Render(button))

	title := inputsTitleLayout.Render(inputsTitleStyle.Render(m.Title))
	inputs := inputsBlockLayout.Render(b.String())

	return lipgloss.JoinVertical(lipgloss.Left, title, inputs)
}

func NewInputsModel() InputsModel {
	m := InputsModel{
		Inputs: make([]InputWithCheck, 4),
	}

	for i := range m.Inputs {
		var iwc InputWithCheck

		iwc.Input = textinput.New()
		iwc.Input.CursorStyle = inputsCursorStyle
		iwc.Input.CharLimit = 128

		switch i {
		case 0:
			iwc.Input.Prompt = "1. SCOPE "
			iwc.Input.Placeholder = "Specifying place of the commit change."
			iwc.Input.PromptStyle = inputsPromptFocusStyle
			iwc.Input.TextStyle = inputsTextFocusStyle
			iwc.Input.Focus()
			iwc.Checker = func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("scope cannot be empty")
				}
				return nil
			}
		case 1:
			iwc.Input.Prompt = "2. SUBJECT "
			iwc.Input.PromptStyle = inputsPromptNormalStyle
			iwc.Input.Placeholder = "A very short description of the change."
			iwc.Checker = func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("subject cannot be empty")
				}
				return nil
			}
		case 2:
			iwc.Input.Prompt = "3. BODY "
			iwc.Input.PromptStyle = inputsPromptNormalStyle
			iwc.Input.Placeholder = "Motivation and contrasts for the change."
		case 3:
			iwc.Input.Prompt = "4. FOOTER "
			iwc.Input.PromptStyle = inputsPromptNormalStyle
			iwc.Input.Placeholder = "Description of the change, justification and migration notes."
		}

		m.Inputs[i] = iwc
	}

	m.ErrSpinner = spinner.New()
	m.ErrSpinner.Spinner = spinner.Spinner{
		Frames: []string{
			// "❯   "
			spinnerMetaFrame1 + "   ",
			// "❯❯  "
			spinnerMetaFrame1 + spinnerMetaFrame2 + "  ",
			// "❯❯❯ "
			spinnerMetaFrame1 + spinnerMetaFrame2 + spinnerMetaFrame3 + " ",
			// " ❯❯❯"
			" " + spinnerMetaFrame1 + spinnerMetaFrame2 + spinnerMetaFrame3,
			// "  ❯❯"
			"  " + spinnerMetaFrame1 + spinnerMetaFrame2,
			// "   ❯"
			"   " + spinnerMetaFrame1,
		},
		FPS: time.Second / 10,
	}

	return m
}