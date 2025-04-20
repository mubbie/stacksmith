package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/spf13/cobra"
)

var fixPrCmd = &cobra.Command{
	Use:   "fix-pr [branch] [target]",
	Short: "🔧 Rebase one branch onto a new base",
	Long:  `Rebase a branch onto a new target and remind to retarget the PR.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		target := args[1]
		
		git := core.NewGitExecutor("")
		
		fmt.Printf("🔧 Reworking %s onto %s... 🪄\n", branch, target)
		
		err := git.CheckoutBranch(branch)
		if err != nil {
			fmt.Printf("Error checking out %s: %s\n", branch, err)
			return
		}
		
		err = git.FetchRemote()
		if err != nil {
			fmt.Printf("Error fetching remote: %s\n", err)
			return
		}
		
		// For fix-pr, we rebase onto origin/target
		err = git.RebaseBranch("origin/" + target)
		if err != nil {
			fmt.Printf("Error rebasing onto origin/%s: %s\n", target, err)
			return
		}
		
		err = git.PushBranch()
		if err != nil {
			fmt.Printf("Error pushing %s: %s\n", branch, err)
			return
		}
		
		fmt.Printf("📢 Don't forget to retarget the PR for %s to %s in Azure DevOps, GitHub, or whatever you are using!\n", branch, target)
	},
}

func init() {
	rootCmd.AddCommand(fixPrCmd)
}