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
