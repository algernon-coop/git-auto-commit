# Contributing to git-auto-commit

Thank you for your interest in contributing to git-auto-commit! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in your interactions with other contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear title and description
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Your environment (OS, Go version, etc.)

### Suggesting Features

Feature suggestions are welcome! Please open an issue with:
- A clear title and description
- Use case for the feature
- Any relevant examples or mockups

### Pull Requests

1. Fork the repository
2. Create a new branch from `main` for your changes
3. Make your changes
4. Add or update tests as needed
5. Ensure all tests pass: `go test ./...`
6. Run the linter: `golangci-lint run`
7. Commit your changes using conventional commits
8. Push to your fork and submit a pull request

#### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(llm): add support for GPT-4 Turbo
fix(config): handle missing config file gracefully
docs(readme): update installation instructions
```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- golangci-lint (for linting)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/git-auto-commit.git
cd git-auto-commit

# Install dependencies
go mod download

# Build
go build -o git-auto-commit .

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Using Dev Container

For a consistent development environment, use the provided dev container:

1. Install Docker and VS Code
2. Install the Dev Containers extension
3. Open the project in VS Code
4. Click "Reopen in Container"

## Testing

- Write tests for new functionality
- Ensure existing tests still pass
- Aim for good test coverage
- Use table-driven tests where appropriate

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detection
go test -race ./...
```

## Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (automatically done by most editors)
- Run `golangci-lint` before submitting PRs
- Write clear, self-documenting code
- Add comments for exported functions and types

## Project Structure

```
.
├── .devcontainer/       # Dev container configuration
├── .github/             # GitHub workflows and configurations
├── internal/
│   ├── cmd/            # CLI commands
│   ├── config/         # Configuration management
│   ├── git/            # Git operations
│   └── llm/            # AI provider implementations
├── main.go             # Entry point
├── go.mod              # Go module definition
└── README.md           # Project documentation
```

## Adding a New AI Provider

To add support for a new AI provider:

1. Create a new file in `internal/llm/` (e.g., `newprovider.go`)
2. Implement the `Provider` interface
3. Add configuration struct in `internal/config/config.go`
4. Update `NewProvider` function in `internal/llm/provider.go`
5. Update the configure command in `internal/cmd/configure.go`
6. Add tests for the new provider
7. Update documentation

## Questions?

If you have questions, feel free to:
- Open an issue for discussion
- Check existing issues and pull requests
- Review the documentation

Thank you for contributing!
