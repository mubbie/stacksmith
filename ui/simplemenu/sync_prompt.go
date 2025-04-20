// ui/simplemenu/sync_prompt.go
package simplemenu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mubbie/stacksmith/core"
)

// branchModel represents a branch in the selection list
type branchModel struct {
	name     string
	selected bool
}

// syncPromptModel handles branch selection for the sync command
type syncPromptModel struct {
	availableBranches []branchModel
	selectedBranches  []string
	cursor           int
	err              string
	step             int // 0 for branch selection, 1 for confirmation
	git              *core.GitExecutor
}

func (m syncPromptModel) Init() tea.Cmd {
	return nil
}

func (m syncPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Quietly signal cancellation
			m.selectedBranches = nil
			m.err = "cancelled" // Internal signal, won't be displayed
			return m, tea.Quit

		case "up", "k":
			if m.step == 0 && m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.step == 0 && m.cursor < len(m.availableBranches)-1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			if m.step == 0 {
				// If no branches are selected, select the current one first
				if len(m.selectedBranches) == 0 {
					m.availableBranches[m.cursor].selected = true
					m.selectedBranches = append(m.selectedBranches, m.availableBranches[m.cursor].name)
				}
				
				// Move to confirmation step
				if len(m.selectedBranches) >= 2 {
					m.err = ""
					m.step = 1
				} else {
					m.err = "Please select at least 2 branches to sync"
				}
			} else if m.step == 1 {
				// Confirmation step - we're done
				return m, tea.Quit
			}
			return m, nil

		case " ":
			if m.step == 0 {
				// Toggle selection of current branch
				currentBranch := m.availableBranches[m.cursor]
				currentName := currentBranch.name
				
				if currentBranch.selected {
					// Remove from selected branches
					for i, name := range m.selectedBranches {
						if name == currentName {
							m.selectedBranches = append(m.selectedBranches[:i], m.selectedBranches[i+1:]...)
							break
						}
					}
				} else {
					// Add to selected branches
					m.selectedBranches = append(m.selectedBranches, currentName)
				}
				
				// Update selected state
				m.availableBranches[m.cursor].selected = !currentBranch.selected
			}
			return m, nil
		}
	}
	
	return m, nil
}

func (m syncPromptModel) View() string {
	// Use title styling but ensure no padding or indentation is applied
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffc27d")).
		PaddingLeft(0).  // Explicitly set padding to 0
		MarginLeft(0).   // Explicitly set margin to 0
		Render("🧽 Sync branch stack")
	
	s := title + "\n\n"

	if m.step == 0 {
		// Branch selection screen - keep text aligned to the left with no indentation
		s += "Select branches to sync in order (parent to child):\n"
		s += "(use space to select, enter to continue)\n\n"

		for i, branch := range m.availableBranches {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checkbox := "[ ]"
			if branch.selected {
				checkbox = "[x]"
			}

			// Add order number if selected
			orderInfo := "   "
			for j, name := range m.selectedBranches {
				if name == branch.name {
					orderInfo = fmt.Sprintf(" %d ", j+1)
					break
				}
			}

			branchStyle := lipgloss.NewStyle()
			if m.cursor == i {
				branchStyle = branchStyle.Foreground(lipgloss.Color("#ffc27d")).Bold(true)
			}

			s += fmt.Sprintf("%s %s %s %s\n", 
				lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc27d")).Render(cursor),
				checkbox, 
				orderInfo,
				branchStyle.Render(branch.name))
		}
	} else if m.step == 1 {
		// Confirmation screen
		s += "Ready to sync branches in this order:\n\n"
		
		for i, name := range m.selectedBranches {
			if i > 0 {
				s += fmt.Sprintf("  %d. %s ← %s\n", i, name, m.selectedBranches[i-1])
			} else {
				s += fmt.Sprintf("  %d. %s (base)\n", i+1, name)
			}
		}
		
		s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("Press Enter to confirm and begin sync")
	}

	// Error message (but don't show "cancelled")
	if m.err != "" && m.err != "cancelled" {
		s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(m.err)
	}

	// Help text
	helpText := ""
	if m.step == 0 {
		helpText = "↑/↓: Navigate • Space: Select • Enter: Continue • Esc: Return to menu"
	} else {
		helpText = "Enter: Confirm • Esc: Return to menu"
	}
	
	s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(helpText)

	return s
}

// RunSyncPrompt shows a prompt for the sync command and returns selected branches
func RunSyncPrompt() ([]string, bool) {
	git := core.NewGitExecutor("")
	
	// Get list of branches
	branchesOutput, err := git.Execute("branch")
	if err != nil {
		fmt.Printf("Error getting branches: %s\n", err)
		return nil, false
	}
	
	branchLines := strings.Split(strings.TrimSpace(branchesOutput), "\n")
	var branches []branchModel
	
	// Parse branch names and filter out the current branch
	for _, line := range branchLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			// Current branch - skip for now
			continue
		}
		
		branchName := strings.TrimSpace(line)
		branches = append(branches, branchModel{
			name:     branchName,
			selected: false,
		})
	}
	
	// Add current branch to the beginning
	currentBranch, err := git.GetCurrentBranch()
	if err == nil {
		branches = append([]branchModel{{
			name:     currentBranch,
			selected: false,
		}}, branches...)
	}
	
	// Initialize model
	initialModel := syncPromptModel{
		availableBranches: branches,
		selectedBranches:  []string{},
		cursor:           0,
		step:             0,
		git:              git,
	}
	
	p := tea.NewProgram(initialModel)
	
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return nil, false
	}
	
	if m, ok := m.(syncPromptModel); ok {
		// Check if user cancelled
		if m.err == "cancelled" {
			return nil, false
		}
		
		if m.step == 1 && len(m.selectedBranches) >= 2 {
			return m.selectedBranches, true
		}
	}
	
	return nil, false
}