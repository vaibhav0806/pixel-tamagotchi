package cmd

import "github.com/spf13/cobra"

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check on Pixel (same as running pixel with no args)",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
