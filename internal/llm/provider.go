package llm

import (
	"context"
	"fmt"

	"github.com/algernon-coop/git-auto-commit/internal/config"
)

// Provider is the interface for AI providers
type Provider interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
}

// NewProvider creates a new provider based on the configuration
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "openai":
		if cfg.OpenAI == nil {
			return nil, fmt.Errorf("OpenAI configuration is required")
		}
		return NewOpenAIProvider(cfg.OpenAI.APIKey, cfg.OpenAI.Model), nil
	case "azure":
		if cfg.Azure == nil {
			return nil, fmt.Errorf("Azure configuration is required")
		}
		return NewAzureOpenAIProvider(cfg.Azure.Endpoint, cfg.Azure.APIKey, cfg.Azure.Deployment), nil
	case "claude":
		if cfg.Claude == nil {
			return nil, fmt.Errorf("Claude configuration is required")
		}
		return NewClaudeProvider(cfg.Claude.APIKey, cfg.Claude.Model), nil
	case "github":
		if cfg.GitHub == nil {
			return nil, fmt.Errorf("GitHub configuration is required")
		}
		return NewGitHubProvider(cfg.GitHub.Token, cfg.GitHub.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

// buildPrompt creates a prompt for generating commit messages
func buildPrompt(diff string) string {
	return fmt.Sprintf(`You are a helpful assistant that generates clear, concise git commit messages following conventional commit format.

Based on the following git diff, generate a commit message that:
1. Uses conventional commit format (e.g., "feat:", "fix:", "docs:", "refactor:", etc.)
2. Has a clear, concise subject line (max 50 characters)
3. Optionally includes a body with more details if the change is complex
4. Focuses on WHAT changed and WHY, not HOW

Git diff:
%s

Generate only the commit message, without any additional explanation or formatting markers.`, diff)
}
