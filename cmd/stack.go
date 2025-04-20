// cmd/stack.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/ui/simplemenu"
	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack [new-branch] [parent-branch]",
	Short: "🪵 Create a new branch atop another",
	Long:  `Forge a new stacked branch on top of an existing parent branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		var newBranch, parentBranch string
		var success bool
		
		if len(args) < 2 {
			// Not enough arguments, launch the interactive prompt
			newBranch, parentBranch, success = simplemenu.RunStackPrompt()
			if !success {
				return // User cancelled
			}
		} else {
			newBranch = args[0]
			parentBranch = args[1]
		}
		
		git := core.NewGitExecutor("")
		
		err := git.CreateBranch(newBranch, parentBranch)
		if err != nil {
			fmt.Printf("Error creating branch: %s\n", err)
			return
		}
		
		fmt.Printf("🪵 Forged new branch %s atop %s. 🌲\n", newBranch, parentBranch)
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	rootCmd.AddCommand(stackCmd)
}