package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mubbie/stacksmith/ui/simplemenu"
	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "dev"
	BuildTime = "unknown"
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
	// Set version information
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(`Stacksmith Version: {{.Version}}
Build Time: ` + BuildTime + `
`)

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
	// Loop until the user selects "quit" or exits with Ctrl+C
	for {
		selectedCmd := simplemenu.RunMenu()

		if selectedCmd == "" || selectedCmd == "quit" {
			// Exit the loop and return to the shell
			return
		}

		// Execute the selected command directly - this allows us to return to the menu properly
		// instead of launching as a separate process
		switch selectedCmd {
		case "stack":
			stackCmd.Run(nil, []string{})
		case "sync":
			syncCmd.Run(nil, []string{})
		case "fix-pr":
			fixPrCmd.Run(nil, []string{})
		case "push":
			pushCmd.Run(nil, []string{})
		case "graph":
			graphCmd.Run(nil, []string{})
		case "tui":
			tuiCmd.Run(nil, []string{})
		default:
			fmt.Printf("Unknown command: %s\n", selectedCmd)
		}

		// Ask user whether to return to menu or exit
		fmt.Println()
		fmt.Print("Would you like to return to the Stacksmith menu? [Y/n]: ")
		var response string
		fmt.Scanln(&response)

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "n" || response == "no" {
			return // Exit the function, which exits the application
		}

		fmt.Println() // Add extra newline for spacing
	}
}
