package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

type tickMsg time.Time
type stateRefreshMsg struct{ state *pet.State }

type Model struct {
	state     *pet.State
	statePath string
	mood      pet.Mood
	message   string
	blinkOpen bool
	width     int
	height    int

	// Animation state.
	frame     int
	animTick  int
	earTwitch bool
	particles ParticleSystem
}

func NewModel() Model {
	statePath := config.DefaultStatePath()
	state, err := pet.LoadState(statePath)
	if err != nil {
		state = &pet.State{
			CreatedAt:    time.Now(),
			LastCommitAt: time.Now(),
		}
	}

	mood := pet.ComputeMood(state.LastCommitAt)

	return Model{
		state:     state,
		statePath: statePath,
		mood:      mood,
		message:   pet.RandomMessage(mood),
		blinkOpen: true,
		particles: NewParticleSystem(mood),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(time.Second),
		refreshState(m.statePath),
		blinkCmd(),
		messageRotateCmd(),
		animTickCmd(),
		frameAdvanceCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		newMood := pet.ComputeMood(m.state.LastCommitAt)
		if newMood != m.mood {
			m.mood = newMood
			m.particles.SetMood(m.mood)
			m.frame = 0
		}
		return m, tickEvery(time.Second)

	case stateRefreshMsg:
		if msg.state != nil {
			oldMood := m.mood
			m.state = msg.state
			m.mood = pet.ComputeMood(m.state.LastCommitAt)
			if m.mood != oldMood {
				m.message = pet.RandomMessage(m.mood)
				m.particles.SetMood(m.mood)
				m.frame = 0
			}
		}
		return m, refreshStateAfter(m.statePath, 3*time.Second)

	case blinkMsg:
		return m, handleBlink(&m)

	case messageRotateMsg:
		return m, handleMessageRotate(&m)

	case animTickMsg:
		return m, handleAnimTick(&m)

	case frameAdvanceMsg:
		return m, handleFrameAdvance(&m)
	}

	return m, nil
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func refreshState(path string) tea.Cmd {
	return func() tea.Msg {
		state, _ := pet.LoadState(path)
		return stateRefreshMsg{state: state}
	}
}

func refreshStateAfter(path string, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		state, _ := pet.LoadState(path)
		return stateRefreshMsg{state: state}
	})
}
