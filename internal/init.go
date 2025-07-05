// internal/init.go
package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// InitHook sets up the git pre-commit hook relative to the given root project path
func InitHook(rootPath string) error {
	hookSource := filepath.Join(rootPath, "pre-commit-hook.sh")
	hookTarget := filepath.Join(rootPath, ".git", "hooks", "pre-commit")

	// Check if source exists
	if _, err := os.Stat(hookSource); os.IsNotExist(err) {
		return fmt.Errorf("❌ %s not found", hookSource)
	}

	// Read file content
	content, err := os.ReadFile(hookSource)
	if err != nil {
		return fmt.Errorf("failed to read hook file: %w", err)
	}

	// Write to .git/hooks/pre-commit
	err = os.WriteFile(hookTarget, content, 0755)
	if err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	// Make executable (only for Unix-based OS)
	if runtime.GOOS != "windows" {
		if err := exec.Command("chmod", "+x", hookTarget).Run(); err != nil {
			return fmt.Errorf("failed to chmod pre-commit hook: %w", err)
		}
	}

	fmt.Println("✅ Pre-commit hook initialized successfully")
	return nil
}
