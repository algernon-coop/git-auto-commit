package cmd

import (
	"fmt"

	"github.com/algernon-coop/git-auto-commit/internal/config"
	"github.com/algernon-coop/git-auto-commit/internal/git"
	"github.com/algernon-coop/git-auto-commit/internal/llm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-auto-commit",
	Short: "Automatically generate commit messages using AI",
	Long: `git-auto-commit is a tool that generates meaningful commit messages
based on your staged changes using various AI providers (OpenAI, Claude, GitHub models).`,
	RunE: runGenerate,
}

var (
	configPath string
	dryRun     bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path (default: $HOME/.git-auto-commit.yaml)")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "generate commit message without committing")
	
	rootCmd.AddCommand(configureCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get staged changes
	gitRepo := git.NewRepository(".")
	diff, err := gitRepo.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("no staged changes found")
	}

	// Generate commit message
	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}

	message, err := provider.GenerateCommitMessage(cmd.Context(), diff)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	fmt.Println("Generated commit message:")
	fmt.Println("---")
	fmt.Println(message)
	fmt.Println("---")

	if dryRun {
		return nil
	}

	// Commit the changes
	if err := gitRepo.Commit(message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	fmt.Println("✓ Changes committed successfully")
	return nil
}
