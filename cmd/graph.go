// cmd/graph.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/internal/core"
	"github.com/mubbie/stacksmith/internal/render"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "🌳 Show commit graph (git log --graph)",
	Long:  `Visualize branch structure using git log --graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer := render.NewPrinter("stacksmith")
		git := core.NewGitExecutor("")

		printer.GraphHeader()
		printer.Divider()

		// Build the branch stack
		stack, err := git.BuildBranchStack()
		if err != nil {
			printer.Error(fmt.Sprintf("Error analyzing branch structure: %s", err))
			
			// Fall back to traditional git graph if stack analysis fails
			graph, _ := git.ShowGraph()
			fmt.Println(graph)
			return
		}

		// No branches or empty repo
		if len(stack.Roots) == 0 {
			printer.Info("No branches found or unable to determine branch relationships.")
			return
		}

		// Render the branch stack
		branchTree := printer.RenderBranchStack(stack)
		fmt.Println(branchTree)
		
		printer.Divider()
		printer.Info("Legend: " + 
             "* HEAD branch • " + 
             "✔ merged into parent • " + 
             "🔁 (+n/-m) ahead/behind counts")
		printer.Info("Tip: For a more detailed view, try 'stacksmith tui' (coming soon)")
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
