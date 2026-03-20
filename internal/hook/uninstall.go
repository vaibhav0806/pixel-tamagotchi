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

	chainFile := filepath.Join(filepath.Dir(hooksDir), "original-hooks-path")
	os.Remove(chainFile)

	return nil
}
