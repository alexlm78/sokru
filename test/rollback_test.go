// Package test
// Description: Unit tests for rollback functionality
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexlm78/sokru/internal/rollback"
)

func TestNewTracker(t *testing.T) {
	tracker := rollback.NewTracker()

	if tracker == nil {
		t.Fatal("NewTracker returned nil")
	}

	if !tracker.IsEnabled() {
		t.Error("Tracker should be enabled by default")
	}

	if tracker.Count() != 0 {
		t.Error("New tracker should have 0 actions")
	}

	if tracker.HasActions() {
		t.Error("New tracker should not have actions")
	}
}

func TestTrackerEnableDisable(t *testing.T) {
	tracker := rollback.NewTracker()

	// Test disable
	tracker.Disable()
	if tracker.IsEnabled() {
		t.Error("Tracker should be disabled")
	}

	// Test enable
	tracker.Enable()
	if !tracker.IsEnabled() {
		t.Error("Tracker should be enabled")
	}
}

func TestTrackCreated(t *testing.T) {
	tracker := rollback.NewTracker()

	tracker.TrackCreated("/test/target", "/test/source")

	if tracker.Count() != 1 {
		t.Errorf("Expected 1 action, got %d", tracker.Count())
	}

	if !tracker.HasActions() {
		t.Error("Tracker should have actions")
	}

	actions := tracker.GetActions()
	if len(actions) != 1 {
		t.Fatalf("Expected 1 action, got %d", len(actions))
	}

	action := actions[0]
	if action.TargetPath != "/test/target" {
		t.Errorf("Expected target '/test/target', got '%s'", action.TargetPath)
	}

	if action.SourcePath != "/test/source" {
		t.Errorf("Expected source '/test/source', got '%s'", action.SourcePath)
	}
}

func TestTrackUpdated(t *testing.T) {
	tracker := rollback.NewTracker()

	tracker.TrackUpdated("/test/target", "/test/new-source", "/test/old-source")

	if tracker.Count() != 1 {
		t.Errorf("Expected 1 action, got %d", tracker.Count())
	}

	actions := tracker.GetActions()
	action := actions[0]

	if action.TargetPath != "/test/target" {
		t.Errorf("Expected target '/test/target', got '%s'", action.TargetPath)
	}

	if action.SourcePath != "/test/new-source" {
		t.Errorf("Expected source '/test/new-source', got '%s'", action.SourcePath)
	}

	if action.PreviousLink != "/test/old-source" {
		t.Errorf("Expected previous link '/test/old-source', got '%s'", action.PreviousLink)
	}

	if !action.WasSymlink {
		t.Error("WasSymlink should be true for updated symlinks")
	}
}

func TestTrackRemoved(t *testing.T) {
	tracker := rollback.NewTracker()

	tracker.TrackRemoved("/test/target", "/test/source")

	if tracker.Count() != 1 {
		t.Errorf("Expected 1 action, got %d", tracker.Count())
	}

	actions := tracker.GetActions()
	action := actions[0]

	if !action.WasSymlink {
		t.Error("WasSymlink should be true for removed symlinks")
	}
}

func TestTrackerClear(t *testing.T) {
	tracker := rollback.NewTracker()

	tracker.TrackCreated("/test/1", "/source/1")
	tracker.TrackCreated("/test/2", "/source/2")
	tracker.TrackCreated("/test/3", "/source/3")

	if tracker.Count() != 3 {
		t.Errorf("Expected 3 actions, got %d", tracker.Count())
	}

	tracker.Clear()

	if tracker.Count() != 0 {
		t.Errorf("Expected 0 actions after clear, got %d", tracker.Count())
	}

	if tracker.HasActions() {
		t.Error("Tracker should not have actions after clear")
	}
}

func TestTrackerDisabledNoTracking(t *testing.T) {
	tracker := rollback.NewTracker()
	tracker.Disable()

	tracker.TrackCreated("/test/target", "/test/source")

	if tracker.Count() != 0 {
		t.Error("Disabled tracker should not track actions")
	}
}

func TestRollbackCreatedSymlink(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sokru-rollback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetPath := filepath.Join(tempDir, "test-link")
	sourcePath := filepath.Join(tempDir, "test-source")

	// Create source file
	if err := os.WriteFile(sourcePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create symlink
	if err := os.Symlink(sourcePath, targetPath); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Track the creation
	tracker := rollback.NewTracker()
	tracker.TrackCreated(targetPath, sourcePath)

	// Verify symlink exists
	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		t.Fatal("Symlink should exist before rollback")
	}

	// Perform rollback
	if err := tracker.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify symlink was removed
	if _, err := os.Lstat(targetPath); !os.IsNotExist(err) {
		t.Error("Symlink should be removed after rollback")
	}
}

func TestRollbackUpdatedSymlink(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sokru-rollback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetPath := filepath.Join(tempDir, "test-link")
	oldSource := filepath.Join(tempDir, "old-source")
	newSource := filepath.Join(tempDir, "new-source")

	// Create source files
	if err := os.WriteFile(oldSource, []byte("old"), 0644); err != nil {
		t.Fatalf("Failed to create old source: %v", err)
	}
	if err := os.WriteFile(newSource, []byte("new"), 0644); err != nil {
		t.Fatalf("Failed to create new source: %v", err)
	}

	// Create initial symlink (old)
	if err := os.Symlink(oldSource, targetPath); err != nil {
		t.Fatalf("Failed to create initial symlink: %v", err)
	}

	// Update symlink to new source
	if err := os.Remove(targetPath); err != nil {
		t.Fatalf("Failed to remove old symlink: %v", err)
	}
	if err := os.Symlink(newSource, targetPath); err != nil {
		t.Fatalf("Failed to create new symlink: %v", err)
	}

	// Track the update
	tracker := rollback.NewTracker()
	tracker.TrackUpdated(targetPath, newSource, oldSource)

	// Verify symlink points to new source
	link, err := os.Readlink(targetPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if link != newSource {
		t.Errorf("Expected symlink to point to %s, got %s", newSource, link)
	}

	// Perform rollback
	if err := tracker.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify symlink was restored to old source
	link, err = os.Readlink(targetPath)
	if err != nil {
		t.Fatalf("Failed to read symlink after rollback: %v", err)
	}
	if link != oldSource {
		t.Errorf("Expected symlink to be restored to %s, got %s", oldSource, link)
	}
}

func TestRollbackMultipleActions(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sokru-rollback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple symlinks
	tracker := rollback.NewTracker()

	for i := 1; i <= 3; i++ {
		targetPath := filepath.Join(tempDir, fmt.Sprintf("link-%d", i))
		sourcePath := filepath.Join(tempDir, fmt.Sprintf("source-%d", i))

		// Create source
		if err := os.WriteFile(sourcePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create source %d: %v", i, err)
		}

		// Create symlink
		if err := os.Symlink(sourcePath, targetPath); err != nil {
			t.Fatalf("Failed to create symlink %d: %v", i, err)
		}

		tracker.TrackCreated(targetPath, sourcePath)
	}

	if tracker.Count() != 3 {
		t.Errorf("Expected 3 actions, got %d", tracker.Count())
	}

	// Perform rollback
	if err := tracker.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify all symlinks were removed
	for i := 1; i <= 3; i++ {
		targetPath := filepath.Join(tempDir, fmt.Sprintf("link-%d", i))
		if _, err := os.Lstat(targetPath); !os.IsNotExist(err) {
			t.Errorf("Symlink %d should be removed after rollback", i)
		}
	}
}

func TestRollbackEmptyTracker(t *testing.T) {
	tracker := rollback.NewTracker()

	// Rollback with no actions should not error
	if err := tracker.Rollback(); err != nil {
		t.Errorf("Rollback of empty tracker should not error: %v", err)
	}
}

func TestRollbackDisabledTracker(t *testing.T) {
	tracker := rollback.NewTracker()
	tracker.Disable()

	// Rollback when disabled should do nothing
	if err := tracker.Rollback(); err != nil {
		t.Errorf("Rollback of disabled tracker should not error: %v", err)
	}
}
