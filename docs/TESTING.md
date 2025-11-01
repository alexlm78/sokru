# Testing Guide

## Overview

Sokru uses Go's built-in testing framework. All tests are centralized in the `test/` directory for easy maintenance and organization.

## Test Structure

All unit tests are located in the `test/` directory:

- **test/config_test.go** - Configuration management tests
- **test/utils_test.go** - Utility function tests
- **test/symlinks_test.go** - Symlink configuration and OS filtering tests
- **test/i18n_test.go** - Internationalization tests

## Test Coverage

Tests cover critical functionality across all packages:

- **Configuration management** - 7 test functions
- **Utility functions** - 3 test functions with 17 test cases
- **Symlink OS filtering** - 2 test functions with 14 scenarios
- **Internationalization** - 10 test functions with full coverage

## Running Tests

### Run All Tests

```bash
# Run all tests
go test ./test/...

# Or from project root
go test ./...
```

### Run Tests with Verbose Output

```bash
go test ./test/... -v
```

### Run Tests with Coverage

```bash
# Coverage for all packages
go test ./... -cover
```

### Generate Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Run Tests with Race Detection

```bash
go test -race ./...
```

### Run Specific Test File

```bash
# Run config tests only
go test ./test -run TestGetDefaultConfig -v

# Run symlinks tests only
go test ./test -run TestSymlinkConfig -v

# Run i18n tests only
go test ./test -run TestGetInstance -v
```

### Run Specific Test

```bash
go test ./test -run TestGetDefaultConfig -v
```

## Test Files

### test/config_test.go

Tests for configuration management (`internal/config` package):

- `TestGetDefaultConfig` - Validates default configuration values
- `TestSaveAndLoadConfig` - Tests saving and loading configuration
- `TestLoadConfigNonExistent` - Tests loading when config doesn't exist
- `TestUpdateConfig` - Tests updating configuration
- `TestGetAndSetConfig` - Tests global config getter/setter
- `TestSaveConfigNil` - Tests error handling for nil config
- `TestConfigPathGeneration` - Tests config path generation

**Total**: 7 test functions

### test/i18n_test.go

Tests for internationalization (`internal/i18n` package):

- `TestGetInstance` - Tests singleton pattern
- `TestSetAndGetLanguage` - Tests language switching
- `TestTranslation` - Tests message translation in multiple languages
- `TestFormattedMessages` - Tests Success/Error/Info/Warning formatting
- `TestGlobalHelperFunctions` - Tests global helper functions
- `TestMissingTranslation` - Tests fallback for missing translations
- `TestLanguageFallback` - Tests fallback to English
- `TestMessagePrefixes` - Tests message prefix constants
- `TestAllEnglishMessagesExist` - Validates English message completeness
- `TestAllSpanishMessagesExist` - Validates Spanish message completeness

**Total**: 10 test functions

### test/utils_test.go

Tests for utility functions (`cmd` package):

- `TestExpandPath` - Tests tilde expansion in paths (7 scenarios)
- `TestValidateOS` - Tests OS validation (10 scenarios)
- `TestExpandPathWithCustomHome` - Tests path expansion with custom HOME (2 scenarios)

**Total**: 3 test functions, 19 test cases

### test/symlinks_test.go

Tests for symlink configuration (`cmd` package):

- `TestSymlinkConfig_GetLinksForOS` - Comprehensive tests for multi-OS filtering
  - Legacy format support
  - Common-only configurations
  - OS-specific configurations (linux, darwin, windows)
  - Priority and override behavior
  - OS filtering
  - Empty configurations
- `TestSymlinkConfig_GetLinksForOS_AllOperatingSystems` - Tests all OS types

**Total**: 2 test functions, 14 test scenarios

## Linting

Configuration file: `.golangci.yml`

Enabled linters:

- errcheck - Check for unchecked errors
- gosimple - Simplify code
- govet - Vet examines Go source code
- ineffassign - Detect ineffectual assignments
- staticcheck - Static analysis
- unused - Check for unused code
- gofmt - Check formatting
- goimports - Check imports
- misspell - Check for misspelled words
- revive - Fast, configurable linter
- gosec - Security checker
- unconvert - Remove unnecessary type conversions
- unparam - Find unused function parameters
- goconst - Find repeated strings
- gocyclo - Check cyclomatic complexity
- dupl - Check for duplicate code

Run linting locally:

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

## Writing Tests

### Test File Naming

Test files should be named `*_test.go` and placed in the same package as the code being tested.

### Test Function Naming

Test functions should start with `Test` followed by the function/feature name:

```go
func TestFunctionName(t *testing.T) {
    // Test code
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestExpandPath(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Path with tilde",
            input:    "~/dotfiles",
            expected: filepath.Join(homeDir, "dotfiles"),
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := expandPath(tt.input)
            if result != tt.expected {
                t.Errorf("expandPath(%q) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        })
    }
}
```

### Testing with Temporary Files

Use `os.MkdirTemp` for tests that need file system operations:

```go
func TestWithTempDir(t *testing.T) {
    tempDir, err := os.MkdirTemp("", "sokru-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    // Test code using tempDir
}
```

### Testing Environment Variables

Save and restore environment variables:

```go
func TestWithEnvVar(t *testing.T) {
    originalHome := os.Getenv("HOME")
    defer os.Setenv("HOME", originalHome)

    os.Setenv("HOME", "/custom/home")
    // Test code
}
```

## Best Practices

### 1. Test Independence

Each test should be independent and not rely on other tests:

```go
// Good
func TestFeatureA(t *testing.T) {
    // Setup
    // Test
    // Cleanup
}

// Bad - depends on TestFeatureA
func TestFeatureB(t *testing.T) {
    // Assumes TestFeatureA ran first
}
```

### 2. Use Subtests

Use `t.Run()` for organizing related tests:

```go
func TestFeature(t *testing.T) {
    t.Run("scenario1", func(t *testing.T) {
        // Test scenario 1
    })

    t.Run("scenario2", func(t *testing.T) {
        // Test scenario 2
    })
}
```

### 3. Clear Error Messages

Provide clear error messages that help debug failures:

```go
// Good
if result != expected {
    t.Errorf("expandPath(%q) = %q, want %q", input, result, expected)
}

// Bad
if result != expected {
    t.Error("test failed")
}
```

### 4. Test Edge Cases

Always test edge cases:

```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"normal case", "~/file", "/home/user/file"},
    {"empty string", "", ""},
    {"just tilde", "~", "~"},
    {"tilde not at start", "/path~/file", "/path~/file"},
}
```

### 5. Use Test Helpers

Extract common test setup into helper functions:

```go
func setupTestConfig(t *testing.T) (*Config, func()) {
    t.Helper()

    tempDir, _ := os.MkdirTemp("", "test-*")
    cfg := &Config{/* ... */}

    cleanup := func() {
        os.RemoveAll(tempDir)
    }

    return cfg, cleanup
}

func TestSomething(t *testing.T) {
    cfg, cleanup := setupTestConfig(t)
    defer cleanup()

    // Test code
}
```

## Coverage Goals

- **Critical packages** (config, i18n): 90%+ coverage
- **Core logic** (symlinks, utils): 80%+ coverage
- **CLI commands**: 50%+ coverage (harder to test, focus on logic)

## Contributing Tests

When contributing new features:

1. Write tests first (TDD approach recommended)
2. Ensure all tests pass: `go test ./...`
3. Check coverage: `go test -cover ./...`
4. Run linter: `golangci-lint run`

## Troubleshooting

### Tests Fail on CI but Pass Locally

- Check OS-specific behavior (path separators, etc.)
- Verify environment variables
- Check file permissions
- Review race conditions with `-race` flag

### Coverage Not Updating

```bash
# Clear test cache
go clean -testcache

# Run tests again
go test -cover ./...
```

### Linter Errors

```bash
# See detailed linter output
golangci-lint run --verbose

# Auto-fix issues where possible
golangci-lint run --fix
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [golangci-lint](https://golangci-lint.run/)

## See Also

- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [Multi-OS Symlinks Guide](MULTI_OS_SYMLINKS.md)
- [Internationalization Guide](I18N.md)
