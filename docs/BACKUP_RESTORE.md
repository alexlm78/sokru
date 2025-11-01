# Backup and Restore Guide

## Overview

Sokru automatically creates backups of existing files and symlinks before making changes. This ensures you can always restore your previous configuration if needed.

## Automatic Backups

### When Backups Are Created

Backups are automatically created when:

- Installing symlinks with `sok symlinks install`
- Applying configuration changes with `sok apply`
- Any existing file or symlink will be replaced

### What Gets Backed Up

- **Regular files** - Complete file content and permissions
- **Symlinks** - Symlink target path and metadata
- **File metadata** - Permissions, timestamps, and file type

### Backup Location

All backups are stored in:

```data
~/.config/sokru/backups/
```

Each backup session creates a subdirectory with a unique ID:

```data
~/.config/sokru/backups/
├── 20241101-143022.123/
│   ├── metadata.json
│   ├── bashrc
│   └── vimrc
├── 20241101-150315.456/
│   ├── metadata.json
│   └── zshrc
```

## Restore Command

### List Available Backups

View all available backup sessions:

```bash
sok restore list
```

**Output:**

```data
Available Backups
================

ID: 20241101-143022.123
  Timestamp: 2024-11-01T14:30:22-06:00
  Command: symlinks install
  Files: 3

ID: 20241101-150315.456
  Timestamp: 2024-11-01T15:03:15-06:00
  Command: apply
  Files: 1

Total backups: 2
```

### Restore a Backup

Restore files from a specific backup:

```bash
sok restore apply <backup-id>
```

**Example:**

```bash
$ sok restore apply 20241101-143022.123
Restoring backup: 20241101-143022.123
Backup created: 2024-11-01T14:30:22-06:00
Files to restore: 3

Files in backup:
  [file]    ~/.bashrc
  [symlink] ~/.vimrc -> ~/.dotfiles/vim/vimrc
  [file]    ~/.zshrc

→ Restoring files...
✓ Restore completed successfully
```

### Delete a Backup

Remove a backup when no longer needed:

```bash
sok restore delete <backup-id>
```

**Example:**

```bash
$ sok restore delete 20241101-143022.123
Deleting backup: 20241101-143022.123
Backup created: 2024-11-01T14:30:22-06:00
Files: 3

✓ Backup deleted: 20241101-143022.123
```

## Backup Metadata

Each backup includes a `metadata.json` file with:

```json
{
  "id": "20241101-143022.123",
  "timestamp": "2024-11-01T14:30:22.123456-06:00",
  "command": "symlinks install",
  "entries": [
    {
      "original_path": "/home/user/.bashrc",
      "backup_path": "/home/user/.config/sokru/backups/20241101-143022.123/bashrc",
      "is_symlink": false,
      "timestamp": "2024-11-01T14:30:22.123456-06:00",
      "file_mode": 420
    },
    {
      "original_path": "/home/user/.vimrc",
      "backup_path": "/home/user/.config/sokru/backups/20241101-143022.123/vimrc",
      "is_symlink": true,
      "symlink_target": "/home/user/.dotfiles/vim/vimrc",
      "timestamp": "2024-11-01T14:30:22.123456-06:00",
      "file_mode": 511
    }
  ]
}
```

## Integration with Rollback

Backups work together with the rollback mechanism:

1. **Before changes**: Backups are created
2. **During changes**: Rollback tracker monitors operations
3. **On error**: Rollback reverts changes immediately
4. **After success**: Backup is saved for future restore

### Example Flow

```bash
$ sok symlinks install --verbose
→ Reading symlinks configuration from: ~/.dotfiles/symlinks.yml
→ Found 3 symlink configuration(s)
→ Backing up: ~/.bashrc
→ Backing up: ~/.vimrc
✓ Symlink created: ~/.bashrc -> ~/.dotfiles/bash/bashrc
✓ Symlink created: ~/.vimrc -> ~/.dotfiles/vim/vimrc
✗ Error creating symlink from ~/.zshrc to ~/.dotfiles/zsh/zshrc: permission denied

⚠ Error occurred, starting rollback of 2 action(s)...
✓ Rollback completed successfully
✓ Backup completed: 20241101-143022.123
```

**Result:**

- Changes were rolled back immediately
- Backup was saved for manual restore if needed
- System is in consistent state

## Use Cases

### 1. Testing New Configurations

```bash
# Install new configuration
sok symlinks install

# If something goes wrong, restore previous state
sok restore list
sok restore apply <backup-id>
```

### 2. Switching Between Configurations

```bash
# Try new dotfiles
sok symlinks install

# Not happy? Restore old ones
sok restore apply <previous-backup-id>
```

### 3. Disaster Recovery

```bash
# List all backups
sok restore list

# Restore from specific point in time
sok restore apply 20241101-120000.000
```

### 4. Cleanup Old Backups

```bash
# List backups
sok restore list

# Delete old backups
sok restore delete 20241001-100000.000
sok restore delete 20241015-140000.000
```

## Best Practices

### 1. Regular Backup Cleanup

Backups accumulate over time. Periodically review and delete old backups:

```bash
# List all backups
sok restore list

# Delete backups older than 30 days
# (manual process - review before deleting)
```

### 2. Keep Recent Backups

Keep at least the last 3-5 backups for safety:

- Latest working configuration
- Previous stable configuration
- Emergency fallback

### 3. Test Restores

Periodically test that restores work:

```bash
# In a test environment
sok restore apply <backup-id>
# Verify files are correct
# Re-apply current configuration
sok symlinks install
```

### 4. Document Important Backups

Keep notes about significant backup IDs:

```data
20241101-143022.123 - Before major refactor
20241015-120000.456 - Stable production config
20241001-090000.789 - Initial setup
```

## Backup Storage

### Disk Space

Backups consume disk space. Monitor usage:

```bash
# Check backup directory size
du -sh ~/.config/sokru/backups/

# List backups with details
sok restore list
```

### Automatic Cleanup

Currently, Sokru does not automatically delete old backups. You must manually manage them using `sok restore delete`.

**Future enhancement**: Automatic cleanup of backups older than N days.

## Troubleshooting

### Backup Failed Warning

If you see:

```data
⚠ Backup failed: permission denied
```

**Actions:**

1. Check file permissions
2. Ensure write access to `~/.config/sokru/backups/`
3. Installation continues (backup is optional)

### Restore Failed

If restore fails:

```data
✗ Restore failed: permission denied
```

**Actions:**

1. Check permissions on target files
2. Ensure backup files exist
3. Try with sudo if needed (not recommended)
4. Manually restore from backup directory

### Backup Directory Full

If disk is full:

```bash
# Check disk space
df -h ~/.config/sokru/

# Delete old backups
sok restore list
sok restore delete <old-backup-id>
```

### Corrupted Backup

If metadata is corrupted:

```bash
# Backup won't appear in list
# Manually inspect backup directory
ls -la ~/.config/sokru/backups/

# Delete corrupted backup directory
rm -rf ~/.config/sokru/backups/<corrupted-id>
```

## Advanced Usage

### Manual Backup Inspection

Backups are stored as regular files and can be inspected manually:

```bash
# View backup metadata
cat ~/.config/sokru/backups/20241101-143022.123/metadata.json | jq

# View backed up file
cat ~/.config/sokru/backups/20241101-143022.123/bashrc

# Check symlink info in metadata
jq '.entries[] | select(.is_symlink == true)' \
  ~/.config/sokru/backups/20241101-143022.123/metadata.json
```

### Selective Restore

Currently, restore is all-or-nothing. For selective restore:

```bash
# Manually copy specific file from backup
cp ~/.config/sokru/backups/20241101-143022.123/bashrc ~/.bashrc

# Or recreate specific symlink from metadata
# (check metadata.json for symlink_target)
ln -s <target-from-metadata> ~/.vimrc
```

## Security Considerations

### Backup Permissions

Backups inherit permissions from original files. Ensure:

- Backup directory has appropriate permissions (0755)
- Sensitive files in backups are protected
- Backup directory is not world-readable if it contains secrets

### Sensitive Data

If backing up files with sensitive data:

```bash
# Check backup directory permissions
ls -la ~/.config/sokru/backups/

# Restrict if needed
chmod 700 ~/.config/sokru/backups/
```

## Limitations

### What Is NOT Backed Up

- Files outside the symlink configuration
- Directory structures (only individual files)
- File ownership (UID/GID)
- Extended attributes
- ACLs (Access Control Lists)

### Backup Size

- Each backup is a full copy (not incremental)
- Large files consume significant space
- No compression (files stored as-is)

## Future Enhancements

Potential improvements:

- [ ] Automatic backup rotation (keep last N backups)
- [ ] Backup compression
- [ ] Incremental backups
- [ ] Selective restore (restore specific files)
- [ ] Backup verification
- [ ] Backup export/import
- [ ] Remote backup storage

## See Also

- [Main README](../README.md)
- [Rollback Guide](ROLLBACK.md)
- [Multi-OS Symlinks Guide](MULTI_OS_SYMLINKS.md)
- [Testing Guide](TESTING.md)
