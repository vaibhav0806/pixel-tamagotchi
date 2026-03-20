package pet

import (
	"encoding/json"
	"os"
	"time"
)

type State struct {
	LastCommitAt  time.Time `json:"last_commit_at"`
	CreatedAt     time.Time `json:"created_at"`
	TotalCommits  int       `json:"total_commits"`
	CurrentStreak int       `json:"current_streak"`
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveState(s *State, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func RecordCommit(s *State) {
	now := time.Now()

	lastDay := s.LastCommitAt.Local().Truncate(24 * time.Hour)
	today := now.Local().Truncate(24 * time.Hour)
	diff := today.Sub(lastDay)

	switch {
	case diff == 24*time.Hour:
		s.CurrentStreak++
	case diff == 0:
		// same day, streak unchanged
	default:
		s.CurrentStreak = 1
	}

	s.LastCommitAt = now
	s.TotalCommits++
}
