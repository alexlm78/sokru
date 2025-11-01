# Rollback Mechanism

## Overview

Sokru includes an automatic rollback mechanism that protects against partial failures during symlink installation and updates. If an error occurs during symlink operations, all changes made in that session are automatically reverted to maintain system consistency.

## How It Works

### Automatic Rollback

When you run `sok symlinks install` or `sok apply`, Sokru:

1. **Tracks all changes** made during the operation
2. **Detects errors** during symlink creation/update
3. **Automatically rolls back** all changes if an error occurs
4. **Restores previous state** to ensure consistency

### What Gets Tracked

The rollback system tracks three types of operations:

1. **Created Symlinks** - New symlinks that were created
2. **Updated Symlinks** - Existing symlinks that were modified
3. **Removed Symlinks** - Symlinks that were deleted

### Rollback Actions

When a rollback is triggered:

- **Created symlinks** → Removed
- **Updated symlinks** → Restored to previous target
- **Removed symlinks** → Recreated with original target

## Example Scenarios

### Scenario 1: Partial Installation Failure

```bash
$ sok symlinks install
→ Reading symlinks configuration from: ~/.dotfiles/symlinks.yml
✓ Symlink created: ~/.bashrc -> ~/.dotfiles/bash/bashrc
✓ Symlink created: ~/.vimrc -> ~/.dotfiles/vim/vimrc
✗ Error creating symlink from ~/.zshrc to ~/.dotfiles/zsh/zshrc: permission denied

⚠ Error occurred, starting rollback of 2 action(s)...
✓ Rollback completed successfully
```

**Result**: The `.bashrc` and `.vimrc` symlinks are removed, system returns to original state.

### Scenario 2: Update Failure with Rollback

```bash
$ sok apply
Applying configuration changes...

=== Changes to Apply ===

📝 To Create (1):
  + ~/.tmux.conf -> ~/.dotfiles/tmux/tmux.conf

🔄 To Update (1):
  ~ ~/.bashrc: ~/.dotfiles/bash/old.bashrc -> ~/.dotfiles/bash/new.bashrc

=== Applying Changes ===
✓ Created: ~/.tmux.conf -> ~/.dotfiles/tmux/tmux.conf
✗ Error creating symlink from ~/.bashrc to ~/.dotfiles/bash/new.bashrc: file exists

⚠ Error occurred, starting rollback of 2 action(s)...
✓ Rollback completed successfully
```

**Result**:

- `~/.tmux.conf` is removed
- `~/.bashrc` is restored to point to `~/.dotfiles/bash/old.bashrc`

### Scenario 3: Successful Operation (No Rollback)

```bash
$ sok symlinks install
→ Reading symlinks configuration from: ~/.dotfiles/symlinks.yml
✓ Symlink created: ~/.bashrc -> ~/.dotfiles/bash/bashrc
✓ Symlink created: ~/.vimrc -> ~/.dotfiles/vim/vimrc
✓ Symlink created: ~/.zshrc -> ~/.dotfiles/zsh/zshrc
```

**Result**: All symlinks created successfully, no rollback needed.

## Technical Details

### Rollback Order

Actions are rolled back in **reverse order** (LIFO - Last In, First Out):

```data
Actions performed:
1. Create ~/.bashrc
2. Create ~/.vimrc
3. Create ~/.zshrc (fails)

Rollback order:
1. Remove ~/.vimrc
2. Remove ~/.bashrc
```

This ensures dependencies are handled correctly.

### State Preservation

For updated symlinks, the system preserves:

- Previous symlink target
- Whether the target was a symlink
- Original file permissions

### Error Handling

If rollback itself encounters errors:

- Rollback continues for remaining actions
- All errors are collected and reported
- Exit code indicates failure
- Partial rollback is better than no rollback

## Dry-Run Mode

Rollback is **not triggered** in dry-run mode:

```bash
$ sok symlinks install --dry-run
[DRY-RUN] Would create symlink: ~/.bashrc -> ~/.dotfiles/bash/bashrc
[DRY-RUN] Would create symlink: ~/.vimrc -> ~/.dotfiles/vim/vimrc
```

No actual changes are made, so no rollback is needed.

## Limitations

### What Rollback Cannot Do

1. **Restore deleted regular files** - Only symlinks are tracked
2. **Restore file permissions** - Original permissions may not be preserved
3. **Handle external changes** - Changes made outside Sokru during operation
4. **Recover from disk failures** - Hardware/filesystem errors may prevent rollback

### When Rollback May Fail

Rollback can fail if:

- Insufficient permissions to remove/create symlinks
- Filesystem is read-only
- Disk is full
- Target paths have been modified externally

In these cases, Sokru reports the rollback errors but continues attempting to rollback remaining actions.

## Best Practices

### 1. Use Dry-Run First

Always test with `--dry-run` before making actual changes:

```bash
sok symlinks install --dry-run
sok apply --dry-run
```

### 2. Check Permissions

Ensure you have write permissions for all target paths:

```bash
# Check permissions before installing
ls -la ~/.config/
```

### 3. Backup Important Files

For critical configurations, maintain backups:

```bash
# Backup before major changes
cp ~/.bashrc ~/.bashrc.backup
```

### 4. Use Verbose Mode

Enable verbose output to see detailed operation logs:

```bash
sok symlinks install --verbose
```

### 5. Test in Stages

For large configurations, test in smaller batches:

```yaml
# Test with a few symlinks first
- common:
    ~/.vimrc: ~/.dotfiles/vim/vimrc
    ~/.gitconfig: ~/.dotfiles/git/gitconfig
```

## Troubleshooting

### Rollback Completed with Errors

If you see:

```output
⚠ Error occurred, starting rollback of 5 action(s)...
✗ Rollback completed with errors: [...]
```

**Actions:**

1. Check the error messages for specific failures
2. Verify file permissions
3. Check disk space
4. Manually inspect affected symlinks
5. Re-run with `--verbose` for more details

### Partial Rollback

If rollback partially fails:

```bash
# List current symlink status
sok symlinks list

# Manually remove problematic symlinks
rm ~/.problematic-link

# Try again
sok symlinks install
```

### System State After Failed Rollback

If rollback fails, your system may be in an inconsistent state:

1. Use `sok symlinks list` to see current status
2. Manually fix any incorrect symlinks
3. Re-run the installation

## Implementation Details

### Rollback Package

Location: `internal/rollback/rollback.go`

Key components:

- `Tracker` - Tracks all symlink operations
- `SymlinkAction` - Represents a single operation
- `ActionType` - Type of operation (Created/Updated/Removed)

### Integration Points

Rollback is integrated into:

- `cmd/symlinks.go` - InstallSymlinksFunc
- `cmd/apply.go` - ApplyFunc

### Testing

Comprehensive tests in `test/rollback_test.go`:

- Tracker functionality
- Created symlink rollback
- Updated symlink rollback
- Multiple action rollback
- Edge cases and error handling

## See Also

- [Main README](../README.md)
- [Multi-OS Symlinks Guide](MULTI_OS_SYMLINKS.md)
- [Testing Guide](TESTING.md)
- [Internationalization Guide](I18N.md)
