// cmd/push.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/render"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "⬆️ Smart push with upstream detection",
	Long:  `Push the current branch with upstream handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer := render.NewPrinter("stacksmith")
		git := core.NewGitExecutor("")
		
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			printer.Error(fmt.Sprintf("Error getting current branch: %s", err))
			return
		}
		
		hasUpstream, err := git.HasUpstream()
		if err != nil {
			printer.Error(fmt.Sprintf("Error checking upstream: %s", err))
			return
		}
		
		// Show a spinner or progress indicator
		printer.Info(fmt.Sprintf("Pushing branch %s...", currentBranch))
		
		if hasUpstream {
			err = git.PushBranch()
			if err != nil {
				printer.Error(fmt.Sprintf("Error pushing branch: %s", err))
				return
			}
			printer.PushSuccess(currentBranch)
		} else {
			err = git.SetUpstreamBranch(currentBranch)
			if err != nil {
				printer.Error(fmt.Sprintf("Error setting upstream: %s", err))
				return
			}
			printer.NewUpstreamSuccess(currentBranch)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}