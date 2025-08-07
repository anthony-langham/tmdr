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

	// Full screen style to enforce dark background
	fullScreenStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(primaryColor)

	// Navigation bar styles
	navBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(secondaryColor).
			Background(bgColor).
			Padding(0, 1)

	navItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(primaryColor)

	navSeparatorStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Padding(0, 1)

	// Content area styles
	contentStyle = lipgloss.NewStyle().
			Background(bgColor).
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
			BorderForeground(secondaryColor).
			Background(bgColor)

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
		Foreground(accentColor).
		Background(bgColor).
		Bold(true)
	
	// Style for keyboard shortcuts - bold and prominent
	keyStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Background(bgColor).
		Bold(true)
	
	// Style for labels - smaller, not bold
	labelStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Background(bgColor)
	
	separatorStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Background(bgColor)
	
	// Build navigation items with bold keys and lighter labels
	nav := tmdrStyle.Render("tmdr") +
		separatorStyle.Render(" │ ") + keyStyle.Render("s") + labelStyle.Render(" search") +
		separatorStyle.Render(" │ ") + keyStyle.Render("b") + labelStyle.Render(" browse") +
		separatorStyle.Render(" │ ") + keyStyle.Render("f") + labelStyle.Render(" feedback") +
		separatorStyle.Render(" │ ") + keyStyle.Render("q") + labelStyle.Render(" quit")

	return nav
}