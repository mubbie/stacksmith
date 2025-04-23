// ui/simplemenu/menu.go
package simplemenu

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mubbie/stacksmith/internal/ui/styles"
)

type MenuItem struct {
	Title   string
	Desc    string
	Emoji   string
	Command string
}

type MenuModel struct {
	Choices  []MenuItem
	Cursor   int
	Selected string
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter", " ":
			m.Selected = m.Choices[m.Cursor].Command
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	s := styles.Title.Render("🧑‍🏭 Stacksmith") + "\n\n"

	for i, choice := range m.Choices {
		cursor := styles.CursorStyle(i == m.Cursor)

		titleStyle := styles.Normal
		descStyle := styles.Subdued

		if m.Cursor == i {
			titleStyle = styles.Selected
			descStyle = styles.Subdued
		}

		title := fmt.Sprintf("%s %s", choice.Emoji, choice.Title)

		choiceLine := fmt.Sprintf("%s %s\n     %s",
			cursor,
			titleStyle.Render(title),
			descStyle.Render(choice.Desc))

		// Add more space between items
		s += choiceLine + "\n\n"
	}

	s += styles.FormatHelpText("↑/↓: Navigate • Enter: Select • q: Quit")

	return s
}

// RunMenu shows a simple command selection menu and returns the selected command
func RunMenu() string {
	choices := []MenuItem{
		{Title: "Stack", Desc: "Create a new branch atop another", Emoji: "🪵", Command: "stack"},
		{Title: "Sync", Desc: "Rebase multiple branches sequentially", Emoji: "🧽", Command: "sync"},
		{Title: "Fix PR", Desc: "Rebase one branch onto a new base", Emoji: "🔧", Command: "fix-pr"},
		{Title: "Push", Desc: "Smart push with upstream detection", Emoji: "⬆️", Command: "push"},
		{Title: "Graph", Desc: "Show commit graph", Emoji: "🌳", Command: "graph"},
		{Title: "TUI", Desc: "Full-screen DAG browser and stack navigator", Emoji: "🖥", Command: "tui"},
		{Title: "Quit", Desc: "Exit Stacksmith", Emoji: "👋", Command: "quit"},
	}

	p := tea.NewProgram(MenuModel{
		Choices: choices,
	})

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running menu: %v\n", err)
		os.Exit(1)
	}

	if m, ok := m.(MenuModel); ok {
		return m.Selected
	}

	return ""
}
