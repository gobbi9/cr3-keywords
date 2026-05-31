package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Script returns the completion/module script for the given shell target.
func Script(target, commandName string) (string, error) {
	if strings.TrimSpace(commandName) == "" {
		commandName = "cr3"
	}

	switch target {
	case "nushell":
		return renderNushellModule(commandName), nil
	case "zsh":
		return renderZshCompletion(commandName), nil
	case "bash":
		return renderBashCompletion(commandName), nil
	default:
		return "", fmt.Errorf("unsupported completion target: %s", target)
	}
}

// InstallPath returns the conventional per-user completion/module path.
func InstallPath(target, commandName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}

	switch target {
	case "nushell":
		dir, err := nushellAutoloadDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, commandName+".nu"), nil
	case "zsh":
		zDotDir := strings.TrimSpace(os.Getenv("ZDOTDIR"))
		if zDotDir != "" {
			return filepath.Join(zDotDir, "completions", "_"+commandName), nil
		}
		return filepath.Join(home, ".zsh", "completions", "_"+commandName), nil
	case "bash":
		if runtime.GOOS == "windows" {
			localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
			if localAppData == "" {
				localAppData = filepath.Join(home, "AppData", "Local")
			}
			return filepath.Join(localAppData, "bash-completion", "completions", commandName), nil
		}
		dataHome := strings.TrimSpace(os.Getenv("XDG_DATA_HOME"))
		if dataHome == "" {
			dataHome = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(dataHome, "bash-completion", "completions", commandName), nil
	default:
		return "", fmt.Errorf("unsupported install target: %s", target)
	}
}
