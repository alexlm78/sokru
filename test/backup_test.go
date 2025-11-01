// Package test
// Description: Unit tests for backup functionality
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexlm78/sokru/internal/backup"
)

func TestNewManager(t *testing.T) {
	backupDir := "/test/backups"
	manager := backup.NewManager(backupDir)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.GetBackupDir() != backupDir {
		t.Errorf("Expected backup dir '%s', got '%s'", backupDir, manager.GetBackupDir())
	}
}

func TestEnsureBackupDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)

	// Ensure directory is created
	if err := manager.EnsureBackupDir(); err != nil {
		t.Fatalf("EnsureBackupDir failed: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		t.Error("Backup directory was not created")
	}

	// Calling again should not error
	if err := manager.EnsureBackupDir(); err != nil {
		t.Errorf("EnsureBackupDir should not error when directory exists: %v", err)
	}
}

func TestCreateBackupRegularFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create backup
	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-backup-001"

	entry, err := manager.CreateBackup(testFile, backupID)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup entry
	if entry.OriginalPath != testFile {
		t.Errorf("Expected original path '%s', got '%s'", testFile, entry.OriginalPath)
	}

	if entry.IsSymlink {
		t.Error("Regular file should not be marked as symlink")
	}

	// Verify backup file exists and has same content
	backupContent, err := os.ReadFile(entry.BackupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}

	if string(backupContent) != string(testContent) {
		t.Errorf("Backup content mismatch: expected '%s', got '%s'", testContent, backupContent)
	}
}

func TestCreateBackupSymlink(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file and symlink
	sourceFile := filepath.Join(tempDir, "source.txt")
	symlinkFile := filepath.Join(tempDir, "link.txt")

	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	if err := os.Symlink(sourceFile, symlinkFile); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Create backup
	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-backup-002"

	entry, err := manager.CreateBackup(symlinkFile, backupID)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup entry
	if !entry.IsSymlink {
		t.Error("Symlink should be marked as symlink")
	}

	if entry.SymlinkTarget != sourceFile {
		t.Errorf("Expected symlink target '%s', got '%s'", sourceFile, entry.SymlinkTarget)
	}
}

func TestSaveAndLoadMetadata(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-backup-003"

	// Create metadata
	metadata := &backup.BackupMetadata{
		ID:        backupID,
		Timestamp: time.Now(),
		Command:   "test command",
		Entries: []backup.BackupEntry{
			{
				OriginalPath: "/test/file1",
				BackupPath:   "/backup/file1",
				IsSymlink:    false,
				Timestamp:    time.Now(),
			},
			{
				OriginalPath:  "/test/link1",
				BackupPath:    "/backup/link1",
				IsSymlink:     true,
				SymlinkTarget: "/test/target",
				Timestamp:     time.Now(),
			},
		},
	}

	// Save metadata
	if err := manager.SaveMetadata(metadata); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}

	// Load metadata
	loadedMetadata, err := manager.LoadMetadata(backupID)
	if err != nil {
		t.Fatalf("LoadMetadata failed: %v", err)
	}

	// Verify loaded metadata
	if loadedMetadata.ID != metadata.ID {
		t.Errorf("ID mismatch: expected '%s', got '%s'", metadata.ID, loadedMetadata.ID)
	}

	if loadedMetadata.Command != metadata.Command {
		t.Errorf("Command mismatch: expected '%s', got '%s'", metadata.Command, loadedMetadata.Command)
	}

	if len(loadedMetadata.Entries) != len(metadata.Entries) {
		t.Errorf("Entries count mismatch: expected %d, got %d",
			len(metadata.Entries), len(loadedMetadata.Entries))
	}
}

func TestListBackups(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)

	// Create multiple backups
	for i := 1; i <= 3; i++ {
		backupID := fmt.Sprintf("backup-%03d", i)
		metadata := &backup.BackupMetadata{
			ID:        backupID,
			Timestamp: time.Now(),
			Command:   "test",
			Entries:   []backup.BackupEntry{},
		}

		if err := manager.SaveMetadata(metadata); err != nil {
			t.Fatalf("Failed to save metadata %d: %v", i, err)
		}
	}

	// List backups
	backups, err := manager.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 3 {
		t.Errorf("Expected 3 backups, got %d", len(backups))
	}
}

func TestListBackupsEmpty(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)

	// List backups when directory doesn't exist
	backups, err := manager.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups should not error when directory doesn't exist: %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("Expected 0 backups, got %d", len(backups))
	}
}

func TestRestoreBackupRegularFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create original file
	originalFile := filepath.Join(tempDir, "original.txt")
	originalContent := []byte("original content")
	if err := os.WriteFile(originalFile, originalContent, 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	// Create backup
	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-restore-001"

	entry, err := manager.CreateBackup(originalFile, backupID)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Save metadata
	metadata := &backup.BackupMetadata{
		ID:        backupID,
		Timestamp: time.Now(),
		Command:   "test",
		Entries:   []backup.BackupEntry{*entry},
	}
	if err := manager.SaveMetadata(metadata); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}

	// Modify original file
	modifiedContent := []byte("modified content")
	if err := os.WriteFile(originalFile, modifiedContent, 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Restore backup
	if err := manager.RestoreBackup(backupID); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	// Verify file was restored
	restoredContent, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}

	if string(restoredContent) != string(originalContent) {
		t.Errorf("Content mismatch: expected '%s', got '%s'", originalContent, restoredContent)
	}
}

func TestRestoreBackupSymlink(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source and symlink
	sourceFile := filepath.Join(tempDir, "source.txt")
	symlinkFile := filepath.Join(tempDir, "link.txt")

	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source: %v", err)
	}

	if err := os.Symlink(sourceFile, symlinkFile); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Create backup
	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-restore-002"

	entry, err := manager.CreateBackup(symlinkFile, backupID)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Save metadata
	metadata := &backup.BackupMetadata{
		ID:        backupID,
		Timestamp: time.Now(),
		Command:   "test",
		Entries:   []backup.BackupEntry{*entry},
	}
	if err := manager.SaveMetadata(metadata); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}

	// Remove symlink
	if err := os.Remove(symlinkFile); err != nil {
		t.Fatalf("Failed to remove symlink: %v", err)
	}

	// Restore backup
	if err := manager.RestoreBackup(backupID); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	// Verify symlink was restored
	restoredTarget, err := os.Readlink(symlinkFile)
	if err != nil {
		t.Fatalf("Failed to read restored symlink: %v", err)
	}

	if restoredTarget != sourceFile {
		t.Errorf("Symlink target mismatch: expected '%s', got '%s'", sourceFile, restoredTarget)
	}
}

func TestDeleteBackup(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-delete-001"

	// Create backup metadata
	metadata := &backup.BackupMetadata{
		ID:        backupID,
		Timestamp: time.Now(),
		Command:   "test",
		Entries:   []backup.BackupEntry{},
	}
	if err := manager.SaveMetadata(metadata); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}

	// Verify backup exists
	backupPath := filepath.Join(backupDir, backupID)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatal("Backup directory should exist")
	}

	// Delete backup
	if err := manager.DeleteBackup(backupID); err != nil {
		t.Fatalf("DeleteBackup failed: %v", err)
	}

	// Verify backup was deleted
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Error("Backup directory should be deleted")
	}
}

func TestGenerateBackupID(t *testing.T) {
	id1 := backup.GenerateBackupID()

	if id1 == "" {
		t.Error("GenerateBackupID should not return empty string")
	}

	// Wait a bit and generate another ID
	time.Sleep(10 * time.Millisecond)
	id2 := backup.GenerateBackupID()

	if id1 == id2 {
		t.Error("GenerateBackupID should generate unique IDs")
	}

	// Verify format (should be timestamp-based with milliseconds)
	if len(id1) != 19 { // Format: 20060102-150405.000
		t.Errorf("Expected ID length 19, got %d", len(id1))
	}
}

func TestGetDefaultBackupDir(t *testing.T) {
	backupDir, err := backup.GetDefaultBackupDir()
	if err != nil {
		t.Fatalf("GetDefaultBackupDir failed: %v", err)
	}

	if backupDir == "" {
		t.Error("GetDefaultBackupDir should not return empty string")
	}

	// Should contain .sokru/backups
	if !filepath.IsAbs(backupDir) {
		t.Error("GetDefaultBackupDir should return absolute path")
	}
}

func TestRestoreMultipleFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sokru-backup-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	manager := backup.NewManager(backupDir)
	backupID := "test-multi-001"

	// Create multiple files
	var entries []backup.BackupEntry
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
		content := []byte(fmt.Sprintf("content %d", i))

		if err := os.WriteFile(filename, content, 0644); err != nil {
			t.Fatalf("Failed to create file %d: %v", i, err)
		}

		entry, err := manager.CreateBackup(filename, backupID)
		if err != nil {
			t.Fatalf("Failed to backup file %d: %v", i, err)
		}
		entries = append(entries, *entry)
	}

	// Save metadata
	metadata := &backup.BackupMetadata{
		ID:        backupID,
		Timestamp: time.Now(),
		Command:   "test",
		Entries:   entries,
	}
	if err := manager.SaveMetadata(metadata); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}

	// Modify all files
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(filename, []byte("modified"), 0644); err != nil {
			t.Fatalf("Failed to modify file %d: %v", i, err)
		}
	}

	// Restore backup
	if err := manager.RestoreBackup(backupID); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	// Verify all files were restored
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("Failed to read restored file %d: %v", i, err)
		}

		expected := fmt.Sprintf("content %d", i)
		if string(content) != expected {
			t.Errorf("File %d content mismatch: expected '%s', got '%s'", i, expected, content)
		}
	}
}
