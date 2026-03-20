package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var updateStateCmd = &cobra.Command{
	Use:    "update-state",
	Short:  "Update pet state (called by git hook)",
	Hidden: true,
	RunE:   runUpdateState,
}

func init() {
	rootCmd.AddCommand(updateStateCmd)
}

func runUpdateState(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		return err
	}

	pet.RecordCommit(state)

	return pet.SaveState(state, statePath)
}
