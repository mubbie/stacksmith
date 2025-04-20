// render/printer.go
package render

import (
	"fmt"
	"strings"

	"github.com/mubbie/stacksmith/internal/core"
)

// Colors for console output
const (
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Blue   = "\033[0;34m"
	Yellow = "\033[1;33m"
	Cyan   = "\033[0;36m"
	Purple = "\033[0;35m"
	Gray   = "\033[0;37m"
	Reset  = "\033[0m" // No Color
	Bold   = "\033[1m"
)

// Printer provides formatted output methods
type Printer struct {
	AppName string
	Verbose bool
}

// NewPrinter creates a new printer instance
func NewPrinter(appName string) *Printer {
	return &Printer{
		AppName: appName,
		Verbose: false,
	}
}

// SetVerbose sets the verbose flag
func (p *Printer) SetVerbose(verbose bool) {
	p.Verbose = verbose
}

// Success prints a success message with icon
func (p *Printer) Success(message string) {
	fmt.Printf("%s%s%s ✅ %s\n", Green, p.AppName, Reset, message)
}

// Info prints an info message with icon
func (p *Printer) Info(message string) {
	fmt.Printf("%s%s%s ℹ️ %s\n", Blue, p.AppName, Reset, message)
}

// Warning prints a warning message with icon
func (p *Printer) Warning(message string) {
	fmt.Printf("%s%s%s ⚠️ %s\n", Yellow, p.AppName, Reset, message)
}

// Error prints an error message with icon
func (p *Printer) Error(message string) {
	fmt.Printf("%s%s%s ❌ %s\n", Red, p.AppName, Reset, message)
}

// ErrorWithSolution prints an error message with a suggested solution
func (p *Printer) ErrorWithSolution(message, solution string) {
	p.Error(message)
	if solution != "" {
		fmt.Printf("   %s↳ %s%s\n", Yellow, solution, Reset)
	}
}

// HandleGitError intelligently formats and prints git errors
func (p *Printer) HandleGitError(err error) {
	switch e := err.(type) {
	case *core.BranchNotFoundError:
		p.ErrorWithSolution(
			fmt.Sprintf("Branch '%s' not found", e.BranchName),
			"Check the branch name or run 'git branch' to see available branches",
		)
	case *core.MergeConflictError:
		p.ErrorWithSolution(
			fmt.Sprintf("Merge conflict when rebasing %s onto %s", e.Branch, e.Target),
			"Resolve the conflicts, then run 'git rebase --continue'",
		)
	case *core.RemoteError:
		p.ErrorWithSolution(
			fmt.Sprintf("Error communicating with remote '%s'", e.Remote),
			"Check your network connection and remote repository access",
		)
	case *core.GitError:
		// Only show the full error details in verbose mode
		if p.Verbose {
			p.Error(e.Error())
		} else {
			// Extract just the first line for a cleaner message
			lines := strings.Split(e.Stderr, "\n")
			message := "Git error"
			if len(lines) > 0 && lines[0] != "" {
				message = fmt.Sprintf("Git error: %s", lines[0])
			}
			p.Error(message)
			p.Info("Run with --verbose for more details")
		}
	default:
		p.Error(err.Error())
	}
}

// Smith prints a message with smith emoji
func (p *Printer) Smith(message string) {
	fmt.Printf("%s%s%s 🧑🏾‍🏭 %s\n", Green, p.AppName, Reset, message)
}

// ForgeSuccess prints a message for branch creation
func (p *Printer) ForgeSuccess(newBranch, parentBranch string) {
	fmt.Printf("%s%s%s 🪵 Forged new branch %s atop %s. 🌲\n",
		Green, p.AppName, Reset, newBranch, parentBranch)
}

// PushSuccess prints a message for pushing to remote
func (p *Printer) PushSuccess(branch string) {
	fmt.Printf("%s%s%s ⬆️ Lifted %s to remote forge.\n",
		Green, p.AppName, Reset, branch)
}

// NewUpstreamSuccess prints a message for setting new upstream
func (p *Printer) NewUpstreamSuccess(branch string) {
	fmt.Printf("%s%s%s 🆕 First lift for %s — set upstream.\n",
		Green, p.AppName, Reset, branch)
}

// GraphHeader prints header for graph view
func (p *Printer) GraphHeader() {
	fmt.Printf("%s%s%s 🌳 Behold your branching masterpiece:\n",
		Green, p.AppName, Reset)
}

// SyncStart prints start message for sync operation
func (p *Printer) SyncStart() {
	fmt.Printf("%s%s%s 🧽 Polishing your branch stack... 🪞\n",
		Green, p.AppName, Reset)
}

// RebaseStart prints start message for rebasing
func (p *Printer) RebaseStart(child, parent string) {
	fmt.Printf("%s%s%s 🔄 Rebasing %s onto %s...\n",
		Green, p.AppName, Reset, child, parent)
}

// FixPrStart prints start message for fix-pr operation
func (p *Printer) FixPrStart(branch, target string) {
	fmt.Printf("%s%s%s 🔧 Reworking %s onto %s... 🪄\n",
		Green, p.AppName, Reset, branch, target)
}

// RetargetReminder prints reminder to retarget PR
func (p *Printer) RetargetReminder(branch, target string) {
	fmt.Printf("%s%s%s 📢 Don't forget to retarget the PR for %s to %s!\n",
		Yellow, p.AppName, Reset, branch, target)
}

// Divider prints a horizontal divider
func (p *Printer) Divider() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// ComingSoon prints a "coming soon" message
func (p *Printer) ComingSoon(feature string) {
	fmt.Printf("%s%s%s 🔧 %s coming soon...\n",
		Cyan, p.AppName, Reset, feature)
}

// BulletPoint prints a bullet point item
func (p *Printer) BulletPoint(text string) {
	fmt.Printf("  • %s\n", text)
}

// Step prints a numbered step
func (p *Printer) Step(number int, text string) {
	fmt.Printf("  %s%d.%s %s\n", Bold, number, Reset, text)
}

// CommandExample prints an example command
func (p *Printer) CommandExample(command string) {
	fmt.Printf("    %s$ %s%s\n", Gray, command, Reset)
}

// StartProgress starts a progress indicator for a long-running operation
func (p *Printer) StartProgress(operation string) {
	fmt.Printf("%s%s%s ⏳ %s...\n",
		Blue, p.AppName, Reset, operation)
}

// EndProgress ends a progress indicator
func (p *Printer) EndProgress(result string) {
	fmt.Printf("%s%s%s ⌛ %s\n",
		Green, p.AppName, Reset, result)
}
