package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FeedbackForm represents our custom feedback form
type FeedbackForm struct {
	fields       []FormField
	currentField int
	width        int
	height       int
	submitted    bool
}

// FormField represents a single form field
type FormField struct {
	label    string
	value    string
	options  []string // For select fields
	isSelect bool
	isBool   bool
	input    textinput.Model // For text input fields
}

// NewFeedbackForm creates a new feedback form
func NewFeedbackForm() *FeedbackForm {
	// Create text input for email field
	emailInput := textinput.New()
	emailInput.Placeholder = "your@email.com (optional)"
	emailInput.Width = 40
	emailInput.CharLimit = 100

	return &FeedbackForm{
		fields: []FormField{
			{
				label:    "Was tmdr useful?",
				options:  []string{"Yes", "No"},
				isSelect: true,
				isBool:   true,
				value:    "Yes",
			},
			{
				label:    "How many times have you used tmdr?",
				options:  []string{"First time", "2-5 times", "6+ times"},
				isSelect: true,
				value:    "First time",
			},
			{
				label:    "Would you use tmdr again?",
				options:  []string{"Definitely", "Probably", "Maybe", "No"},
				isSelect: true,
				value:    "Definitely",
			},
			{
				label:    "How likely are you to recommend tmdr? (1-5)",
				options:  []string{"1 - Not at all", "2", "3", "4", "5 - Extremely likely"},
				isSelect: true,
				value:    "3",
			},
			{
				label:    "What's your role?",
				options:  []string{"Engineer", "DevOps", "Data Scientist", "Healthcare", "Other"},
				isSelect: true,
				value:    "Engineer",
			},
			{
				label:    "Email for updates",
				isSelect: false,
				input:    emailInput,
				value:    "",
			},
		},
		currentField: 0,
		submitted:    false,
	}
}

// Init initializes the form
func (f *FeedbackForm) Init() tea.Cmd {
	// Focus email input if we're on that field
	if f.currentField == len(f.fields)-1 {
		return f.fields[f.currentField].input.Focus()
	}
	return nil
}

// Update handles form updates
func (f *FeedbackForm) Update(msg tea.Msg) (*FeedbackForm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle ESC first, before text input can consume it
		if msg.String() == "esc" {
			// ESC should always be handled by parent, not consumed here
			return f, nil
		}
		
		switch msg.String() {
		case "up", "shift+tab":
			if f.currentField > 0 {
				// Blur current field if it's text input
				if !f.fields[f.currentField].isSelect {
					f.fields[f.currentField].value = f.fields[f.currentField].input.Value()
					f.fields[f.currentField].input.Blur()
				}
				f.currentField--
				// Focus new field if it's text input
				if !f.fields[f.currentField].isSelect {
					return f, f.fields[f.currentField].input.Focus()
				}
			}
			return f, nil
			
		case "down", "tab":
			if f.currentField < len(f.fields)-1 {
				// Save value if text input before moving
				if !f.fields[f.currentField].isSelect {
					f.fields[f.currentField].value = f.fields[f.currentField].input.Value()
					f.fields[f.currentField].input.Blur()
				}
				f.currentField++
				// Focus new field if it's text input
				if !f.fields[f.currentField].isSelect {
					return f, f.fields[f.currentField].input.Focus()
				}
			}
			return f, nil
			
		case "left":
			// For select fields, go to previous option
			if f.fields[f.currentField].isSelect {
				field := &f.fields[f.currentField]
				currentIndex := 0
				for i, opt := range field.options {
					if opt == field.value {
						currentIndex = i
						break
					}
				}
				if currentIndex > 0 {
					field.value = field.options[currentIndex-1]
				}
			}
			return f, nil
			
		case "right":
			// For select fields, go to next option
			if f.fields[f.currentField].isSelect {
				field := &f.fields[f.currentField]
				currentIndex := 0
				for i, opt := range field.options {
					if opt == field.value {
						currentIndex = i
						break
					}
				}
				if currentIndex < len(field.options)-1 {
					field.value = field.options[currentIndex+1]
				}
			}
			return f, nil
			
		case "enter":
			// If on last field, submit the form
			if f.currentField == len(f.fields)-1 {
				// Save email value from input
				if !f.fields[f.currentField].isSelect {
					f.fields[f.currentField].value = f.fields[f.currentField].input.Value()
				}
				f.submitted = true
				return f, nil
			}
			// Move to next field
			if f.currentField < len(f.fields)-1 {
				// Save value if text input
				if !f.fields[f.currentField].isSelect {
					f.fields[f.currentField].value = f.fields[f.currentField].input.Value()
					f.fields[f.currentField].input.Blur()
				}
				f.currentField++
				if !f.fields[f.currentField].isSelect {
					return f, f.fields[f.currentField].input.Focus()
				}
			}
			return f, nil
			
		default:
			// For text input fields, pass through other keys
			if !f.fields[f.currentField].isSelect {
				var cmd tea.Cmd
				f.fields[f.currentField].input, cmd = f.fields[f.currentField].input.Update(msg)
				return f, cmd
			}
		}
	}
	
	// Update text input for non-key messages (like tick for cursor blink)
	if !f.fields[f.currentField].isSelect {
		var cmd tea.Cmd
		f.fields[f.currentField].input, cmd = f.fields[f.currentField].input.Update(msg)
		return f, cmd
	}
	
	return f, nil
}

// View renders the form
func (f *FeedbackForm) View() string {
	if f.submitted {
		return ""
	}
	
	var s strings.Builder
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6600")).
		MarginBottom(1)
	
	s.WriteString(titleStyle.Render("📝 Product Feedback"))
	s.WriteString("\n\n")
	
	for i, field := range f.fields {
		// Field label
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		if i == f.currentField {
			labelStyle = labelStyle.Bold(true).Foreground(lipgloss.Color("#FF6600"))
		}
		
		s.WriteString(labelStyle.Render(fmt.Sprintf("%d. %s", i+1, field.label)))
		s.WriteString("\n")
		
		// Field value/options
		if field.isSelect {
			for _, opt := range field.options {
				optionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
				prefix := "   ○ "
				
				if opt == field.value {
					if i == f.currentField {
						optionStyle = optionStyle.Foreground(lipgloss.Color("#FF6600")).Bold(true)
						prefix = " ▸ ● "
					} else {
						optionStyle = optionStyle.Foreground(lipgloss.Color("#FFFFFF"))
						prefix = "   ● "
					}
				} else if i == f.currentField {
					optionStyle = optionStyle.Foreground(lipgloss.Color("#AAAAAA"))
				}
				
				s.WriteString(optionStyle.Render(prefix + opt))
				s.WriteString("\n")
			}
		} else {
			// Text input field
			s.WriteString("   ")
			s.WriteString(field.input.View())
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}
	
	// Instructions
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).MarginTop(1)
	help := "↑/↓ or tab: navigate • ←/→: change selection • enter: next/submit • esc: cancel"
	s.WriteString(helpStyle.Render(help))
	
	return s.String()
}

// GetValues returns the form values
func (f *FeedbackForm) GetValues() map[string]string {
	// Update email value from input if needed
	if !f.fields[5].isSelect && f.fields[5].input.Value() != "" {
		f.fields[5].value = f.fields[5].input.Value()
	}
	
	return map[string]string{
		"useful":        f.fields[0].value,
		"usage":         f.fields[1].value,
		"wouldUseAgain": f.fields[2].value,
		"npsScore":      f.fields[3].value,
		"role":          f.fields[4].value,
		"email":         f.fields[5].value,
	}
}

// IsSubmitted returns whether the form was submitted
func (f *FeedbackForm) IsSubmitted() bool {
	return f.submitted
}

// Reset resets the form to initial state
func (f *FeedbackForm) Reset() {
	f.currentField = 0
	f.submitted = false
	f.fields[0].value = "Yes"
	f.fields[1].value = "First time"
	f.fields[2].value = "Definitely"
	f.fields[3].value = "3"
	f.fields[4].value = "Engineer"
	f.fields[5].value = ""
	f.fields[5].input.Reset()
}