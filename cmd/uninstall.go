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

	chainFile := filepath.Join(baseDir, "original-hooks-path")
	if data, err := os.ReadFile(chainFile); err == nil {
		originalPath := strings.TrimSpace(string(data))
		exec.Command("git", "config", "--global", "core.hooksPath", originalPath).Run()
	} else {
		exec.Command("git", "config", "--global", "--unset", "core.hooksPath").Run()
	}

	hook.Uninstall(hooksDir)
	os.RemoveAll(baseDir)

	fmt.Println("Pixel has left the terminal. Goodbye! 😿")
	return nil
}
