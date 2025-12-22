package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Repository represents a git repository
type Repository struct {
	path string
}

// NewRepository creates a new Repository instance
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// GetStagedDiff returns the diff of staged changes
func (r *Repository) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = r.path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w (stderr: %s)", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Commit commits the staged changes with the given message
func (r *Repository) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}
