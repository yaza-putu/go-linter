package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// InitHook sets up the git pre-commit hook
func InitHook(rootPath string) error {
	hookSource := filepath.Join(rootPath, "pre-commit-hook.sh")
	hookTarget := filepath.Join(rootPath, ".git", "hooks", "pre-commit")

	var content []byte
	var err error

	// 1. Try local first
	if _, statErr := os.Stat(hookSource); statErr == nil {
		content, err = os.ReadFile(hookSource)
		if err != nil {
			return fmt.Errorf("failed to read local hook file: %w", err)
		}
	} else {
		// 2. Fallback to remote GitHub URL
		url := "https://raw.githubusercontent.com/yaza-putu/golinter/main/pre-commit-hook.sh"
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch hook from GitHub: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("hook file not found at GitHub (%d)", resp.StatusCode)
		}

		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read hook content from response: %w", err)
		}
	}

	// 3. Write to .git/hooks/pre-commit
	err = os.WriteFile(hookTarget, content, 0755)
	if err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	// 4. Make it executable (Unix only)
	if runtime.GOOS != "windows" {
		if err := exec.Command("chmod", "+x", hookTarget).Run(); err != nil {
			return fmt.Errorf("failed to chmod pre-commit hook: %w", err)
		}
	}

	fmt.Println("✅ Pre-commit hook initialized successfully")
	return nil
}
