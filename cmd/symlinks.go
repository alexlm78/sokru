// Package cmd
// Description: This file contains the symlinks command for the cli tool. It is used to manage the symlinks in the system.
// (c) 2023 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type SymlinkConfig struct {
	Link map[string]string `yaml:"link"`
}

// symlinksCmd represents the symlinks command
var symlinksCmd = &cobra.Command{
	Use:   "symlinks",
	Short: "Manage the symlinks",
	Long:  `This command will allow you to manage the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Main sumlinks function called")
	},
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the symlinks",
	Long:  `This command will install the symlinks in the system.`,
	Run:   InstallSymlinksFunc,
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the symlinks",
	Long:  `This command will uninstall the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling symlinks")
	},
}

// listCMd represents the list symlink configured (symlinks.yaml)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the symlinks",
	Long:  `This command will list the symlinks in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing symlinks")
	},
}

// helpCmd represents the help command
var symhelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help for the symlinks",
	Long:  `This command will show the help for the symlinks command.`,
	Run:   HelpSymlinksFunc,
}

func InstallSymlinksFunc(*cobra.Command, []string) {
	// Read the YAML file
	data, err := ioutil.ReadFile("/Users/alexlm78/.dotfiles/symlinks.yml")
	if err != nil {
		log.Fatalf("Error leyendo el archivo YAML: %v", err)
	}

	// Unmarshal del YAML al struct
	var config []SymlinkConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parseando el YAML: %v", err)
	}

	// Itera sobre los items y crea los enlaces simbólicos
	for _, entry := range config {
		for target, source := range entry.Link {
			targetPath := expandPath(target)
			sourcePath := expandPath(source)

			existingLink, err := os.Readlink(targetPath)
			if err == nil {
				if existingLink == sourcePath {
					fmt.Printf("El enlace simbólico ya existe y es correcto para %s -> %s\n", targetPath, sourcePath)
					continue
				} else {
					// Eliminar el enlace existente si apunta a un destino diferente
					err = os.Remove(targetPath)
					if err != nil {
						log.Printf("Error eliminando el symlink existente en %s: %v", targetPath, err)
						continue
					}
					fmt.Printf("Enlace simbólico existente eliminado para %s\n", targetPath)
				}
			} else if !os.IsNotExist(err) {
				log.Printf("Error verificando el symlink en %s: %v", targetPath, err)
				continue
			}

			// Crear el nuevo enlace simbólico
			err = os.Symlink(sourcePath, targetPath)
			if err != nil {
				log.Printf("Error creando el symlink de %s a %s: %v", targetPath, sourcePath, err)
			} else {
				fmt.Printf("Enlace simbólico creado de %s a %s\n", targetPath, sourcePath)
			}
		}
	}
}

func init() {
	symlinksCmd.AddCommand(installCmd)
	symlinksCmd.AddCommand(uninstallCmd)
	symlinksCmd.AddCommand(listCmd)
	symlinksCmd.AddCommand(symhelpCmd)
	rootCmd.AddCommand(symlinksCmd)
}

// expandPath reemplaza ~ con el directorio HOME
func expandPath(path string) string {
	if path[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
