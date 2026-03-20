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
	"github.com/vaibhav0806/pixel-tamagotchi/internal/config"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/hook"
	"github.com/vaibhav0806/pixel-tamagotchi/internal/pet"
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
	statePath := config.DefaultStatePath()
	if _, err := os.Stat(statePath); err == nil {
		fmt.Println("Pixel is already here! If you want to start fresh, run 'pixel uninstall' first.")
		return nil
	}

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
	hooksDir := filepath.Join(baseDir, "hooks")
	existingPath := ""
	out, err := exec.Command("git", "config", "--global", "core.hooksPath").Output()
	if err == nil {
		existingPath = strings.TrimSpace(string(out))
		// Don't chain to our own hooks dir (happens on re-init)
		if existingPath == hooksDir {
			existingPath = ""
		}
		if existingPath != "" {
			fmt.Printf("Found existing global hooks path: %s\n", existingPath)
			fmt.Println("Pixel will chain to your existing hooks.")
		}
	}
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
	fmt.Println("Scanning for recent commits...")
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
	fmt.Println("Pixel has arrived! She'll track your commits automatically.")
	fmt.Println("Run 'pixel' anytime to check on her.")

	return nil
}

func expandHome(path string) string {
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
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
