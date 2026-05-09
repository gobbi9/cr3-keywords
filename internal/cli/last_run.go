package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SaveLastRunCR3Path stores the absolute CR3 directory path from the last
// successful run in the local state file under ~/.cr3-keywords.
func SaveLastRunCR3Path(cr3Path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}

	abs, err := filepath.Abs(cr3Path)
	if err != nil {
		return fmt.Errorf("resolve absolute CR3 path: %w", err)
	}

	stateFile := lastRunPathFile(home)
	if err := os.MkdirAll(filepath.Dir(stateFile), 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}
	if err := os.WriteFile(stateFile, []byte(abs+"\n"), 0o644); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}

	return nil
}

// ClearLastRunXMP loads the last-run CR3 directory, lists .xmp files in that
// directory, asks for confirmation, and deletes them when confirmed.
//
// It returns the number of deleted files.
func ClearLastRunXMP(in io.Reader, out io.Writer) (int, error) {
	cr3Path, err := loadLastRunCR3Path()
	if err != nil {
		return 0, err
	}

	paths, err := os.ReadDir(cr3Path)
	if err != nil {
		return 0, fmt.Errorf("read CR3 path: %w", err)
	}

	xmpFiles := make([]string, 0)
	for _, path := range paths {
		if !path.IsDir() && strings.HasSuffix(strings.ToLower(path.Name()), ".xmp") {
			xmpFiles = append(xmpFiles, filepath.Join(cr3Path, path.Name()))
		}
	}

	fmt.Fprintf(out, "Found %d .xmp file(s) in %s\n", len(xmpFiles), cr3Path)
	if len(xmpFiles) == 0 {
		return 0, nil
	}

	fmt.Fprintf(out, "Delete all %d file(s)? [y/N]: ", len(xmpFiles))
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("read confirmation: %w", err)
		}
		fmt.Fprintln(out, "Clear cancelled.")
		return 0, nil
	}

	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Fprintln(out, "Clear cancelled.")
		return 0, nil
	}

	deleted := 0
	for _, xmpFile := range xmpFiles {
		if err := os.Remove(xmpFile); err != nil {
			return deleted, fmt.Errorf("delete %s: %w", xmpFile, err)
		}
		deleted++
	}

	return deleted, nil
}

func loadLastRunCR3Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}

	lastRunPathRaw, err := os.ReadFile(lastRunPathFile(home))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no previous run found; execute cr3-keyword with a CR3 path first")
		}
		return "", fmt.Errorf("read state file: %w", err)
	}

	lastRunPath := strings.TrimSpace(string(lastRunPathRaw))
	if lastRunPath == "" {
		return "", fmt.Errorf("last run CR3 path is empty")
	}
	if st, err := os.Stat(lastRunPath); err != nil || !st.IsDir() {
		if err != nil {
			return "", fmt.Errorf("last run CR3 path is not accessible: %w", err)
		}
		return "", fmt.Errorf("last run CR3 path is not a directory: %s", lastRunPath)
	}

	return lastRunPath, nil
}

func lastRunPathFile(home string) string {
	return filepath.Join(home, ".cr3-keywords", "last-run-cr3-path.txt")
}
