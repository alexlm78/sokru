// Package cmd
// Description: This file contains the restore command for the cli tool
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/alexlm78/sokru/internal/backup"
	"github.com/alexlm78/sokru/internal/i18n"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore files from backups",
	Long:  `This command allows you to restore files from previous backups.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Restore command - use subcommands: list, apply, delete")
	},
}

// restoreListCmd lists all available backups
var restoreListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available backups",
	Long:  `This command lists all available backup sessions.`,
	Run:   RestoreListFunc,
}

// restoreApplyCmd restores a specific backup
var restoreApplyCmd = &cobra.Command{
	Use:   "apply [backup-id]",
	Short: "Restore files from a specific backup",
	Long:  `This command restores files from a specific backup session.`,
	Args:  cobra.ExactArgs(1),
	Run:   RestoreApplyFunc,
}

// restoreDeleteCmd deletes a backup
var restoreDeleteCmd = &cobra.Command{
	Use:   "delete [backup-id]",
	Short: "Delete a specific backup",
	Long:  `This command deletes a specific backup session.`,
	Args:  cobra.ExactArgs(1),
	Run:   RestoreDeleteFunc,
}

func RestoreListFunc(cmd *cobra.Command, args []string) {
	backupDir, err := backup.GetDefaultBackupDir()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorGettingBackupDir, err))
	}

	manager := backup.NewManager(backupDir)
	backups, err := manager.ListBackups()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorListingBackups, err))
	}

	if len(backups) == 0 {
		fmt.Println(i18n.Info(i18n.MsgNoBackupsFound))
		return
	}

	// Sort backups by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	fmt.Println(i18n.T(i18n.MsgAvailableBackups))
	fmt.Println("================")
	fmt.Println()

	for _, bk := range backups {
		fmt.Printf("ID: %s\n", bk.ID)
		fmt.Printf("  %s: %s\n", i18n.T(i18n.MsgTimestamp), bk.Timestamp.Format(time.RFC3339))
		fmt.Printf("  %s: %s\n", i18n.T(i18n.MsgCommand), bk.Command)
		fmt.Printf("  %s: %d\n", i18n.T(i18n.MsgFiles), len(bk.Entries))
		fmt.Println()
	}

	fmt.Printf("%s: %d\n", i18n.T(i18n.MsgTotalBackups), len(backups))
}

func RestoreApplyFunc(cmd *cobra.Command, args []string) {
	backupID := args[0]

	backupDir, err := backup.GetDefaultBackupDir()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorGettingBackupDir, err))
	}

	manager := backup.NewManager(backupDir)

	// Load metadata to show what will be restored
	metadata, err := manager.LoadMetadata(backupID)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorLoadingBackup, backupID, err))
	}

	fmt.Printf("%s: %s\n", i18n.T(i18n.MsgRestoringBackup), backupID)
	fmt.Printf("%s: %s\n", i18n.T(i18n.MsgBackupCreated), metadata.Timestamp.Format(time.RFC3339))
	fmt.Printf("%s: %d\n\n", i18n.T(i18n.MsgFilesToRestore), len(metadata.Entries))

	// Show files that will be restored
	fmt.Println(i18n.T(i18n.MsgFilesInBackup))
	for _, entry := range metadata.Entries {
		if entry.IsSymlink {
			fmt.Printf("  [symlink] %s -> %s\n", entry.OriginalPath, entry.SymlinkTarget)
		} else {
			fmt.Printf("  [file]    %s\n", entry.OriginalPath)
		}
	}
	fmt.Println()

	// Perform restore
	fmt.Println(i18n.Info(i18n.MsgRestoringFiles))
	if err := manager.RestoreBackup(backupID); err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgRestoreFailed, err))
	}

	fmt.Println(i18n.Success(i18n.MsgRestoreComplete))
}

func RestoreDeleteFunc(cmd *cobra.Command, args []string) {
	backupID := args[0]

	backupDir, err := backup.GetDefaultBackupDir()
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorGettingBackupDir, err))
	}

	manager := backup.NewManager(backupDir)

	// Load metadata to show what will be deleted
	metadata, err := manager.LoadMetadata(backupID)
	if err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorLoadingBackup, backupID, err))
	}

	fmt.Printf("%s: %s\n", i18n.T(i18n.MsgDeletingBackup), backupID)
	fmt.Printf("%s: %s\n", i18n.T(i18n.MsgBackupCreated), metadata.Timestamp.Format(time.RFC3339))
	fmt.Printf("%s: %d\n\n", i18n.T(i18n.MsgFiles), len(metadata.Entries))

	// Delete backup
	if err := manager.DeleteBackup(backupID); err != nil {
		log.Fatalf("%s", i18n.Error(i18n.MsgErrorDeletingBackup, err))
	}

	fmt.Println(i18n.Success(i18n.MsgBackupDeleted, backupID))
}

func init() {
	restoreCmd.AddCommand(restoreListCmd)
	restoreCmd.AddCommand(restoreApplyCmd)
	restoreCmd.AddCommand(restoreDeleteCmd)
	rootCmd.AddCommand(restoreCmd)
}
