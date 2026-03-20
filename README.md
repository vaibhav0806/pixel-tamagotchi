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

```bash
go install github.com/vaibhav/terminal-pet@latest
```

Or build from source:

```bash
git clone https://github.com/vaibhav/terminal-pet.git
cd terminal-pet
go build -o terminal-pet .
```

## Setup

```bash
terminal-pet init
```

This sets up a global git post-commit hook so Pixel knows every time you commit, in any repo. If you already have a global `core.hooksPath`, it chains to your existing hooks.

## Usage

**Check on Pixel:**

```bash
terminal-pet
```

**Open the animated dashboard:**

```bash
terminal-pet watch
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
| `terminal-pet` | Quick status check |
| `terminal-pet init` | Set up Pixel and install git hook |
| `terminal-pet watch` | Animated TUI dashboard |
| `terminal-pet reset` | Wake Pixel up, reset streak |
| `terminal-pet uninstall` | Remove hooks and clean up |

## Dev

Test different moods without waiting:

```bash
terminal-pet watch --mood happy
terminal-pet watch --mood hungry
terminal-pet watch --mood sad
terminal-pet watch --mood asleep
```

## Built with

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components

## License

MIT
