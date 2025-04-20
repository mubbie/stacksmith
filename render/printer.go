// render/printer.go
package render

import (
	"fmt"
)

// Colors for console output
const (
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Blue   = "\033[0;34m"
	Yellow = "\033[1;33m"
	Reset  = "\033[0m" // No Color
)

// Printer provides formatted output methods
type Printer struct {
	AppName string
}

// NewPrinter creates a new printer instance
func NewPrinter(appName string) *Printer {
	return &Printer{
		AppName: appName,
	}
}

// Success prints a success message with icon
func (p *Printer) Success(message string) {
	fmt.Printf("%s%s%s ✓ %s\n", Green, p.AppName, Reset, message)
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
	fmt.Printf("%s%s%s ✗ %s\n", Red, p.AppName, Reset, message)
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