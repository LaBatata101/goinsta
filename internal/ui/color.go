package ui

import "github.com/charmbracelet/lipgloss"

var BoldText = lipgloss.NewStyle().Bold(true)
var GreenText = lipgloss.NewStyle().Foreground(lipgloss.Color("#728F66"))
var GreenText2Underlined = lipgloss.NewStyle().Underline(true).Inherit(greenText2)
var greenText2 = lipgloss.NewStyle().Foreground(lipgloss.Color("#658E84"))
var RedText = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF3F42"))
var lineNumberColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9588"))
var YellowText = lipgloss.NewStyle().Foreground(lipgloss.Color("#BA9E6B"))
var grayText = lipgloss.NewStyle().Foreground(lipgloss.Color("#94907E"))
