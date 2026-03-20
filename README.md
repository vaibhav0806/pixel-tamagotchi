# pixel-tamagotchi

A tamagotchi cat that lives in your terminal. Meet **Pixel** — she tracks your git commits and gets hungry if you haven't coded in a while.

![pixel-tamagotchi demo](demo.gif)

## Install

**npm (easiest — works everywhere):**

```bash
npm install -g pixel-tamagotchi
```

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
pixel init
```

This sets up a global git post-commit hook so Pixel knows every time you commit, in any repo. If you already have a global `core.hooksPath`, it chains to your existing hooks.

## Usage

**Check on Pixel:**

```bash
pixel
```

**Open the animated dashboard:**

```bash
pixel watch
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
| `pixel` / `pixel status` | Quick status check |
| `pixel init` | Set up Pixel and install git hook |
| `pixel watch` | Animated TUI dashboard (bell on mood changes) |
| `pixel stats` | Lifetime stats — commits, streaks, days together |
| `pixel reset` | Wake Pixel up, reset streak |
| `pixel uninstall` | Remove hooks and clean up |

`pixel-tamagotchi` also works as an alias for all commands.

## Dev

Test different moods without waiting:

```bash
pixel watch --mood happy
pixel watch --mood hungry
pixel watch --mood sad
pixel watch --mood asleep
```

## How it works

Pixel uses a global git `post-commit` hook to detect when you commit — in any repo on your machine. No background daemon, no network calls, no data leaves your computer. Your commit timestamps are stored locally in `~/.pixel-tamagotchi/state.json`.

## Built with

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components

## License

MIT
