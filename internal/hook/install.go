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
