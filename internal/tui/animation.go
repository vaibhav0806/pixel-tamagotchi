package tui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

type blinkMsg struct{}
type messageRotateMsg struct{}

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
