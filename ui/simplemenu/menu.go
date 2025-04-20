// ui/simplemenu/menu.go
package simplemenu

import (
	"fmt"
	"os"

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
	s := "🧑‍🏭 Stacksmith\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		title := fmt.Sprintf("%s %s", choice.emoji, choice.title)
		desc := choice.desc

		titleStyle := lipgloss.NewStyle()
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

		if m.cursor == i {
			titleStyle = titleStyle.Foreground(lipgloss.Color("#ffc27d")).Bold(true)
			descStyle = descStyle.Foreground(lipgloss.Color("#999999"))
		}

		choiceLine := fmt.Sprintf("%s %s\n     %s", 
			lipgloss.NewStyle().Foreground(lipgloss.Color("#ffc27d")).Render(cursor),
			titleStyle.Render(title),
			descStyle.Render(desc))
		
		// Add more space between items
		s += choiceLine + "\n\n"
	}

	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("↑/↓: Navigate • Enter: Select • q: Quit")

	// Return the UI as a string without extra newlines
	return s
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