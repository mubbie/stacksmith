// ui/simplemenu/fixpr_prompt.go
package simplemenu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mubbie/stacksmith/internal/core"
	"github.com/mubbie/stacksmith/internal/ui/styles"
)

type branchSelectionState int

const (
	selectingBranch branchSelectionState = iota
	selectingTarget
	confirming
)

// FixPrPromptModel handles user input for the fix-pr command
type FixPrPromptModel struct {
	BasePrompt
	BranchList     *SelectableList
	TargetList     *SelectableList
	SelectedBranch string
	TargetBranch   string
	State          branchSelectionState
	Git            *core.GitExecutor
}

// NewFixPrPromptModel creates a new fix-pr prompt model
func NewFixPrPromptModel(branches []string, git *core.GitExecutor) FixPrPromptModel {
	return FixPrPromptModel{
		BasePrompt: BasePrompt{
			Title: "🔧 Fix PR branch target",
		},
		BranchList: NewSelectableList(branches),
		TargetList: NewSelectableList(branches),
		State:      selectingBranch,
		Git:        git,
	}
}

func (m FixPrPromptModel) Init() tea.Cmd {
	return nil
}

func (m FixPrPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancel()
			return m, tea.Quit

		case "up", "k":
			if m.State == selectingBranch {
				m.BranchList.MoveUp()
			} else if m.State == selectingTarget {
				m.TargetList.MoveUp()
			}
			m.ClearError()
			return m, nil

		case "down", "j":
			if m.State == selectingBranch {
				m.BranchList.MoveDown()
			} else if m.State == selectingTarget {
				m.TargetList.MoveDown()
			}
			m.ClearError()
			return m, nil

		case "enter":
			switch m.State {
			case selectingBranch:
				m.SelectedBranch = m.BranchList.Items[m.BranchList.Cursor]
				m.State = selectingTarget
				return m, nil

			case selectingTarget:
				m.TargetBranch = m.TargetList.Items[m.TargetList.Cursor]

				// Don't allow a branch to target itself
				if m.TargetBranch == m.SelectedBranch {
					m.SetError("Branch cannot target itself")
					return m, nil
				}

				m.ClearError()
				m.State = confirming
				return m, nil

			case confirming:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m FixPrPromptModel) View() string {
	s := m.RenderTitle()

	switch m.State {
	case selectingBranch:
		s += "Select the branch to retarget:\n\n"
		// Don't show checkboxes (false) or order (false)
		s += m.BranchList.Render(false, false)

	case selectingTarget:
		s += styles.Selected.Render("Selected branch: ") + m.SelectedBranch + "\n\n"
		s += "Select new target branch:\n\n"

		// Special case: highlight the current branch to avoid selecting it
		customList := ""
		for i, branch := range m.TargetList.Items {
			cursor := styles.CursorStyle(i == m.TargetList.Cursor)

			itemStyle := styles.Normal
			if i == m.TargetList.Cursor {
				itemStyle = styles.Selected
			}

			// Gray out the current branch to indicate it's not a valid target
			if branch == m.SelectedBranch {
				itemStyle = styles.Subdued
			}

			customList += fmt.Sprintf("%s %s\n", cursor, itemStyle.Render(branch))
		}
		s += customList

	case confirming:
		s += "Ready to rebase:\n\n"
		s += fmt.Sprintf("  Branch: %s\n", styles.Selected.Render(m.SelectedBranch))
		s += fmt.Sprintf("  New target: %s\n\n", styles.Selected.Render(m.TargetBranch))
		s += styles.Success.Render("Press Enter to confirm and begin rebase")
	}

	// Error message
	s += m.RenderError()

	// Help text
	helpText := ""
	if m.State == confirming {
		helpText = "Enter: Confirm • Esc: Return to menu"
	} else {
		helpText = "↑/↓: Navigate • Enter: Select • Esc: Return to menu"
	}

	s += m.RenderHelpText(helpText)

	return s
}

// RunFixPrPrompt shows a prompt for the fix-pr command
func RunFixPrPrompt() (string, string, bool) {
	git := core.NewGitExecutor("")

	// Get list of branches
	branchesOutput, err := git.Execute("branch")
	if err != nil {
		fmt.Printf("Error getting branches: %v\n", err)
		return "", "", false
	}

	branchLines := strings.Split(strings.TrimSpace(branchesOutput), "\n")
	var branches []string

	// Parse branch names
	for _, line := range branchLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "* ")
		}

		branches = append(branches, strings.TrimSpace(line))
	}

	// Add remote branches as possibilities for targets
	remoteOutput, err := git.Execute("branch", "-r")
	if err == nil {
		remoteLines := strings.Split(strings.TrimSpace(remoteOutput), "\n")
		for _, line := range remoteLines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "origin/") && !strings.Contains(line, "HEAD") {
				remoteBranch := line
				branches = append(branches, remoteBranch)
			}
		}
	}

	p := tea.NewProgram(NewFixPrPromptModel(branches, git))

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return "", "", false
	}

	if m, ok := m.(FixPrPromptModel); ok {
		// Check if user cancelled
		if m.IsCancelled() {
			return "", "", false
		}

		if m.State == confirming && m.SelectedBranch != "" && m.TargetBranch != "" {
			return m.SelectedBranch, m.TargetBranch, true
		}
	}

	return "", "", false
}
