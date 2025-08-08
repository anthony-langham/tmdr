package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors - using adaptive colors for better terminal compatibility
	primaryColor   = lipgloss.AdaptiveColor{Light: "0", Dark: "15"}    // Black in light mode, White in dark mode
	secondaryColor = lipgloss.AdaptiveColor{Light: "8", Dark: "7"}     // Gray
	accentColor    = lipgloss.AdaptiveColor{Light: "202", Dark: "202"} // Orange (ANSI 256 color)
	bgColor        = lipgloss.NoColor{}                                // No explicit background

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Full screen style
	fullScreenStyle = lipgloss.NewStyle().
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
			Foreground(accentColor)  // Orange color for title

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

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
			Foreground(lipgloss.AdaptiveColor{Light: "1", Dark: "9"})
)

func renderNavBar() string {
	// Special style for "tmdr" with orange and bold
	tmdrStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	// Style for keyboard shortcuts - bold and prominent
	keyStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
	
	// Style for labels - smaller, not bold
	labelStyle := lipgloss.NewStyle().
		Foreground(secondaryColor)
	
	separatorStyle := lipgloss.NewStyle().
		Foreground(secondaryColor)
	
	// Build navigation items with bold keys and lighter labels
	nav := tmdrStyle.Render("tmdr") +
		separatorStyle.Render(" │ ") + keyStyle.Render("s") + labelStyle.Render(" search") +
		separatorStyle.Render(" │ ") + keyStyle.Render("b") + labelStyle.Render(" browse") +
		separatorStyle.Render(" │ ") + keyStyle.Render("f") + labelStyle.Render(" feedback") +
		separatorStyle.Render(" │ ") + keyStyle.Render("q") + labelStyle.Render(" quit")

	return nav
}