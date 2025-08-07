package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Build the navigation bar
	navBar := renderNavBar()

	// Build the content based on state
	var content string
	switch m.state {
	case StateHome:
		content = m.viewHome()
	case StateBrowse:
		content = m.viewBrowse()
	case StateSearch:
		content = m.viewSearch()
	case StateRandom:
		content = m.viewRandom()
	case StateFeedback:
		content = m.viewFeedback()
	}

	// Combine navigation and content
	fullView := lipgloss.JoinVertical(
		lipgloss.Left,
		navBar,
		content,
	)

	// Apply container styling
	return containerStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(fullView)
}

func (m Model) viewHome() string {
	title := titleStyle.Render("🩺  Welcome to tmdr")
	subtitle := subtitleStyle.Render("Too Medical; Didn't Read.")
	tagline := subtitleStyle.Render("Your terminal-native tool for instant medical acronym help.")

	instructions := []string{
		"📖  Quick Start:",
		"    • Press 's' to search for an acronym",
		"    • Press 'r' for a random medical term",
		"    • Press 'b' to browse all acronyms",
		"    • Press 'f' to send feedback",
	}

	dataInfo := fmt.Sprintf("💾  Data version: v0.1  |  Acronyms loaded: %d", len(m.acronyms))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		subtitle,
		tagline,
		"",
		strings.Join(instructions, "\n"),
		"",
		"",
		dataInfo,
	)

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(content)
}

func (m Model) viewBrowse() string {
	var listBuilder strings.Builder
	
	// Display list of acronyms
	start := 0
	if m.cursor > 5 {
		start = m.cursor - 5
	}
	end := start + 10
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
			strings.Repeat("─", m.width-6),
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
	searchLine := searchPromptStyle.Render("Search: ") + 
		searchInputStyle.Render(m.searchQuery) + "_"

	var results strings.Builder
	if m.searchQuery != "" {
		results.WriteString("\nResults:\n")
		
		displayCount := 8
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
		strings.Repeat("─", m.width-6),
		help,
	)

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(content)
}

func (m Model) viewRandom() string {
	title := titleStyle.Render("🎲  Random Medical Acronym")

	var content string
	if m.selected != nil {
		acronymLine := fmt.Sprintf("%s → %s",
			acronymStyle.Render(m.selected.Acronym),
			fullFormStyle.Render(m.selected.FullForm))

		definition := definitionStyle.Width(m.width - 8).Render(m.selected.Definition)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			title,
			"",
			acronymLine,
			"",
			definition,
			"",
			strings.Repeat("─", m.width-6),
			helpStyle.Render("Press 'r' for another random acronym"),
		)
	} else {
		content = errorStyle.Render("Error loading random acronym")
	}

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Align(lipgloss.Center).
		Render(content)
}

func (m Model) viewFeedback() string {
	title := titleStyle.Render("📝  Send Feedback")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		"Found a missing acronym? Have a suggestion?",
		"Visit: github.com/anthonylangham/tmdr/issues",
		"",
		"Or press Enter to open in your browser",
		"",
		strings.Repeat("─", m.width-6),
		helpStyle.Render("Press any other key to go back"),
	)

	return contentStyle.
		Width(m.width - 4).
		Height(m.height - 6).
		Align(lipgloss.Center).
		Render(content)
}