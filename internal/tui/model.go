package tui

import (
	"github.com/anthonylangham/tmdr/internal/acronym"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateHome State = iota
	StateBrowse
	StateSearch
	StateFeedback
)

type Model struct {
	state        State
	repo         acronym.Repository
	acronyms     []acronym.Acronym
	filtered     []acronym.Acronym
	cursor       int
	searchInput  textinput.Model
	selected     *acronym.Acronym
	width        int
	height       int
	err          error
}

func NewModel(repo acronym.Repository) Model {
	acronyms, _ := repo.All()
	
	// Initialize text input with orange cursor
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50
	
	// Set cursor to orange using the Cursor properties
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6600")).Background(lipgloss.Color("#FF6600"))
	ti.Cursor.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6600"))
	
	return Model{
		state:       StateHome,
		repo:        repo,
		acronyms:    acronyms,
		filtered:    acronyms,
		searchInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	// Request window size and start textinput blinking
	return tea.Batch(
		tea.WindowSize(),
		textinput.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ALWAYS handle window size updates first, regardless of state
	// This ensures the UI can always render properly
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsMsg.Width
		m.height = wsMsg.Height
		// Don't return, continue processing in case we're in search mode
	}
	
	// If we're in search mode, handle textinput updates for ALL message types
	// This ensures tick messages for cursor blinking are processed
	if m.state == StateSearch {
		// Handle special keys first before textinput gets them
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc":
				if m.searchInput.Value() != "" {
					// Clear search if there's a query
					m.searchInput.SetValue("")
					m.filterAcronyms()
					return m, nil
				} else {
					// Go home if search is empty
					m.state = StateHome
					m.cursor = 0
					return m, nil
				}
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					if m.cursor < len(m.filtered) {
						m.selected = &m.filtered[m.cursor]
					}
				}
				return m, nil
			case "down", "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					if m.cursor < len(m.filtered) {
						m.selected = &m.filtered[m.cursor]
					}
				}
				return m, nil
			case "enter":
				if m.cursor < len(m.filtered) {
					m.selected = &m.filtered[m.cursor]
					m.state = StateBrowse
				}
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				// In search mode, 'q' should be typed, not quit
				// Let it fall through to textinput update
			}
		}
		
		// Update textinput for all other messages (including ticks for blinking)
		var cmd tea.Cmd
		prevValue := m.searchInput.Value()
		var updatedInput textinput.Model
		updatedInput, cmd = m.searchInput.Update(msg)
		m.searchInput = updatedInput
		
		// If the value changed, filter the acronyms
		if m.searchInput.Value() != prevValue {
			m.filterAcronyms()
		}
		
		return m, cmd
	}

	// Handle messages for other states
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle escape key with context-aware behavior
		if msg.String() == "esc" {
			switch m.state {
			case StateHome:
				// From home, quit the app
				return m, tea.Quit
			default:
				// From browse/feedback, go home
				m.state = StateHome
				m.searchInput.SetValue("")
				m.cursor = 0
				return m, nil
			}
		}

		// Handle quit keys (ctrl+c always quits, 'q' only quits outside of search)
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}

		// Handle navigation keys
		// Handle home keys 'h' or 't'
		if msg.String() == "h" || msg.String() == "t" {
			m.state = StateHome
			m.searchInput.SetValue("")
			m.cursor = 0
			return m, nil
		}

		// Handle state transition keys
		switch msg.String() {
		case "s":
			m.state = StateSearch
			m.searchInput.SetValue("")
			m.searchInput.Focus()
			m.filtered = m.acronyms
			m.cursor = 0
			// Start the cursor blinking when entering search
			return m, textinput.Blink

		case "b":
			m.state = StateBrowse
			m.cursor = 0
			m.filtered = m.acronyms
			if len(m.filtered) > 0 {
				m.selected = &m.filtered[0]
			}
			return m, nil

		case "f":
			m.state = StateFeedback
			return m, nil
		}

		// State-specific key handling
		switch m.state {
		case StateBrowse:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					if m.cursor < len(m.filtered) {
						m.selected = &m.filtered[m.cursor]
					}
				}
			case "down", "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					if m.cursor < len(m.filtered) {
						m.selected = &m.filtered[m.cursor]
					}
				}
			}
		}
	}

	return m, nil
}

func (m *Model) filterAcronyms() {
	query := m.searchInput.Value()
	if query == "" {
		m.filtered = m.acronyms
		m.cursor = 0
		return
	}

	filtered := []acronym.Acronym{}

	for _, a := range m.acronyms {
		if contains(a.Acronym, query) || contains(a.FullForm, query) {
			filtered = append(filtered, a)
		}
	}

	m.filtered = filtered
	m.cursor = 0
	if len(m.filtered) > 0 {
		m.selected = &m.filtered[0]
	} else {
		m.selected = nil
	}
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toUpper(s[i+j]) != toUpper(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}