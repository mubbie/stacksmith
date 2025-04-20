// ui/simplemenu/stack_prompt.go
package simplemenu

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StackPromptModel handles user input for the stack command
type stackPromptModel struct {
	newBranch    string
	parentBranch string
	currentField int
	cursorPos    int
	err          string
}

func (m stackPromptModel) Init() tea.Cmd {
	return nil
}

func (m stackPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Don't set an error message, just quietly return a signal
			m.newBranch = ""
			m.parentBranch = ""
			m.err = "cancelled" // Internal signal, won't be displayed
			return m, tea.Quit

		case "enter":
			if m.currentField == 0 {
				if m.newBranch == "" {
					m.err = "New branch name cannot be empty"
					return m, nil
				}
				m.currentField = 1
				m.cursorPos = len(m.parentBranch)
				return m, nil
			}

			if m.parentBranch == "" {
				m.err = "Parent branch name cannot be empty"
				return m, nil
			}

			// Both fields are filled, time to return
			return m, tea.Quit

		case "tab":
			if m.currentField == 0 && m.newBranch != "" {
				m.currentField = 1
				m.cursorPos = len(m.parentBranch)
			} else if m.currentField == 1 && m.parentBranch != "" {
				m.currentField = 0
				m.cursorPos = len(m.newBranch)
			}
			return m, nil

		case "backspace":
			if m.currentField == 0 && len(m.newBranch) > 0 && m.cursorPos > 0 {
				m.newBranch = m.newBranch[:m.cursorPos-1] + m.newBranch[m.cursorPos:]
				m.cursorPos--
			} else if m.currentField == 1 && len(m.parentBranch) > 0 && m.cursorPos > 0 {
				m.parentBranch = m.parentBranch[:m.cursorPos-1] + m.parentBranch[m.cursorPos:]
				m.cursorPos--
			}
			m.err = ""
			return m, nil

		case "left":
			if m.currentField == 0 && m.cursorPos > 0 {
				m.cursorPos--
			} else if m.currentField == 1 && m.cursorPos > 0 {
				m.cursorPos--
			}
			return m, nil

		case "right":
			if m.currentField == 0 && m.cursorPos < len(m.newBranch) {
				m.cursorPos++
			} else if m.currentField == 1 && m.cursorPos < len(m.parentBranch) {
				m.cursorPos++
			}
			return m, nil

		default:
			// Handle regular character input
			if msg.String() != "tab" && msg.String() != "enter" && len(msg.String()) == 1 {
				if m.currentField == 0 {
					m.newBranch = m.newBranch[:m.cursorPos] + msg.String() + m.newBranch[m.cursorPos:]
					m.cursorPos++
				} else if m.currentField == 1 {
					m.parentBranch = m.parentBranch[:m.cursorPos] + msg.String() + m.parentBranch[m.cursorPos:]
					m.cursorPos++
				}
				m.err = ""
			}
			return m, nil
		}
	}
	return m, nil
}

func (m stackPromptModel) View() string {
	// Use title styling but ensure no padding or indentation is applied
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffc27d")).
		PaddingLeft(0). // Explicitly set padding to 0
		MarginLeft(0).  // Explicitly set margin to 0
		Render("🪵 Create a new branch")

	s := title + "\n\n"

	// New branch field
	newBranchLabel := "New branch name: "
	if m.currentField == 0 {
		s += lipgloss.NewStyle().Bold(true).Render(newBranchLabel)

		// Insert cursor at the current position
		if m.cursorPos == len(m.newBranch) {
			s += m.newBranch + lipgloss.NewStyle().Background(lipgloss.Color("#666666")).Render(" ")
		} else {
			s += m.newBranch[:m.cursorPos] +
				lipgloss.NewStyle().Background(lipgloss.Color("#666666")).Render(string(m.newBranch[m.cursorPos])) +
				m.newBranch[m.cursorPos+1:]
		}
	} else {
		s += newBranchLabel + m.newBranch
	}

	s += "\n"

	// Parent branch field
	parentBranchLabel := "Parent branch: "
	if m.currentField == 1 {
		s += lipgloss.NewStyle().Bold(true).Render(parentBranchLabel)

		// Insert cursor at the current position
		if m.cursorPos == len(m.parentBranch) {
			s += m.parentBranch + lipgloss.NewStyle().Background(lipgloss.Color("#666666")).Render(" ")
		} else {
			s += m.parentBranch[:m.cursorPos] +
				lipgloss.NewStyle().Background(lipgloss.Color("#666666")).Render(string(m.parentBranch[m.cursorPos])) +
				m.parentBranch[m.cursorPos+1:]
		}
	} else {
		s += parentBranchLabel + m.parentBranch
	}

	// Error message
	if m.err != "" {
		s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.err)
	}

	// Help text
	helpText := "Tab: Switch fields • Enter: Confirm • Esc: Return to menu"

	s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(helpText)

	return s
}

// RunStackPrompt shows a prompt for the stack command
func RunStackPrompt() (string, string, bool) {
	p := tea.NewProgram(stackPromptModel{
		currentField: 0,
		parentBranch: "main", // Default parent branch
	})

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return "", "", false
	}

	if m, ok := m.(stackPromptModel); ok {
		// Check if user cancelled
		if m.err == "cancelled" {
			return "", "", false
		}

		// Check if both fields are filled
		if m.newBranch != "" && m.parentBranch != "" {
			return m.newBranch, m.parentBranch, true
		}
	}

	return "", "", false
}
