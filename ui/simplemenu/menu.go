// ui/simplemenu/menu.go
package simplemenu

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc, emoji, command string
}

type model struct {
	choices  []item
	cursor   int
	selected string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.choices[m.cursor].command
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	var builder strings.Builder

	builder.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")). // purple-pinkish
		Render("🧑‍🏭 Stacksmith") + "\n\n")

	for i, choice := range m.choices {
		isSelected := m.cursor == i

		cursor := " "
		if isSelected {
			cursor = ">"
		}

		// Styles
		cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc27d"))
		titleStyle := lipgloss.NewStyle()
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

		if isSelected {
			titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("#ffc27d"))
			descStyle = descStyle.Foreground(lipgloss.Color("#999999"))
		}

		// Build lines
		title := fmt.Sprintf("%s %s", choice.emoji, choice.title)
		titleLine := lipgloss.JoinHorizontal(1,
			cursorStyle.Render(cursor),
			titleStyle.Render(title),
		)

		descLine := lipgloss.NewStyle().
			MarginLeft(4). // indent description
			Render(descStyle.Render(choice.desc))

		// Combine with vertical spacing
		itemBlock := lipgloss.JoinVertical(lipgloss.Top, titleLine, descLine)
		builder.WriteString(itemBlock + "\n\n") // add spacing between items
	}

	builder.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Render("↑/↓: Navigate • Enter: Select • q: Quit"),
	)

	return builder.String()
}

// RunMenu shows a simple command selection menu and returns the selected command
func RunMenu() string {
	choices := []item{
		{title: "Stack", desc: "Create a new branch atop another", emoji: "🪵", command: "stack"},
		{title: "Sync", desc: "Rebase multiple branches sequentially", emoji: "🧽", command: "sync"},
		{title: "Fix PR", desc: "Rebase one branch onto a new base", emoji: "🔧", command: "fix-pr"},
		{title: "Push", desc: "Smart push with upstream detection", emoji: "⬆️", command: "push"},
		{title: "Graph", desc: "Show commit graph", emoji: "🌳", command: "graph"},
		{title: "TUI", desc: "Full-screen DAG browser and stack navigator", emoji: "🖥", command: "tui"},
		{title: "Quit", desc: "Exit Stacksmith", emoji: "👋", command: "quit"},
	}

	p := tea.NewProgram(model{
		choices: choices,
	})

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running menu: %v\n", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok {
		return m.selected
	}

	return ""
}