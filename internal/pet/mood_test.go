package pet_test

import (
	"testing"
	"time"

	"github.com/vaibhav/terminal-pet/internal/pet"
)

func TestComputeMood(t *testing.T) {
	tests := []struct {
		name     string
		ago      time.Duration
		expected pet.Mood
	}{
		{"just committed", 1 * time.Hour, pet.MoodHappy},
		{"23 hours ago", 23 * time.Hour, pet.MoodHappy},
		{"25 hours ago", 25 * time.Hour, pet.MoodHungry},
		{"47 hours ago", 47 * time.Hour, pet.MoodHungry},
		{"49 hours ago", 49 * time.Hour, pet.MoodSad},
		{"71 hours ago", 71 * time.Hour, pet.MoodSad},
		{"73 hours ago", 73 * time.Hour, pet.MoodAsleep},
		{"a week ago", 168 * time.Hour, pet.MoodAsleep},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastCommit := time.Now().Add(-tt.ago)
			mood := pet.ComputeMood(lastCommit)
			if mood != tt.expected {
				t.Errorf("got %v, want %v", mood, tt.expected)
			}
		})
	}
}

func TestMoodString(t *testing.T) {
	if pet.MoodHappy.String() != "Happy" {
		t.Errorf("got %q, want Happy", pet.MoodHappy.String())
	}
}
