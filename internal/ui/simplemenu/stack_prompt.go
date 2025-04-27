// ui/simplemenu/stack_prompt.go
package simplemenu

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mubbie/stacksmith/internal/core"
	"github.com/mubbie/stacksmith/internal/ui/styles"
)

// StackPromptModel handles user input for the stack command
type StackPromptModel struct {
	BasePrompt
	NewBranch    string
	ParentBranch string
	CurrentField int
	CursorPos    int
}

// NewStackPromptModel creates a new stack prompt model
func NewStackPromptModel() StackPromptModel {
	// Create Git executor
	git := core.NewGitExecutor("")

	// Get current branch as default parent
	defaultParent := "main" // Fallback default
	currentBranch, err := git.GetCurrentBranch()
	if err == nil && currentBranch != "" {
		defaultParent = currentBranch
	}

	return StackPromptModel{
		BasePrompt: BasePrompt{
			Title: "🪵 Create a new branch",
		},
		ParentBranch: defaultParent, // Use current branch as default
		CurrentField: 0,
		CursorPos:    0,
	}
}

func (m StackPromptModel) Init() tea.Cmd {
	return nil
}

func (m StackPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancel()
			return m, tea.Quit

		case "enter":
			if m.CurrentField == 0 {
				if m.NewBranch == "" {
					m.SetError("New branch name cannot be empty")
					return m, nil
				}
				m.CurrentField = 1
				m.CursorPos = len(m.ParentBranch)
				m.ClearError()
				return m, nil
			}

			if m.ParentBranch == "" {
				m.SetError("Parent branch name cannot be empty")
				return m, nil
			}

			// Both fields are filled, time to return
			return m, tea.Quit

		case "tab":
			if m.CurrentField == 0 && m.NewBranch != "" {
				m.CurrentField = 1
				m.CursorPos = len(m.ParentBranch)
			} else if m.CurrentField == 1 && m.ParentBranch != "" {
				m.CurrentField = 0
				m.CursorPos = len(m.NewBranch)
			}
			m.ClearError()
			return m, nil

		case "backspace":
			if m.CurrentField == 0 && len(m.NewBranch) > 0 && m.CursorPos > 0 {
				m.NewBranch = m.NewBranch[:m.CursorPos-1] + m.NewBranch[m.CursorPos:]
				m.CursorPos--
			} else if m.CurrentField == 1 && len(m.ParentBranch) > 0 && m.CursorPos > 0 {
				m.ParentBranch = m.ParentBranch[:m.CursorPos-1] + m.ParentBranch[m.CursorPos:]
				m.CursorPos--
			}
			m.ClearError()
			return m, nil

		case "left":
			if m.CurrentField == 0 && m.CursorPos > 0 {
				m.CursorPos--
			} else if m.CurrentField == 1 && m.CursorPos > 0 {
				m.CursorPos--
			}
			return m, nil

		case "right":
			if m.CurrentField == 0 && m.CursorPos < len(m.NewBranch) {
				m.CursorPos++
			} else if m.CurrentField == 1 && m.CursorPos < len(m.ParentBranch) {
				m.CursorPos++
			}
			return m, nil

		default:
			// Handle regular character input (critical for typing branch names)
			if msg.String() != "tab" && msg.String() != "enter" && len(msg.String()) == 1 {
				if m.CurrentField == 0 {
					m.NewBranch = m.NewBranch[:m.CursorPos] + msg.String() + m.NewBranch[m.CursorPos:]
					m.CursorPos++
				} else if m.CurrentField == 1 {
					m.ParentBranch = m.ParentBranch[:m.CursorPos] + msg.String() + m.ParentBranch[m.CursorPos:]
					m.CursorPos++
				}
				m.ClearError()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m StackPromptModel) View() string {
	s := m.RenderTitle()

	// New branch field
	newBranchLabel := "New branch name: "
	if m.CurrentField == 0 {
		s += styles.Selected.Render(newBranchLabel)

		// Insert cursor at the current position
		if m.CursorPos == len(m.NewBranch) {
			s += m.NewBranch + styles.Normal.Background(lipgloss.Color(styles.ColorSubdued)).Render(" ")
		} else {
			s += m.NewBranch[:m.CursorPos] +
				styles.Normal.Background(lipgloss.Color(styles.ColorSubdued)).Render(string(m.NewBranch[m.CursorPos])) +
				m.NewBranch[m.CursorPos+1:]
		}
	} else {
		s += newBranchLabel + m.NewBranch
	}

	s += "\n"

	// Parent branch field
	parentBranchLabel := "Parent branch: "
	if m.CurrentField == 1 {
		s += styles.Selected.Render(parentBranchLabel)

		// Insert cursor at the current position
		if m.CursorPos == len(m.ParentBranch) {
			s += m.ParentBranch + styles.Normal.Background(lipgloss.Color(styles.ColorSubdued)).Render(" ")
		} else {
			s += m.ParentBranch[:m.CursorPos] +
				styles.Normal.Background(lipgloss.Color(styles.ColorSubdued)).Render(string(m.ParentBranch[m.CursorPos])) +
				m.ParentBranch[m.CursorPos+1:]
		}
	} else {
		s += parentBranchLabel + m.ParentBranch
	}

	// Error message
	s += m.RenderError()

	// Help text
	s += m.RenderHelpText("Tab: Switch fields • Enter: Confirm • Esc: Return to menu")

	return s
}

// RunStackPrompt shows a prompt for the stack command
func RunStackPrompt() (string, string, bool) {
	p := tea.NewProgram(NewStackPromptModel())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return "", "", false
	}

	if m, ok := m.(StackPromptModel); ok {
		// Check if user cancelled
		if m.IsCancelled() {
			return "", "", false
		}

		// Check if both fields are filled
		if m.NewBranch != "" && m.ParentBranch != "" {
			return m.NewBranch, m.ParentBranch, true
		}
	}

	return "", "", false
}
