package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetStagedDiff(t *testing.T) {
	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create and commit a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage initial file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Modify the file and stage it
	if err := os.WriteFile(testFile, []byte("modified content\n"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Test GetStagedDiff
	repo := NewRepository(tmpDir)
	diff, err := repo.GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	if diff == "" {
		t.Error("Expected non-empty diff")
	}

	if !strings.Contains(diff, "modified content") {
		t.Error("Diff should contain 'modified content'")
	}
}

func TestGetStagedDiff_NoChanges(t *testing.T) {
	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Test GetStagedDiff with no staged changes
	repo := NewRepository(tmpDir)
	diff, err := repo.GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	if diff != "" {
		t.Errorf("Expected empty diff, got: %s", diff)
	}
}

func TestCommit(t *testing.T) {
	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create and stage a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test Commit
	repo := NewRepository(tmpDir)
	commitMessage := "test: add test file"
	if err := repo.Commit(commitMessage); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Verify commit was created
	cmd = exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}

	if !strings.Contains(string(output), commitMessage) {
		t.Errorf("Commit message not found in log: %s", string(output))
	}
}

func TestGetCommitGuidelines_WithContributing(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a CONTRIBUTING.md file with commit guidelines
	contributingContent := `# Contributing Guide

## Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

` + "```" + `
<type>(<scope>): <subject>

<body>

<footer>
` + "```" + `

Types:
- feat: A new feature
- fix: A bug fix
- docs: Documentation changes

Examples:
` + "```" + `
feat(llm): add support for GPT-4 Turbo
fix(config): handle missing config file gracefully
` + "```" + `

## Other Guidelines

Some other unrelated content here.
`
	contributingPath := filepath.Join(tmpDir, "CONTRIBUTING.md")
	if err := os.WriteFile(contributingPath, []byte(contributingContent), 0644); err != nil {
		t.Fatalf("Failed to create CONTRIBUTING.md: %v", err)
	}

	// Test GetCommitGuidelines
	repo := NewRepository(tmpDir)
	guidelines := repo.GetCommitGuidelines()

	if guidelines == "" {
		t.Error("Expected non-empty guidelines")
	}

	if !strings.Contains(guidelines, "Conventional Commits") {
		t.Error("Guidelines should contain 'Conventional Commits'")
	}

	if !strings.Contains(guidelines, "feat:") {
		t.Error("Guidelines should contain commit types like 'feat:'")
	}
}

func TestGetCommitGuidelines_WithGitmessage(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a .gitmessage template file
	gitmessageContent := `# Title: Summary, imperative, start upper case, don't end with a period
# No more than 50 chars. #### 50 chars is here:  #

# Remember blank line between title and body.

# Body: Explain *what* and *why* (not *how*). Include task ID (Jira issue).
# Wrap at 72 chars. ################################## which is here:  #

# At the end: Include Co-authored-by for all contributors.
# Include at least one empty line before it. Format:
# Co-authored-by: name <user@users.noreply.github.com>
`
	gitmessagePath := filepath.Join(tmpDir, ".gitmessage")
	if err := os.WriteFile(gitmessagePath, []byte(gitmessageContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitmessage: %v", err)
	}

	// Test GetCommitGuidelines
	repo := NewRepository(tmpDir)
	guidelines := repo.GetCommitGuidelines()

	if guidelines == "" {
		t.Error("Expected non-empty guidelines")
	}

	if !strings.Contains(guidelines, "50 chars") {
		t.Error("Guidelines should contain template content")
	}
}

func TestGetCommitGuidelines_NoGuidelines(t *testing.T) {
	// Create a temporary directory with no guidelines
	tmpDir := t.TempDir()

	// Test GetCommitGuidelines
	repo := NewRepository(tmpDir)
	guidelines := repo.GetCommitGuidelines()

	if guidelines != "" {
		t.Errorf("Expected empty guidelines, got: %s", guidelines)
	}
}

func TestGetCommitGuidelines_WithCopilotInstructions(t *testing.T) {
	testCases := []struct {
		name     string
		filepath string
	}{
		{
			name:     "GitHub directory",
			filepath: ".github/copilot-instructions.md",
		},
		{
			name:     "Root directory",
			filepath: "copilot-instructions.md",
		},
		{
			name:     "Hidden file",
			filepath: ".copilot-instructions.md",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create Copilot instructions file
			copilotContent := `# Copilot Instructions

## Commit Message Format

Always use conventional commits with the following types:
- feat: New features
- fix: Bug fixes
- docs: Documentation updates

Keep messages concise and descriptive.
`
			copilotPath := filepath.Join(tmpDir, tc.filepath)
			// Create directory if needed
			if dir := filepath.Dir(copilotPath); dir != tmpDir {
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
			}
			if err := os.WriteFile(copilotPath, []byte(copilotContent), 0644); err != nil {
				t.Fatalf("Failed to create copilot instructions: %v", err)
			}

			// Test GetCommitGuidelines
			repo := NewRepository(tmpDir)
			guidelines := repo.GetCommitGuidelines()

			if guidelines == "" {
				t.Error("Expected non-empty guidelines")
			}

			if !strings.Contains(guidelines, "Copilot Instructions") {
				t.Error("Guidelines should contain 'Copilot Instructions'")
			}

			if !strings.Contains(guidelines, "conventional commits") {
				t.Error("Guidelines should contain 'conventional commits'")
			}
		})
	}
}

func TestGetCommitGuidelines_WithMultipleSources(t *testing.T) {
	tmpDir := t.TempDir()

	// Create CONTRIBUTING.md
	contributingContent := `# Contributing

## Commit Message Format

Use conventional commits.
`
	contributingPath := filepath.Join(tmpDir, "CONTRIBUTING.md")
	if err := os.WriteFile(contributingPath, []byte(contributingContent), 0644); err != nil {
		t.Fatalf("Failed to create CONTRIBUTING.md: %v", err)
	}

	// Create Copilot instructions
	copilotContent := `# Copilot Instructions

## Commit Format

Keep commits under 50 characters.
`
	githubDir := filepath.Join(tmpDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatalf("Failed to create .github directory: %v", err)
	}
	copilotPath := filepath.Join(githubDir, "copilot-instructions.md")
	if err := os.WriteFile(copilotPath, []byte(copilotContent), 0644); err != nil {
		t.Fatalf("Failed to create copilot instructions: %v", err)
	}

	// Create .gitmessage
	gitmessageContent := `# Template for commits`
	gitmessagePath := filepath.Join(tmpDir, ".gitmessage")
	if err := os.WriteFile(gitmessagePath, []byte(gitmessageContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitmessage: %v", err)
	}

	// Test GetCommitGuidelines
	repo := NewRepository(tmpDir)
	guidelines := repo.GetCommitGuidelines()

	// Should contain all three sources
	if !strings.Contains(guidelines, "Repository Contributing Guidelines") {
		t.Error("Guidelines should contain contributing guidelines")
	}

	if !strings.Contains(guidelines, "Copilot Instructions") {
		t.Error("Guidelines should contain copilot instructions")
	}

	if !strings.Contains(guidelines, "Repository Commit Template") {
		t.Error("Guidelines should contain commit template")
	}
}

func TestExtractCommitSection(t *testing.T) {
	content := `# Contributing Guide

## Getting Started

Some content here.

## Commit Message Format

We use conventional commits.

Types:
- feat: new feature
- fix: bug fix

## Code Style

Other content.
`

	result := extractCommitSection(content)

	if result == "" {
		t.Error("Expected non-empty result")
	}

	if !strings.Contains(result, "Commit Message Format") {
		t.Error("Result should contain the commit section header")
	}

	if !strings.Contains(result, "conventional commits") {
		t.Error("Result should contain commit section content")
	}

	if strings.Contains(result, "Code Style") {
		t.Error("Result should not contain the next section")
	}
}
