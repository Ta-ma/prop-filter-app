/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ta-ma/prop-filter-app/internal/config"
)

var cfg *config.Cli

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prop-filter-app",
	Short: "Command-line tool to filter real estate data",
	Long: `This application allows you to query real estate data using a variety
of operators. Use the command query -h for more information.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(config *config.Cli) {
	cfg = config
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
