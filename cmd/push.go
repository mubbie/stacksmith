package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "⬆️ Smart push with upstream detection",
	Long:  `Push the current branch with upstream handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		git := core.NewGitExecutor("")
		
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %s\n", err)
			return
		}
		
		hasUpstream, err := git.HasUpstream()
		if err != nil {
			fmt.Printf("Error checking upstream: %s\n", err)
			return
		}
		
		if hasUpstream {
			fmt.Printf("⬆️ Lifting %s to remote forge...\n", currentBranch)
			err = git.PushBranch()
			if err != nil {
				fmt.Printf("Error pushing branch: %s\n", err)
				return
			}
		} else {
			fmt.Printf("🆕 First lift for %s — setting upstream...\n", currentBranch)
			err = git.SetUpstreamBranch(currentBranch)
			if err != nil {
				fmt.Printf("Error setting upstream: %s\n", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}