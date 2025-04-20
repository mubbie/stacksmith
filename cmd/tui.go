// cmd/tui.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/internal/render"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "🖥 Full-screen DAG browser and stack navigator",
	Long:  `Launch the full-screen terminal UI for managing git branch stacks.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer := render.NewPrinter("stacksmith")
		
		printer.ComingSoon("Full-screen TUI")
		printer.Divider()
		fmt.Println("The TUI will feature:")
		printer.BulletPoint("Interactive branch DAG visualization")
		printer.BulletPoint("One-key branch operations (rebase, push, retarget)")
		printer.BulletPoint("Visual commit history exploration")
		printer.BulletPoint("Stack health monitoring")
		printer.Divider()
		printer.Info("For now, try 'stacksmith graph' to see a basic branch structure")
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}