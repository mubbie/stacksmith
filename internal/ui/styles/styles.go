// ui/styles/styles.go
package styles

import "github.com/charmbracelet/lipgloss"

// Color constants
const (
	ColorPrimary    = "#ffc27d" // Orange/amber for brand
	ColorSecondary  = "#50fa7b" // Green for success
	ColorError      = "#ff5555" // Red for errors
	ColorSubdued    = "#666666" // Gray for help text
	ColorHighlight  = "#bd93f9" // Purple for highlights
	ColorBackground = "#333333" // Background color for selections
)

// Text styles
var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorPrimary))

	Selected = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	Normal = lipgloss.NewStyle()

	Subdued = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubdued))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary))

	HelpText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSubdued))

	Cursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary))

	Highlight = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorBackground))
)

// CursorStyle returns a styled cursor character
func CursorStyle(active bool) string {
	if active {
		return Cursor.Render(">")
	}
	return " "
}

// FormatHelpText renders help text for command navigation
func FormatHelpText(text string) string {
	return HelpText.Render(text)
}
