package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	// Create test config
	cfg := &Config{
		Provider: "openai",
		OpenAI: &OpenAIConfig{
			APIKey: "test-key",
			Model:  "gpt-4",
		},
	}

	// Save config
	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify
	if loaded.Provider != cfg.Provider {
		t.Errorf("Provider mismatch: got %s, want %s", loaded.Provider, cfg.Provider)
	}

	if loaded.OpenAI == nil {
		t.Fatal("OpenAI config is nil")
	}

	if loaded.OpenAI.APIKey != cfg.OpenAI.APIKey {
		t.Errorf("APIKey mismatch: got %s, want %s", loaded.OpenAI.APIKey, cfg.OpenAI.APIKey)
	}

	if loaded.OpenAI.Model != cfg.OpenAI.Model {
		t.Errorf("Model mismatch: got %s, want %s", loaded.OpenAI.Model, cfg.OpenAI.Model)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/non/existent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestSaveAllProviders(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name   string
		config *Config
	}{
		{
			name: "OpenAI",
			config: &Config{
				Provider: "openai",
				OpenAI: &OpenAIConfig{
					APIKey: "test-key",
					Model:  "gpt-4",
				},
			},
		},
		{
			name: "Azure",
			config: &Config{
				Provider: "azure",
				Azure: &AzureOpenAIConfig{
					Endpoint:   "https://example.openai.azure.com",
					APIKey:     "test-key",
					Deployment: "gpt-4",
				},
			},
		},
		{
			name: "Claude",
			config: &Config{
				Provider: "claude",
				Claude: &ClaudeConfig{
					APIKey: "test-key",
					Model:  "claude-3-5-sonnet-20241022",
				},
			},
		},
		{
			name: "GitHub",
			config: &Config{
				Provider: "github",
				GitHub: &GitHubConfig{
					Token: "test-token",
					Model: "gpt-4o",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, tc.name+"-config.yaml")

			if err := Save(tc.config, configPath); err != nil {
				t.Fatalf("Failed to save config: %v", err)
			}

			loaded, err := Load(configPath)
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			if loaded.Provider != tc.config.Provider {
				t.Errorf("Provider mismatch: got %s, want %s", loaded.Provider, tc.config.Provider)
			}
		})
	}
}

func TestConfigFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	cfg := &Config{
		Provider: "openai",
		OpenAI: &OpenAIConfig{
			APIKey: "secret-key",
			Model:  "gpt-4",
		},
	}

	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Check that file permissions are restrictive (0600)
	mode := info.Mode().Perm()
	expected := os.FileMode(0600)
	if mode != expected {
		t.Errorf("File permissions mismatch: got %o, want %o", mode, expected)
	}
}
