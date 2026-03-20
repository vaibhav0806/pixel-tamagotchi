package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(hooksDir string, existingHooksPath string) error {
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

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
	// Resolve the full path to the binary so the hook works in
	// non-interactive shells (e.g. git hooks) that don't have
	// nvm/homebrew/go paths in $PATH.
	binPath := "pixel-tamagotchi"
	if resolved, err := exec.LookPath("pixel-tamagotchi"); err == nil {
		if abs, err := filepath.Abs(resolved); err == nil {
			binPath = abs
		}
	} else {
		// Fallback: use the current executable's path
		if exe, err := os.Executable(); err == nil {
			binPath = exe
		}
	}

	return fmt.Sprintf(`#!/bin/sh
# pixel-tamagotchi post-commit hook
BIN="%s"
if [ ! -x "$BIN" ]; then BIN="pixel-tamagotchi"; fi
"$BIN" update-state 2>/dev/null &

# Chain existing hooks if present
ORIGINAL_HOOKS="$HOME/.pixel-tamagotchi/original-hooks-path"
if [ -f "$ORIGINAL_HOOKS" ]; then
    ORIG=$(cat "$ORIGINAL_HOOKS")
    if [ -x "$ORIG/post-commit" ]; then
        exec "$ORIG/post-commit" "$@"
    fi
fi
`, binPath)
}
