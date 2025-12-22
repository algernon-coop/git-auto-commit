package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Provider string             `yaml:"provider"`
	OpenAI   *OpenAIConfig      `yaml:"openai,omitempty"`
	Azure    *AzureOpenAIConfig `yaml:"azure,omitempty"`
	Claude   *ClaudeConfig      `yaml:"claude,omitempty"`
	GitHub   *GitHubConfig      `yaml:"github,omitempty"`
}

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// AzureOpenAIConfig represents Azure OpenAI configuration
type AzureOpenAIConfig struct {
	Endpoint   string `yaml:"endpoint"`
	APIKey     string `yaml:"api_key"`
	Deployment string `yaml:"deployment"`
}

// ClaudeConfig represents Anthropic Claude configuration
type ClaudeConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// GitHubConfig represents GitHub Models configuration
type GitHubConfig struct {
	Token string `yaml:"token"`
	Model string `yaml:"model"`
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = home + "/.git-auto-commit.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to a file
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
