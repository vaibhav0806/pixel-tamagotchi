package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Wake Pixel up and reset streak",
	RunE:  runReset,
}

func init() {
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' first.")
		return nil
	}

	state.LastCommitAt = time.Now()
	state.CurrentStreak = 0

	if err := pet.SaveState(state, statePath); err != nil {
		return err
	}

	fmt.Println(pet.Render(pet.MoodHappy))
	fmt.Println()
	fmt.Println("Pixel is awake! Streak reset to 0.")

	return nil
}
