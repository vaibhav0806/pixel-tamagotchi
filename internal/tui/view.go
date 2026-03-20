package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

const (
	canvasWidth  = 25
	canvasHeight = 12
	// Cat art is placed starting at this row/col in the canvas.
	catStartRow = 3
	catStartCol = 7
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

	canvas := m.renderCanvas()

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("🐱 Pixel"),
		artStyle.Render(canvas),
		messageStyle.Render(fmt.Sprintf("%q", m.message)),
		statsStyle.Render(fmt.Sprintf("Mood: %s  |  Streak: %d 🔥", m.mood.String(), m.state.CurrentStreak)),
		statsStyle.Render(fmt.Sprintf("Last commit: %s ago", elapsedStr)),
		moodBar,
		helpStyle.Render("q: quit"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
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

	// Place cat art on the grid.
	catLines := strings.Split(catArt, "\n")
	for dy, line := range catLines {
		row := catStartRow + dy
		if row >= canvasHeight {
			break
		}
		col := catStartCol
		for _, r := range line {
			if col >= canvasWidth {
				break
			}
			grid[row][col] = string(r)
			col++
		}
	}

	// Place particles on the grid (particles behind cat art stay hidden
	// because we place particles after, but only on empty spaces).
	for _, p := range m.particles.particles {
		if p.x >= 0 && p.x < canvasWidth && p.y >= 0 && p.y < canvasHeight {
			// Only place particle if the cell is a space (don't overwrite cat).
			if grid[p.y][p.x] == " " {
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
	lines := strings.Split(art, "\n")
	for i, line := range lines {
		// The face line contains "( ... )" pattern.
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			openParen := strings.Index(line, "(")
			closeParen := strings.LastIndex(line, ")")
			if closeParen > openParen {
				// Count runes between parens to maintain width.
				inner := line[openParen+1 : closeParen]
				runeCount := utf8.RuneCountInString(inner)
				// Build replacement: space + -.- + spaces to fill.
				var face string
				if runeCount >= 5 {
					padding := runeCount - 5
					face = " -.- " + strings.Repeat(" ", padding)
				} else {
					face = " -.- "[:runeCount]
				}
				lines[i] = line[:openParen+1] + face + line[closeParen:]
			}
		}
	}
	return strings.Join(lines, "\n")
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
