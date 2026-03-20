package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
)

const (
	canvasWidth  = 26
	canvasHeight = 8
	catStartRow  = 1
	catStartCol  = 7
)

func (m Model) View() string {
	color := lipgloss.Color(pet.ColorForMood(m.mood))

	titleStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	artStyle := lipgloss.NewStyle().
		Foreground(color)

	messageStyle := lipgloss.NewStyle().
		Foreground(color).
		Italic(true)

	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))

	elapsed := time.Since(m.state.LastCommitAt)

	// Compute progress value and label for mood bar.
	var progressValue float64
	var label string

	switch m.mood {
	case pet.MoodHappy:
		progressValue = 1.0 - (float64(elapsed) / float64(24*time.Hour))
		remaining := 24*time.Hour - elapsed
		label = fmt.Sprintf("%s until hungry", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodHungry:
		progressValue = 1.0 - (float64(elapsed-24*time.Hour) / float64(24*time.Hour))
		remaining := 48*time.Hour - elapsed
		label = fmt.Sprintf("%s until sad", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodSad:
		progressValue = 1.0 - (float64(elapsed-48*time.Hour) / float64(24*time.Hour))
		remaining := 72*time.Hour - elapsed
		label = fmt.Sprintf("%s until asleep", pet.FormatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodAsleep:
		progressValue = 0
		label = "Pixel is asleep... commit to wake him up"
	}

	if progressValue < 0 {
		progressValue = 0
	}

	prog := progress.New(
		progress.WithSolidFill(pet.ColorForMood(m.mood)),
		progress.WithWidth(30),
		progress.WithoutPercentage(),
	)
	moodBar := prog.ViewAs(progressValue) + " " + label

	elapsedStr := pet.FormatDuration(elapsed.Truncate(time.Minute))

	canvas := m.renderCanvas()

	// Compact stats line
	statsLine := fmt.Sprintf("%s  %s  |  Streak: %d 🔥  |  Last commit: %s ago",
		m.mood.Emoji(), m.mood.String(), m.state.CurrentStreak, elapsedStr)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("🐱 Pixel"),
		"",
		artStyle.Render(canvas),
		messageStyle.Render(fmt.Sprintf("%q", m.message)),
		"",
		statsStyle.Render(statsLine),
		moodBar,
		"",
		dimStyle.Render("q: quit"),
	)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(1, 3)

	boxed := borderStyle.Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxed)
}

// renderCanvas builds a character grid, places the cat art and particles,
// then returns it as a single string.
func (m Model) renderCanvas() string {
	// Initialize grid with spaces.
	grid := make([][]string, canvasHeight)
	for y := 0; y < canvasHeight; y++ {
		grid[y] = make([]string, canvasWidth)
		for x := 0; x < canvasWidth; x++ {
			grid[y][x] = " "
		}
	}

	// Track which cells belong to the cat (including internal spaces).
	catZone := make([][]bool, canvasHeight)
	for y := 0; y < canvasHeight; y++ {
		catZone[y] = make([]bool, canvasWidth)
	}

	// Get the current cat frame.
	var catArt string
	if m.earTwitch {
		catArt = pet.EarTwitchFrame(m.mood, m.frame)
	} else {
		frames := pet.AnimationFrames(m.mood)
		frameIdx := m.frame % len(frames)
		catArt = frames[frameIdx]
	}

	// If blinking, override eyes.
	if !m.blinkOpen {
		catArt = overrideEyes(catArt)
	}

	// Place cat art on the grid and mark the cat zone.
	catLines := strings.Split(catArt, "\n")
	for dy, line := range catLines {
		row := catStartRow + dy
		if row >= canvasHeight {
			break
		}
		col := catStartCol
		lineStart := -1
		lineEnd := -1
		for _, r := range line {
			if col >= canvasWidth {
				break
			}
			grid[row][col] = string(r)
			if r != ' ' {
				if lineStart == -1 {
					lineStart = col
				}
				lineEnd = col
			}
			col++
		}
		// Mark the entire span from first non-space to last non-space as cat zone.
		// This prevents particles from landing inside the cat body.
		if lineStart >= 0 && lineEnd >= 0 {
			for x := lineStart; x <= lineEnd; x++ {
				catZone[row][x] = true
			}
		}
	}

	// Place particles on the grid, but only outside the cat zone.
	for _, p := range m.particles.particles {
		if p.x >= 0 && p.x < canvasWidth && p.y >= 0 && p.y < canvasHeight {
			if !catZone[p.y][p.x] {
				grid[p.y][p.x] = p.char
			}
		}
	}

	// Render grid to string.
	var sb strings.Builder
	for y := 0; y < canvasHeight; y++ {
		for x := 0; x < canvasWidth; x++ {
			sb.WriteString(grid[y][x])
		}
		if y < canvasHeight-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// overrideEyes replaces the face expression in a cat art string with closed eyes.
func overrideEyes(art string) string {
	// Simple string replacement for known face patterns.
	result := art
	for _, face := range []string{"o.o", "^.^", "T.T", ";_;", ">.<"} {
		result = strings.Replace(result, face, "-.-", 1)
	}
	return result
}
