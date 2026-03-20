package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/config"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show Pixel's lifetime stats",
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	state, err := pet.LoadState(config.DefaultStatePath())
	if err != nil {
		fmt.Println("Pixel isn't here yet! Run 'pixel init' first.")
		return nil
	}

	mood := pet.ComputeMood(state.LastCommitAt)
	color := lipgloss.Color(pet.ColorForMood(mood))

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		MarginBottom(1)

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(20)

	value := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	dim := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))

	// Calculate days since adoption
	daysSinceAdoption := int(time.Since(state.CreatedAt).Hours() / 24)
	if daysSinceAdoption == 0 {
		daysSinceAdoption = 1 // at least today
	}

	// Commits per day
	commitsPerDay := float64(state.TotalCommits) / float64(daysSinceAdoption)

	// Time since last commit
	elapsed := time.Since(state.LastCommitAt)
	lastCommitStr := pet.FormatDuration(elapsed.Truncate(time.Minute))

	// Best streak — use current if it's the best
	bestStreak := state.BestStreak
	if state.CurrentStreak > bestStreak {
		bestStreak = state.CurrentStreak
	}

	fmt.Println()
	fmt.Println(title.Render("📊 Pixel's Stats"))
	fmt.Println()
	fmt.Println(label.Render("  Mood") + value.Render(fmt.Sprintf("%s %s", mood.Emoji(), mood.String())))
	fmt.Println(label.Render("  Last commit") + value.Render(lastCommitStr+" ago"))
	fmt.Println()
	fmt.Println(label.Render("  Total commits") + value.Render(fmt.Sprintf("%d", state.TotalCommits)))
	fmt.Println(label.Render("  Current streak") + value.Render(fmt.Sprintf("%d days 🔥", state.CurrentStreak)))
	fmt.Println(label.Render("  Best streak") + value.Render(fmt.Sprintf("%d days 🏆", bestStreak)))
	fmt.Println(label.Render("  Commits/day") + value.Render(fmt.Sprintf("%.1f", commitsPerDay)))
	fmt.Println()
	fmt.Println(label.Render("  Adopted") + dim.Render(state.CreatedAt.Local().Format("Jan 2, 2006")))
	fmt.Println(label.Render("  Days together") + dim.Render(fmt.Sprintf("%d", daysSinceAdoption)))
	fmt.Println()

	return nil
}
