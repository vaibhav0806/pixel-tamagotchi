# pixel-tamagotchi

A tamagotchi cat that lives in your terminal. Meet **Pixel** — she tracks your git commits and gets hungry if you haven't coded in a while.

```
      💖
  /\_/\  ~
=( ^.^ )=
 /     \
( u   u )
```

## Install

**Homebrew (macOS/Linux):**

```bash
brew install vaibhav0806/tap/pixel-tamagotchi
```

**Go:**

```bash
go install github.com/vaibhav0806/pixel-tamagotchi@latest
```

**From source:**

```bash
git clone https://github.com/vaibhav0806/pixel-tamagotchi.git
cd pixel-tamagotchi
go build -o pixel-tamagotchi .
```

**Binary download:**

Grab the latest release from [GitHub Releases](https://github.com/vaibhav0806/pixel-tamagotchi/releases).

## Setup

```bash
pixel-tamagotchi init
```

This sets up a global git post-commit hook so Pixel knows every time you commit, in any repo. If you already have a global `core.hooksPath`, it chains to your existing hooks.

## Usage

**Check on Pixel:**

```bash
pixel-tamagotchi
```

**Open the animated dashboard:**

```bash
pixel-tamagotchi watch
```

The dashboard shows Pixel with live animations, floating particles, a mood bar, commit streak, and rotating status messages. Press `q` to quit.

## Moods

Pixel's mood depends on how recently you committed:

| Mood | Threshold | What happens |
|------|-----------|-------------|
| Happy | < 24h | Purring, hearts floating, tail wagging |
| Hungry | 24-48h | Reaching for food, meowing at you |
| Sad | 48-72h | Crying, guilt-tripping you |
| Asleep | 72h+ | Curled up sleeping — commit to wake her up |

No permadeath. Pixel always comes back.

## Commands

| Command | Description |
|---------|-------------|
| `pixel-tamagotchi` | Quick status check |
| `pixel-tamagotchi init` | Set up Pixel and install git hook |
| `pixel-tamagotchi watch` | Animated TUI dashboard |
| `pixel-tamagotchi reset` | Wake Pixel up, reset streak |
| `pixel-tamagotchi uninstall` | Remove hooks and clean up |

## Dev

Test different moods without waiting:

```bash
pixel-tamagotchi watch --mood happy
pixel-tamagotchi watch --mood hungry
pixel-tamagotchi watch --mood sad
pixel-tamagotchi watch --mood asleep
```

## Built with

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components

## License

MIT
