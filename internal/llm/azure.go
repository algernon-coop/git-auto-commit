package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AzureOpenAIProvider implements the Provider interface for Azure OpenAI
type AzureOpenAIProvider struct {
	endpoint   string
	apiKey     string
	deployment string
}

// NewAzureOpenAIProvider creates a new Azure OpenAI provider
func NewAzureOpenAIProvider(endpoint, apiKey, deployment string) *AzureOpenAIProvider {
	return &AzureOpenAIProvider{
		endpoint:   endpoint,
		apiKey:     apiKey,
		deployment: deployment,
	}
}

// GenerateCommitMessage generates a commit message using Azure OpenAI
func (p *AzureOpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	req := openAIRequest{
		Messages: []openAIMessage{
			{
				Role:    "user",
				Content: buildPrompt(diff),
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Construct Azure OpenAI URL
	endpoint := strings.TrimSuffix(p.endpoint, "/")
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=2024-02-15-preview", endpoint, p.deployment)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api-key", p.apiKey)

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
		return "", fmt.Errorf("azure OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from Azure OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
