// cmd/sync.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/render"
	"github.com/mubbie/stacksmith/ui/simplemenu"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [branch1] [branch2] ...",
	Short: "🧽 Rebase multiple branches sequentially",
	Long:  `Rebase and push a stack of branches in sequence.`,
	Run: func(cmd *cobra.Command, args []string) {
		var branches []string
		var success bool
		
		printer := render.NewPrinter("stacksmith")
		
		if len(args) < 2 {
			// Not enough arguments, launch the interactive prompt
			branches, success = simplemenu.RunSyncPrompt()
			if !success {
				// Command was cancelled or failed, just return silently
				return
			}
		} else {
			branches = args
		}
		
		git := core.NewGitExecutor("")
		
		printer.SyncStart()
		
		for i := 1; i < len(branches); i++ {
			child := branches[i]
			parent := branches[i-1]
			
			printer.RebaseStart(child, parent)
			
			err := git.CheckoutBranch(child)
			if err != nil {
				printer.Error(fmt.Sprintf("Error checking out %s: %s", child, err))
				return
			}
			
			err = git.FetchRemote()
			if err != nil {
				printer.Error(fmt.Sprintf("Error fetching remote: %s", err))
				return
			}
			
			err = git.RebaseBranch(parent)
			if err != nil {
				printer.Error(fmt.Sprintf("Error rebasing %s onto %s: %s", child, parent, err))
				return
			}
			
			err = git.PushBranch()
			if err != nil {
				printer.Error(fmt.Sprintf("Error pushing %s: %s", child, err))
				return
			}
			
			printer.PushSuccess(child)
		}
		
		printer.Success("Stack sync complete!")
	},
	Args: cobra.MaximumNArgs(100), // Allow multiple branches
}

func init() {
	rootCmd.AddCommand(syncCmd)
}