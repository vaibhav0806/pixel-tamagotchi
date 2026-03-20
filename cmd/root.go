package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terminal-pet",
	Short: "A tamagotchi cat that lives in your terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Pixel says hi! Run 'terminal-pet init' to get started.")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
