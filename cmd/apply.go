// Package cmd
// Description: This file contains the apply command for the cli tool. It is used to apply the changes in memory and reload the symlinks and dotfiles.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alexlm78/sokru/internal/config"
	"github.com/alexlm78/sokru/internal/i18n"
	"github.com/alexlm78/sokru/internal/rollback"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the changes in memory and reload the symlinks and dotfiles",
	Long:  `This command will apply the changes in memory and reload the symlinks and dotfiles.`,
	Run:   ApplyFunc,
}

func ApplyFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Applying configuration changes...")
	fmt.Println()

	// Get current config to preserve flags
	currentCfg, _ := config.GetConfig()
	preserveDryRun := currentCfg.DryRun
	preserveVerbose := currentCfg.Verbose

	// 1. Reload configuration from disk
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Preserve command-line flags
	cfg.DryRun = preserveDryRun
	cfg.Verbose = preserveVerbose
	config.SetConfig(cfg)

	if cfg.Verbose {
		fmt.Println("✓ Configuration reloaded from disk")
	}

	// 2. Get current symlinks from YAML
	symlinkFile := expandPath(cfg.SymlinksFile)

	if _, err := os.Stat(symlinkFile); os.IsNotExist(err) {
		log.Fatalf("Symlinks file not found: %s\nPlease create the file or update the configuration with: sok config symlinkfile <path>", symlinkFile)
	}

	data, err := os.ReadFile(symlinkFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var symlinkConfigs []SymlinkConfig
	err = yaml.Unmarshal(data, &symlinkConfigs)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	if cfg.Verbose {
		fmt.Printf("✓ Read %d symlink configuration(s) from: %s\n", len(symlinkConfigs), symlinkFile)
	}

	// 3. Build a map of configured symlinks (filtered by OS)
	configuredSymlinks := make(map[string]string)
	for _, entry := range symlinkConfigs {
		// Get links for current OS
		links := entry.getLinksForOS(cfg.OS)

		for target, source := range links {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)
			configuredSymlinks[targetPath] = sourcePath
		}
	}

	// 4. Check existing symlinks and track changes
	var toCreate, toUpdate []string
	var alreadyCorrect int

	// Check what needs to be created or updated
	for targetPath, sourcePath := range configuredSymlinks {
		existingLink, err := os.Readlink(targetPath)
		if os.IsNotExist(err) {
			// Symlink doesn't exist, needs to be created
			toCreate = append(toCreate, fmt.Sprintf("%s -> %s", targetPath, sourcePath))
		} else if err == nil {
			if existingLink != sourcePath {
				// Symlink exists but points to wrong source
				toUpdate = append(toUpdate, fmt.Sprintf("%s: %s -> %s", targetPath, existingLink, sourcePath))
			} else {
				// Symlink is correct
				alreadyCorrect++
			}
		}
	}

	// 5. Show what will change
	fmt.Println("=== Changes to Apply ===")

	if len(toCreate) > 0 {
		fmt.Printf("\n📝 To Create (%d):\n", len(toCreate))
		for _, item := range toCreate {
			fmt.Printf("  + %s\n", item)
		}
	}

	if len(toUpdate) > 0 {
		fmt.Printf("\n🔄 To Update (%d):\n", len(toUpdate))
		for _, item := range toUpdate {
			fmt.Printf("  ~ %s\n", item)
		}
	}

	if alreadyCorrect > 0 {
		fmt.Printf("\n✅ Already Correct: %d\n", alreadyCorrect)
	}

	if len(toCreate) == 0 && len(toUpdate) == 0 {
		fmt.Println("\n✓ No changes needed - all symlinks are up to date!")
		return
	}

	// 6. Apply changes (unless dry-run)
	if cfg.DryRun {
		fmt.Println("\n[DRY-RUN] No changes were made")
		return
	}

	fmt.Println("\n=== Applying Changes ===")

	// Create rollback tracker
	tracker := rollback.NewTracker()
	var created, updated, failed int
	var hasError bool

	// Create new symlinks
	for targetPath, sourcePath := range configuredSymlinks {
		existingLink, err := os.Readlink(targetPath)

		if os.IsNotExist(err) {
			// Create new symlink
			if err := os.Symlink(sourcePath, targetPath); err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCreatingSymlink, targetPath, sourcePath, err))
				failed++
				hasError = true
				break
			} else {
				if cfg.Verbose {
					fmt.Println(i18n.Success(i18n.MsgCreated, targetPath, sourcePath))
				}
				created++
				tracker.TrackCreated(targetPath, sourcePath)
			}
		} else if err == nil && existingLink != sourcePath {
			// Update existing symlink
			if err := os.Remove(targetPath); err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorRemovingSymlink, targetPath, err))
				failed++
				hasError = true
				break
			}
			if err := os.Symlink(sourcePath, targetPath); err != nil {
				log.Printf("%s", i18n.Error(i18n.MsgErrorCreatingSymlink, targetPath, sourcePath, err))
				failed++
				hasError = true
				break
			} else {
				if cfg.Verbose {
					fmt.Println(i18n.Success(i18n.MsgUpdated, targetPath, sourcePath))
				}
				updated++
				tracker.TrackUpdated(targetPath, sourcePath, existingLink)
			}
		}
	}

	// Perform rollback if there was an error
	if hasError && tracker.HasActions() {
		fmt.Println()
		fmt.Println(i18n.Warning(i18n.MsgRollbackStarting, tracker.Count()))

		if err := tracker.Rollback(); err != nil {
			log.Printf("%s", i18n.Error(i18n.MsgRollbackFailed, err))
		} else {
			fmt.Println(i18n.Success(i18n.MsgRollbackComplete))
		}

		os.Exit(1)
	}

	// 7. Summary
	fmt.Printf("\n=== %s ===\n", i18n.T(i18n.MsgApplySummary))
	if created > 0 {
		fmt.Printf("Created: %d symlink(s)\n", created)
	}
	if updated > 0 {
		fmt.Printf("Updated: %d symlink(s)\n", updated)
	}
	if alreadyCorrect > 0 {
		fmt.Printf("Already correct: %d symlink(s)\n", alreadyCorrect)
	}
	if failed > 0 {
		fmt.Printf("Failed: %d symlink(s)\n", failed)
	}

	fmt.Println(i18n.Success(i18n.MsgConfigApplied))
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
