# sokru

### Considerations
- The symlinks file contais the symlinks that will be created in the system. a part for every OS soported.
- The configuration options are stored in a file called `config.yml` in the ~/.sokru/ dir.




### Options
#### Main commands
```bash
sok init                    # Initialize the configuration
sok apply                   # Apply the changes in memory and reload the symlinks and dotfiles
sok version                 # Show the version
sok config                  # Show the configuration options
sok symlinks                # Show the symlinks options
sok help                    # Show this help
```
#### Config
```bash
sok config dotDir <dir>     # Set the directory where the dotfiles are stored (default: ~/.dotfiles)
sok config symlinks <file>  # Set the file that contains the symlinks (default: ~/.dotfiles/symlinks.yml)
sok config os <os>          # Set the OS to use (default: linux)
sok config verbose <bool>   # Set the verbose mode (default: false)
sok config dryRun <bool>    # Set the dry run mode (default: false)
sok config help             # Show this help
```
#### Symlinks
```bash
sok symlinks install        # Install the symlinks
sok symlinks uninstall      # Uninstall the symlinks
sok symlinks list           # List the symlinks
sok symlinks help           # Show this help
```
