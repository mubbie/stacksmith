package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "🖥 Full-screen DAG browser and stack navigator",
	Long:  `Launch the full-screen terminal UI for managing git branch stacks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This will eventually launch the Bubble Tea TUI
		// For now, just show a placeholder message
		fmt.Println("Full-screen TUI coming in Phase 3...")
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}