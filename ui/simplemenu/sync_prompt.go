// ui/simplemenu/sync_prompt.go
package simplemenu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/ui/styles"
)

type syncStep int

const (
	selectingBranches syncStep = iota
	confirmingSelection
)

// SyncPromptModel handles branch selection for the sync command
type SyncPromptModel struct {
	BasePrompt
	BranchList *SelectableList
	Step       syncStep
	Git        *core.GitExecutor
}

// NewSyncPromptModel creates a new sync prompt model
func NewSyncPromptModel(branches []string, git *core.GitExecutor) SyncPromptModel {
	return SyncPromptModel{
		BasePrompt: BasePrompt{
			Title: "🧽 Sync branch stack",
		},
		BranchList: NewSelectableList(branches),
		Step:       selectingBranches,
		Git:        git,
	}
}

func (m SyncPromptModel) Init() tea.Cmd {
	return nil
}

func (m SyncPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancel()
			return m, tea.Quit

		case "up", "k":
			if m.Step == selectingBranches {
				m.BranchList.MoveUp()
				m.ClearError()
			}
			return m, nil

		case "down", "j":
			if m.Step == selectingBranches {
				m.BranchList.MoveDown()
				m.ClearError()
			}
			return m, nil

		case "enter":
			if m.Step == selectingBranches {
				// If no branches are selected yet, select the current one first
				if m.BranchList.GetSelectedCount() == 0 {
					m.BranchList.ToggleSelected()
				}

				// Validate we have enough branches selected
				if m.BranchList.GetSelectedCount() >= 2 {
					m.ClearError()
					m.Step = confirmingSelection
				} else {
					m.SetError("Please select at least 2 branches to sync")
				}
			} else if m.Step == confirmingSelection {
				// Confirmation step - we're done
				return m, tea.Quit
			}
			return m, nil

		case " ":
			if m.Step == selectingBranches {
				m.BranchList.ToggleSelected()
			}
			return m, nil
		}
	}

	return m, nil
}

func (m SyncPromptModel) View() string {
    s := m.RenderTitle()

    if m.Step == selectingBranches {
        // Branch selection screen
        s += "Select branches to sync in order (parent to child):\n"
        s += "(use space to select, enter to continue)\n\n"
        // Show checkboxes (true) and order numbers (true)
        s += m.BranchList.Render(true, true)
        
    } else if m.Step == confirmingSelection {
        // Confirmation screen
        s += "Ready to sync branches in this order:\n\n"
        
        selectedBranches := m.BranchList.GetSelectedItems()
        
        for i, name := range selectedBranches {
            if i > 0 {
                s += fmt.Sprintf("  %d. %s ← %s\n", i+1, name, selectedBranches[i-1])
            } else {
                s += fmt.Sprintf("  %d. %s (base)\n", i+1, name)
            }
        }
        
        s += "\n" + styles.Success.Render("Press Enter to confirm and begin sync")
    }

    // Error message
    s += m.RenderError()

    // Help text
    helpText := ""
    if m.Step == selectingBranches {
        helpText = "↑/↓: Navigate • Space: Select • Enter: Continue • Esc: Return to menu"
    } else {
        helpText = "Enter: Confirm • Esc: Return to menu"
    }
    
    s += m.RenderHelpText(helpText)

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
	var branches []string

	// Parse branch names and filter out the current branch
	for _, line := range branchLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			// Current branch - mark it but keep it
			line = strings.TrimPrefix(line, "* ")
		}

		branchName := strings.TrimSpace(line)
		branches = append(branches, branchName)
	}

	// Initialize model with branches
	p := tea.NewProgram(NewSyncPromptModel(branches, git))

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return nil, false
	}

	if m, ok := m.(SyncPromptModel); ok {
		// Check if user cancelled
		if m.IsCancelled() {
			return nil, false
		}

		if m.Step == confirmingSelection && m.BranchList.GetSelectedCount() >= 2 {
			return m.BranchList.GetSelectedItems(), true
		}
	}

	return nil, false
}
