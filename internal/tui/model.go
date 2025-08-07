package tui

import (
	"github.com/anthonylangham/tmdr/internal/acronym"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateHome State = iota
	StateBrowse
	StateSearch
	StateRandom
	StateFeedback
)

type Model struct {
	state        State
	repo         acronym.Repository
	acronyms     []acronym.Acronym
	filtered     []acronym.Acronym
	cursor       int
	searchQuery  string
	selected     *acronym.Acronym
	width        int
	height       int
	err          error
}

func NewModel(repo acronym.Repository) Model {
	acronyms, _ := repo.All()
	return Model{
		state:    StateHome,
		repo:     repo,
		acronyms: acronyms,
		filtered: acronyms,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "s":
			m.state = StateSearch
			m.searchQuery = ""
			m.filtered = m.acronyms
			m.cursor = 0
			return m, nil

		case "r":
			m.state = StateRandom
			if random, err := m.repo.Random(); err == nil {
				m.selected = random
			}
			return m, nil

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

		case "escape":
			if m.state != StateHome {
				m.state = StateHome
				m.searchQuery = ""
				m.cursor = 0
			}
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

		case StateSearch:
			switch msg.String() {
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.filterAcronyms()
				}
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
			case "enter":
				if m.cursor < len(m.filtered) {
					m.selected = &m.filtered[m.cursor]
					m.state = StateBrowse
				}
			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.filterAcronyms()
				}
			}
		}
	}

	return m, nil
}

func (m *Model) filterAcronyms() {
	if m.searchQuery == "" {
		m.filtered = m.acronyms
		m.cursor = 0
		return
	}

	filtered := []acronym.Acronym{}
	query := m.searchQuery

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