package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func nushellAutoloadDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "nushell", "vendor", "autoload"), nil
	case "windows":
		appData := strings.TrimSpace(os.Getenv("APPDATA"))
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "nushell", "vendor", "autoload"), nil
	default:
		config := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME"))
		if config == "" {
			config = filepath.Join(home, ".config")
		}
		return filepath.Join(config, "nushell", "vendor", "autoload"), nil
	}
}

func renderNushellModule(commandName string) (string, error) {
	return renderTemplate("nushell.tmpl.nu", commandName)
}
