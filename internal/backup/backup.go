// Package backup
// Description: Backup and restore functionality for files and symlinks
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupEntry represents a single backed up file or symlink
type BackupEntry struct {
	OriginalPath  string      `json:"original_path"`
	BackupPath    string      `json:"backup_path"`
	IsSymlink     bool        `json:"is_symlink"`
	SymlinkTarget string      `json:"symlink_target,omitempty"`
	Timestamp     time.Time   `json:"timestamp"`
	FileMode      os.FileMode `json:"file_mode"`
}

// BackupMetadata contains information about a backup session
type BackupMetadata struct {
	ID        string        `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Entries   []BackupEntry `json:"entries"`
	Command   string        `json:"command"`
}

// Manager handles backup operations
type Manager struct {
	backupDir string
}

// NewManager creates a new backup manager
func NewManager(backupDir string) *Manager {
	return &Manager{
		backupDir: backupDir,
	}
}

// GetBackupDir returns the backup directory path
func (m *Manager) GetBackupDir() string {
	return m.backupDir
}

// EnsureBackupDir creates the backup directory if it doesn't exist
func (m *Manager) EnsureBackupDir() error {
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	return nil
}

// CreateBackup creates a backup of a file or symlink
func (m *Manager) CreateBackup(originalPath string, backupID string) (*BackupEntry, error) {
	// Ensure backup directory exists
	if err := m.EnsureBackupDir(); err != nil {
		return nil, err
	}

	// Get file info
	fileInfo, err := os.Lstat(originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Create backup subdirectory for this backup session
	sessionDir := filepath.Join(m.backupDir, backupID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	// Generate backup filename
	backupFilename := filepath.Base(originalPath)
	backupPath := filepath.Join(sessionDir, backupFilename)

	// Handle symlinks
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(originalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read symlink: %w", err)
		}

		entry := &BackupEntry{
			OriginalPath:  originalPath,
			BackupPath:    backupPath,
			IsSymlink:     true,
			SymlinkTarget: target,
			Timestamp:     time.Now(),
			FileMode:      fileInfo.Mode(),
		}

		return entry, nil
	}

	// Handle regular files - copy the file
	if err := m.copyFile(originalPath, backupPath); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	entry := &BackupEntry{
		OriginalPath: originalPath,
		BackupPath:   backupPath,
		IsSymlink:    false,
		Timestamp:    time.Now(),
		FileMode:     fileInfo.Mode(),
	}

	return entry, nil
}

// copyFile copies a file from src to dst
func (m *Manager) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// SaveMetadata saves backup metadata to a JSON file
func (m *Manager) SaveMetadata(metadata *BackupMetadata) error {
	sessionDir := filepath.Join(m.backupDir, metadata.ID)

	// Ensure session directory exists
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	metadataPath := filepath.Join(sessionDir, "metadata.json")

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// LoadMetadata loads backup metadata from a JSON file
func (m *Manager) LoadMetadata(backupID string) (*BackupMetadata, error) {
	metadataPath := filepath.Join(m.backupDir, backupID, "metadata.json")

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata BackupMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

// ListBackups returns a list of all available backups
func (m *Manager) ListBackups() ([]BackupMetadata, error) {
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupMetadata{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupMetadata
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadata, err := m.LoadMetadata(entry.Name())
		if err != nil {
			// Skip invalid backups
			continue
		}

		backups = append(backups, *metadata)
	}

	return backups, nil
}

// RestoreBackup restores files from a backup
func (m *Manager) RestoreBackup(backupID string) error {
	metadata, err := m.LoadMetadata(backupID)
	if err != nil {
		return err
	}

	var errors []error

	for _, entry := range metadata.Entries {
		if entry.IsSymlink {
			// Restore symlink
			// Remove current file/symlink if exists
			if err := os.Remove(entry.OriginalPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("failed to remove %s: %w", entry.OriginalPath, err))
				continue
			}

			// Recreate symlink
			if err := os.Symlink(entry.SymlinkTarget, entry.OriginalPath); err != nil {
				errors = append(errors, fmt.Errorf("failed to restore symlink %s: %w", entry.OriginalPath, err))
			}
		} else {
			// Restore regular file
			// Remove current file if exists
			if err := os.Remove(entry.OriginalPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("failed to remove %s: %w", entry.OriginalPath, err))
				continue
			}

			// Copy backup file back
			if err := m.copyFile(entry.BackupPath, entry.OriginalPath); err != nil {
				errors = append(errors, fmt.Errorf("failed to restore file %s: %w", entry.OriginalPath, err))
				continue
			}

			// Restore permissions
			if err := os.Chmod(entry.OriginalPath, entry.FileMode); err != nil {
				errors = append(errors, fmt.Errorf("failed to restore permissions for %s: %w", entry.OriginalPath, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("restore completed with %d error(s): %v", len(errors), errors)
	}

	return nil
}

// DeleteBackup removes a backup directory
func (m *Manager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(m.backupDir, backupID)
	if err := os.RemoveAll(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}
	return nil
}

// GenerateBackupID generates a unique backup ID based on timestamp
func GenerateBackupID() string {
	return time.Now().Format("20060102-150405.000")
}

// GetDefaultBackupDir returns the default backup directory path
func GetDefaultBackupDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "sokru", "backups"), nil
}
