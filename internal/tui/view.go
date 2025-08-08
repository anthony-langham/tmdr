package tui

import (
	"fmt"
	"strings"

	"github.com/anthonylangham/tmdr/internal/version"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	
	// Build update notification if applicable
	var updateNotification string
	if m.updateReady {
		updateNotification = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "2", Dark: "10"}).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
			Padding(0, 1).
			Width(m.width).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("✨ tmdr v%s downloaded! Restart to apply update", m.updateInfo.Version))
	} else if m.updateDownloading {
		updateNotification = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "202", Dark: "202"}).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
			Padding(0, 1).
			Width(m.width).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("⬇️ Downloading tmdr v%s...", m.updateInfo.Version))
	} else if m.updateInfo.Available && m.updateError == nil {
		updateNotification = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "4", Dark: "12"}).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
			Padding(0, 1).
			Width(m.width).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("🆕 tmdr v%s is available!", m.updateInfo.Version))
	}

	// Check minimum terminal size
	minWidth := 50
	minHeight := 16
	
	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal too small!\n\nMinimum size: %dx%d\nCurrent size: %dx%d\n\nPlease resize your terminal.",
			minWidth, minHeight, m.width, m.height)
		return fullScreenStyle.
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(msg)
	}

	// Build the navigation bar
	navContent := renderNavBar()
	navBar := lipgloss.NewStyle().
		Width(m.width - 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "8", Dark: "7"}).
		Padding(0, 1).
		Render(navContent)

	// Build the content based on state
	var content string
	switch m.state {
	case StateHome:
		content = m.viewHome()
	case StateBrowse:
		content = m.viewBrowse()
	case StateSearch:
		content = m.viewSearch()
	case StateFeedback:
		content = m.viewFeedback()
	}

	// Combine navigation and content
	fullView := lipgloss.JoinVertical(
		lipgloss.Left,
		navBar,
		content,
	)

	// Apply container styling with full dark background
	containerView := containerStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(fullView)

	// Wrap in full screen style to ensure dark background fills entire terminal
	finalView := fullScreenStyle.
		Width(m.width).
		Height(m.height).
		Render(containerView)
	
	// Add update notification at the top if present
	if updateNotification != "" {
		return lipgloss.JoinVertical(lipgloss.Top, updateNotification, finalView)
	}
	
	return finalView
}

func (m Model) viewHome() string {
	// Use PlaceHorizontal for centering without background artifacts
	width := m.width - 4
	if width < 1 {
		width = 1
	}
	
	title := lipgloss.PlaceHorizontal(width, lipgloss.Center, titleStyle.Render("too medical; didn't read"))
	subtitle := lipgloss.PlaceHorizontal(width, lipgloss.Center, subtitleStyle.Render("Your terminal-native tool for instant medical acronym help."))
	
	// Adjust content based on available height
	var content string
	if m.height > 20 {
		// Full version for larger terminals
		instructions := []string{
			"📖  Quick Start:",
			"    • Press 's' to search for an acronym",
			"    • Press 'b' to browse all acronyms",
			"    • Press 'f' to send feedback",
			"    • Press 'h' or 't' to return home",
		}
		
		// Center each instruction line
		centeredInstructions := make([]string, len(instructions))
		for i, line := range instructions {
			centeredInstructions[i] = lipgloss.PlaceHorizontal(width, lipgloss.Center, line)
		}
		
		dataInfo := fmt.Sprintf("⚙️  version: v%s", version.Version)
		centeredDataInfo := lipgloss.PlaceHorizontal(width, lipgloss.Center, dataInfo)
		
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			title,
			subtitle,
			"",
			strings.Join(centeredInstructions, "\n"),
			"",
			centeredDataInfo,
		)
	} else {
		// Compact version for smaller terminals
		shortcuts := "s: search | b: browse | f: feedback"
		centeredShortcuts := lipgloss.PlaceHorizontal(width, lipgloss.Center, shortcuts)
		
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			centeredShortcuts,
		)
	}

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(content)
}

func (m Model) viewBrowse() string {
	var listBuilder strings.Builder
	
	// Calculate available height for list (accounting for borders, nav, details)
	// Cap at 10 results for better UX and to prevent navbar issues
	availableHeight := 10
	if m.height-8 < availableHeight {
		availableHeight = m.height - 8
	}
	if availableHeight < 3 {
		availableHeight = 3
	}
	
	// Display list of acronyms
	start := 0
	if m.cursor > availableHeight/2 {
		start = m.cursor - availableHeight/2
	}
	end := start + availableHeight
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		item := m.filtered[i]
		line := fmt.Sprintf("%-6s %s", item.Acronym, item.FullForm)
		
		if i == m.cursor {
			listBuilder.WriteString(selectedItemStyle.Render("> " + line))
		} else {
			listBuilder.WriteString(listItemStyle.Render(line))
		}
		if i < end-1 {
			listBuilder.WriteString("\n")
		}
	}

	// Display selected acronym details
	var details string
	if m.selected != nil {
		acronymLine := fmt.Sprintf("%s → %s", 
			acronymStyle.Render(m.selected.Acronym),
			fullFormStyle.Render(m.selected.FullForm))
		
		definition := definitionStyle.Render(m.selected.Definition)
		
		details = lipgloss.JoinVertical(
			lipgloss.Left,
			strings.Repeat("─", 60),
			acronymLine,
			definition,
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		listBuilder.String(),
		"",
		details,
	)

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(content)
}

func (m Model) viewSearch() string {
	// Use the textinput component with blinking cursor
	searchLine := searchPromptStyle.Render("Search: ") + m.searchInput.View()

	var results strings.Builder
	if m.searchInput.Value() != "" {
		results.WriteString("\nResults:\n")
		
		// Calculate available height for results
		// Cap at 10 results for better UX and to prevent navbar issues
		availableHeight := 10
		if m.height-7 < availableHeight {
			availableHeight = m.height - 7
		}
		if availableHeight < 2 {
			availableHeight = 2
		}
		
		displayCount := availableHeight
		if len(m.filtered) < displayCount {
			displayCount = len(m.filtered)
		}

		if displayCount == 0 {
			results.WriteString(helpStyle.Render("  No results found"))
		} else {
			for i := 0; i < displayCount; i++ {
				item := m.filtered[i]
				line := fmt.Sprintf("%-6s %s", item.Acronym, item.FullForm)
				
				if i == m.cursor {
					results.WriteString(selectedItemStyle.Render("> " + line))
				} else {
					results.WriteString(listItemStyle.Render(line))
				}
				if i < displayCount-1 {
					results.WriteString("\n")
				}
			}
		}
	}

	help := helpStyle.Render("Type to search • ↑↓ navigate • Enter select • Esc cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		searchLine,
		results.String(),
		"",
		strings.Repeat("─", 60),
		help,
	)

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(content)
}

func (m Model) viewFeedback() string {
	if m.formSubmitted {
		// Show success message briefly
		title := titleStyle.Render("✅  Thank You!")
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			title,
			"",
			"Your feedback has been prepared in your email client.",
			"",
			"Returning to home...",
		)
		
		return contentStyle.
			Width(m.width - 4).
			Height(m.height - 6).
			Align(lipgloss.Center).
			Render(content)
	}
	
	// Show the feedback form directly
	// Get the form view
	formView := ""
	if m.feedbackForm != nil {
		formView = m.feedbackForm.View()
	}
	
	// Debug: If form is empty, show a message
	if formView == "" {
		formView = "Error: Feedback form not initialized"
	}
	
	content := formView

	// Don't restrict height for feedback form to prevent truncation
	return contentStyle.
		Width(m.width - 4).
		Padding(1, 2).
		Render(content)
}