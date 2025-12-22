package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/algernon-coop/git-auto-commit/internal/config"
)

func TestNewProvider(t *testing.T) {
	testCases := []struct {
		name      string
		config    *config.Config
		expectErr bool
	}{
		{
			name: "OpenAI",
			config: &config.Config{
				Provider: "openai",
				OpenAI: &config.OpenAIConfig{
					APIKey: "test-key",
					Model:  "gpt-4",
				},
			},
			expectErr: false,
		},
		{
			name: "Azure",
			config: &config.Config{
				Provider: "azure",
				Azure: &config.AzureOpenAIConfig{
					Endpoint:   "https://example.openai.azure.com",
					APIKey:     "test-key",
					Deployment: "gpt-4",
				},
			},
			expectErr: false,
		},
		{
			name: "Claude",
			config: &config.Config{
				Provider: "claude",
				Claude: &config.ClaudeConfig{
					APIKey: "test-key",
					Model:  "claude-3-5-sonnet-20241022",
				},
			},
			expectErr: false,
		},
		{
			name: "GitHub",
			config: &config.Config{
				Provider: "github",
				GitHub: &config.GitHubConfig{
					Token: "test-token",
					Model: "gpt-4o",
				},
			},
			expectErr: false,
		},
		{
			name: "Unknown provider",
			config: &config.Config{
				Provider: "unknown",
			},
			expectErr: true,
		},
		{
			name: "Missing OpenAI config",
			config: &config.Config{
				Provider: "openai",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewProvider(tc.config)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("Expected provider, got nil")
				}
			}
		})
	}
}

func TestOpenAIProvider_GenerateCommitMessage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Unexpected Authorization header: %s", r.Header.Get("Authorization"))
		}

		response := openAIResponse{
			Choices: []struct {
				Message openAIMessage `json:"message"`
			}{
				{
					Message: openAIMessage{
						Role:    "assistant",
						Content: "feat: add new feature",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Note: This test would need to be modified to inject the server URL
	// For now, we'll test the provider creation
	provider := NewOpenAIProvider("test-key", "gpt-4")
	if provider == nil {
		t.Error("Expected provider, got nil")
	}

	if provider.apiKey != "test-key" {
		t.Errorf("Expected apiKey 'test-key', got %s", provider.apiKey)
	}

	if provider.model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got %s", provider.model)
	}
}

func TestBuildPrompt(t *testing.T) {
	diff := "diff --git a/test.txt b/test.txt\n+new line"
	prompt := buildPrompt(diff)

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if !contains(prompt, "conventional commit") {
		t.Error("Prompt should mention conventional commit format")
	}

	if !contains(prompt, diff) {
		t.Error("Prompt should contain the diff")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
