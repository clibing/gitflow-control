package main

import (
	"github.com/mattn/go-runewidth"
)
var (
	version     string
	buildDate   string
	buildCommit string
)

func main() {

}
// See also: https://github.com/charmbracelet/lipgloss/issues/40#issuecomment-891167509
func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}
