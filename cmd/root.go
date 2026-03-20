package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/config"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
)

// Version is set at build time via ldflags. Falls back to module version
// from debug.ReadBuildInfo() for go install users.
var Version = func() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}()

var rootCmd = &cobra.Command{
	Use:     "pixel",
	Aliases: []string{"pixel-tamagotchi"},
	Short:   "A tamagotchi cat that lives in your terminal",
	Long: `pixel-tamagotchi — a tamagotchi cat that lives in your terminal.

Meet Pixel! She tracks your git commits and reflects your coding
activity as her mood. Commit regularly to keep her happy.

  pixel           Check on Pixel
  pixel init      Set up Pixel and install git hook
  pixel watch     Open the animated dashboard
  pixel reset     Wake Pixel up and reset streak
  pixel uninstall Remove hooks and clean up`,
	Version: Version,
	RunE:    runStatus,
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
		fmt.Println(pet.Render(pet.MoodAsleep))
		fmt.Println()
		fmt.Println("Pixel isn't here yet! Run 'pixel init' to adopt her.")
		return nil
	}

	mood := pet.ComputeMood(state.LastCommitAt)
	art := pet.Render(mood)
	msg := pet.RandomMessage(mood)

	color := lipgloss.Color(pet.ColorForMood(mood))
	style := lipgloss.NewStyle().Foreground(color)

	elapsed := time.Since(state.LastCommitAt).Truncate(time.Minute)
	stats := fmt.Sprintf("Last commit: %s ago | Streak: %d days", pet.FormatDuration(elapsed), state.CurrentStreak)

	fmt.Println(style.Render(art) + "   " + msg)
	fmt.Println(style.Render("          " + stats))

	if state.WelcomeBack {
		wb := lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24")).Bold(true)
		fmt.Println()
		fmt.Println(wb.Render("  🎉 Pixel woke up! She missed you so much!"))
		// Clear the flag
		state.WelcomeBack = false
		_ = pet.SaveState(state, statePath)
	}

	return nil
}
