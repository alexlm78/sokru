# Sokru Architecture

## Overview

Sokru is a command-line interface (CLI) tool designed for managing dotfiles across multiple operating systems. This document describes the architectural design, components, and patterns used in the project.

## High-Level Architecture

```diagram
┌─────────────────────────────────────────────────────────┐
│                     CLI Interface                        │
│                  (Cobra Framework)                       │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────┐
│                  Command Layer                           │
│  (init, apply, config, symlinks, restore, version)       │
└─────┬──────────┬──────────┬──────────┬─────────────────┘
      │          │          │          │
┌─────▼────┐ ┌──▼─────┐ ┌──▼────┐ ┌──▼─────────┐
│ Config   │ │  i18n  │ │Backup │ │  Rollback  │
│ Package  │ │Package │ │Package│ │  Package   │
└──────────┘ └────────┘ └───────┘ └────────────┘
      │          │          │          │
      └──────────┴──────────┴──────────┘
                  │
          ┌───────▼────────┐
          │  File System   │
          └────────────────┘
```

## Directory Structure

```tree
sokru/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and CLI entry point
│   ├── init.go            # Initialize configuration
│   ├── apply.go           # Apply configuration changes
│   ├── config.go          # Configuration management commands
│   ├── symlinks.go        # Symlink management commands
│   ├── restore.go         # Backup restore commands
│   ├── version.go         # Version information
│   ├── help.go            # Help text and utilities
│   └── utils.go           # Utility functions
│
├── internal/              # Internal packages (not exported)
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── i18n/             # Internationalization
│   │   ├── i18n.go
│   │   ├── messages_en.go
│   │   └── messages_es.go
│   ├── backup/           # Backup and restore system
│   │   └── backup.go
│   └── rollback/         # Rollback mechanism
│       └── rollback.go
│
├── test/                 # Test files (centralized)
│   ├── config_test.go
│   ├── i18n_test.go
│   ├── symlinks_test.go
│   ├── utils_test.go
│   ├── backup_test.go
│   └── rollback_test.go
│
├── docs/                 # Documentation
│   ├── ARCHITECTURE.md   # This file
│   ├── I18N.md
│   ├── MULTI_OS_SYMLINKS.md
│   ├── ROLLBACK.md
│   ├── BACKUP_RESTORE.md
│   └── TESTING.md
│
├── main.go              # Application entry point
├── Makefile            # Build automation
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── README.md           # Main documentation
└── CLAUDE.md           # AI assistant context
```

## Core Components

### 1. Command Layer (`cmd/`)

The command layer implements the CLI interface using the [Cobra](https://github.com/spf13/cobra) framework.

**Key files:**

- **`root.go`**: Entry point that registers all commands and global flags
- **`init.go`**: Initializes Sokru configuration directory and default config
- **`apply.go`**: Applies configuration changes with backup and rollback
- **`config.go`**: Manages configuration settings (get/set operations)
- **`symlinks.go`**: Manages symlink operations (install/uninstall/list)
- **`restore.go`**: Manages backup restore operations (list/apply/delete)
- **`version.go`**: Displays version information
- **`utils.go`**: Shared utility functions (path expansion, OS validation)

**Design Pattern**: Command Pattern (via Cobra framework)

Each command follows this structure:

```go
var commandCmd = &cobra.Command{
    Use:   "command",
    Short: "Brief description",
    Long:  "Detailed description",
    Run:   commandFunc,
}

func commandFunc(cmd *cobra.Command, args []string) {
    // Command implementation
}

func init() {
    rootCmd.AddCommand(commandCmd)
}
```

### 2. Configuration Package (`internal/config/`)

Manages application configuration with YAML persistence.

**Key responsibilities:**

- Load configuration from `~/.sokru/config.yaml`
- Save configuration changes
- Provide default configuration values
- Global configuration singleton

**Configuration structure:**

```go
type Config struct {
    DotfilesDir  string // Dotfiles directory location
    SymlinksFile string // Symlinks configuration file
    OS           string // Target OS (linux/darwin/windows)
    Language     string // UI language (en/es)
    Verbose      bool   // Verbose output flag
    DryRun       bool   // Dry-run mode flag
}
```

**Design Pattern**: Singleton Pattern

### 3. Internationalization Package (`internal/i18n/`)

Provides multi-language support for UI messages.

**Key features:**

- Type-safe message keys
- Language switching at runtime
- Message formatting with parameters
- Standard message prefixes (✓, ✗, →, ⚠)

**Supported languages:**

- English (en) - Default
- Spanish (es)

**Design Pattern**: Singleton Pattern with Factory Methods

**Usage example:**

```go
fmt.Println(i18n.Success(i18n.MsgSymlinkCreated, target, source))
fmt.Println(i18n.Error(i18n.MsgErrorLoadingConfig, err))
```

### 4. Backup Package (`internal/backup/`)

Implements automatic backup and restore functionality.

**Key features:**

- Timestamped backup sessions
- Metadata tracking (files, symlinks, permissions)
- Restore from any backup point
- Backup deletion

**Backup structure:**

```tree
~/.sokru/backups/
├── 20241101-143022.123/
│   ├── metadata.json
│   ├── file1
│   └── file2
└── 20241101-150315.456/
    ├── metadata.json
    └── file3
```

**Design Pattern**: Factory Pattern (BackupManager)

### 5. Rollback Package (`internal/rollback/`)

Implements automatic rollback on errors during symlink operations.

**Key features:**

- Tracks all operations (create/update/remove)
- Automatic rollback in reverse order (LIFO)
- Preserves previous state for updates

**Design Pattern**: Memento Pattern (state preservation)

**Action types:**

```go
type ActionType int
const (
    ActionCreated  ActionType = iota
    ActionUpdated
    ActionRemoved
)
```

## Data Flow

### Symlink Installation Flow

```flow
1. User runs: sok symlinks install
       │
       ▼
2. Load configuration from ~/.sokru/config.yaml
       │
       ▼
3. Read symlinks from configured file (e.g., ~/dotfiles/symlinks.yaml)
       │
       ▼
4. Filter symlinks by configured OS
       │
       ▼
5. Create backup session
       │
       ▼
6. Initialize rollback tracker
       │
       ▼
7. For each symlink:
   ├─ Check if target exists
   ├─ Backup existing file/symlink
   ├─ Track operation in rollback
   ├─ Create/update symlink
   └─ On error: trigger rollback
       │
       ▼
8. On success: complete backup
   On error: rollback all changes
```

### Configuration Change Flow

```flow
1. User runs: sok config <key> <value>
       │
       ▼
2. Load current configuration
       │
       ▼
3. Validate new value
       │
       ▼
4. Update configuration in memory
       │
       ▼
5. Save configuration to ~/.sokru/config.yaml
       │
       ▼
6. Update global configuration instance
```

### Restore Flow

```flow
1. User runs: sok restore apply <backup-id>
       │
       ▼
2. Load backup metadata
       │
       ▼
3. Display files to be restored
       │
       ▼
4. For each backup entry:
   ├─ If symlink: recreate symlink with saved target
   └─ If file: restore file content and permissions
       │
       ▼
5. Report results
```

## Design Patterns

### 1. Singleton Pattern

**Used in:**

- Configuration (`internal/config/`)
- Internationalization (`internal/i18n/`)

**Purpose:** Ensure single instance of configuration and i18n manager throughout application lifecycle.

### 2. Command Pattern

**Used in:**

- All CLI commands (`cmd/`)

**Purpose:** Encapsulate commands as objects, enabling parameterization and queuing.

### 3. Factory Pattern

**Used in:**

- Backup Manager (`internal/backup/`)

**Purpose:** Create backup instances without exposing creation logic.

### 4. Builder Pattern

**Used in:**

- Rollback Tracker (`internal/rollback/`)

**Purpose:** Construct complex rollback operations step by step.

### 5. Strategy Pattern

**Used in:**

- OS-specific symlink filtering (`cmd/symlinks.go`)

**Purpose:** Select behavior (symlink selection) based on configured OS.

### 6. Memento Pattern

**Used in:**

- Rollback system (`internal/rollback/`)

**Purpose:** Capture and restore object state for rollback.

## Error Handling Strategy

### 1. Early Validation

Validate inputs before making changes:

```go
if err := validateOS(os); err != nil {
    return err
}
```

### 2. Automatic Rollback

On any error during symlink operations:

```go
if err != nil {
    rollbackErrors := tracker.Rollback()
    // Handle rollback
}
```

### 3. User-Friendly Messages

Use i18n for consistent, translated error messages:

```go
fmt.Fprintf(os.Stderr, "%s", i18n.Error(i18n.MsgErrorLoadingConfig, err))
```

### 4. Dry-Run Mode

Allow users to preview changes without applying them:

```go
if cfg.DryRun {
    fmt.Printf("[DRY-RUN] Would create symlink: %s -> %s\n", target, source)
    return nil
}
```

## Security Considerations

### 1. Path Validation

- All paths are validated and expanded
- Symlinks are not followed during validation
- Absolute paths are preferred

### 2. Backup Permissions

- Backups inherit original file permissions
- Backup directory has restrictive permissions (0755)
- Sensitive files should be protected by user

### 3. No Privilege Escalation

- Sokru runs with user permissions
- No sudo or privilege escalation
- User must have write permissions to target paths

### 4. Configuration File Security

- Config file stored in user home directory
- Standard file permissions (0644)
- No passwords or secrets stored

## Performance Considerations

### 1. Lazy Loading

- Configuration loaded only when needed
- i18n messages loaded once at startup

### 2. Minimal File Operations

- Single read for symlinks configuration
- Batch operations where possible
- Efficient YAML parsing with gopkg.in/yaml.v3

### 3. Efficient Backups

- Only backup files that will be changed
- No compression (trade-off for simplicity)
- Metadata stored separately for fast listing

## Testing Strategy

### Test Organization

All tests are centralized in `test/` directory:

- **Unit tests**: Test individual functions and packages
- **Integration tests**: Test command workflows
- **Table-driven tests**: Test multiple scenarios efficiently

### Coverage Goals

- **Critical packages** (config, i18n): 90%+ coverage
- **Core logic** (symlinks, utils): 80%+ coverage
- **CLI commands**: 50%+ coverage

### Test Patterns

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected interface{}
    }{
        // Test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Build System

### Makefile Targets

```makefile
mac      # Build for macOS ARM64
macx86   # Build for macOS Intel
lin      # Build for Linux AMD64
win      # Build for Windows AMD64
release  # Build for all platforms
clean    # Remove build artifacts
```

### Cross-Compilation

Uses Go's built-in cross-compilation:

```bash
GOOS=linux GOARCH=amd64 go build -o build/sok-linux-amd64
```

## Future Architecture Improvements

### Potential Enhancements

1. **Plugin System**
   - Allow custom commands via plugins
   - Plugin discovery and loading

2. **Configuration Profiles**
   - Multiple named configurations (work, home, etc.)
   - Easy switching between profiles

3. **Remote Configuration**
   - Fetch symlinks config from git repository
   - Auto-update dotfiles

4. **Hooks System**
   - Pre/post install hooks
   - Custom validation hooks

5. **Template Engine**
   - Generate dotfiles from templates
   - Variable substitution

6. **Dependency Management**
   - Define dependencies between dotfiles
   - Ordered installation

7. **Conflict Resolution**
   - Interactive conflict resolution
   - Merge strategies

## References

- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [YAML v3 Library](https://github.com/go-yaml/yaml)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## Glossary

- **Dotfiles**: Configuration files typically stored in home directory, often starting with `.`
- **Symlink**: Symbolic link, a file that points to another file or directory
- **Rollback**: Reverting changes to previous state after an error
- **Backup**: Copy of files before modification for safety
- **Dry-Run**: Preview mode that shows what would happen without making actual changes
- **i18n**: Internationalization, supporting multiple languages

## See Also

- [Main README](../README.md) - Project overview and quick start
- [Testing Guide](TESTING.md) - Test suite documentation
- [Multi-OS Symlinks](MULTI_OS_SYMLINKS.md) - Multi-OS configuration format
- [Backup & Restore](BACKUP_RESTORE.md) - Backup system documentation
- [Rollback Mechanism](ROLLBACK.md) - Rollback system documentation
- [Internationalization](I18N.md) - i18n system documentation
