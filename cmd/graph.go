// cmd/graph.go
package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/mubbie/stacksmith/render"
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
		
		graph, err := git.ShowGraph()
		if err != nil {
			printer.Error(fmt.Sprintf("Error showing graph: %s", err))
			return
		}
		
		fmt.Println(graph)
		printer.Divider()
		printer.Info("Tip: For a more detailed view, try 'stacksmith tui' (coming soon)")
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}