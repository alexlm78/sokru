# Multi-OS Symlinks Configuration

## Overview

Sokru supports multi-OS symlink configurations, allowing you to define symlinks that work across different operating systems (Linux, macOS/Darwin, Windows) in a single configuration file.

## Supported Operating Systems

- `linux` - Linux distributions
- `darwin` - macOS
- `windows` - Windows

## Configuration Formats

### Format 1: Legacy (Backward Compatible)

The original format continues to work without any changes:

```yaml
- link:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    ~/.zshrc: ~/.dotfiles/zsh/zshrc
```

### Format 2: OS-Specific Sections

Define different symlinks for each operating system:

```yaml
- linux:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    /usr/local/bin/myapp: ~/.dotfiles/bin/myapp-linux
  darwin:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    /usr/local/bin/myapp: ~/.dotfiles/bin/myapp-macos
  windows:
    ~/AppData/Roaming/.bashrc: ~/.dotfiles/bash/bashrc
    ~/bin/myapp.exe: ~/.dotfiles/bin/myapp-windows.exe
```

### Format 3: Common + OS-Specific

Use `common` for symlinks that work on all OS, and OS-specific sections for platform-specific ones:

```yaml
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
  linux:
    ~/.config/i3/config: ~/.dotfiles/i3/config
  darwin:
    ~/Library/Application Support/Code/User/settings.json: ~/.dotfiles/vscode/settings.json
```

### Format 4: OS Filter

Use the `os` field to specify that an entire entry should only apply to a specific OS:

```yaml
- os: linux
  link:
    ~/.config/i3/config: ~/.dotfiles/i3/config
    ~/.xinitrc: ~/.dotfiles/x11/xinitrc

- os: darwin
  link:
    ~/Library/Preferences/com.apple.Terminal.plist: ~/.dotfiles/macos/Terminal.plist
```

## Priority Rules

When multiple formats are combined, symlinks are applied in this priority order (highest to lowest):

1. **Legacy `link` field** - Highest priority
2. **OS-specific sections** (`linux`, `darwin`, `windows`)
3. **Common section** - Lowest priority

This means OS-specific symlinks can override common ones.

## Complete Example

```yaml
# Common symlinks for all operating systems
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
    ~/.tmux.conf: ~/.dotfiles/tmux/tmux.conf
  linux:
    ~/.config/i3/config: ~/.dotfiles/i3/config
    ~/.xinitrc: ~/.dotfiles/x11/xinitrc
  darwin:
    ~/Library/Application Support/Code/User/settings.json: ~/.dotfiles/vscode/settings.json
    ~/Library/Preferences/com.googlecode.iterm2.plist: ~/.dotfiles/iterm2/com.googlecode.iterm2.plist

# Linux-only configuration
- os: linux
  link:
    ~/.config/polybar/config: ~/.dotfiles/polybar/config
    ~/.config/rofi/config.rasi: ~/.dotfiles/rofi/config.rasi

# macOS-only configuration
- os: darwin
  link:
    ~/.hammerspoon/init.lua: ~/.dotfiles/hammerspoon/init.lua
    ~/.yabairc: ~/.dotfiles/yabai/yabairc

# Shell configurations with OS-specific overrides
- common:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.common
    ~/.bash_profile: ~/.dotfiles/bash/bash_profile.common
  darwin:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.macos
    ~/.bash_profile: ~/.dotfiles/bash/bash_profile.macos
  linux:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.linux
```

## How It Works

1. **OS Detection**: Sokru reads the `os` field from your configuration (`~/.config/sokru/config.yaml`)
2. **Filtering**: When processing symlinks, it filters entries based on:
   - If an entry has an `os` field, it only processes it if it matches the current OS
   - Within each entry, it combines `common` and OS-specific sections
3. **Merging**: OS-specific symlinks override common ones if they have the same target path

## Usage Examples

### Viewing Symlinks for Current OS

```bash
# List symlinks (automatically filtered by configured OS)
sok symlinks list

# With verbose output to see filtering details
sok symlinks list --verbose
```

### Installing Symlinks

```bash
# Install symlinks for current OS
sok symlinks install

# Dry-run to see what would be installed
sok symlinks install --dry-run
```

### Changing OS Configuration

```bash
# Check current OS setting
sok config os

# Change OS (useful for testing or multi-boot systems)
sok config os linux
sok config os darwin
sok config os windows
```

## Best Practices

### 1. Use Common for Shared Configurations

Put configurations that work across all platforms in the `common` section:

```yaml
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
```

### 2. OS-Specific Paths

Use OS-specific sections for platform-dependent paths:

```yaml
- linux:
    ~/.config/systemd/user/myservice.service: ~/.dotfiles/systemd/myservice.service
  darwin:
    ~/Library/LaunchAgents/com.example.myservice.plist: ~/.dotfiles/launchd/myservice.plist
```

### 3. Organize by Category

Group related symlinks together:

```yaml
# Editor configurations
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.config/nvim/init.vim: ~/.dotfiles/nvim/init.vim

# Shell configurations
- common:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    ~/.zshrc: ~/.dotfiles/zsh/zshrc

# Window manager (Linux only)
- os: linux
  link:
    ~/.config/i3/config: ~/.dotfiles/i3/config
```

### 4. Test with Dry-Run

Always test your configuration with `--dry-run` first:

```bash
sok symlinks install --dry-run --verbose
```

## Migration from Legacy Format

If you have an existing `symlinks.yml` using the legacy format, it will continue to work without any changes. To migrate to the multi-OS format:

### Before (Legacy)

```yaml
- link:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    ~/.vimrc: ~/.dotfiles/vim/vimrc
```

### After (Multi-OS)

```yaml
- common:
    ~/.bashrc: ~/.dotfiles/bash/bashrc
    ~/.vimrc: ~/.dotfiles/vim/vimrc
```

Or with OS-specific configurations:

```yaml
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
  linux:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.linux
  darwin:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.macos
```

## Troubleshooting

### Symlinks Not Being Created

1. Check your OS configuration:

   ```bash
   sok config show
   ```

2. Verify the OS field in your symlinks.yml matches your system

3. Use verbose mode to see filtering details:

   ```bash
   sok symlinks list --verbose
   ```

### Wrong Symlinks Being Created

1. Check the priority rules - OS-specific sections override common
2. Verify you don't have conflicting entries
3. Use `sok symlinks list` to see what will be created

### Testing on Different OS

You can temporarily change the OS setting to test configurations:

```bash
# Save current OS
CURRENT_OS=$(sok config os)

# Test Linux configuration
sok config os linux
sok symlinks install --dry-run

# Test macOS configuration
sok config os darwin
sok symlinks install --dry-run

# Restore original OS
sok config os $CURRENT_OS
```

## Advanced Examples

### Multi-Boot System

For systems that dual-boot Linux and Windows:

```yaml
- common:
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
  linux:
    ~/.bashrc: ~/.dotfiles/bash/bashrc.linux
    ~/.config/i3/config: ~/.dotfiles/i3/config
  windows:
    ~/AppData/Roaming/.bashrc: ~/.dotfiles/bash/bashrc.windows
    ~/AppData/Roaming/Code/User/settings.json: ~/.dotfiles/vscode/settings.json
```

### Development Tools

Different paths for development tools on different platforms:

```yaml
- common:
    ~/.config/nvim/init.vim: ~/.dotfiles/nvim/init.vim
  linux:
    /usr/local/bin/mydev: ~/.dotfiles/bin/mydev-linux
  darwin:
    /usr/local/bin/mydev: ~/.dotfiles/bin/mydev-macos
  windows:
    ~/bin/mydev.exe: ~/.dotfiles/bin/mydev-windows.exe
```

## See Also

- [Main README](../README.md)
- [Configuration Guide](../README.md#configuration)
- [Internationalization Guide](I18N.md)
