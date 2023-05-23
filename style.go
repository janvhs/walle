package main

import "github.com/charmbracelet/lipgloss"

var messageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(""))

var scanningStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))

var langStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("2"))

var dirListStyle = lipgloss.NewStyle().
	PaddingLeft(2).
	Foreground(lipgloss.Color("6"))

var destructiveQuestion = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("5"))

var strongWarningStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("1"))

var errorStyle = strongWarningStyle

func clearLine() {
	lipgloss.DefaultRenderer().Output().ClearLines(1)
}
