package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mubbie/stacksmith/ui/simplemenu"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stacksmith",
	Short: "🧑🏾‍🏭 Stacksmith - Artisan Git Stacking Tool",
	Long: `Stacksmith is a lightweight, expressive CLI for managing stacked Git branches
			using vanilla Git. Whether you're crafting one-liner PRs or sculpting a majestic stack,
			Stacksmith helps you move fast and stay clean — artisan-style.`,
	Run: func(cmd *cobra.Command, args []string) {
		// When no command is given, launch the Bubble Tea UI menu
		launchMainMenu()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Hide the completion command from help
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	
	// Replace the default help command with our custom one
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "ℹ️ Help about any command",
		Long: `Help provides help for any command in the application.
				Simply type stacksmith help [path to command] for full details.`,
		Run: func(c *cobra.Command, args []string) {
			cmdToShow := rootCmd
			if len(args) > 0 {
				cmdToShow, _, _ = rootCmd.Find(args)
			}
			cmdToShow.HelpFunc()(cmdToShow, args)
		},
	})
}


// launchMainMenu starts the main menu and handles the selected command
func launchMainMenu() {
	selectedCmd := simplemenu.RunMenu()
	
	if selectedCmd == "" || selectedCmd == "quit" {
		return
	}
	
	// Execute the selected command
	cmd := exec.Command(os.Args[0], selectedCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Error executing command: %v\n", err)
	}
}