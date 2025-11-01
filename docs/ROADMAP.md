# Sokru Roadmap

This document outlines the future development plans for Sokru, including features that could enhance dotfiles management.

## Current Status (v1.0)

### Implemented Features ✅

- Multi-OS support (Linux, macOS, Windows)
- Automatic backup system
- Automatic rollback on errors
- Internationalization (English, Spanish)
- OS-specific and common symlink configurations
- Dry-run mode
- Comprehensive test suite (90%+ coverage)
- YAML-based configuration
- Command-line interface with Cobra

## Roadmap to v2.0

### Phase 1: Enhanced Safety & Recovery (v1.1)

#### Automatic Backup Rotation

**Priority:** High
**Complexity:** Low

- Automatically clean up old backups (configurable retention)
- Keep last N backups or backups from last X days
- Warning when backups consume too much space

```yaml
# Configuration example
backup:
  retention_days: 30
  max_backups: 10
  max_size_mb: 500
```

#### Backup Compression

**Priority:** Medium
**Complexity:** Medium

- Compress backups to save disk space
- Use gzip or similar for backup files
- Transparent decompression on restore

#### Selective Restore

**Priority:** Medium
**Complexity:** Medium

- Restore specific files from a backup instead of all
- Interactive selection mode

```bash
sok restore apply <backup-id> --select
# Interactive prompt to choose files

sok restore apply <backup-id> --files ~/.bashrc,~/.vimrc
# Restore only specific files
```

### Phase 2: Git Integration (v1.2)

#### Git Repository Sync

**Priority:** High
**Complexity:** Medium

- Initialize/clone dotfiles from git repository
- Automatic commit and push after changes
- Pull latest changes before applying

```bash
sok git init https://github.com/user/dotfiles.git
sok git sync    # Pull, apply, commit, push
sok git status  # Show git status of dotfiles
```

#### Remote Configuration

**Priority:** Medium
**Complexity:** Medium

- Fetch symlinks configuration from remote repository
- Support for multiple remote configs
- Auto-update from remote

```yaml
# Configuration
dotfiles:
  remote: https://github.com/user/dotfiles.git
  auto_sync: true
  branch: main
```

### Phase 3: Advanced Configuration (v1.3)

#### Configuration Profiles

**Priority:** High
**Complexity:** Medium

- Multiple named profiles (work, home, minimal, full)
- Easy switching between profiles
- Profile inheritance

```bash
sok profile create work
sok profile use work
sok profile list
sok profile delete minimal

# In symlinks.yaml
profiles:
  minimal:
    - common:
        ~/.bashrc: ~/dotfiles/bash/bashrc
  work:
    inherits: minimal
    - common:
        ~/.ssh/config: ~/dotfiles/ssh/work_config
```

#### Template System

**Priority:** Medium
**Complexity:** High

- Generate dotfiles from templates
- Variable substitution
- Conditional sections

```yaml
# templates/bashrc.tmpl
export USER="{{ .Username }}"
export EMAIL="{{ .Email }}"

{{if .WorkProfile}}
export WORK_DIR="{{ .WorkDir }}"
{{end}}

# Configuration
templates:
  variables:
    Username: john
    Email: john@example.com
    WorkProfile: true
    WorkDir: /work/projects
```

#### Environment Variables Management

**Priority:** Medium
**Complexity:** Medium

- Manage `.env` files securely
- Encrypted storage for secrets
- Different environments (dev, staging, prod)

```bash
sok env set DATABASE_URL postgres://localhost/db --encrypt
sok env list
sok env export > .env
```

### Phase 4: Hooks & Validation (v1.4)

#### Pre/Post Hooks

**Priority:** High
**Complexity:** Medium

- Run custom scripts before/after operations
- Validation hooks
- Notification hooks

```yaml
hooks:
  pre_install:
    - script: ~/dotfiles/scripts/check_system.sh
  post_install:
    - script: ~/dotfiles/scripts/reload_shell.sh
    - notify: "Dotfiles installed successfully"
  on_error:
    - notify: "Installation failed"
    - script: ~/dotfiles/scripts/send_alert.sh
```

#### Syntax Validation

**Priority:** Medium
**Complexity:** Medium

- Validate dotfile syntax before applying
- Plugin system for validators (bash, vim, zsh, etc.)
- Lint configuration files

```bash
sok validate ~/.bashrc
sok validate --all
```

#### Conflict Detection

**Priority:** High
**Complexity:** High

- Detect conflicts before applying
- Interactive conflict resolution
- Merge strategies

```bash
sok symlinks install --interactive
# Shows conflicts and asks for resolution:
# Conflict: ~/.bashrc
#   [1] Keep existing
#   [2] Replace with new
#   [3] Show diff
#   [4] Merge
```

### Phase 5: Dependencies & Package Management (v1.5)

#### Dependency Management

**Priority:** Medium
**Complexity:** High

- Define dependencies between dotfiles
- Ordered installation based on dependencies
- Conditional dependencies (OS-specific)

```yaml
dependencies:
  ~/.vimrc:
    requires:
      - ~/.vim/colors/theme.vim
      - packages: [vim, git]
  ~/.zshrc:
    requires:
      - ~/.oh-my-zsh
    optional:
      - ~/.zsh-autosuggestions
```

#### Package Integration

**Priority:** Low
**Complexity:** High

- Detect and install required packages
- Support for multiple package managers (apt, brew, pacman, chocolatey)
- Optional package installation

```yaml
packages:
  linux:
    apt:
      - vim
      - tmux
      - git
  darwin:
    brew:
      - neovim
      - tmux
      - ripgrep
```

### Phase 6: Plugin System (v1.6)

#### Plugin Architecture

**Priority:** Low
**Complexity:** High

- Custom commands via plugins
- Plugin discovery and loading
- Community plugins

```bash
sok plugin install sok-docker
sok plugin list
sok docker generate-config  # Custom plugin command
```

#### Plugin API

- Well-defined plugin interface
- Plugin hooks into Sokru lifecycle
- Plugin configuration

### Phase 7: UI/UX Improvements (v1.7)

#### Interactive Mode

**Priority:** Medium
**Complexity:** Medium

- TUI (Terminal UI) with keyboard navigation
- Visual diff viewer
- Interactive symlink management

```bash
sok interactive  # Launch TUI
```

#### Better Diff Visualization

**Priority:** Low
**Complexity:** Low

- Side-by-side diff view
- Syntax highlighting in diffs
- Integration with diff tools

#### Progress Indicators

**Priority:** Low
**Complexity:** Low

- Progress bars for long operations
- Real-time status updates
- Estimated time remaining

### Phase 8: Additional Languages (v1.8)

#### More i18n Languages

**Priority:** Low
**Complexity:** Low

- French (fr)
- German (de)
- Portuguese (pt)
- Japanese (ja)
- Chinese (zh)

Community contributions welcome!

### Phase 9: Cloud & Sync (v2.0)

#### Cloud Backup

**Priority:** Low
**Complexity:** High

- Optional cloud backup storage
- Support for S3, Dropbox, Google Drive
- Encrypted cloud backups

```yaml
cloud:
  provider: s3
  bucket: my-dotfiles-backup
  encrypt: true
```

#### Multi-Device Sync

**Priority:** Low
**Complexity:** High

- Sync dotfiles across multiple devices
- Conflict resolution for concurrent changes
- Device-specific configurations

## Feature Requests

We track feature requests in GitHub Issues. Vote on features you'd like to see!

### Community Requested Features

1. **Shell completion** - Bash/Zsh/Fish completion scripts
2. **Watch mode** - Auto-apply on dotfiles changes
3. **Import from existing tools** - Import from GNU Stow, yadm, etc.
4. **Windows symlink alternatives** - Support for Windows junction points
5. **Dotfile discovery** - Auto-detect common dotfiles to manage

## Contributing

Want to work on a feature from this roadmap?

1. Check if there's an existing issue/PR
2. Comment on the issue to claim it
3. Discuss implementation approach
4. Submit a PR with tests and documentation

See [Contributing Guide](../CONTRIBUTING.md) for details.

## Version History

### v1.0 (Current)

- Multi-OS symlink management
- Automatic backup/restore
- Automatic rollback
- Internationalization (en/es)
- Comprehensive test suite

## Priorities Explained

- **High Priority**: Essential features for dotfiles management
- **Medium Priority**: Nice-to-have features that improve usability
- **Low Priority**: Advanced features for power users

## Complexity Estimates

- **Low**: 1-2 weeks of development
- **Medium**: 1-2 months of development
- **High**: 2-3+ months of development

## Deprecation Policy

We follow semantic versioning (SemVer):

- **Major version** (2.0): Breaking changes allowed
- **Minor version** (1.1): New features, backward compatible
- **Patch version** (1.0.1): Bug fixes only

We maintain backward compatibility within major versions.

## Related Projects

Similar tools to consider for inspiration:

- [GNU Stow](https://www.gnu.org/software/stow/) - Symlink farm manager
- [yadm](https://yadm.io/) - Yet Another Dotfiles Manager
- [chezmoi](https://www.chezmoi.io/) - Manage dotfiles across multiple machines
- [dotbot](https://github.com/anishathalye/dotbot) - Tool that bootstraps dotfiles
- [rcm](https://github.com/thoughtbot/rcm) - Management suite for dotfiles

## Feedback

Have ideas not listed here? Open an issue on GitHub!

## See Also

- [Main README](../README.md) - Project overview
- [Architecture](ARCHITECTURE.md) - Technical architecture
- [Contributing](../CONTRIBUTING.md) - How to contribute
- [Testing Guide](TESTING.md) - Test suite documentation
