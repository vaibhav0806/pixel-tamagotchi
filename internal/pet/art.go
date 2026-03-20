package pet

import (
	"fmt"
	"math/rand"
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

func RandomMessage(mood Mood) string {
	pool := messages[mood]
	if len(pool) == 0 {
		return ""
	}
	return pool[rand.Intn(len(pool))]
}
