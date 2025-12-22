package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/algernon-coop/git-auto-commit/internal/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure AI provider settings",
	Long:  `Interactive configuration wizard to set up your preferred AI provider and API credentials.`,
	RunE:  runConfigure,
}

func runConfigure(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	cfg := &config.Config{}

	fmt.Println("Git Auto-Commit Configuration")
	fmt.Println("==============================")
	fmt.Println()
	fmt.Println("Select AI Provider:")
	fmt.Println("1. OpenAI (native)")
	fmt.Println("2. OpenAI (Azure)")
	fmt.Println("3. Anthropic Claude")
	fmt.Println("4. GitHub Models")
	fmt.Print("\nEnter choice (1-4): ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read choice: %w", err)
	}
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		cfg.Provider = "openai"
		fmt.Print("Enter OpenAI API Key: ")
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		cfg.OpenAI = &config.OpenAIConfig{
			APIKey: strings.TrimSpace(apiKey),
			Model:  "gpt-4",
		}
		fmt.Print("Enter model (default: gpt-4): ")
		model, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read model: %w", err)
		}
		model = strings.TrimSpace(model)
		if model != "" {
			cfg.OpenAI.Model = model
		}

	case "2":
		cfg.Provider = "azure"
		fmt.Print("Enter Azure OpenAI Endpoint: ")
		endpoint, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read endpoint: %w", err)
		}
		fmt.Print("Enter Azure OpenAI API Key: ")
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		fmt.Print("Enter Deployment Name: ")
		deployment, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read deployment: %w", err)
		}

		cfg.Azure = &config.AzureOpenAIConfig{
			Endpoint:   strings.TrimSpace(endpoint),
			APIKey:     strings.TrimSpace(apiKey),
			Deployment: strings.TrimSpace(deployment),
		}

	case "3":
		cfg.Provider = "claude"
		fmt.Print("Enter Anthropic API Key: ")
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		cfg.Claude = &config.ClaudeConfig{
			APIKey: strings.TrimSpace(apiKey),
			Model:  "claude-3-5-sonnet-20241022",
		}
		fmt.Print("Enter model (default: claude-3-5-sonnet-20241022): ")
		model, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read model: %w", err)
		}
		model = strings.TrimSpace(model)
		if model != "" {
			cfg.Claude.Model = model
		}

	case "4":
		cfg.Provider = "github"
		fmt.Print("Enter GitHub Token: ")
		token, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		cfg.GitHub = &config.GitHubConfig{
			Token: strings.TrimSpace(token),
			Model: "gpt-4o",
		}
		fmt.Print("Enter model (default: gpt-4o): ")
		model, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read model: %w", err)
		}
		model = strings.TrimSpace(model)
		if model != "" {
			cfg.GitHub.Model = model
		}

	default:
		return fmt.Errorf("invalid choice")
	}

	// Save configuration
	path := configPath
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = home + "/.git-auto-commit.yaml"
	}

	if err := config.Save(cfg, path); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("\n✓ Configuration saved to %s\n", path)
	return nil
}
