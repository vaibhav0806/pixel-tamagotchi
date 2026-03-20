package tui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

// Message types for animation timing.
type blinkMsg struct{}
type messageRotateMsg struct{}
type animTickMsg struct{}
type frameAdvanceMsg struct{}

// Particle represents a floating character in the animation canvas.
type Particle struct {
	x    int
	y    int
	char string
	life int
}

// ParticleSystem manages spawning, updating, and removing particles.
type ParticleSystem struct {
	particles []Particle
	mood      pet.Mood
}

// Mood-specific particle characters.
var moodParticles = map[pet.Mood][]string{
	pet.MoodHappy:  {"♡", "✦", "♪", "·"},
	pet.MoodHungry: {"·", "?"},
	pet.MoodSad:    {"·", ",", "'", "."},
	pet.MoodAsleep: {"z", "Z", "·"},
}

// Mood-specific spawn rates (particles per tick).
var moodSpawnRate = map[pet.Mood]int{
	pet.MoodHappy:  2,
	pet.MoodHungry: 1,
	pet.MoodSad:    1,
	pet.MoodAsleep: 1,
}

// NewParticleSystem creates a particle system for the given mood.
func NewParticleSystem(mood pet.Mood) ParticleSystem {
	return ParticleSystem{mood: mood}
}

// Update advances all particles by one tick: moves them, ages them,
// removes dead ones, and spawns new ones.
func (ps *ParticleSystem) Update() {
	// Move and age existing particles.
	alive := ps.particles[:0]
	for i := range ps.particles {
		p := &ps.particles[i]
		p.life--
		if p.life <= 0 {
			continue
		}
		// Sad tears fall down; everything else floats up.
		if ps.mood == pet.MoodSad {
			p.y++
		} else {
			p.y--
		}
		// Small horizontal drift.
		p.x += rand.Intn(3) - 1
		alive = append(alive, *p)
	}
	ps.particles = alive

	// Spawn new particles.
	chars := moodParticles[ps.mood]
	if len(chars) == 0 {
		return
	}
	rate := moodSpawnRate[ps.mood]
	for i := 0; i < rate; i++ {
		// Spawn within a ~18-char band centered on the cat in a 26-wide canvas.
		spawnX := 3 + rand.Intn(20)
		// Spawn near the cat vertically in an 8-tall canvas (cat at rows 1-5).
		var spawnY int
		if ps.mood == pet.MoodSad {
			spawnY = 2 + rand.Intn(2) // start near cat face, fall down
		} else {
			spawnY = 5 + rand.Intn(2) // start below the cat body, float up
		}
		ps.particles = append(ps.particles, Particle{
			x:    spawnX,
			y:    spawnY,
			char: chars[rand.Intn(len(chars))],
			life: 5 + rand.Intn(4), // 5-8 ticks
		})
	}
}

// SetMood changes the mood and clears existing particles.
func (ps *ParticleSystem) SetMood(mood pet.Mood) {
	if ps.mood != mood {
		ps.mood = mood
		ps.particles = nil
	}
}

// Commands for scheduling ticks.

func blinkCmd() tea.Cmd {
	delay := time.Duration(3+rand.Intn(5)) * time.Second
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func messageRotateCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return messageRotateMsg{}
	})
}

func animTickCmd() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return animTickMsg{}
	})
}

func frameAdvanceCmd() tea.Cmd {
	return tea.Tick(600*time.Millisecond, func(t time.Time) tea.Msg {
		return frameAdvanceMsg{}
	})
}

// Handlers called from Update.

func handleBlink(m *Model) tea.Cmd {
	m.blinkOpen = !m.blinkOpen
	if !m.blinkOpen {
		return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
			return blinkMsg{}
		})
	}
	return blinkCmd()
}

func handleMessageRotate(m *Model) tea.Cmd {
	m.message = pet.RandomMessage(m.mood)
	return messageRotateCmd()
}

func handleAnimTick(m *Model) tea.Cmd {
	m.animTick++
	m.particles.Update()
	return animTickCmd()
}

func handleFrameAdvance(m *Model) tea.Cmd {
	frames := pet.AnimationFrames(m.mood)
	if len(frames) > 0 {
		m.frame = (m.frame + 1) % len(frames)
	}
	// Occasionally trigger an ear twitch (~10% chance).
	if rand.Intn(10) == 0 {
		m.earTwitch = true
	} else {
		m.earTwitch = false
	}
	return frameAdvanceCmd()
}
