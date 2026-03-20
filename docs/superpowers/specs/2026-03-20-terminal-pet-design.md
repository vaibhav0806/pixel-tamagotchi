# terminal-pet: Design Spec

A tamagotchi cat named Pixel that lives in your terminal. Gets hungry if you haven't committed in a while, falls asleep if you're gone too long.

## Overview

`terminal-pet` is a Go CLI tool that tracks your git commit activity and reflects it as the mood of an ASCII cat named Pixel. Built with Charm libraries (Bubble Tea, Lip Gloss, Harmonica, Bubbles) for a polished terminal UI.

## Pet: Pixel

**Fixed name.** No name prompt, no customization — the cat is always called Pixel.

### Mood States

Mood is computed on-read from the time since the last commit. Never stored directly.

| State   | Threshold        | Face    | Emoji | Description                                      |
|---------|------------------|---------|-------|--------------------------------------------------|
| Happy   | < 24h            | `o.o`   | 💖    | Purring, content                                 |
| Hungry  | 24–48h           | `o.o`   | 🍕    | Meowing, begging for attention                   |
| Sad     | 48–72h           | `T.T`   | 😿    | Crying, guilt-tripping you                       |
| Asleep  | 72h+             | `-.-`   | 💤    | Curled up sleeping. First commit wakes him up.   |

No permadeath. Asleep is the worst state — Pixel always comes back on the next commit with a welcome-back message.

### ASCII Art

```
   [emoji]
 /\_/\
( [face] )
 > ^ <       (happy)
 > ~ <       (hungry/sad/asleep)
```

Pure ASCII with monospace rendering. Emoji floats above the head. Face expression and paw character change per mood.

## Architecture

### Components

1. **CLI (`terminal-pet`)** — Cobra-based with subcommands
2. **State file (`~/.terminal-pet/state.json`)** — Pet state persisted to disk
3. **Config file (`~/.terminal-pet/config.json`)** — User preferences (scan directories)
4. **Global git hook (`~/.terminal-pet/hooks/post-commit`)** — Updates state on every commit

### Data Flow

```
User commits in any repo
  → global post-commit hook fires
  → runs `terminal-pet update-state` in background
  → updates ~/.terminal-pet/state.json (last_commit_at, total_commits, streak)

User runs `terminal-pet`
  → reads state.json
  → computes mood from last_commit_at
  → renders Pixel with appropriate face/emoji/message

User runs `terminal-pet watch`
  → Bubble Tea TUI launches
  → re-reads state.json periodically
  → live mood updates, animations, status messages
```

No background daemon. The git hook is the event source, the CLI is the reader.

## State File

`~/.terminal-pet/state.json`:

```json
{
  "last_commit_at": "2026-03-20T10:30:00Z",
  "created_at": "2026-03-18T08:00:00Z",
  "total_commits": 42,
  "current_streak": 3
}
```

- `last_commit_at`: Updated by the post-commit hook. Used to compute mood.
- `created_at`: Set during `init`. For display purposes.
- `total_commits`: Incremented by the hook. Lifetime counter.
- `current_streak`: Consecutive days with at least one commit. Resets if a full calendar day passes with no commit.

## CLI Commands

### `terminal-pet` (default: status)

Quick status view. Renders Pixel + mood + stats in a few lines, then exits.

```
   💖
 /\_/\
( o.o )   Pixel is purring!
 > ^ <    Last commit: 2h ago | Streak: 3 days
```

Colored with Lip Gloss: green (happy), yellow (hungry), blue (sad), gray (asleep).

### `terminal-pet init`

Setup wizard:

1. Ask for directories to scan (used only for initial state seeding — the hook handles ongoing tracking)
2. Check if `git config --global core.hooksPath` is already set
   - If set: save existing path to `~/.terminal-pet/original-hooks-path`, chain our hook to it
   - If not set: proceed
3. Create `~/.terminal-pet/hooks/post-commit`
4. Run `git config --global core.hooksPath ~/.terminal-pet/hooks`
5. Scan configured directories for most recent commit to seed `state.json`

### `terminal-pet watch`

Full-screen Bubble Tea TUI dashboard:

```
┌──────────────────────────────────┐
│         🐱 Pixel                 │
│                                  │
│           💖                     │
│         /\_/\                    │
│        ( o.o )                   │
│         > ^ <                    │
│                                  │
│   "Pixel is chasing a bug!"     │
│                                  │
│   Mood: Happy  |  Streak: 3 🔥  │
│   Last commit: 2h ago            │
│                                  │
│   ██████████░░░░ 16h until       │
│                  hungry          │
│                                  │
│   q: quit                        │
└──────────────────────────────────┘
```

**Animations (Harmonica for spring physics):**
- Pixel blinks every few seconds (eyes swap `o.o` → `-.-` → `o.o`)
- Status message rotates from a pool every ~10 seconds
- Mood bar decreases in real-time
- Brief animation on mood transitions

**Live updates:** Re-reads `state.json` every few seconds. If you commit in another terminal, Pixel reacts live.

**Status message pools per mood:**
- Happy: "Pixel is purring!", "Pixel loves your commits!", "Pixel is chasing a bug (for fun)"
- Hungry: "Pixel is eyeing your keyboard...", "Pixel meows at you expectantly"
- Sad: "Pixel misses you...", "Pixel is staring out the terminal window"
- Asleep: "Pixel is curled up sleeping... commit to wake him up"

### `terminal-pet reset`

Reset Pixel — wake from sleep, reset streak. Keeps `total_commits` and `created_at`.

### `terminal-pet update-state`

Hidden subcommand called by the post-commit hook. Not user-facing.

- Reads `state.json`
- Sets `last_commit_at` to now
- Increments `total_commits`
- Updates `current_streak`
- Writes back

### `terminal-pet uninstall`

Clean removal:
- Restore original `core.hooksPath` if we chained
- Remove `~/.terminal-pet/` directory
- Unset `core.hooksPath` if we own it

## Global Git Hook

`~/.terminal-pet/hooks/post-commit`:

```bash
#!/bin/sh
# terminal-pet post-commit hook
terminal-pet update-state 2>/dev/null &

# Chain existing hooks if present
ORIGINAL_HOOKS="$HOME/.terminal-pet/original-hooks-path"
if [ -f "$ORIGINAL_HOOKS" ]; then
    ORIG=$(cat "$ORIGINAL_HOOKS")
    if [ -x "$ORIG/post-commit" ]; then
        exec "$ORIG/post-commit" "$@"
    fi
fi
```

Runs `update-state` in the background so it doesn't slow down commits. Chains to any pre-existing hook path.

## Hook Chaining Strategy

During `init`:

1. Read `git config --global core.hooksPath`
2. If already set → save that path to `~/.terminal-pet/original-hooks-path`. Our `post-commit` hook calls the original after updating state.
3. If not set → no chaining needed. Per-repo hooks (husky, etc.) use local `.git/hooks` which is separate from `core.hooksPath` — they continue to work.
4. Set `core.hooksPath` to `~/.terminal-pet/hooks/`

## Tech Stack

- **Language:** Go
- **CLI framework:** `github.com/spf13/cobra`
- **TUI framework:** `github.com/charmbracelet/bubbletea`
- **Terminal styling:** `github.com/charmbracelet/lipgloss`
- **Spring animations:** `github.com/charmbracelet/harmonica`
- **TUI components:** `github.com/charmbracelet/bubbles` (progress bar, spinner)

## Project Structure

```
terminal-pet/
├── main.go
├── cmd/
│   ├── root.go           # default status command
│   ├── init.go           # setup wizard
│   ├── watch.go          # launch TUI
│   ├── reset.go          # reset pet
│   ├── update_state.go   # hidden cmd called by hook
│   └── uninstall.go      # cleanup
├── internal/
│   ├── pet/
│   │   ├── state.go      # read/write state.json
│   │   ├── mood.go       # mood calculation from timestamp
│   │   └── art.go        # ASCII art per mood
│   ├── hook/
│   │   ├── install.go    # hook installation + chaining
│   │   └── uninstall.go  # hook removal
│   ├── config/
│   │   └── config.go     # read/write config.json
│   └── tui/
│       ├── model.go      # Bubble Tea model
│       ├── view.go       # rendering
│       └── animation.go  # blink, message rotation, transitions
├── go.mod
└── go.sum
```
