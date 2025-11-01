# Contributing to Sokru

Thank you for your interest in contributing to Sokru! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive experience for everyone. We expect all contributors to:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment, trolling, or discriminatory comments
- Personal attacks or political arguments
- Publishing others' private information without permission
- Other conduct that would be inappropriate in a professional setting

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Git
- Make (recommended)
- golangci-lint (for linting)

### Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/sokru.git
cd sokru

# Add upstream remote
git remote add upstream https://github.com/alexlm78/sokru.git
```

## Development Setup

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install golangci-lint (optional but recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build and Test

```bash
# Build for your platform
make mac      # macOS ARM64
make macx86   # macOS Intel
make lin      # Linux
make win      # Windows

# Run tests
go test ./test/...

# Run tests with coverage
go test -cover ./test/...

# Run linter
golangci-lint run
```

### Project Structure

```proj
sokru/
├── cmd/               # Command implementations
├── internal/          # Internal packages
│   ├── config/       # Configuration management
│   ├── i18n/         # Internationalization
│   ├── backup/       # Backup system
│   └── rollback/     # Rollback mechanism
├── test/             # Test files
├── docs/             # Documentation
├── main.go           # Entry point
└── Makefile          # Build automation
```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug reports** - Found a bug? Report it!
- **Feature requests** - Have an idea? Suggest it!
- **Bug fixes** - Fix a bug and submit a PR
- **New features** - Implement a feature from the roadmap
- **Documentation** - Improve or add documentation
- **Translations** - Add new language translations
- **Tests** - Improve test coverage

### Reporting Bugs

Before submitting a bug report:

1. Check if the bug has already been reported in [Issues](https://github.com/alexlm78/sokru/issues)
2. Ensure you're using the latest version
3. Try to reproduce with minimal configuration

When reporting, include:

- Sokru version (`sok version`)
- Operating system and version
- Go version (if building from source)
- Steps to reproduce the bug
- Expected behavior vs actual behavior
- Error messages or logs
- Configuration files (sanitize sensitive data)

**Example bug report:**

```markdown
**Environment:**
- Sokru version: 1.0.0
- OS: macOS 13.0 (M1)
- Go version: 1.20

**Steps to reproduce:**
1. Run `sok symlinks install`
2. Error occurs when...

**Expected behavior:**
Symlinks should be created successfully.

**Actual behavior:**
Error: permission denied

**Error message:**

```output
✗ Error creating symlink: permission denied
```

**Configuration:**
(attach relevant config files)

```feat

### Suggesting Features

Before suggesting a feature:

1. Check the [Roadmap](docs/ROADMAP.md) - it might already be planned
2. Search [Issues](https://github.com/alexlm78/sokru/issues) for similar suggestions
3. Consider if it aligns with Sokru's goals (simple, safe dotfiles management)

When suggesting, include:

- Clear use case - why is this needed?
- Proposed solution - how should it work?
- Alternatives considered
- Willingness to implement

**Example feature request:**

```markdown
**Feature:** Git integration for dotfiles synchronization

**Use case:**
As a user with multiple machines, I want to sync my dotfiles via git
automatically when I run `sok apply`.

**Proposed solution:**
Add `sok git sync` command that:
1. Pulls latest changes from remote
2. Applies symlinks
3. Commits and pushes any local changes

**Alternatives:**
- Manual git operations before/after sok commands
- Shell scripts wrapping sok commands

**Willing to implement:** Yes
```

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

Key points:

- Use `gofmt` to format code
- Follow Go naming conventions
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors properly - don't ignore them

### Code Style Examples

**Good:**

```go
// ExpandPath expands ~ in paths to the user's home directory
func expandPath(path string) (string, error) {
    if !strings.HasPrefix(path, "~") {
        return path, nil
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get home directory: %w", err)
    }

    return filepath.Join(homeDir, path[1:]), nil
}
```

**Bad:**

```go
// expand path
func exp(p string) string {
    h, _ := os.UserHomeDir()  // Error ignored!
    return strings.Replace(p, "~", h, 1)  // Brittle string replacement
}
```

### Error Handling

Always handle errors properly:

```go
// Good
config, err := config.LoadConfig()
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Bad - never ignore errors
config, _ := config.LoadConfig()
```

### Comments

- Use complete sentences in comments
- Start comments with the name of the thing being described
- Document why, not what (code shows what)

```go
// LoadConfig reads the configuration from ~/.config/sokru/config.yaml.
// If the file doesn't exist, it returns the default configuration.
func LoadConfig() (*Config, error) {
    // Implementation
}
```

## Testing Guidelines

### Writing Tests

- Write tests for all new features
- Write tests when fixing bugs
- Use table-driven tests for multiple scenarios
- Test edge cases and error conditions
- Aim for 80%+ code coverage

### Test Structure

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "",
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Feature(tt.input)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }

            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./test/...

# Run specific test
go test ./test -run TestFeature

# Run with coverage
go test -cover ./test/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detection
go test -race ./test/...
```

See [Testing Guide](docs/TESTING.md) for more details.

## Documentation

### Types of Documentation

1. **Code comments** - Document exported functions and types
2. **README** - User-facing documentation
3. **docs/** - Detailed guides and references
4. **CLAUDE.md** - AI assistant context

### Writing Documentation

- Use clear, concise language
- Include examples
- Update documentation when changing code
- Check for broken links
- Use proper markdown formatting

### Documentation Style

- Use headings for structure (##, ###)
- Use code blocks with language tags (\`\`\`bash, \`\`\`go)
- Use lists for steps or options
- Include command examples with expected output

## Pull Request Process

### Before Submitting

1. **Sync with upstream:**

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests:**

   ```bash
   go test ./test/...
   golangci-lint run
   ```

3. **Update documentation:**
   - Update README if adding user-facing features
   - Update relevant docs/ files
   - Add/update code comments

4. **Commit messages:**

   ```bash
   # Good commit messages
   feat: add git integration for dotfiles sync
   fix: correct symlink path expansion on Windows
   docs: update installation guide for Windows
   test: add tests for backup compression

   # Include details in commit body
   git commit -m "feat: add git integration" -m "
   - Add sok git sync command
   - Implement auto-commit on changes
   - Add configuration for git remote

   Closes #123
   "
   ```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `test:` - Adding or updating tests
- `refactor:` - Code change that neither fixes a bug nor adds a feature
- `perf:` - Performance improvement
- `chore:` - Maintenance tasks

### Creating Pull Request

1. **Push to your fork:**

   ```bash
   git push origin feature-branch
   ```

2. **Open PR on GitHub:**
   - Clear title describing the change
   - Reference related issues (Closes #123)
   - Describe what changed and why
   - Include screenshots if UI changes
   - Mention any breaking changes

3. **PR template:**

   ```markdown
   ## Description
   Brief description of changes

   ## Motivation
   Why is this change needed?

   ## Changes
   - Added feature X
   - Fixed bug Y
   - Updated documentation

   ## Testing
   - [ ] Unit tests pass
   - [ ] Manual testing completed
   - [ ] Documentation updated

   ## Related Issues
   Closes #123

   ## Breaking Changes
   None / List any breaking changes
   ```

### Review Process

1. Maintainers will review your PR
2. Address feedback and push updates
3. Maintainer will merge when approved
4. PR may be squashed or rebased

### After Merge

1. Delete your feature branch
2. Update your fork:

   ```bash
   git checkout main
   git pull upstream main
   git push origin main
   ```

## Issue Guidelines

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `question` - Further information requested
- `wontfix` - Will not be worked on

### Claiming Issues

Comment on an issue to claim it:

```markdown
I'd like to work on this. My approach would be:
1. ...
2. ...

Estimated time: 1 week
```

## Adding Translations

To add a new language (example: French):

1. **Create message file:**

   ```bash
   touch internal/i18n/messages_fr.go
   ```

2. **Implement translations:**

   ```go
   package i18n

   func getFrenchMessages() map[MessageKey]string {
       return map[MessageKey]string{
           MsgErrorLoadingConfig: "Erreur lors du chargement: %v",
           // ... all message keys
       }
   }
   ```

3. **Register language:**
   - Add constant in `internal/i18n/i18n.go`
   - Register in `loadMessages()`
   - Update validation in `cmd/config.go`

4. **Add tests:**
   - Update `test/i18n_test.go`
   - Verify all messages are translated

5. **Update docs:**
   - Add to `docs/I18N.md`
   - Update README

See [Internationalization Guide](docs/I18N.md) for details.

## Getting Help

Need help with contribution?

- **Discord:** (coming soon)
- **GitHub Discussions:** Ask questions
- **Email:** <alejandro@kreaker.dev>

## Recognition

Contributors will be:

- Listed in CONTRIBUTORS file
- Mentioned in release notes
- Thanked in pull requests

Significant contributors may be invited as maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You

Your contributions help make Sokru better for everyone. Thank you for taking the time to contribute!

## See Also

- [Main README](README.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Testing Guide](docs/TESTING.md)
- [Roadmap](docs/ROADMAP.md)
