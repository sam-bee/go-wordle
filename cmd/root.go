package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wordle",
	Short: "A Wordle player written in Go",
	Long: `This Wordle player is written in Go.
It has knowledge of the valid guesses and solutions for the game, from the New York Times website's Javascript.

Play with "wordle play SPARE", for example.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
