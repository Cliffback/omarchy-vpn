package main

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Semantic color variables — ANSI terminal colors, set by initColors()
var (
	green     color.Color // connected, active states
	red       color.Color // errors
	yellow    color.Color // warnings
	accent    color.Color // titles, active borders, highlights, shortcuts
	textCol   color.Color // primary text
	dimCol    color.Color // labels, inactive text, dim elements
	borderCol color.Color // panel borders, dividers
	base      color.Color // background
)

// Style variables — set by initStyles()
var (
	panelBorderStyle  lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	labelStyle        lipgloss.Style
	valueStyle        lipgloss.Style
	connectedStyle    lipgloss.Style
	errorStyle        lipgloss.Style
	warnStyle         lipgloss.Style
	dimStyle          lipgloss.Style
	helpOverlayStyle  lipgloss.Style
	helpTitleStyle    lipgloss.Style
	inputPromptStyle  lipgloss.Style
	spinnerStyle      lipgloss.Style
)

func initStyles() {
	panelBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderCol)

	itemStyle = lipgloss.NewStyle().
		Foreground(textCol)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	labelStyle = lipgloss.NewStyle().
		Foreground(dimCol).
		Width(12)

	valueStyle = lipgloss.NewStyle().
		Foreground(textCol)

	connectedStyle = lipgloss.NewStyle().Foreground(green)
	errorStyle = lipgloss.NewStyle().Foreground(red)
	warnStyle = lipgloss.NewStyle().Foreground(yellow)
	dimStyle = lipgloss.NewStyle().Foreground(dimCol)

	helpOverlayStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 3)

	helpTitleStyle = lipgloss.NewStyle().
		Foreground(textCol)

	inputPromptStyle = lipgloss.NewStyle().
		Foreground(accent)

	spinnerStyle = lipgloss.NewStyle().
		Foreground(accent)
}
