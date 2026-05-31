package shell

import (
	"fmt"
	"os"
	"path/filepath"
)

// Install writes a shell completion/module file to the conventional
// per-user location and returns the installed file path.
func Install(target, commandName string) (string, error) {
	script, err := Script(target, commandName)
	if err != nil {
		return "", err
	}

	path, err := InstallPath(target, commandName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create completion directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(script), 0o644); err != nil {
		return "", fmt.Errorf("write completion file: %w", err)
	}
	return path, nil
}
