# git-auto-commit

[![CI](https://github.com/algernon-coop/git-auto-commit/workflows/CI/badge.svg)](https://github.com/algernon-coop/git-auto-commit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/algernon-coop/git-auto-commit)](https://goreportcard.com/report/github.com/algernon-coop/git-auto-commit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A cross-platform CLI tool that automatically generates meaningful commit messages from your staged changes using AI.

## Features

- 🤖 **Multiple AI Providers**: Supports OpenAI (native & Azure), Anthropic Claude, and GitHub Models
- 🌍 **Cross-Platform**: Works on Windows, Linux, and macOS
- 📝 **Conventional Commits**: Generates commit messages following the conventional commit format
- 🔒 **Secure**: API keys stored in local config file with restricted permissions
- 🎯 **Simple to Use**: Interactive configuration wizard and straightforward CLI
- ⚡ **Fast**: Written in Go for optimal performance

## Installation

### Windows (MSI Installer)

Download and run the MSI installer from the [releases page](https://github.com/algernon-coop/git-auto-commit/releases):

- `git-auto-commit-windows-amd64.msi` for 64-bit Intel/AMD processors
- `git-auto-commit-windows-arm64.msi` for ARM64 processors

The MSI installer will:
- Install the binary to `C:\Program Files\Algernon Coop\Git Auto Commit\`
- Automatically add the installation directory to your system PATH
- Allow easy uninstallation via Windows Settings

### From Source

```bash
go install github.com/algernon-coop/git-auto-commit@latest
```

### Using Go Get

```bash
git clone https://github.com/algernon-coop/git-auto-commit.git
cd git-auto-commit
go build -o git-auto-commit .
```

### Pre-built Binaries

Download pre-built binaries or MSI installers for Windows from the [releases page](https://github.com/algernon-coop/git-auto-commit/releases).

## Quick Start

### 1. Configure AI Provider

Run the interactive configuration wizard:

```bash
git-auto-commit configure
```

This will prompt you to select an AI provider and enter your credentials:

- **OpenAI (native)**: Requires an OpenAI API key
- **OpenAI (Azure)**: Requires Azure OpenAI endpoint, API key, and deployment name
- **Anthropic Claude**: Requires an Anthropic API key
- **GitHub Models**: Requires a GitHub token

### 2. Stage Your Changes

```bash
git add .
```

### 3. Generate and Commit

```bash
git-auto-commit
```

The tool will:
1. Analyze your staged changes
2. Generate a meaningful commit message using AI
3. Display the message for review
4. Commit the changes

### Dry Run Mode

To preview the commit message without committing:

```bash
git-auto-commit --dry-run
```

## Configuration

The configuration is stored in `~/.git-auto-commit.yaml` by default. You can specify a custom config file:

```bash
git-auto-commit --config /path/to/config.yaml
```

### Configuration File Format

```yaml
provider: openai
openai:
  api_key: your-api-key-here
  model: gpt-4
```

## Supported AI Providers

### OpenAI (Native)

```yaml
provider: openai
openai:
  api_key: sk-...
  model: gpt-4  # or gpt-3.5-turbo, gpt-4-turbo, etc.
```

### Azure OpenAI

```yaml
provider: azure
azure:
  endpoint: https://your-resource.openai.azure.com
  api_key: your-azure-api-key
  deployment: your-deployment-name
```

### Anthropic Claude

```yaml
provider: claude
claude:
  api_key: sk-ant-...
  model: claude-3-5-sonnet-20241022  # or other Claude models
```

### GitHub Models

```yaml
provider: github
github:
  token: ghp_...
  model: gpt-4o  # or other available models
```

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
git clone https://github.com/algernon-coop/git-auto-commit.git
cd git-auto-commit
go mod download
go build -o git-auto-commit .
```

### Running Tests

```bash
go test ./... -v
```

### Using Dev Container

This project includes a dev container configuration for VS Code:

1. Install the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
2. Open the project in VS Code
3. Click "Reopen in Container" when prompted

### Linting

```bash
golangci-lint run
```

## CI/CD

The project uses GitHub Actions for continuous integration:

- **Test**: Runs tests on Linux, Windows, and macOS with multiple Go versions
- **Lint**: Runs golangci-lint to ensure code quality
- **Build**: Creates binaries for all supported platforms

## Dependencies

Dependencies are managed with Go modules and automatically updated weekly via Dependabot.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by the need for better commit messages
- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses various AI providers for intelligent commit message generation