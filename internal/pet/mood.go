package pet

import "time"

type Mood int

const (
	MoodHappy Mood = iota
	MoodHungry
	MoodSad
	MoodAsleep
)

func (m Mood) String() string {
	switch m {
	case MoodHappy:
		return "Happy"
	case MoodHungry:
		return "Hungry"
	case MoodSad:
		return "Sad"
	case MoodAsleep:
		return "Asleep"
	default:
		return "Unknown"
	}
}

func (m Mood) Emoji() string {
	switch m {
	case MoodHappy:
		return "\U0001f496"
	case MoodHungry:
		return "\U0001f355"
	case MoodSad:
		return "\U0001f63f"
	case MoodAsleep:
		return "\U0001f4a4"
	default:
		return "\u2753"
	}
}

func ComputeMood(lastCommitAt time.Time) Mood {
	elapsed := time.Since(lastCommitAt)

	switch {
	case elapsed < 24*time.Hour:
		return MoodHappy
	case elapsed < 48*time.Hour:
		return MoodHungry
	case elapsed < 72*time.Hour:
		return MoodSad
	default:
		return MoodAsleep
	}
}
