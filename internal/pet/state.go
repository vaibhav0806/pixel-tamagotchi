package pet

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LastCommitAt  time.Time `json:"last_commit_at"`
	CreatedAt     time.Time `json:"created_at"`
	TotalCommits  int       `json:"total_commits"`
	CurrentStreak int       `json:"current_streak"`
	BestStreak    int       `json:"best_streak"`
	WelcomeBack   bool      `json:"welcome_back,omitempty"`
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
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func RecordCommit(s *State) {
	now := time.Now()

	// Check if Pixel was sad or asleep before this commit
	oldMood := ComputeMood(s.LastCommitAt)
	if oldMood == MoodSad || oldMood == MoodAsleep {
		s.WelcomeBack = true
	}

	lastYear, lastMonth, lastDay := s.LastCommitAt.Local().Date()
	lastMidnight := time.Date(lastYear, lastMonth, lastDay, 0, 0, 0, 0, time.Local)

	nowYear, nowMonth, nowDay := now.Local().Date()
	todayMidnight := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local)

	daysDiff := int(todayMidnight.Sub(lastMidnight).Hours() / 24)

	switch {
	case daysDiff == 1:
		s.CurrentStreak++
	case daysDiff == 0:
		// same day, streak unchanged
	default:
		s.CurrentStreak = 1
	}

	if s.CurrentStreak > s.BestStreak {
		s.BestStreak = s.CurrentStreak
	}

	s.LastCommitAt = now
	s.TotalCommits++
}

func LoadAndUpdate(path string, fn func(*State)) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := lockFile(f); err != nil {
		return err
	}
	defer unlockFile(f)

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	fn(&s)

	out, err := json.MarshalIndent(&s, "", "  ")
	if err != nil {
		return err
	}

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	_, err = f.Write(out)
	return err
}
