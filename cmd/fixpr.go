// cmd/fixpr.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/render"
	"github.com/mubbie/stacksmith/ui/simplemenu"
	"github.com/spf13/cobra"
)

var fixPrCmd = &cobra.Command{
	Use:   "fix-pr [branch] [target]",
	Short: "🔧 Rebase one branch onto a new base",
	Long:  `Rebase a branch onto a new target and remind to retarget the PR.`,
	Run: func(cmd *cobra.Command, args []string) {
		var branch, target string
		var success bool

		printer := render.NewPrinter("stacksmith")

		if len(args) < 2 {
			// Not enough arguments, launch the interactive prompt
			branch, target, success = simplemenu.RunFixPrPrompt()
			if !success {
				// Command was cancelled or failed, just return silently
				return
			}
		} else {
			branch = args[0]
			target = args[1]
		}

		git := core.NewGitExecutor("")

		printer.FixPrStart(branch, target)

		err := git.CheckoutBranch(branch)
		if err != nil {
			printer.Error(fmt.Sprintf("Error checking out %s: %s", branch, err))
			return
		}

		err = git.FetchRemote()
		if err != nil {
			printer.Error(fmt.Sprintf("Error fetching remote: %s", err))
			return
		}

		// For fix-pr, we rebase onto the target (which might be a remote branch)
		rebaseTarget := target
		if !strings.HasPrefix(target, "origin/") {
			rebaseTarget = "origin/" + target
		}

		err = git.RebaseBranch(rebaseTarget)
		if err != nil {
			printer.Error(fmt.Sprintf("Error rebasing onto %s: %s", rebaseTarget, err))
			return
		}

		err = git.PushBranch()
		if err != nil {
			printer.Error(fmt.Sprintf("Error pushing %s: %s", branch, err))
			return
		}

		printer.Success(fmt.Sprintf("Successfully rebased %s onto %s", branch, target))
		printer.RetargetReminder(branch, target)
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	rootCmd.AddCommand(fixPrCmd)
}
