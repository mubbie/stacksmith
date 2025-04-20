package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [branch1] [branch2] ...",
	Short: "🧽 Rebase multiple branches sequentially",
	Long:  `Rebase and push a stack of branches in sequence.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		branches := args
		
		git := core.NewGitExecutor("")
		
		fmt.Println("🧽 Polishing your branch stack... 🪞")
		
		for i := 1; i < len(branches); i++ {
			child := branches[i]
			parent := branches[i-1]
			
			fmt.Printf("🔄 Rebasing %s onto %s...\n", child, parent)
			
			err := git.CheckoutBranch(child)
			if err != nil {
				fmt.Printf("Error checking out %s: %s\n", child, err)
				return
			}
			
			err = git.FetchRemote()
			if err != nil {
				fmt.Printf("Error fetching remote: %s\n", err)
				return
			}
			
			err = git.RebaseBranch(parent)
			if err != nil {
				fmt.Printf("Error rebasing %s onto %s: %s\n", child, parent, err)
				return
			}
			
			err = git.PushBranch()
			if err != nil {
				fmt.Printf("Error pushing %s: %s\n", child, err)
				return
			}
			
			fmt.Printf("🚀 Pushed %s upstream.\n", child)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}