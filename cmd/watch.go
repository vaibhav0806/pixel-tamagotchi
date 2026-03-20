package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
	"github.com/vaibhav/terminal-pet/internal/tui"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Open Pixel's dashboard",
	RunE:  runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()
	if _, err := pet.LoadState(statePath); err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' first.")
		return nil
	}

	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
