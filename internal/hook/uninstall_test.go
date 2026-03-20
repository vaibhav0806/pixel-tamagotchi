package hook_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaibhav0806/pixel-tamagotchi/internal/hook"
)

func TestUninstall_RemovesHookDir(t *testing.T) {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, "hooks")

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
