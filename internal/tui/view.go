package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

func (m Model) View() string {
	color := lipgloss.Color(pet.ColorForMood(m.mood))

	titleStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		MarginBottom(1)

	artStyle := lipgloss.NewStyle().
		Foreground(color)

	messageStyle := lipgloss.NewStyle().
		Foreground(color).
		Italic(true).
		MarginTop(1)

	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		MarginTop(1)

	elapsed := time.Since(m.state.LastCommitAt)
	moodBar := renderMoodBar(elapsed, m.mood, color)

	elapsedStr := pet.FormatDuration(elapsed.Truncate(time.Minute))

	art := pet.RenderWithBlink(m.mood, m.blinkOpen)

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("🐱 Pixel"),
		artStyle.Render(art),
		messageStyle.Render(fmt.Sprintf("%q", m.message)),
		statsStyle.Render(fmt.Sprintf("Mood: %s  |  Streak: %d 🔥", m.mood.String(), m.state.CurrentStreak)),
		statsStyle.Render(fmt.Sprintf("Last commit: %s ago", elapsedStr)),
		moodBar,
		helpStyle.Render("q: quit"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func renderMoodBar(elapsed time.Duration, mood pet.Mood, color lipgloss.Color) string {
	var progress float64
	var label string

	switch mood {
	case pet.MoodHappy:
		progress = 1.0 - (float64(elapsed) / float64(24*time.Hour))
		remaining := 24*time.Hour - elapsed
		label = fmt.Sprintf("%s until hungry", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodHungry:
		progress = 1.0 - (float64(elapsed-24*time.Hour) / float64(24*time.Hour))
		remaining := 48*time.Hour - elapsed
		label = fmt.Sprintf("%s until sad", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodSad:
		progress = 1.0 - (float64(elapsed-48*time.Hour) / float64(24*time.Hour))
		remaining := 72*time.Hour - elapsed
		label = fmt.Sprintf("%s until asleep", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodAsleep:
		progress = 0
		label = "Pixel is asleep... commit to wake him up"
	}

	if progress < 0 {
		progress = 0
	}

	barWidth := 20
	filled := int(progress * float64(barWidth))
	empty := barWidth - filled

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

	bar := ""
	for i := 0; i < filled; i++ {
		bar += filledStyle.Render("█")
	}
	for i := 0; i < empty; i++ {
		bar += emptyStyle.Render("░")
	}

	return fmt.Sprintf("%s %s", bar, label)
}
