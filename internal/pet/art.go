package pet

import (
	"fmt"
	"math/rand"
	"time"
)

func Render(mood Mood) string {
	face := faceFor(mood)
	paws := pawsFor(mood)
	emoji := mood.Emoji()

	return fmt.Sprintf("   %s\n /\\_/\\\n( %s )\n %s", emoji, face, paws)
}

func faceFor(mood Mood) string {
	switch mood {
	case MoodHappy, MoodHungry:
		return "o.o"
	case MoodSad:
		return "T.T"
	case MoodAsleep:
		return "-.-"
	default:
		return "o.o"
	}
}

func pawsFor(mood Mood) string {
	if mood == MoodHappy {
		return "> ^ <"
	}
	return "> ~ <"
}

var messages = map[Mood][]string{
	MoodHappy: {
		"Pixel is purring!",
		"Pixel loves your commits!",
		"Pixel is chasing a bug (for fun)",
	},
	MoodHungry: {
		"Pixel is eyeing your keyboard...",
		"Pixel meows at you expectantly",
	},
	MoodSad: {
		"Pixel misses you...",
		"Pixel is staring out the terminal window",
	},
	MoodAsleep: {
		"Pixel is curled up sleeping... commit to wake him up",
	},
}

func RenderWithBlink(mood Mood, blinkOpen bool) string {
	face := faceFor(mood)
	if !blinkOpen {
		face = "-.-"
	}
	paws := pawsFor(mood)
	emoji := mood.Emoji()

	return fmt.Sprintf("   %s\n /\\_/\\\n( %s )\n %s", emoji, face, paws)
}

func RandomMessage(mood Mood) string {
	pool := messages[mood]
	if len(pool) == 0 {
		return ""
	}
	return pool[rand.Intn(len(pool))]
}

func ColorForMood(mood Mood) string {
	switch mood {
	case MoodHappy:
		return "#4ade80"
	case MoodHungry:
		return "#facc15"
	case MoodSad:
		return "#60a5fa"
	case MoodAsleep:
		return "#94a3b8"
	default:
		return "#e0e0e0"
	}
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
