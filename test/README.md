# Test Directory

This directory contains all unit tests for the Sokru project.

## Structure

- `config_test.go` - Tests for configuration management
- `utils_test.go` - Tests for utility functions
- `symlinks_test.go` - Tests for symlink configuration and OS filtering
- `i18n_test.go` - Tests for internationalization

## Running Tests

```bash
# Run all tests
go test ./test/...

# Run with verbose output
go test ./test/... -v

# Run with coverage (note: coverage is calculated per-package)
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Coverage

To get accurate coverage metrics, run tests from the project root:

```bash
go test ./... -cover
```

This will show coverage for each package:

- `cmd` - Command-line interface
- `internal/config` - Configuration management
- `internal/i18n` - Internationalization

## Note on Test Package

Tests are in a separate `test` package to:

1. Centralize all tests in one location
2. Test the public API of packages
3. Ensure proper encapsulation
4. Make it easier to find and maintain tests

Some functions are exported with `ForTesting` suffix to allow testing from external packages while keeping the internal API clean.
