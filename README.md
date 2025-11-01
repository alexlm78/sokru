# Sokru

> A powerful CLI tool for managing dotfiles across multiple operating systems

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](test/)
[![Coverage](https://img.shields.io/badge/coverage-90%25-brightgreen.svg)](test/)

Sokru (command: `sok`) helps you manage your dotfiles configuration across Linux, macOS, and Windows. It handles symlink creation, automatic backups, rollback on errors, and supports multi-OS configurations in a single file.

## Features

- **Multi-OS Support** - Define different symlinks for Linux, macOS (Intel/ARM), and Windows
- **Automatic Backups** - Every change creates a timestamped backup you can restore
- **Automatic Rollback** - If anything fails, changes are automatically reverted
- **Internationalization** - Full support for English and Spanish (more coming!)
- **Dry-Run Mode** - Test changes before applying them
- **OS-Specific & Common Configs** - Share common dotfiles and override per-OS
- **Type-Safe Configuration** - YAML-based configuration with validation
- **Comprehensive Testing** - 90%+ test coverage with extensive test suite

## Quick Start

### Installation

```bash
# Download latest release from GitHub releases
# Or build from source:
git clone https://github.com/alexlm78/sokru.git
cd sokru
make mac      # macOS ARM64
make macx86   # macOS Intel
make lin      # Linux AMD64
make win      # Windows AMD64
```

### Initialize Configuration

```bash
# Create configuration directory and default config
sok init

# Configure your dotfiles directory
sok config dotDir ~/dotfiles

# Set your symlinks file location
sok config symlinks ~/dotfiles/symlinks.yaml

# Set your language preference (en/es)
sok config language en
```

### Create Your Symlinks Configuration

Create `~/dotfiles/symlinks.yaml`:

```yaml
# Common dotfiles for all operating systems
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
    ~/.bashrc: ~/.dotfiles/bash/bashrc

  # Linux-specific
  linux:
    ~/.config/i3/config: ~/.dotfiles/i3/config
    ~/.xinitrc: ~/.dotfiles/x11/xinitrc

  # macOS-specific
  darwin:
    ~/Library/Application Support/Code/User/settings.json: ~/.dotfiles/vscode/settings.json
    ~/.hammerspoon/init.lua: ~/.dotfiles/hammerspoon/init.lua
```

### Install Symlinks

```bash
# Preview changes (dry-run)
sok symlinks install --dry-run

# Install symlinks
sok symlinks install

# List current symlinks
sok symlinks list
```

## Commands

### Core Commands

```bash
sok init                      # Initialize configuration
sok apply                     # Apply configuration changes
sok version                   # Show version information
sok help                      # Show help message
```

### Configuration Management

```bash
sok config show               # Show current configuration
sok config dotDir <path>      # Set dotfiles directory
sok config symlinks <path>    # Set symlinks configuration file
sok config os <os>            # Set target OS (linux/darwin/windows)
sok config language <lang>    # Set language (en/es)
sok config verbose <bool>     # Enable/disable verbose output
sok config dryRun <bool>      # Enable/disable dry-run mode
```

### Symlink Management

```bash
sok symlinks install          # Install symlinks from configuration
sok symlinks uninstall        # Remove all managed symlinks
sok symlinks list             # List configured symlinks (filtered by OS)
```

### Backup & Restore

```bash
sok restore list              # List all available backups
sok restore apply <id>        # Restore from a specific backup
sok restore delete <id>       # Delete a specific backup
```

## Configuration

Configuration is stored in `~/.sokru/config.yaml`:

```yaml
dotfiles_dir: /home/user/dotfiles
symlinks_file: /home/user/dotfiles/symlinks.yaml
os: linux                     # linux, darwin, or windows
language: en                  # en or es
verbose: false
dry_run: false
```

## Symlinks Configuration Formats

### Format 1: Legacy (Backward Compatible)

```yaml
- link:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    ~/.vimrc: ~/.dotfiles/vim/vimrc
```

### Format 2: Common + OS-Specific

```yaml
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
  linux:
    ~/.config/i3/config: ~/.dotfiles/i3/config
  darwin:
    ~/Library/Preferences/plist: ~/.dotfiles/macos/plist
```

### Format 3: OS Filter

```yaml
- os: linux
  link:
    ~/.config/polybar/config: ~/.dotfiles/polybar/config

- os: darwin
  link:
    ~/.yabairc: ~/.dotfiles/yabai/yabairc
```

See [docs/MULTI_OS_SYMLINKS.md](docs/MULTI_OS_SYMLINKS.md) for complete documentation.

## Safety Features

### Automatic Backups

Before making any changes, Sokru automatically backs up existing files and symlinks:

```bash
$ sok symlinks install
→ Backing up: ~/.bashrc
→ Backing up: ~/.vimrc
✓ Symlink created: ~/.bashrc -> ~/.dotfiles/bash/bashrc
✓ Symlink created: ~/.vimrc -> ~/.dotfiles/vim/vimrc
✓ Backup completed: 20241101-143022.123
```

### Automatic Rollback

If an error occurs, all changes are automatically reverted:

```bash
$ sok symlinks install
✓ Symlink created: ~/.bashrc -> ~/.dotfiles/bash/bashrc
✓ Symlink created: ~/.vimrc -> ~/.dotfiles/vim/vimrc
✗ Error creating symlink: permission denied

⚠ Error occurred, starting rollback of 2 action(s)...
✓ Rollback completed successfully
```

See [docs/ROLLBACK.md](docs/ROLLBACK.md) and [docs/BACKUP_RESTORE.md](docs/BACKUP_RESTORE.md) for details.

## Internationalization

Sokru supports multiple languages:

```bash
# Set to Spanish
sok config language es

# Set to English
sok config language en
```

See [docs/I18N.md](docs/I18N.md) for complete i18n documentation.

## Development

### Building

```bash
# Build for current platform
make mac       # macOS ARM64 (M1/M2/M3)
make macx86    # macOS Intel
make lin       # Linux AMD64
make win       # Windows AMD64

# Build for all platforms
make release

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
go test ./test/...

# Run tests with coverage
go test -cover ./test/...

# Run tests with verbose output
go test -v ./test/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

See [docs/TESTING.md](docs/TESTING.md) for complete testing documentation.

## Documentation

- [Multi-OS Symlinks Guide](docs/MULTI_OS_SYMLINKS.md) - Configure dotfiles for multiple operating systems
- [Backup & Restore Guide](docs/BACKUP_RESTORE.md) - Automatic backups and restore system
- [Rollback Mechanism](docs/ROLLBACK.md) - Automatic rollback on errors
- [Internationalization](docs/I18N.md) - Multi-language support
- [Testing Guide](docs/TESTING.md) - Test suite and coverage

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Run linter: `golangci-lint run`
7. Submit a pull request

## License

MIT License - see LICENSE file for details

## Author

Alejandro Lopez Monzon ([@alexlm78](https://github.com/alexlm78))

---

**Note:** This project was previously named "Sokru" and the command is `sok` for brevity.
