package pet_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/vaibhav/terminal-pet/internal/pet"
)

func TestSaveAndLoadState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	now := time.Now().Truncate(time.Second)
	state := &pet.State{
		LastCommitAt:  now,
		CreatedAt:     now.Add(-48 * time.Hour),
		TotalCommits:  42,
		CurrentStreak: 3,
	}

	if err := pet.SaveState(state, path); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	loaded, err := pet.LoadState(path)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}

	if !loaded.LastCommitAt.Equal(state.LastCommitAt) {
		t.Errorf("LastCommitAt: got %v, want %v", loaded.LastCommitAt, state.LastCommitAt)
	}
	if loaded.TotalCommits != 42 {
		t.Errorf("TotalCommits: got %d, want 42", loaded.TotalCommits)
	}
	if loaded.CurrentStreak != 3 {
		t.Errorf("CurrentStreak: got %d, want 3", loaded.CurrentStreak)
	}
}

func TestLoadState_FileNotFound(t *testing.T) {
	_, err := pet.LoadState("/nonexistent/state.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRecordCommit(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	state := &pet.State{
		LastCommitAt:  now.Add(-2 * time.Hour),
		CreatedAt:     now.Add(-48 * time.Hour),
		TotalCommits:  10,
		CurrentStreak: 1,
	}

	pet.RecordCommit(state)

	if state.TotalCommits != 11 {
		t.Errorf("TotalCommits: got %d, want 11", state.TotalCommits)
	}
	if time.Since(state.LastCommitAt) > time.Second {
		t.Error("LastCommitAt was not updated to now")
	}
}
