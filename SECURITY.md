# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within git-auto-commit, please send an email to the maintainers. All security vulnerabilities will be promptly addressed.

Please do not open public issues for security vulnerabilities.

### What to Include

- Type of vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### What to Expect

- Acknowledgment within 48 hours
- Regular updates on the progress
- Credit in the security advisory (if desired)

## Security Best Practices

When using git-auto-commit:

1. **API Keys**: Store your API keys securely in the config file (`~/.git-auto-commit.yaml`), which is created with restricted permissions (0600)
2. **Never Commit API Keys**: The default `.gitignore` excludes config files. Never commit API keys to version control
3. **Keep Updated**: Regularly update to the latest version to get security patches
4. **Review Generated Messages**: Always review AI-generated commit messages before they are committed
5. **Use Environment Variables**: For CI/CD or shared environments, consider using environment variables instead of config files

## Known Security Considerations

- API keys are stored in plaintext in the config file (encrypted storage may be added in future versions)
- All API communications use HTTPS
- No telemetry or analytics data is collected
- The tool only reads staged git changes and does not modify files
