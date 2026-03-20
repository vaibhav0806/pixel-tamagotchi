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
