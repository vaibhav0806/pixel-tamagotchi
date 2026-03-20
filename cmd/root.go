package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var rootCmd = &cobra.Command{
	Use:   "terminal-pet",
	Short: "A tamagotchi cat that lives in your terminal",
	RunE:  runStatus,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' to adopt him.")
		return nil
	}

	mood := pet.ComputeMood(state.LastCommitAt)
	art := pet.Render(mood)
	msg := pet.RandomMessage(mood)

	color := colorForMood(mood)
	style := lipgloss.NewStyle().Foreground(color)

	elapsed := time.Since(state.LastCommitAt).Truncate(time.Minute)
	stats := fmt.Sprintf("Last commit: %s ago | Streak: %d days", formatDuration(elapsed), state.CurrentStreak)

	fmt.Println(style.Render(art) + "   " + msg)
	fmt.Println(style.Render("          " + stats))

	return nil
}

func colorForMood(mood pet.Mood) lipgloss.Color {
	switch mood {
	case pet.MoodHappy:
		return lipgloss.Color("#4ade80")
	case pet.MoodHungry:
		return lipgloss.Color("#facc15")
	case pet.MoodSad:
		return lipgloss.Color("#60a5fa")
	case pet.MoodAsleep:
		return lipgloss.Color("#94a3b8")
	default:
		return lipgloss.Color("#e0e0e0")
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
