# Sokru Installation Guide

This guide covers multiple methods to install Sokru on different operating systems.

## Quick Install

### Using Pre-built Binaries (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/alexlm78/sokru/releases):

```bash
# Linux (AMD64)
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-linux-amd64
chmod +x sok-linux-amd64
sudo mv sok-linux-amd64 /usr/local/bin/sok

# macOS (ARM64 - M1/M2/M3)
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-darwin-arm64
chmod +x sok-darwin-arm64
sudo mv sok-darwin-arm64 /usr/local/bin/sok

# macOS (Intel)
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-darwin-amd64
chmod +x sok-darwin-amd64
sudo mv sok-darwin-amd64 /usr/local/bin/sok

# Windows (AMD64)
# Download from GitHub releases and add to PATH
```

### Using curl/wget One-liner

```bash
# Linux/macOS (auto-detect architecture)
curl -sSL https://raw.githubusercontent.com/alexlm78/sokru/main/install.sh | bash

# Or with wget
wget -qO- https://raw.githubusercontent.com/alexlm78/sokru/main/install.sh | bash
```

## Platform-Specific Installation

### Linux

#### Debian/Ubuntu (APT)

```bash
# Add repository (coming soon)
# sudo add-apt-repository ppa:alexlm78/sokru
# sudo apt update
# sudo apt install sokru

# For now, use pre-built binary
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-linux-amd64
chmod +x sok-linux-amd64
sudo mv sok-linux-amd64 /usr/local/bin/sok
```

#### Fedora/RHEL/CentOS (DNF/YUM)

```bash
# Using pre-built binary
sudo wget -O /usr/local/bin/sok https://github.com/alexlm78/sokru/releases/latest/download/sok-linux-amd64
sudo chmod +x /usr/local/bin/sok
```

#### Arch Linux (AUR)

```bash
# Coming soon: AUR package
# yay -S sokru

# For now, use pre-built binary
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-linux-amd64
chmod +x sok-linux-amd64
sudo mv sok-linux-amd64 /usr/local/bin/sok
```

#### Snap Package

```bash
# Coming soon
# sudo snap install sokru
```

### macOS

#### Homebrew

```bash
# Coming soon: Homebrew formula
# brew install alexlm78/tap/sokru

# For now, use pre-built binary
# M1/M2/M3 (ARM)
brew install wget
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-darwin-arm64
chmod +x sok-darwin-arm64
sudo mv sok-darwin-arm64 /usr/local/bin/sok

# Intel
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-darwin-amd64
chmod +x sok-darwin-amd64
sudo mv sok-darwin-amd64 /usr/local/bin/sok
```

#### MacPorts

```bash
# Coming soon
# sudo port install sokru
```

### Windows

#### Chocolatey

```powershell
# Coming soon
# choco install sokru

# For now, manual installation:
# 1. Download sok.exe from GitHub releases
# 2. Add to PATH environment variable
```

#### Scoop

```powershell
# Coming soon
# scoop bucket add alexlm78 https://github.com/alexlm78/scoop-bucket
# scoop install sokru
```

#### Manual Installation

1. Download `sok-windows-amd64.exe` from [GitHub Releases](https://github.com/alexlm78/sokru/releases)
2. Rename to `sok.exe`
3. Move to a directory in your PATH (e.g., `C:\Windows\System32` or `C:\Program Files\sok\`)

## Build from Source

### Prerequisites

- Go 1.19 or higher
- Git
- Make (optional, but recommended)

### Clone and Build

```bash
# Clone repository
git clone https://github.com/alexlm78/sokru.git
cd sokru

# Build for your platform
make mac       # macOS ARM64
make macx86    # macOS Intel
make lin       # Linux AMD64
make win       # Windows AMD64

# Or build all platforms
make release

# Built binaries will be in build/ directory
```

### Manual Build (without Make)

```bash
# Clone repository
git clone https://github.com/alexlm78/sokru.git
cd sokru

# Build for current platform
go build -o sok

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o sok-linux-amd64
GOOS=darwin GOARCH=arm64 go build -o sok-darwin-arm64
GOOS=darwin GOARCH=amd64 go build -o sok-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o sok-windows-amd64.exe

# Install to /usr/local/bin
sudo cp sok /usr/local/bin/

# Or add to PATH manually
export PATH=$PATH:$(pwd)
```

### Development Build

```bash
git clone https://github.com/alexlm78/sokru.git
cd sokru

# Install dependencies
go mod download

# Run tests
go test ./test/...

# Build with debug info
go build -gcflags="all=-N -l" -o sok

# Run directly without installing
go run main.go --help
```

## Docker Installation

```bash
# Coming soon: Official Docker image
# docker pull alexlm78/sokru:latest
# docker run -v ~/.dotfiles:/dotfiles -v ~/.config/sokru:/root/.config/sokru alexlm78/sokru sok --help
```

## Verify Installation

After installation, verify Sokru is working:

```bash
# Check version
sok version

# Check help
sok help

# Initialize configuration
sok init
```

Expected output:

```shell
$ sok version
Sokru version 1.0.0
```

## Post-Installation Setup

### 1. Initialize Configuration

```bash
sok init
```

This creates:

- `~/.config/sokru/` directory
- `~/.config/sokru/config.yaml` with default settings
- `~/.config/sokru/backups/` directory

### 2. Configure Dotfiles Directory

```bash
sok config dotDir ~/dotfiles
sok config symlinks ~/dotfiles/symlinks.yaml
```

### 3. Set Language (Optional)

```bash
# English (default)
sok config language en

# Spanish
sok config language es
```

### 4. Create Symlinks Configuration

Create `~/dotfiles/symlinks.yaml`:

```yaml
- common:
    ~/.bashrc: ~/dotfiles/bash/bashrc
    ~/.vimrc: ~/dotfiles/vim/vimrc
    ~/.gitconfig: ~/dotfiles/git/gitconfig
```

### 5. Test Configuration

```bash
# Dry-run to preview changes
sok symlinks install --dry-run

# If everything looks good, install
sok symlinks install
```

## Shell Completion (Optional)

### Bash

```bash
# Generate completion script
sok completion bash > /etc/bash_completion.d/sok

# Or for user only
sok completion bash > ~/.bash_completion.d/sok

# Add to ~/.bashrc
echo "source ~/.bash_completion.d/sok" >> ~/.bashrc
```

### Zsh

```bash
# Generate completion script
sok completion zsh > "${fpath[1]}/_sok"

# Or add to ~/.zshrc
echo 'eval "$(sok completion zsh)"' >> ~/.zshrc
```

### Fish

```bash
# Generate completion script
sok completion fish > ~/.config/fish/completions/sok.fish
```

Note: Shell completion is a planned feature for v1.1.

## Updating Sokru

### Using Pre-built Binaries

Download and replace the binary:

```bash
# Linux
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-linux-amd64
chmod +x sok-linux-amd64
sudo mv sok-linux-amd64 /usr/local/bin/sok

# macOS
wget https://github.com/alexlm78/sokru/releases/latest/download/sok-darwin-arm64
chmod +x sok-darwin-arm64
sudo mv sok-darwin-arm64 /usr/local/bin/sok
```

### From Source

```bash
cd sokru
git pull
make mac  # or appropriate target for your platform
sudo cp build/sok /usr/local/bin/
```

## Uninstalling

### Remove Binary

```bash
# Linux/macOS
sudo rm /usr/local/bin/sok

# Or if installed elsewhere
which sok  # Find location
sudo rm $(which sok)
```

### Remove Configuration (Optional)

```bash
# Remove configuration and backups
rm -rf ~/.config/sokru

# This removes:
# - Configuration file (~/.config/sokru/config.yaml)
# - All backups (~/.config/sokru/backups/)
```

### Uninstall Symlinks

Before uninstalling, you may want to uninstall managed symlinks:

```bash
sok symlinks uninstall
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Check if binary is in PATH:

   ```bash
   echo $PATH
   ls -l /usr/local/bin/sok
   ```

2. Add to PATH if needed:

   ```bash
   export PATH=$PATH:/usr/local/bin
   echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
   ```

3. Restart shell or source rc file:

   ```bash
   source ~/.bashrc
   ```

### Permission Denied

If you get "permission denied":

```bash
chmod +x /usr/local/bin/sok

# Or if you don't have sudo access, install to user directory
mkdir -p ~/bin
mv sok ~/bin/
export PATH=$PATH:~/bin
echo 'export PATH=$PATH:~/bin' >> ~/.bashrc
```

### macOS Gatekeeper Issues

If macOS blocks the binary:

```bash
xattr -d com.apple.quarantine /usr/local/bin/sok

# Or in System Preferences:
# Security & Privacy → Allow apps downloaded from: Anywhere
```

### Build Errors

If building from source fails:

1. Check Go version:

   ```bash
   go version  # Should be 1.19+
   ```

2. Update dependencies:

   ```bash
   go mod tidy
   go mod download
   ```

3. Clean and rebuild:

   ```bash
   make clean
   make mac
   ```

### Missing Dependencies

If you see "missing module" errors:

```bash
cd sokru
go mod download
go mod verify
```

## Getting Help

If you encounter issues:

1. Check [GitHub Issues](https://github.com/alexlm78/sokru/issues)
2. Read [Documentation](../README.md)
3. Open a new issue with:
   - Your OS and version
   - Go version (if building from source)
   - Installation method used
   - Error messages

## Next Steps

After installation:

1. Read [Quick Start Guide](../README.md#quick-start)
2. Learn about [Multi-OS Symlinks](MULTI_OS_SYMLINKS.md)
3. Understand [Backup & Restore](BACKUP_RESTORE.md)
4. Explore [Rollback Mechanism](ROLLBACK.md)

## See Also

- [Main README](../README.md) - Project overview
- [Architecture](ARCHITECTURE.md) - Technical details
- [Testing Guide](TESTING.md) - Development setup
- [Roadmap](ROADMAP.md) - Future plans
