// cmd/stack.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/internal/core"
    "github.com/mubbie/stacksmith/internal/render"
    "github.com/mubbie/stacksmith/internal/ui/simplemenu"
	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack [new-branch] [parent-branch]",
	Short: "🪵 Create a new branch atop another",
	Long:  `Forge a new stacked branch on top of an existing parent branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		var newBranch, parentBranch string
		var success bool
		
		printer := render.NewPrinter("stacksmith")
		
		if len(args) < 2 {
			// Not enough arguments, launch the interactive prompt
			newBranch, parentBranch, success = simplemenu.RunStackPrompt()
			if !success {
				// Command was cancelled or failed, just return silently
				return
			}
		} else {
			newBranch = args[0]
			parentBranch = args[1]
		}
		
		git := core.NewGitExecutor("")
		
		err := git.CreateBranch(newBranch, parentBranch)
		if err != nil {
			printer.Error(fmt.Sprintf("Error creating branch: %s", err))
			return
		}
		
		printer.ForgeSuccess(newBranch, parentBranch)
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	rootCmd.AddCommand(stackCmd)
}