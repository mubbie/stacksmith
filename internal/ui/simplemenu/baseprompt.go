package simplemenu

import (
	"github.com/mubbie/stacksmith/internal/ui/styles"
)

// BasePrompt provides common functionality for all prompts
type BasePrompt struct {
	Title     string
	ErrorMsg  string
	Cancelled bool
}

// RenderTitle renders the title with consistent styling
func (b *BasePrompt) RenderTitle() string {
	return styles.Title.Render(b.Title) + "\n\n"
}

// RenderError renders error message if present
func (b *BasePrompt) RenderError() string {
	if b.ErrorMsg != "" && !b.Cancelled {
		return "\n\n" + styles.Error.Render(b.ErrorMsg)
	}
	return ""
}

// RenderHelpText renders help text with consistent styling
func (b *BasePrompt) RenderHelpText(text string) string {
	return "\n\n" + styles.HelpText.Render(text)
}

// Cancel marks the prompt as cancelled
func (b *BasePrompt) Cancel() {
	b.Cancelled = true
	b.ErrorMsg = "cancelled" // Internal signal
}

// IsCancelled returns true if the prompt was cancelled
func (b *BasePrompt) IsCancelled() bool {
	return b.Cancelled
}

// SetError sets an error message
func (b *BasePrompt) SetError(msg string) {
	b.ErrorMsg = msg
}

// ClearError clears the error message
func (b *BasePrompt) ClearError() {
	b.ErrorMsg = ""
}
