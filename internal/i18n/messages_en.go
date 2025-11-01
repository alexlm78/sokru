// Package i18n
// Description: English translations
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package i18n

// Message keys
const (
	// Common messages
	MsgErrorLoadingConfig   MessageKey = "error_loading_config"
	MsgErrorReadingFile     MessageKey = "error_reading_file"
	MsgErrorParsingYAML     MessageKey = "error_parsing_yaml"
	MsgErrorCreatingSymlink MessageKey = "error_creating_symlink"
	MsgErrorRemovingSymlink MessageKey = "error_removing_symlink"
	MsgErrorCheckingFile    MessageKey = "error_checking_file"
	MsgErrorReadingSymlink  MessageKey = "error_reading_symlink"
	MsgFileNotFound         MessageKey = "file_not_found"
	MsgInvalidBoolValue     MessageKey = "invalid_bool_value"
	MsgInvalidOS            MessageKey = "invalid_os"

	// Symlinks messages
	MsgSymlinkFileNotFound    MessageKey = "symlink_file_not_found"
	MsgReadingSymlinksFrom    MessageKey = "reading_symlinks_from"
	MsgFoundConfigurations    MessageKey = "found_configurations"
	MsgDryRunWouldCreate      MessageKey = "dry_run_would_create"
	MsgSymlinkAlreadyExists   MessageKey = "symlink_already_exists"
	MsgExistingSymlinkRemoved MessageKey = "existing_symlink_removed"
	MsgSymlinkCreated         MessageKey = "symlink_created"
	MsgSymlinkNotFound        MessageKey = "symlink_not_found"
	MsgNotSymlink             MessageKey = "not_symlink"
	MsgSymlinkWrongTarget     MessageKey = "symlink_wrong_target"
	MsgDryRunWouldRemove      MessageKey = "dry_run_would_remove"
	MsgSymlinkRemoved         MessageKey = "symlink_removed"
	MsgUninstallSummary       MessageKey = "uninstall_summary"
	MsgWouldRemove            MessageKey = "would_remove"
	MsgRemoved                MessageKey = "removed"
	MsgNotFound               MessageKey = "not_found"
	MsgNotSymlinks            MessageKey = "not_symlinks"
	MsgSkipped                MessageKey = "skipped"
	MsgSymlinksStatus         MessageKey = "symlinks_status"
	MsgStatus                 MessageKey = "status"
	MsgTarget                 MessageKey = "target"
	MsgSource                 MessageKey = "source"
	MsgSummary                MessageKey = "summary"
	MsgInstalledCorrectly     MessageKey = "installed_correctly"
	MsgWrongTarget            MessageKey = "wrong_target"
	MsgNotInstalled           MessageKey = "not_installed"
	MsgRegularFileExists      MessageKey = "regular_file_exists"
	MsgTotalSymlinks          MessageKey = "total_symlinks"
	MsgLegend                 MessageKey = "legend"
	MsgLegendInstalled        MessageKey = "legend_installed"
	MsgLegendWrongTarget      MessageKey = "legend_wrong_target"
	MsgLegendNotInstalled     MessageKey = "legend_not_installed"
	MsgLegendRegularFile      MessageKey = "legend_regular_file"
	MsgFilteringByOS          MessageKey = "filtering_by_os"
	MsgUsingCommonLinks       MessageKey = "using_common_links"
	MsgUsingOSSpecificLinks   MessageKey = "using_os_specific_links"
	
	// Rollback messages
	MsgRollbackStarting       MessageKey = "rollback_starting"
	MsgRollbackComplete       MessageKey = "rollback_complete"
	MsgRollbackFailed         MessageKey = "rollback_failed"
	MsgRollbackAction         MessageKey = "rollback_action"
	MsgRollbackRemoved        MessageKey = "rollback_removed"
	MsgRollbackRestored       MessageKey = "rollback_restored"
	MsgRollbackRecreated      MessageKey = "rollback_recreated"

	// Apply messages
	MsgApplyingChanges    MessageKey = "applying_changes"
	MsgConfigReloaded     MessageKey = "config_reloaded"
	MsgChangesToApply     MessageKey = "changes_to_apply"
	MsgToCreate           MessageKey = "to_create"
	MsgToUpdate           MessageKey = "to_update"
	MsgAlreadyCorrect     MessageKey = "already_correct"
	MsgNoChangesNeeded    MessageKey = "no_changes_needed"
	MsgDryRunNoChanges    MessageKey = "dry_run_no_changes"
	MsgApplyingChangesNow MessageKey = "applying_changes_now"
	MsgCreated            MessageKey = "created"
	MsgUpdated            MessageKey = "updated"
	MsgFailed             MessageKey = "failed"
	MsgApplySummary       MessageKey = "apply_summary"
	MsgConfigApplied      MessageKey = "config_applied"

	// Init messages
	MsgInitializing        MessageKey = "initializing"
	MsgCreatedDirectory    MessageKey = "created_directory"
	MsgDirectoryExists     MessageKey = "directory_exists"
	MsgCreatedConfigFile   MessageKey = "created_config_file"
	MsgConfigFileExists    MessageKey = "config_file_exists"
	MsgCreatedSymlinksFile MessageKey = "created_symlinks_file"
	MsgSymlinksFileExists  MessageKey = "symlinks_file_exists"
	MsgInitSummary         MessageKey = "init_summary"
	MsgNextSteps           MessageKey = "next_steps"
	MsgAddDotfiles         MessageKey = "add_dotfiles"
	MsgEditSymlinks        MessageKey = "edit_symlinks"
	MsgInstallSymlinks     MessageKey = "install_symlinks"
	MsgViewConfig          MessageKey = "view_config"
	MsgInitialized         MessageKey = "initialized"

	// Config messages
	MsgConfigFile          MessageKey = "config_file"
	MsgCurrentConfig       MessageKey = "current_config"
	MsgDotfilesDirectory   MessageKey = "dotfiles_directory"
	MsgSymlinksFile        MessageKey = "symlinks_file"
	MsgOperatingSystem     MessageKey = "operating_system"
	MsgVerbose             MessageKey = "verbose"
	MsgDryRun              MessageKey = "dry_run"
	MsgCurrentDotfilesDir  MessageKey = "current_dotfiles_dir"
	MsgDotfilesDirSet      MessageKey = "dotfiles_dir_set"
	MsgCurrentSymlinksFile MessageKey = "current_symlinks_file"
	MsgSymlinksFileSet     MessageKey = "symlinks_file_set"
	MsgCurrentVerbose      MessageKey = "current_verbose"
	MsgVerboseSet          MessageKey = "verbose_set"
	MsgCurrentDryRun       MessageKey = "current_dry_run"
	MsgDryRunSet           MessageKey = "dry_run_set"
	MsgCurrentOS           MessageKey = "current_os"
	MsgOSSet               MessageKey = "os_set"

	// Help messages
	MsgHelpSymlinks MessageKey = "help_symlinks"
	MsgHelpConfig   MessageKey = "help_config"
	MsgHelp         MessageKey = "help"

	// Version messages
	MsgVersion MessageKey = "version"
)

func getEnglishMessages() map[MessageKey]string {
	return map[MessageKey]string{
		// Common messages
		MsgErrorLoadingConfig:   "Error loading configuration: %v",
		MsgErrorReadingFile:     "Error reading file: %v",
		MsgErrorParsingYAML:     "Error parsing YAML: %v",
		MsgErrorCreatingSymlink: "Error creating symlink from %s to %s: %v",
		MsgErrorRemovingSymlink: "Error removing symlink at %s: %v",
		MsgErrorCheckingFile:    "Error checking file at %s: %v",
		MsgErrorReadingSymlink:  "Error reading symlink at %s: %v",
		MsgFileNotFound:         "File not found: %s",
		MsgInvalidBoolValue:     "Invalid boolean value. Use 'true' or 'false'",
		MsgInvalidOS:            "Invalid OS '%s'. Valid options are: linux, darwin, windows",

		// Symlinks messages
		MsgSymlinkFileNotFound:    "Symlinks file not found: %s\nPlease create the file or update the configuration with: sok config symlinkfile <path>",
		MsgReadingSymlinksFrom:    "Reading symlinks configuration from: %s",
		MsgFoundConfigurations:    "Found %d symlink configuration(s)",
		MsgDryRunWouldCreate:      "[DRY-RUN] Would create symlink: %s -> %s",
		MsgSymlinkAlreadyExists:   "Symlink already exists and is correct: %s -> %s",
		MsgExistingSymlinkRemoved: "Existing symlink removed: %s",
		MsgSymlinkCreated:         "Symlink created: %s -> %s",
		MsgSymlinkNotFound:        "Symlink not found (already removed): %s",
		MsgNotSymlink:             "Warning: %s is not a symlink, skipping (will not remove regular files)",
		MsgSymlinkWrongTarget:     "Symlink points to different source: %s -> %s (expected: %s), skipping",
		MsgDryRunWouldRemove:      "[DRY-RUN] Would remove symlink: %s -> %s",
		MsgSymlinkRemoved:         "Symlink removed: %s -> %s",
		MsgUninstallSummary:       "Uninstall Summary",
		MsgWouldRemove:            "Would remove: %d symlink(s)",
		MsgRemoved:                "Removed: %d symlink(s)",
		MsgNotFound:               "Not found: %d symlink(s)",
		MsgNotSymlinks:            "Not symlinks (skipped): %d file(s)",
		MsgSkipped:                "Skipped: %d symlink(s)",
		MsgSymlinksStatus:         "Symlinks Status:",
		MsgStatus:                 "Status",
		MsgTarget:                 "Target",
		MsgSource:                 "Source",
		MsgSummary:                "Summary:",
		MsgInstalledCorrectly:     "Installed correctly:    %d",
		MsgWrongTarget:            "Wrong target:          %d",
		MsgNotInstalled:           "Not installed:         %d",
		MsgRegularFileExists:      "Regular file exists:   %d",
		MsgTotalSymlinks:          "Total symlinks configured: %d",
		MsgLegend:                 "Legend:",
		MsgLegendInstalled:        "✅ = Symlink installed and points to correct source",
		MsgLegendWrongTarget:      "⚠️  = Symlink exists but points to different source",
		MsgLegendNotInstalled:     "❌ = Symlink not installed",
		MsgLegendRegularFile:      "⛔ = Regular file exists at target location (not a symlink)",
		MsgFilteringByOS:          "Filtering symlinks for OS: %s",
		MsgUsingCommonLinks:       "Using common links (all OS)",
		MsgUsingOSSpecificLinks:   "Using OS-specific links for: %s",
		
		// Rollback messages
		MsgRollbackStarting:       "Error occurred, starting rollback of %d action(s)...",
		MsgRollbackComplete:       "Rollback completed successfully",
		MsgRollbackFailed:         "Rollback completed with errors: %v",
		MsgRollbackAction:         "Rolling back action %d/%d",
		MsgRollbackRemoved:        "Removed created symlink: %s",
		MsgRollbackRestored:       "Restored previous symlink: %s -> %s",
		MsgRollbackRecreated:      "Recreated removed symlink: %s -> %s",

		// Apply messages
		MsgApplyingChanges:    "Applying configuration changes...",
		MsgConfigReloaded:     "Configuration reloaded from disk",
		MsgChangesToApply:     "Changes to Apply",
		MsgToCreate:           "To Create (%d):",
		MsgToUpdate:           "To Update (%d):",
		MsgAlreadyCorrect:     "Already Correct: %d",
		MsgNoChangesNeeded:    "No changes needed - all symlinks are up to date!",
		MsgDryRunNoChanges:    "[DRY-RUN] No changes were made",
		MsgApplyingChangesNow: "Applying Changes",
		MsgCreated:            "Created: %s -> %s",
		MsgUpdated:            "Updated: %s -> %s",
		MsgFailed:             "Failed: %d symlink(s)",
		MsgApplySummary:       "Apply Summary",
		MsgConfigApplied:      "Configuration applied successfully!",

		// Init messages
		MsgInitializing:        "Initializing sokru...",
		MsgCreatedDirectory:    "Created directory: %s",
		MsgDirectoryExists:     "Directory already exists: %s",
		MsgCreatedConfigFile:   "Created config file: %s",
		MsgConfigFileExists:    "Config file already exists: %s",
		MsgCreatedSymlinksFile: "Created symlinks file: %s",
		MsgSymlinksFileExists:  "Symlinks file already exists: %s",
		MsgInitSummary:         "Initialization Summary",
		MsgNextSteps:           "Next Steps",
		MsgAddDotfiles:         "1. Add your dotfiles to the shell directories:",
		MsgEditSymlinks:        "2. Edit the symlinks configuration:",
		MsgInstallSymlinks:     "3. Install your symlinks:",
		MsgViewConfig:          "4. View your configuration:",
		MsgInitialized:         "Sokru initialized successfully!",

		// Config messages
		MsgConfigFile:          "Configuration file: %s",
		MsgCurrentConfig:       "Current configuration:",
		MsgDotfilesDirectory:   "Dotfiles Directory: %s",
		MsgSymlinksFile:        "Symlinks File:      %s",
		MsgOperatingSystem:     "Operating System:   %s",
		MsgVerbose:             "Verbose:            %v",
		MsgDryRun:              "Dry Run:            %v",
		MsgCurrentDotfilesDir:  "Current dotfiles directory: %s",
		MsgDotfilesDirSet:      "Dotfiles directory set to: %s",
		MsgCurrentSymlinksFile: "Current symlinks file: %s",
		MsgSymlinksFileSet:     "Symlinks file set to: %s",
		MsgCurrentVerbose:      "Current verbose setting: %v",
		MsgVerboseSet:          "Verbose mode set to: %v",
		MsgCurrentDryRun:       "Current dry-run setting: %v",
		MsgDryRunSet:           "Dry-run mode set to: %v",
		MsgCurrentOS:           "Current OS setting: %s",
		MsgOSSet:               "OS set to: %s",

		// Help messages
		MsgHelpSymlinks: "SymLinks command help::",
		MsgHelpConfig:   "Config command help::",
		MsgHelp:         "Sokru help::",

		// Version messages
		MsgVersion: "Sok v0.1 -- HEAD",
	}
}
