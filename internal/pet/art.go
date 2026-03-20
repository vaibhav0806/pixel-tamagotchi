package pet

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Render(mood Mood) string {
	frames := AnimationFrames(mood)
	return frames[0]
}

// AnimationFrames returns all animation frames for a given mood.
// Each frame is 5 lines: emoji, ears, face, body, paws.
func AnimationFrames(mood Mood) []string {
	e := mood.Emoji()
	switch mood {
	case MoodHappy:
		return []string{
			"      " + e + "\n  /\\_/\\  ~\n=( o.o )=\n /     \\\n( u   u )",
			"      " + e + "\n  /\\_/\\  )\n=( ^.^ )=\n /     \\\n( u   u )",
			"   " + e + "\n(  /\\_/\\\n=( o.o )=\n /     \\\n( u   u )",
		}
	case MoodHungry:
		return []string{
			"      " + e + "\n  /\\_/\\\n=( o.o )=\n /     \\\n( u   u )",
			"      " + e + "\n  /\\_/\\\n=( o.o )=" + "\u3064" + "\n /     \\\n( u   u )",
			"      " + e + "\n  /\\_/\\\n=( >.< )=\n /     \\\n( u   u )",
		}
	case MoodSad:
		return []string{
			"      " + e + "\n  /\\_/\\\n=( T.T )=\n /     \\\n( u   u )",
			"      " + e + "\n  /\\_/\\\n=( ;_; )=\n /     \\\n( u   u )",
		}
	case MoodAsleep:
		return []string{
			"      " + e + "\n  /\\_/\\\n=( -.- )=\n /     \\\n( u _ u )",
			"      " + e + "\n  /\\_/\\\n=( -.- )=\n /     \\\n(  u_u  )",
		}
	default:
		return []string{
			"      " + e + "\n  /\\_/\\\n=( o.o )=\n /     \\\n( u   u )",
		}
	}
}

// EarTwitchFrame returns a frame with a twitched ear for variety.
func EarTwitchFrame(mood Mood, frameIdx int) string {
	frames := AnimationFrames(mood)
	if frameIdx >= len(frames) {
		frameIdx = 0
	}
	frame := frames[frameIdx]
	return strings.Replace(frame, `/\_/\`, `/\~/\`, 1)
}

var messages = map[Mood][]string{
	MoodHappy: {
		"Pixel is purring!",
		"Pixel loves your commits!",
		"Pixel is chasing a bug (for fun)",
		"Pixel knocked your cursor off the desk",
		"Pixel is sitting on your keyboard... productively",
		"Pixel is code reviewing your last commit... looks good!",
		"Pixel found a yarn ball in your node_modules",
		"Pixel deployed to production and it worked first try",
		"Pixel is napping in a sunbeam on your monitor",
		"Pixel just mass-approved all your PRs",
		"Pixel is batting at your blinking cursor",
		"Pixel wrote a unit test while you weren't looking",
	},
	MoodHungry: {
		"Pixel is eyeing your keyboard...",
		"Pixel meows at you expectantly",
		"Pixel is pawing at your git log",
		"Pixel just knocked your coffee over... commit something!",
		"Pixel is dramatically flopping on your terminal",
		"Pixel is sitting on your backspace key in protest",
		"Pixel brought you a dead mouse... the peripheral kind",
		"Pixel is chewing on your ethernet cable",
		"Pixel is giving you the slow blink of disappointment",
	},
	MoodSad: {
		"Pixel misses you...",
		"Pixel is staring out the terminal window",
		"Pixel pushed an empty commit just to feel something",
		"Pixel is scrolling through old commit messages",
		"Pixel is watching your GitHub activity go grey",
		"Pixel tried to rebase but there's nothing to rebase",
		"Pixel is writing TODO comments with no TODOs",
		"Pixel is listening to lo-fi beats and crying",
	},
	MoodAsleep: {
		"Pixel is curled up sleeping... commit to wake her up",
		"Pixel is dreaming about merge conflicts",
		"Pixel is hibernating in /dev/null",
		"Pixel's contribution graph has flatlined",
		"Pixel.exe has stopped responding",
		"Pixel went to sleep mode to save energy",
	},
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
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
