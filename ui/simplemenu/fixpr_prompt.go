// ui/simplemenu/fixpr_prompt.go
package simplemenu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mubbie/stacksmith/core"
)

type branchSelectionState int

const (
	selectingBranch branchSelectionState = iota
	selectingTarget
	confirming
)

// fixPrPromptModel handles user input for the fix-pr command
type fixPrPromptModel struct {
	availableBranches []string
	selectedBranch    string
	targetBranch      string
	cursor            int
	err               string
	state             branchSelectionState
	git               *core.GitExecutor
}

func (m fixPrPromptModel) Init() tea.Cmd {
	return nil
}

func (m fixPrPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Quietly signal cancellation
			m.selectedBranch = ""
			m.targetBranch = ""
			m.err = "cancelled" // Internal signal, won't be displayed
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.availableBranches)-1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			switch m.state {
			case selectingBranch:
				m.selectedBranch = m.availableBranches[m.cursor]
				m.cursor = 0 // Reset cursor for target selection
				m.state = selectingTarget
				return m, nil

			case selectingTarget:
				m.targetBranch = m.availableBranches[m.cursor]

				// Don't allow a branch to target itself
				if m.targetBranch == m.selectedBranch {
					m.err = "Branch cannot target itself"
					return m, nil
				}

				m.err = ""
				m.state = confirming
				return m, nil

			case confirming:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m fixPrPromptModel) View() string {
	// Use title styling but ensure no padding or indentation is applied
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffc27d")).
		PaddingLeft(0). // Explicitly set padding to 0
		MarginLeft(0).  // Explicitly set margin to 0
		Render("🔧 Fix PR branch target")

	s := title + "\n\n"

	switch m.state {
	case selectingBranch:
		s += "Select the branch to retarget:\n\n"

		for i, branch := range m.availableBranches {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			branchStyle := lipgloss.NewStyle()
			if m.cursor == i {
				branchStyle = branchStyle.Foreground(lipgloss.Color("#ffc27d")).Bold(true)
			}

			s += fmt.Sprintf("%s %s\n",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc27d")).Render(cursor),
				branchStyle.Render(branch))
		}

	case selectingTarget:
		s += fmt.Sprintf("Selected branch: %s\n\n",
			lipgloss.NewStyle().Bold(true).Render(m.selectedBranch))
		s += "Select new target branch:\n\n"

		for i, branch := range m.availableBranches {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			branchStyle := lipgloss.NewStyle()
			if m.cursor == i {
				branchStyle = branchStyle.Foreground(lipgloss.Color("#ffc27d")).Bold(true)
			}

			// Highlight currently selected branch
			if branch == m.selectedBranch {
				branchStyle = branchStyle.Foreground(lipgloss.Color("#999999"))
			}

			s += fmt.Sprintf("%s %s\n",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc27d")).Render(cursor),
				branchStyle.Render(branch))
		}

	case confirming:
		s += "Ready to rebase:\n\n"
		s += fmt.Sprintf("  Branch: %s\n",
			lipgloss.NewStyle().Bold(true).Render(m.selectedBranch))
		s += fmt.Sprintf("  New target: %s\n\n",
			lipgloss.NewStyle().Bold(true).Render(m.targetBranch))

		s += lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("Press Enter to confirm and begin rebase")
	}

	// Error message (but don't show "cancelled")
	if m.err != "" && m.err != "cancelled" {
		s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(m.err)
	}

	// Help text
	helpText := ""
	if m.state == confirming {
		helpText = "Enter: Confirm • Esc: Return to menu"
	} else {
		helpText = "↑/↓: Navigate • Enter: Select • Esc: Return to menu"
	}

	s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(helpText)

	return s
}

// RunFixPrPrompt shows a prompt for the fix-pr command
func RunFixPrPrompt() (string, string, bool) {
	git := core.NewGitExecutor("")

	// Get list of branches
	branchesOutput, err := git.Execute("branch")
	if err != nil {
		fmt.Printf("Error getting branches: %s\n", err)
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

	initialModel := fixPrPromptModel{
		availableBranches: branches,
		cursor:            0,
		state:             selectingBranch,
		git:               git,
	}

	p := tea.NewProgram(initialModel)

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return "", "", false
	}

	if m, ok := m.(fixPrPromptModel); ok {
		// Check if user cancelled
		if m.err == "cancelled" {
			return "", "", false
		}

		if m.state == confirming && m.selectedBranch != "" && m.targetBranch != "" {
			return m.selectedBranch, m.targetBranch, true
		}
	}

	return "", "", false
}
