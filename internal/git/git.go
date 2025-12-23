package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// GetCommitGuidelines returns commit message guidelines from the repository
// It checks for CONTRIBUTING.md and .gitmessage files
func (r *Repository) GetCommitGuidelines() string {
	var guidelines strings.Builder

	// Check for CONTRIBUTING.md
	contributingPath := filepath.Join(r.path, "CONTRIBUTING.md")
	if content, err := os.ReadFile(contributingPath); err == nil {
		contributingContent := string(content)
		// Extract commit message related sections
		if extracted := extractCommitSection(contributingContent); extracted != "" {
			guidelines.WriteString("Repository Contributing Guidelines:\n")
			guidelines.WriteString(extracted)
			guidelines.WriteString("\n\n")
		}
	}

	// Check for .gitmessage template
	gitmessagePath := filepath.Join(r.path, ".gitmessage")
	if content, err := os.ReadFile(gitmessagePath); err == nil {
		guidelines.WriteString("Repository Commit Template:\n")
		guidelines.WriteString(string(content))
		guidelines.WriteString("\n\n")
	}

	return strings.TrimSpace(guidelines.String())
}

// extractCommitSection extracts commit-related sections from CONTRIBUTING.md
func extractCommitSection(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	inCommitSection := false
	sectionLevel := 0

	for i, line := range lines {
		// Detect section headers
		if strings.HasPrefix(line, "#") {
			headerLevel := strings.Count(strings.Split(line, " ")[0], "#")
			lowerLine := strings.ToLower(line)

			// Check if this is a commit-related section
			if strings.Contains(lowerLine, "commit") && (strings.Contains(lowerLine, "message") || strings.Contains(lowerLine, "format")) {
				inCommitSection = true
				sectionLevel = headerLevel
				result.WriteString(line)
				result.WriteString("\n")
				continue
			}

			// If we're in a commit section and encounter a header of equal or higher level, end the section
			if inCommitSection && headerLevel <= sectionLevel {
				break
			}
		}

		// Add content if we're in a commit section
		if inCommitSection {
			result.WriteString(line)
			result.WriteString("\n")

			// Stop if we've gathered enough content (to avoid including unrelated sections)
			if i > 0 && result.Len() > 1500 {
				break
			}
		}
	}

	return strings.TrimSpace(result.String())
}
