# terminal-pet Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a CLI tamagotchi cat (Pixel) that tracks git commits and reflects coding activity as mood states in the terminal.

**Architecture:** Cobra CLI with subcommands reads/writes a state file at `~/.terminal-pet/state.json`. A global git post-commit hook updates the state on every commit. Mood is computed on-read from the time since last commit. A Bubble Tea TUI provides an animated dashboard.

**Tech Stack:** Go, Cobra, Bubble Tea, Lip Gloss, Harmonica, Bubbles

**Spec:** `docs/superpowers/specs/2026-03-20-terminal-pet-design.md`

---

## File Structure

```
terminal-pet/
├── main.go                     # Entry point, executes root command
├── cmd/
│   ├── root.go                 # Root command: default status display
│   ├── init_cmd.go             # Init command: setup wizard
│   ├── watch.go                # Watch command: launches TUI
│   ├── reset.go                # Reset command: wake Pixel / reset streak
│   ├── update_state.go         # Hidden command: called by post-commit hook
│   └── uninstall.go            # Uninstall command: cleanup hooks + config
├── internal/
│   ├── pet/
│   │   ├── state.go            # State struct, Load/Save to JSON
│   │   ├── state_test.go       # Tests for state persistence
│   │   ├── mood.go             # Mood enum, ComputeMood(lastCommit) → Mood
│   │   ├── mood_test.go        # Tests for mood thresholds
│   │   ├── art.go              # Render(mood) → string, status messages
│   │   └── art_test.go         # Tests for art rendering
│   ├── hook/
│   │   ├── install.go          # InstallHook(), chain detection
│   │   ├── install_test.go     # Tests for hook installation
│   │   ├── uninstall.go        # UninstallHook(), restore original
│   │   └── uninstall_test.go   # Tests for hook uninstall
│   ├── config/
│   │   ├── config.go           # Config struct, Load/Save, default paths
│   │   └── config_test.go      # Tests for config
│   └── tui/
│       ├── model.go            # Bubble Tea Model, Init, Update, View
│       ├── view.go             # View rendering with Lip Gloss
│       └── animation.go        # Blink timer, message rotation, mood bar
├── go.mod
└── go.sum
```

---

### Task 1: Project Scaffolding

**Files:**
- Create: `main.go`
- Create: `cmd/root.go`
- Create: `go.mod`

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/vaibhav/Documents/projects/tamagotchi
go mod init github.com/vaibhav/terminal-pet
```

- [ ] **Step 2: Install dependencies**

```bash
go get github.com/spf13/cobra@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/harmonica@latest
go get github.com/charmbracelet/bubbles@latest
```

- [ ] **Step 3: Create main.go**

```go
package main

import "github.com/vaibhav/terminal-pet/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 4: Create cmd/root.go with placeholder**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terminal-pet",
	Short: "A tamagotchi cat that lives in your terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Pixel says hi! Run 'terminal-pet init' to get started.")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 5: Verify it builds and runs**

```bash
go build -o terminal-pet .
./terminal-pet
```

Expected: `Pixel says hi! Run 'terminal-pet init' to get started.`

- [ ] **Step 6: Commit**

```bash
git add main.go cmd/root.go go.mod go.sum
git commit -m "feat: scaffold project with cobra CLI"
```

---

### Task 2: State Management

**Files:**
- Create: `internal/pet/state.go`
- Create: `internal/pet/state_test.go`

- [ ] **Step 1: Write failing tests for state Load/Save**

```go
package pet_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaibhav/terminal-pet/internal/pet"
)

func TestSaveAndLoadState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	now := time.Now().Truncate(time.Second)
	state := &pet.State{
		LastCommitAt:  now,
		CreatedAt:     now.Add(-48 * time.Hour),
		TotalCommits:  42,
		CurrentStreak: 3,
	}

	if err := pet.SaveState(state, path); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	loaded, err := pet.LoadState(path)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}

	if !loaded.LastCommitAt.Equal(state.LastCommitAt) {
		t.Errorf("LastCommitAt: got %v, want %v", loaded.LastCommitAt, state.LastCommitAt)
	}
	if loaded.TotalCommits != 42 {
		t.Errorf("TotalCommits: got %d, want 42", loaded.TotalCommits)
	}
	if loaded.CurrentStreak != 3 {
		t.Errorf("CurrentStreak: got %d, want 3", loaded.CurrentStreak)
	}
}

func TestLoadState_FileNotFound(t *testing.T) {
	_, err := pet.LoadState("/nonexistent/state.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRecordCommit(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	state := &pet.State{
		LastCommitAt:  now.Add(-2 * time.Hour),
		CreatedAt:     now.Add(-48 * time.Hour),
		TotalCommits:  10,
		CurrentStreak: 1,
	}

	pet.RecordCommit(state)

	if state.TotalCommits != 11 {
		t.Errorf("TotalCommits: got %d, want 11", state.TotalCommits)
	}
	if time.Since(state.LastCommitAt) > time.Second {
		t.Error("LastCommitAt was not updated to now")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/pet/ -v
```

Expected: compilation error — `pet` package doesn't exist yet.

- [ ] **Step 3: Implement state.go**

```go
package pet

import (
	"encoding/json"
	"os"
	"time"
)

type State struct {
	LastCommitAt  time.Time `json:"last_commit_at"`
	CreatedAt     time.Time `json:"created_at"`
	TotalCommits  int       `json:"total_commits"`
	CurrentStreak int       `json:"current_streak"`
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveState(s *State, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func RecordCommit(s *State) {
	now := time.Now()

	// Update streak: if last commit was yesterday (local time), continue streak.
	// If last commit was today, keep streak. Otherwise reset to 1.
	lastDay := s.LastCommitAt.Local().Truncate(24 * time.Hour)
	today := now.Local().Truncate(24 * time.Hour)
	diff := today.Sub(lastDay)

	switch {
	case diff == 24*time.Hour:
		s.CurrentStreak++
	case diff == 0:
		// same day, streak unchanged
	default:
		s.CurrentStreak = 1
	}

	s.LastCommitAt = now
	s.TotalCommits++
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/pet/ -v
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/pet/state.go internal/pet/state_test.go
git commit -m "feat: add state persistence and commit recording"
```

---

### Task 3: Mood Calculation

**Files:**
- Create: `internal/pet/mood.go`
- Create: `internal/pet/mood_test.go`

- [ ] **Step 1: Write failing tests for mood computation**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/pet/ -v -run TestComputeMood
```

Expected: compilation error.

- [ ] **Step 3: Implement mood.go**

```go
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
		return "💖"
	case MoodHungry:
		return "🍕"
	case MoodSad:
		return "😿"
	case MoodAsleep:
		return "💤"
	default:
		return "❓"
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/pet/ -v -run TestComputeMood
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/pet/mood.go internal/pet/mood_test.go
git commit -m "feat: add mood computation from commit timestamp"
```

---

### Task 4: ASCII Art Rendering

**Files:**
- Create: `internal/pet/art.go`
- Create: `internal/pet/art_test.go`

- [ ] **Step 1: Write failing tests for art rendering**

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/pet/ -v -run TestRender
```

Expected: compilation error.

- [ ] **Step 3: Implement art.go**

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/pet/ -v -run "TestRender|TestRandom"
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/pet/art.go internal/pet/art_test.go
git commit -m "feat: add ASCII art rendering and status messages"
```

---

### Task 5: Config Management

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests for config**

```go
package config_test

import (
	"path/filepath"
	"testing"

	"github.com/vaibhav/terminal-pet/internal/config"
)

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := &config.Config{
		ScanDirs: []string{"~/projects", "~/work"},
	}

	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.ScanDirs) != 2 {
		t.Fatalf("ScanDirs: got %d, want 2", len(loaded.ScanDirs))
	}
	if loaded.ScanDirs[0] != "~/projects" {
		t.Errorf("ScanDirs[0]: got %q, want ~/projects", loaded.ScanDirs[0])
	}
}

func TestDefaultDir(t *testing.T) {
	dir := config.DefaultDir()
	if dir == "" {
		t.Error("DefaultDir returned empty string")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/config/ -v
```

Expected: compilation error.

- [ ] **Step 3: Implement config.go**

```go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ScanDirs []string `json:"scan_dirs"`
}

func DefaultDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".terminal-pet")
}

func DefaultStatePath() string {
	return filepath.Join(DefaultDir(), "state.json")
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultDir(), "config.json")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Save(c *Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/config/ -v
```

Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add config management"
```

---

### Task 6: Hook Installation

**Files:**
- Create: `internal/hook/install.go`
- Create: `internal/hook/install_test.go`
- Create: `internal/hook/uninstall.go`
- Create: `internal/hook/uninstall_test.go`

- [ ] **Step 1: Write failing tests for hook installation**

```go
package hook_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaibhav/terminal-pet/internal/hook"
)

func TestInstallHook_CreatesPostCommit(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	err := hook.Install(hooksDir, "")
	if err != nil {
		t.Fatalf("Install: %v", err)
	}

	hookPath := filepath.Join(hooksDir, "post-commit")
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("hook file not created: %v", err)
	}

	// Check executable
	if info.Mode()&0111 == 0 {
		t.Error("hook file is not executable")
	}
}

func TestInstallHook_ChainsExisting(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")
	existingPath := "/some/existing/hooks"

	err := hook.Install(hooksDir, existingPath)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}

	chainFile := filepath.Join(filepath.Dir(hooksDir), "original-hooks-path")
	data, err := os.ReadFile(chainFile)
	if err != nil {
		t.Fatalf("chain file not created: %v", err)
	}
	if string(data) != existingPath {
		t.Errorf("chain file: got %q, want %q", string(data), existingPath)
	}
}

func TestGenerateHookScript(t *testing.T) {
	script := hook.GenerateHookScript()
	if script == "" {
		t.Error("expected non-empty hook script")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/hook/ -v
```

Expected: compilation error.

- [ ] **Step 3: Implement install.go**

```go
package hook

import (
	"fmt"
	"os"
	"path/filepath"
)

func Install(hooksDir string, existingHooksPath string) error {
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	// Save existing hooks path for chaining
	if existingHooksPath != "" {
		chainFile := filepath.Join(filepath.Dir(hooksDir), "original-hooks-path")
		if err := os.WriteFile(chainFile, []byte(existingHooksPath), 0644); err != nil {
			return fmt.Errorf("save original hooks path: %w", err)
		}
	}

	hookPath := filepath.Join(hooksDir, "post-commit")
	script := GenerateHookScript()
	if err := os.WriteFile(hookPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	return nil
}

func GenerateHookScript() string {
	return `#!/bin/sh
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
`
}
```

- [ ] **Step 4: Write failing tests for hook uninstall**

```go
package hook_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaibhav/terminal-pet/internal/hook"
)

func TestUninstall_RemovesHookDir(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

	// Install first
	if err := hook.Install(hooksDir, ""); err != nil {
		t.Fatalf("Install: %v", err)
	}

	if err := hook.Uninstall(hooksDir); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	if _, err := os.Stat(hooksDir); !os.IsNotExist(err) {
		t.Error("hooks dir still exists after uninstall")
	}
}
```

- [ ] **Step 5: Implement uninstall.go**

```go
package hook

import (
	"fmt"
	"os"
	"path/filepath"
)

func Uninstall(hooksDir string) error {
	if err := os.RemoveAll(hooksDir); err != nil {
		return fmt.Errorf("remove hooks dir: %w", err)
	}

	// Remove chain file if present
	chainFile := filepath.Join(filepath.Dir(hooksDir), "original-hooks-path")
	os.Remove(chainFile) // ignore error — may not exist

	return nil
}
```

- [ ] **Step 6: Run all hook tests**

```bash
go test ./internal/hook/ -v
```

Expected: all PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/hook/install.go internal/hook/install_test.go internal/hook/uninstall.go internal/hook/uninstall_test.go
git commit -m "feat: add git hook installation and uninstallation"
```

---

### Task 7: Status Command (Root)

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 1: Update root.go with actual status display**

```go
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var rootCmd = &cobra.Command{
	Use:   "terminal-pet",
	Short: "A tamagotchi cat that lives in your terminal",
	RunE:  runStatus,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' to adopt him.")
		return nil
	}

	mood := pet.ComputeMood(state.LastCommitAt)
	art := pet.Render(mood)
	msg := pet.RandomMessage(mood)

	color := colorForMood(mood)
	style := lipgloss.NewStyle().Foreground(color)

	elapsed := time.Since(state.LastCommitAt).Truncate(time.Minute)
	stats := fmt.Sprintf("Last commit: %s ago | Streak: %d days", formatDuration(elapsed), state.CurrentStreak)

	fmt.Println(style.Render(art) + "   " + msg)
	fmt.Println(style.Render("          " + stats))

	return nil
}

func colorForMood(mood pet.Mood) lipgloss.Color {
	switch mood {
	case pet.MoodHappy:
		return lipgloss.Color("#4ade80")
	case pet.MoodHungry:
		return lipgloss.Color("#facc15")
	case pet.MoodSad:
		return lipgloss.Color("#60a5fa")
	case pet.MoodAsleep:
		return lipgloss.Color("#94a3b8")
	default:
		return lipgloss.Color("#e0e0e0")
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet
```

Expected: "Pixel isn't here yet! Run 'terminal-pet init' to adopt him."

- [ ] **Step 3: Commit**

```bash
git add cmd/root.go
git commit -m "feat: add status display as default command"
```

---

### Task 8: Init Command

**Files:**
- Create: `cmd/init_cmd.go`

- [ ] **Step 1: Implement init command**

```go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/hook"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up Pixel — your terminal pet",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	baseDir := config.DefaultDir()
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("create base dir: %w", err)
	}

	// Ask for scan directories
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Which directories should Pixel scan for git repos? (default: ~/): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		home, _ := os.UserHomeDir()
		input = home
	}

	scanDirs := strings.Split(input, ",")
	for i := range scanDirs {
		scanDirs[i] = strings.TrimSpace(scanDirs[i])
	}

	// Save config
	cfg := &config.Config{ScanDirs: scanDirs}
	if err := config.Save(cfg, config.DefaultConfigPath()); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Check existing hooks path
	existingPath := ""
	out, err := exec.Command("git", "config", "--global", "core.hooksPath").Output()
	if err == nil {
		existingPath = strings.TrimSpace(string(out))
		if existingPath != "" {
			fmt.Printf("Found existing global hooks path: %s\n", existingPath)
			fmt.Println("Pixel will chain to your existing hooks.")
		}
	}

	// Install hook
	hooksDir := filepath.Join(baseDir, "hooks")
	if err := hook.Install(hooksDir, existingPath); err != nil {
		return fmt.Errorf("install hook: %w", err)
	}

	// Set global hooks path
	if err := exec.Command("git", "config", "--global", "core.hooksPath", hooksDir).Run(); err != nil {
		return fmt.Errorf("set global hooks path: %w", err)
	}

	// Seed initial state — start with zero time so any found commit wins
	now := time.Now()
	state := &pet.State{
		CreatedAt:     now,
		LastCommitAt:  time.Time{},
		TotalCommits:  0,
		CurrentStreak: 0,
	}

	// Try to find most recent commit across scan dirs
	for _, dir := range scanDirs {
		expanded := expandHome(dir)
		latestCommit := findLatestCommit(expanded)
		if !latestCommit.IsZero() && latestCommit.After(state.LastCommitAt) {
			state.LastCommitAt = latestCommit
		}
	}

	// If no commits found, default to now (Pixel starts Happy)
	if state.LastCommitAt.IsZero() {
		state.LastCommitAt = now
	}

	if err := pet.SaveState(state, config.DefaultStatePath()); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	fmt.Println()
	fmt.Println(pet.Render(pet.ComputeMood(state.LastCommitAt)))
	fmt.Println()
	fmt.Println("Pixel has arrived! He'll track your commits automatically.")
	fmt.Println("Run 'terminal-pet' anytime to check on him.")

	return nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func findLatestCommit(dir string) time.Time {
	var latest time.Time

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == ".git" {
			repoDir := filepath.Dir(path)
			out, err := exec.Command("git", "-C", repoDir, "log", "-1", "--format=%aI").Output()
			if err == nil {
				t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
				if err == nil && t.After(latest) {
					latest = t
				}
			}
			return filepath.SkipDir
		}
		// Skip deep directories
		if info.IsDir() {
			rel, _ := filepath.Rel(dir, path)
			if strings.Count(rel, string(filepath.Separator)) > 3 {
				return filepath.SkipDir
			}
		}
		return nil
	})

	return latest
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet init --help
```

Expected: shows help text for init command.

- [ ] **Step 3: Commit**

```bash
git add cmd/init_cmd.go
git commit -m "feat: add init command with hook setup and repo scanning"
```

---

### Task 9: Update State Command (Hidden)

**Files:**
- Create: `cmd/update_state.go`

- [ ] **Step 1: Implement update-state command**

```go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var updateStateCmd = &cobra.Command{
	Use:    "update-state",
	Short:  "Update pet state (called by git hook)",
	Hidden: true,
	RunE:   runUpdateState,
}

func init() {
	rootCmd.AddCommand(updateStateCmd)
}

func runUpdateState(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		return err
	}

	pet.RecordCommit(state)

	return pet.SaveState(state, statePath)
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet update-state --help
```

Expected: shows help text. Verify it's hidden: `./terminal-pet --help` should NOT list update-state.

- [ ] **Step 3: Commit**

```bash
git add cmd/update_state.go
git commit -m "feat: add hidden update-state command for git hook"
```

---

### Task 10: Reset Command

**Files:**
- Create: `cmd/reset.go`

- [ ] **Step 1: Implement reset command**

```go
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Wake Pixel up and reset streak",
	RunE:  runReset,
}

func init() {
	rootCmd.AddCommand(resetCmd)
}

func runReset(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()

	state, err := pet.LoadState(statePath)
	if err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' first.")
		return nil
	}

	state.LastCommitAt = time.Now()
	state.CurrentStreak = 0

	if err := pet.SaveState(state, statePath); err != nil {
		return err
	}

	fmt.Println(pet.Render(pet.MoodHappy))
	fmt.Println()
	fmt.Println("Pixel is awake! Streak reset to 0.")

	return nil
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet reset --help
```

Expected: shows help text.

- [ ] **Step 3: Commit**

```bash
git add cmd/reset.go
git commit -m "feat: add reset command"
```

---

### Task 11: Uninstall Command

**Files:**
- Create: `cmd/uninstall.go`

- [ ] **Step 1: Implement uninstall command**

```go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/hook"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Pixel and clean up hooks",
	RunE:  runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to uninstall Pixel? (y/N): ")
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("Pixel is relieved!")
		return nil
	}

	baseDir := config.DefaultDir()
	hooksDir := filepath.Join(baseDir, "hooks")

	// Restore original hooks path if we chained
	chainFile := filepath.Join(baseDir, "original-hooks-path")
	if data, err := os.ReadFile(chainFile); err == nil {
		originalPath := strings.TrimSpace(string(data))
		exec.Command("git", "config", "--global", "core.hooksPath", originalPath).Run()
	} else {
		// We own the hooks path — unset it
		exec.Command("git", "config", "--global", "--unset", "core.hooksPath").Run()
	}

	// Remove hooks
	hook.Uninstall(hooksDir)

	// Remove entire base dir
	os.RemoveAll(baseDir)

	fmt.Println("Pixel has left the terminal. Goodbye! 😿")
	return nil
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet uninstall --help
```

Expected: shows help text.

- [ ] **Step 3: Commit**

```bash
git add cmd/uninstall.go
git commit -m "feat: add uninstall command with hook cleanup"
```

---

### Task 12: Bubble Tea TUI — Model & View

**Files:**
- Create: `internal/tui/model.go`
- Create: `internal/tui/view.go`

- [ ] **Step 1: Implement model.go**

```go
package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

type tickMsg time.Time
type stateRefreshMsg struct{ state *pet.State }

type Model struct {
	state     *pet.State
	statePath string
	mood      pet.Mood
	message   string
	blinkOpen bool
	width     int
	height    int
}

func NewModel() Model {
	statePath := config.DefaultStatePath()
	state, err := pet.LoadState(statePath)
	if err != nil {
		state = &pet.State{
			CreatedAt:    time.Now(),
			LastCommitAt: time.Now(),
		}
	}

	mood := pet.ComputeMood(state.LastCommitAt)

	return Model{
		state:     state,
		statePath: statePath,
		mood:      mood,
		message:   pet.RandomMessage(mood),
		blinkOpen: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(time.Second),
		refreshState(m.statePath),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.mood = pet.ComputeMood(m.state.LastCommitAt)
		return m, tickEvery(time.Second)

	case stateRefreshMsg:
		if msg.state != nil {
			oldMood := m.mood
			m.state = msg.state
			m.mood = pet.ComputeMood(m.state.LastCommitAt)
			if m.mood != oldMood {
				m.message = pet.RandomMessage(m.mood)
			}
		}
		return m, refreshStateAfter(m.statePath, 3*time.Second)
	}

	return m, nil
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func refreshState(path string) tea.Cmd {
	return func() tea.Msg {
		state, _ := pet.LoadState(path)
		return stateRefreshMsg{state: state}
	}
}

func refreshStateAfter(path string, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		state, _ := pet.LoadState(path)
		return stateRefreshMsg{state: state}
	})
}
```

- [ ] **Step 2: Implement view.go**

```go
package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

func (m Model) View() string {
	color := colorForMood(m.mood)

	titleStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		MarginBottom(1)

	artStyle := lipgloss.NewStyle().
		Foreground(color)

	messageStyle := lipgloss.NewStyle().
		Foreground(color).
		Italic(true).
		MarginTop(1)

	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		MarginTop(1)

	// Build mood bar
	elapsed := time.Since(m.state.LastCommitAt)
	moodBar := renderMoodBar(elapsed, m.mood, color)

	// Format elapsed time
	elapsedStr := formatDuration(elapsed.Truncate(time.Minute))

	art := pet.Render(m.mood)

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("🐱 Pixel"),
		artStyle.Render(art),
		messageStyle.Render(fmt.Sprintf("%q", m.message)),
		statsStyle.Render(fmt.Sprintf("Mood: %s  |  Streak: %d 🔥", m.mood.String(), m.state.CurrentStreak)),
		statsStyle.Render(fmt.Sprintf("Last commit: %s ago", elapsedStr)),
		moodBar,
		helpStyle.Render("q: quit"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func renderMoodBar(elapsed time.Duration, mood pet.Mood, color lipgloss.Color) string {
	// Calculate progress to next mood
	var progress float64
	var label string

	switch mood {
	case pet.MoodHappy:
		progress = 1.0 - (float64(elapsed) / float64(24*time.Hour))
		remaining := 24*time.Hour - elapsed
		label = fmt.Sprintf("%s until hungry", formatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodHungry:
		progress = 1.0 - (float64(elapsed-24*time.Hour) / float64(24*time.Hour))
		remaining := 48*time.Hour - elapsed
		label = fmt.Sprintf("%s until sad", formatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodSad:
		progress = 1.0 - (float64(elapsed-48*time.Hour) / float64(24*time.Hour))
		remaining := 72*time.Hour - elapsed
		label = fmt.Sprintf("%s until asleep", formatDuration(remaining.Truncate(time.Minute)))
	case pet.MoodAsleep:
		progress = 0
		label = "Pixel is asleep... commit to wake him up"
	}

	if progress < 0 {
		progress = 0
	}

	barWidth := 20
	filled := int(progress * float64(barWidth))
	empty := barWidth - filled

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

	bar := ""
	for i := 0; i < filled; i++ {
		bar += filledStyle.Render("█")
	}
	for i := 0; i < empty; i++ {
		bar += emptyStyle.Render("░")
	}

	return fmt.Sprintf("%s %s", bar, label)
}

func colorForMood(mood pet.Mood) lipgloss.Color {
	switch mood {
	case pet.MoodHappy:
		return lipgloss.Color("#4ade80")
	case pet.MoodHungry:
		return lipgloss.Color("#facc15")
	case pet.MoodSad:
		return lipgloss.Color("#60a5fa")
	case pet.MoodAsleep:
		return lipgloss.Color("#94a3b8")
	default:
		return lipgloss.Color("#e0e0e0")
	}
}

func formatDuration(d time.Duration) string {
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
```

- [ ] **Step 3: Build and verify compilation**

```bash
go build -o terminal-pet .
```

Expected: compiles without errors.

- [ ] **Step 4: Commit**

```bash
git add internal/tui/model.go internal/tui/view.go
git commit -m "feat: add Bubble Tea TUI model and view"
```

---

### Task 13: TUI Animations

**Files:**
- Create: `internal/tui/animation.go`
- Modify: `internal/tui/model.go`

- [ ] **Step 1: Implement animation.go**

```go
package tui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaibhav/terminal-pet/internal/pet"
)

type blinkMsg struct{}
type messageRotateMsg struct{}

func blinkCmd() tea.Cmd {
	delay := time.Duration(3+rand.Intn(5)) * time.Second
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func messageRotateCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return messageRotateMsg{}
	})
}

func handleBlink(m *Model) tea.Cmd {
	m.blinkOpen = !m.blinkOpen
	if !m.blinkOpen {
		// Close eyes briefly, then reopen
		return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
			return blinkMsg{}
		})
	}
	return blinkCmd()
}

func handleMessageRotate(m *Model) tea.Cmd {
	m.message = pet.RandomMessage(m.mood)
	return messageRotateCmd()
}
```

- [ ] **Step 2: Add blink rendering to art.go**

Add to `internal/pet/art.go`:

```go
func RenderWithBlink(mood Mood, blinkOpen bool) string {
	face := faceFor(mood)
	if !blinkOpen {
		face = "-.-"
	}
	paws := pawsFor(mood)
	emoji := mood.Emoji()

	return fmt.Sprintf("   %s\n /\\_/\\\n( %s )\n %s", emoji, face, paws)
}
```

- [ ] **Step 3: Update model.go to handle animation messages**

Add to the `Init()` method's `tea.Batch`:

```go
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(time.Second),
		refreshState(m.statePath),
		blinkCmd(),
		messageRotateCmd(),
	)
}
```

Add to the `Update()` switch:

```go
	case blinkMsg:
		return m, handleBlink(&m)

	case messageRotateMsg:
		return m, handleMessageRotate(&m)
```

Update `View()` to use `RenderWithBlink`:

```go
	art := pet.RenderWithBlink(m.mood, m.blinkOpen)
```

- [ ] **Step 4: Build and verify compilation**

```bash
go build -o terminal-pet .
```

Expected: compiles without errors.

- [ ] **Step 5: Commit**

```bash
git add internal/tui/animation.go internal/tui/model.go internal/tui/view.go internal/pet/art.go
git commit -m "feat: add blink and message rotation animations"
```

---

### Task 14: Watch Command

**Files:**
- Create: `cmd/watch.go`

- [ ] **Step 1: Implement watch command**

```go
package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vaibhav/terminal-pet/internal/config"
	"github.com/vaibhav/terminal-pet/internal/pet"
	"github.com/vaibhav/terminal-pet/internal/tui"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Open Pixel's dashboard",
	RunE:  runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	statePath := config.DefaultStatePath()
	if _, err := pet.LoadState(statePath); err != nil {
		fmt.Println("Pixel isn't here yet! Run 'terminal-pet init' first.")
		return nil
	}

	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
```

- [ ] **Step 2: Build and verify**

```bash
go build -o terminal-pet .
./terminal-pet watch --help
```

Expected: shows help text.

- [ ] **Step 3: Commit**

```bash
git add cmd/watch.go
git commit -m "feat: add watch command for TUI dashboard"
```

---

### Task 15: Integration Test & Polish

**Files:**
- Modify: `cmd/root.go` (remove duplicate `colorForMood` and `formatDuration` — use from `tui` or extract to shared)

- [ ] **Step 1: Fix duplicate functions**

Move `colorForMood` and `formatDuration` from `cmd/root.go` to a shared location. Since `cmd` imports `internal/tui` would cause issues, move these helpers to `internal/pet/art.go`:

Add to `internal/pet/art.go`:

```go
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
```

Update `cmd/root.go` to use `pet.ColorForMood()` and `pet.FormatDuration()`.
Update `internal/tui/view.go` to use `pet.ColorForMood()` and `pet.FormatDuration()`, removing local duplicates.

- [ ] **Step 2: Run all tests**

```bash
go test ./... -v
```

Expected: all PASS.

- [ ] **Step 3: Build final binary**

```bash
go build -o terminal-pet .
```

Expected: compiles cleanly.

- [ ] **Step 4: Commit**

```bash
git add .
git commit -m "refactor: extract shared helpers, fix duplicates"
```

---

### Task 16: End-to-End Manual Test

- [ ] **Step 1: Run the full flow**

```bash
# Build
go build -o terminal-pet .

# Check status before init
./terminal-pet

# Initialize
./terminal-pet init

# Check status after init
./terminal-pet

# Simulate a commit
./terminal-pet update-state
./terminal-pet

# Open dashboard
./terminal-pet watch
# Press q to quit

# Reset
./terminal-pet reset
./terminal-pet
```

- [ ] **Step 2: Verify hook works in a real repo**

```bash
cd /tmp && mkdir test-repo && cd test-repo && git init
echo "test" > test.txt && git add . && git commit -m "test commit"
cd /Users/vaibhav/Documents/projects/tamagotchi
./terminal-pet
```

Expected: Pixel should be Happy (commit was just made).

- [ ] **Step 3: Commit any fixes**

```bash
git add -A
git commit -m "fix: address issues found during manual testing"
```

---

### Task 17: Add .gitignore

- [ ] **Step 1: Create .gitignore**

```
terminal-pet
.superpowers/
```

- [ ] **Step 2: Commit**

```bash
git add .gitignore
git commit -m "chore: add gitignore"
```
