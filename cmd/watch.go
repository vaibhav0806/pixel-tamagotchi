package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/tui"
)

var moodFlag string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Open Pixel's dashboard",
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().StringVar(&moodFlag, "mood", "", "Override mood for testing (happy, hungry, sad, asleep)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	if moodFlag != "" {
		switch strings.ToLower(moodFlag) {
		case "happy":
			tui.MoodOverride = pet.MoodHappy
		case "hungry":
			tui.MoodOverride = pet.MoodHungry
		case "sad":
			tui.MoodOverride = pet.MoodSad
		case "asleep":
			tui.MoodOverride = pet.MoodAsleep
		default:
			fmt.Printf("Unknown mood %q. Use: happy, hungry, sad, asleep\n", moodFlag)
			return nil
		}
	}

	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
