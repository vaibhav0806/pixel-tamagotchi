package pet_test

import (
	"strings"
	"testing"

	"github.com/vaibhav/terminal-pet/internal/pet"
)

func TestRenderContainsFace(t *testing.T) {
	tests := []struct {
		mood pet.Mood
		face string
	}{
		{pet.MoodHappy, "o.o"},
		{pet.MoodHungry, "o.o"},
		{pet.MoodSad, "T.T"},
		{pet.MoodAsleep, "-.-"},
	}

	for _, tt := range tests {
		t.Run(tt.mood.String(), func(t *testing.T) {
			art := pet.Render(tt.mood)
			if !strings.Contains(art, tt.face) {
				t.Errorf("render for %v missing face %q:\n%s", tt.mood, tt.face, art)
			}
		})
	}
}

func TestRenderContainsEars(t *testing.T) {
	art := pet.Render(pet.MoodHappy)
	if !strings.Contains(art, `/\_/\`) {
		t.Errorf("render missing ears:\n%s", art)
	}
}

func TestRandomMessage(t *testing.T) {
	msg := pet.RandomMessage(pet.MoodHappy)
	if msg == "" {
		t.Error("expected non-empty message for Happy mood")
	}
}
