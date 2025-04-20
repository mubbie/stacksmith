package cmd

import (
	"fmt"

	"github.com/mubbie/stacksmith/core"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "🌳 Show commit graph (git log --graph)",
	Long:  `Visualize branch structure using git log --graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		git := core.NewGitExecutor("")
		
		fmt.Println("🌳 Behold your branching masterpiece:")
		graph, err := git.ShowGraph()
		if err != nil {
			fmt.Printf("Error showing graph: %s\n", err)
			return
		}
		
		fmt.Println(graph)
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}