package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GitHubProvider implements the Provider interface for GitHub Models
type GitHubProvider struct {
	token string
	model string
}

// NewGitHubProvider creates a new GitHub Models provider
func NewGitHubProvider(token, model string) *GitHubProvider {
	return &GitHubProvider{
		token: token,
		model: model,
	}
}

// GenerateCommitMessage generates a commit message using GitHub Models
func (p *GitHubProvider) GenerateCommitMessage(ctx context.Context, diff string, guidelines string) (string, error) {
	req := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{
				Role:    "user",
				Content: buildPromptWithGuidelines(diff, guidelines),
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://models.inference.ai.azure.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("GitHub Models API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from GitHub Models")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
