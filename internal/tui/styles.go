package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#FFFFFF")
	secondaryColor = lipgloss.Color("#888888")
	accentColor    = lipgloss.Color("#FF6600")
	bgColor        = lipgloss.Color("#000000")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(primaryColor)

	// Navigation bar styles
	navBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	navItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(primaryColor)

	navSeparatorStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Padding(0, 1)

	// Content area styles
	contentStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Align(lipgloss.Center)

	// List styles
	listItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(accentColor).
				Bold(true)

	// Acronym display styles
	acronymStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	fullFormStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	definitionStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			PaddingTop(1)

	// Search styles
	searchPromptStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Border styles for main container
	containerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)

func renderNavBar() string {
	// Special style for "tmdr" with orange and bold
	tmdrStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(accentColor).
		Bold(true)
	
	items := []string{
		tmdrStyle.Render("tmdr"),
		navItemStyle.Render("s search"),
		navItemStyle.Render("r random"),
		navItemStyle.Render("b browse"),
		navItemStyle.Render("f feedback"),
		navItemStyle.Render("q quit"),
	}

	bar := items[0]
	for i := 1; i < len(items); i++ {
		bar += navSeparatorStyle.Render("│") + items[i]
	}

	return navBarStyle.Render(bar)
}