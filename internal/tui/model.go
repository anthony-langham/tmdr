package tui

import (
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/anthonylangham/tmdr/internal/acronym"
	"github.com/anthonylangham/tmdr/internal/update"
	"github.com/anthonylangham/tmdr/internal/version"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
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
	
	// Feedback form fields
	feedbackForm      *huh.Form
	useful            bool
	usage             string
	wouldUseAgain     string
	npsScore          int
	role              string
	email             string
	formSubmitted     bool
	formInteracted    bool  // Track if user has interacted with the form
	formFieldIndex    int   // Track current field position
	
	// Update state
	updateInfo        update.UpdateInfo
	updateDownloading bool
	updateProgress    float64
	updateError       error
	updateReady       bool
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
	
	m := Model{
		state:       StateHome,
		repo:        repo,
		acronyms:    acronyms,
		filtered:    acronyms,
		searchInput: ti,
		npsScore:    3, // Default to middle score
	}
	
	// Initialize the feedback form
	m.createFeedbackForm()
	
	return m
}

func (m Model) Init() tea.Cmd {
	// Request window size, start textinput blinking, and check for updates
	return tea.Batch(
		tea.WindowSize(),
		textinput.Blink,
		m.checkForUpdate(),
	)
}

// Custom messages for update process
type updateAvailableMsg update.UpdateInfo
type updateProgressMsg float64
type updateCompleteMsg string
type updateErrorMsg error

func (m Model) checkForUpdate() tea.Cmd {
	return func() tea.Msg {
		info := update.CheckForUpdateWithAssets()
		return updateAvailableMsg(info)
	}
}

func (m Model) downloadUpdate() tea.Cmd {
	return func() tea.Msg {
		if m.updateInfo.DownloadURL == "" {
			return updateErrorMsg(fmt.Errorf("no download URL available"))
		}

		tempFile, err := update.DownloadUpdate(m.updateInfo.DownloadURL, func(downloaded, total int64) {
			// Progress callback - in production would send progress messages
		})
		
		if err != nil {
			return updateErrorMsg(err)
		}
		
		// Try to install the update
		if err := update.InstallUpdate(tempFile); err != nil {
			// If install fails, still mark as ready for manual install
			return updateCompleteMsg(tempFile)
		}
		
		return updateCompleteMsg(tempFile)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ALWAYS handle window size updates first, regardless of state
	// This ensures the UI can always render properly
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsMsg.Width
		m.height = wsMsg.Height
		// Don't return, continue processing in case we're in search mode
	}
	
	// Handle update messages
	switch msg := msg.(type) {
	case updateAvailableMsg:
		m.updateInfo = update.UpdateInfo(msg)
		if m.updateInfo.Available && !m.updateDownloading && !m.updateReady {
			// Start downloading automatically
			m.updateDownloading = true
			return m, m.downloadUpdate()
		}
		return m, nil
		
	case updateProgressMsg:
		m.updateProgress = float64(msg)
		return m, nil
		
	case updateCompleteMsg:
		m.updateDownloading = false
		m.updateReady = true
		return m, nil
		
	case updateErrorMsg:
		m.updateDownloading = false
		m.updateError = msg
		return m, nil
	}
	
	// Handle StateFeedback for ALL message types (not just KeyMsg)
	// This ensures form renders on initial load and gets all necessary updates
	if m.state == StateFeedback {
		// Intercept ESC and handle it specially
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "esc" {
				// If we're at field index 0 (first question), exit to home
				if m.formFieldIndex == 0 {
					m.state = StateHome
					m.resetFeedbackForm()
					return m, nil
				}
				
				// Otherwise try to go back with Shift+Tab
				prevFieldMsg := tea.KeyMsg{
					Type: tea.KeyShiftTab,
					Runes: []rune{},
					Alt: false,
				}
				
				// Navigate back and decrement field index
				form, cmd := m.feedbackForm.Update(prevFieldMsg)
				if f, ok := form.(*huh.Form); ok {
					m.feedbackForm = f
					if m.formFieldIndex > 0 {
						m.formFieldIndex--
					}
				}
				return m, cmd
			}
			
			// Track forward navigation
			if keyMsg.String() == "enter" || keyMsg.String() == "tab" {
				m.formInteracted = true
				// Increment field index when moving forward
				m.formFieldIndex++
			} else if keyMsg.String() == "shift+tab" {
				// Decrement field index when moving backward
				if m.formFieldIndex > 0 {
					m.formFieldIndex--
				}
			}
		}
		
		// Update form with the original message (if not ESC)
		form, cmd := m.feedbackForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.feedbackForm = f
			
			// Check if form is completed
			if m.feedbackForm.State == huh.StateCompleted {
				m.submitFeedback()
				// Return to home after submission
				m.state = StateHome
				m.resetFeedbackForm() // Reset form after submission
				return m, nil
			}
			
			// Check if form was aborted
			if m.feedbackForm.State == huh.StateAborted {
				// User explicitly aborted the form
				m.state = StateHome
				m.resetFeedbackForm()
				return m, nil
			}
		}
		
		return m, cmd
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
			// Reset field index when entering feedback
			m.formFieldIndex = 0
			m.formInteracted = false
			// Send an Init command to ensure form renders immediately
			return m, m.feedbackForm.Init()
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

func (m *Model) createFeedbackForm() {
	// Create a custom theme with orange for focused elements
	theme := huh.ThemeBase()
	
	// Orange color for selections and focus
	orange := lipgloss.Color("#FF6600")
	
	// Customize focused state
	theme.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(orange).SetString("▸ ")
	theme.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("  ")
	
	// Style for Yes/No buttons in Confirm
	theme.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(orange).
		Padding(0, 2).
		Bold(true)
	theme.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#999999")).
		Padding(0, 2)
	
	// Additional styling
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(orange).Background(orange)
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	theme.Focused.SelectedOption = lipgloss.NewStyle().Foreground(orange)
	theme.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	
	// Also set the cursor/prompt for input fields
	theme.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(orange)
	theme.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	
	m.feedbackForm = huh.NewForm(
		// Group 1: Basic feedback questions
		huh.NewGroup(
			huh.NewConfirm().
				Title("Was tmdr useful?").
				Value(&m.useful).
				Affirmative("Yes").
				Negative("No"),

			huh.NewSelect[string]().
				Title("How many times have you used tmdr?").
				Options(
					huh.NewOption("First time", "first"),
					huh.NewOption("2-5 times", "2-5"),
					huh.NewOption("6+ times", "6+"),
				).
				Value(&m.usage),

			huh.NewSelect[string]().
				Title("Would you use tmdr again?").
				Options(
					huh.NewOption("Definitely", "definitely"),
					huh.NewOption("Probably", "probably"),
					huh.NewOption("Maybe", "maybe"),
					huh.NewOption("No", "no"),
				).
				Value(&m.wouldUseAgain),
		),
		
		// Group 2: NPS and demographics
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How likely are you to recommend tmdr? (1-5)").
				Options(
					huh.NewOption("1 - Not at all likely", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5 - Extremely likely", 5),
				).
				Value(&m.npsScore).
				Height(7), // Limit height for scrollable list

			huh.NewSelect[string]().
				Title("What's your role?").
				Options(
					huh.NewOption("Engineer", "engineer"),
					huh.NewOption("DevOps", "devops"),
					huh.NewOption("Data Scientist", "data_scientist"),
					huh.NewOption("Healthcare Professional", "healthcare"),
					huh.NewOption("Other", "other"),
				).
				Value(&m.role),

			huh.NewInput().
				Title("Email for updates (optional)").
				Value(&m.email).
				Placeholder("your@email.com"),
		),
	).WithWidth(60).WithShowHelp(false).WithShowErrors(true).WithTheme(theme)
}

func (m *Model) formatFeedbackEmail() string {
	usefulStr := "No"
	if m.useful {
		usefulStr = "Yes"
	}
	
	timestamp := time.Now().Format(time.RFC3339)
	osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	
	body := fmt.Sprintf(`=== TMDR Product Feedback ===
Useful: %s
Usage: %s
Would use again: %s
NPS Score: %d/5
Role: %s
Email: %s
---
Version: v%s
OS: %s
Timestamp: %s`,
		usefulStr,
		m.usage,
		m.wouldUseAgain,
		m.npsScore,
		m.role,
		m.email,
		version.Version,
		osInfo,
		timestamp,
	)
	
	return body
}

func (m *Model) submitFeedback() {
	body := m.formatFeedbackEmail()
	encodedBody := url.QueryEscape(body)
	subject := url.QueryEscape("TMDR Product Feedback")
	
	mailtoURL := fmt.Sprintf("mailto:hello@tmdr.sh?subject=%s&body=%s", subject, encodedBody)
	
	// Open the mailto link in the default email client
	browser.OpenURL(mailtoURL)
	
	m.formSubmitted = true
}

func (m *Model) resetFeedbackForm() {
	m.useful = false
	m.usage = ""
	m.wouldUseAgain = ""
	m.npsScore = 3
	m.role = ""
	m.email = ""
	m.formSubmitted = false
	m.formInteracted = false  // Reset the interaction flag
	m.formFieldIndex = 0     // Reset field index
	// Don't recreate the form - it's already bound to our fields
}