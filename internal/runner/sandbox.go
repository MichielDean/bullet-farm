package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// EnsureSandbox guarantees a persistent sandbox directory exists for a worker.
// On first call it clones the repo; on subsequent calls it fetches and resets
// to origin/main so the sandbox is clean and current.
func EnsureSandbox(dir, repoURL string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return cloneSandbox(dir, repoURL)
	}
	return resetSandbox(dir)
}

// cloneSandbox performs a fresh git clone into dir.
func cloneSandbox(dir, repoURL string) error {
	if repoURL == "" {
		return fmt.Errorf("repo URL is required for initial clone")
	}
	cmd := exec.Command("git", "clone", repoURL, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone %s: %w", repoURL, err)
	}
	return nil
}

// resetSandbox fetches latest from origin and resets to origin/main.
// This gives each step a clean working tree without re-cloning.
func resetSandbox(dir string) error {
	fetch := exec.Command("git", "fetch", "origin")
	fetch.Dir = dir
	if out, err := fetch.CombinedOutput(); err != nil {
		return fmt.Errorf("git fetch in %s: %w: %s", dir, err, out)
	}

	checkout := exec.Command("git", "checkout", "main")
	checkout.Dir = dir
	if out, err := checkout.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout main in %s: %w: %s", dir, err, out)
	}

	pull := exec.Command("git", "pull", "--ff-only", "origin", "main")
	pull.Dir = dir
	if out, err := pull.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull in %s: %w: %s", dir, err, out)
	}

	// Clean untracked files left by previous runs.
	clean := exec.Command("git", "clean", "-fd")
	clean.Dir = dir
	if out, err := clean.CombinedOutput(); err != nil {
		return fmt.Errorf("git clean in %s: %w: %s", dir, err, out)
	}

	return nil
}
