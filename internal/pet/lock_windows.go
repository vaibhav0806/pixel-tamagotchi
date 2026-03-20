//go:build windows

package pet

import "os"

// Windows doesn't support flock. No-op for now — concurrent hooks are
// unlikely on Windows since git hooks run sequentially.
func lockFile(f *os.File) error   { return nil }
func unlockFile(f *os.File) error { return nil }
