// Package rollback
// Description: Rollback mechanism for symlink operations
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package rollback

import (
	"fmt"
	"os"
)

// SymlinkAction represents an action performed on a symlink
type SymlinkAction struct {
	Type         ActionType
	TargetPath   string
	SourcePath   string
	PreviousLink string // For updates, stores the previous symlink target
	WasSymlink   bool   // Whether the target was a symlink before
}

// ActionType represents the type of action performed
type ActionType int

const (
	ActionCreated ActionType = iota // New symlink created
	ActionUpdated                   // Existing symlink updated
	ActionRemoved                   // Symlink removed
)

// Tracker tracks symlink operations for rollback
type Tracker struct {
	actions []SymlinkAction
	enabled bool
}

// NewTracker creates a new rollback tracker
func NewTracker() *Tracker {
	return &Tracker{
		actions: make([]SymlinkAction, 0),
		enabled: true,
	}
}

// Enable enables rollback tracking
func (t *Tracker) Enable() {
	t.enabled = true
}

// Disable disables rollback tracking
func (t *Tracker) Disable() {
	t.enabled = false
}

// IsEnabled returns whether rollback tracking is enabled
func (t *Tracker) IsEnabled() bool {
	return t.enabled
}

// TrackCreated records a symlink creation
func (t *Tracker) TrackCreated(targetPath, sourcePath string) {
	if !t.enabled {
		return
	}

	t.actions = append(t.actions, SymlinkAction{
		Type:       ActionCreated,
		TargetPath: targetPath,
		SourcePath: sourcePath,
	})
}

// TrackUpdated records a symlink update
func (t *Tracker) TrackUpdated(targetPath, sourcePath, previousLink string) {
	if !t.enabled {
		return
	}

	t.actions = append(t.actions, SymlinkAction{
		Type:         ActionUpdated,
		TargetPath:   targetPath,
		SourcePath:   sourcePath,
		PreviousLink: previousLink,
		WasSymlink:   true,
	})
}

// TrackRemoved records a symlink removal
func (t *Tracker) TrackRemoved(targetPath, sourcePath string) {
	if !t.enabled {
		return
	}

	t.actions = append(t.actions, SymlinkAction{
		Type:       ActionRemoved,
		TargetPath: targetPath,
		SourcePath: sourcePath,
		WasSymlink: true,
	})
}

// GetActions returns all tracked actions
func (t *Tracker) GetActions() []SymlinkAction {
	return t.actions
}

// Clear clears all tracked actions
func (t *Tracker) Clear() {
	t.actions = make([]SymlinkAction, 0)
}

// Rollback reverts all tracked actions in reverse order
func (t *Tracker) Rollback() error {
	if !t.enabled {
		return nil
	}

	var errors []error

	// Process actions in reverse order
	for i := len(t.actions) - 1; i >= 0; i-- {
		action := t.actions[i]

		switch action.Type {
		case ActionCreated:
			// Remove the created symlink
			if err := os.Remove(action.TargetPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("failed to remove %s: %w", action.TargetPath, err))
			}

		case ActionUpdated:
			// Restore the previous symlink
			// First remove the current one
			if err := os.Remove(action.TargetPath); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("failed to remove %s: %w", action.TargetPath, err))
				continue
			}

			// Recreate the previous symlink
			if action.PreviousLink != "" {
				if err := os.Symlink(action.PreviousLink, action.TargetPath); err != nil {
					errors = append(errors, fmt.Errorf("failed to restore %s -> %s: %w",
						action.TargetPath, action.PreviousLink, err))
				}
			}

		case ActionRemoved:
			// Recreate the removed symlink
			if err := os.Symlink(action.SourcePath, action.TargetPath); err != nil {
				errors = append(errors, fmt.Errorf("failed to recreate %s -> %s: %w",
					action.TargetPath, action.SourcePath, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("rollback completed with %d error(s): %v", len(errors), errors)
	}

	return nil
}

// Count returns the number of tracked actions
func (t *Tracker) Count() int {
	return len(t.actions)
}

// HasActions returns whether any actions have been tracked
func (t *Tracker) HasActions() bool {
	return len(t.actions) > 0
}
