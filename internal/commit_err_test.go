package internal

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"testing"
)

func TestStyle(t *testing.T) {
	layOutStyle1 := lipgloss.NewStyle().
		Padding(1, 0, 1, 2)

	errorStyle1 := lipgloss.NewStyle().
		Bold(true).
		Width(64).
		Foreground(lipgloss.AdaptiveColor{Light: "#E11C9C", Dark: "#FF62DA"}).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#E11C9C", Dark: "#FF62DA"}).
		Padding(1, 3, 1, 3)
	fmt.Println(layOutStyle1.Render(errorStyle1.Render("err msg")))
}
