package internal

import "github.com/charmbracelet/lipgloss"

const (
	BinName = "gitflow-control"
)

var (
	Cyan  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	Green = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	Gray  = lipgloss.NewStyle().Foreground(lipgloss.Color("#696969"))
	Gold  = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8860B"))
)
