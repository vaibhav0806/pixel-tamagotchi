package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/config"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
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
	return pet.LoadAndUpdate(config.DefaultStatePath(), func(s *pet.State) {
		pet.RecordCommit(s)
	})
}
