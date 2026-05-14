package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ClearMatchedFiles removes generated temporary files (.jpg/.txt) and sidecar
// .xmp files for CR3 files selected by CLI input.
//
// Selection rules:
//   - If files is empty, all .CR3 files in cr3Path are targeted.
//   - If files is non-empty, only those CR3 files are targeted.
//
// It returns the number of deleted files.
func ClearMatchedFiles(cr3Path string, files []string) (int, error) {
	resolvedPath := cr3Path
	if strings.TrimSpace(resolvedPath) == "" {
		var err error
		resolvedPath, err = loadLastRunCR3Path()
		if err != nil {
			return 0, err
		}
	}

	if err := ensureDir(resolvedPath); err != nil {
		return 0, fmt.Errorf("CR3 path not found: %w", err)
	}

	cr3Path = resolvedPath

	bases, err := resolveTargetCR3Bases(cr3Path, files)
	if err != nil {
		return 0, err
	}

	folderName := filepath.Base(filepath.Clean(cr3Path))
	jpgDir := filepath.Join(os.TempDir(), "cr3-keywords", "jpgs", folderName)
	outputDir := filepath.Join(os.TempDir(), "cr3-keywords", "outputs", folderName)

	deleted := 0
	for _, base := range bases {
		for _, target := range []string{
			filepath.Join(cr3Path, base+".xmp"),
			filepath.Join(jpgDir, base+".jpg"),
			filepath.Join(outputDir, base+".txt"),
		} {
			if err := os.Remove(target); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return deleted, fmt.Errorf("delete %s: %w", target, err)
			}
			deleted++
		}
	}

	return deleted, nil
}

func resolveTargetCR3Bases(cr3Path string, files []string) ([]string, error) {
	baseSet := map[string]struct{}{}

	if len(files) == 0 {
		entries, err := os.ReadDir(cr3Path)
		if err != nil {
			return nil, fmt.Errorf("read CR3 path: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(strings.ToLower(name), ".cr3") {
				base := strings.TrimSuffix(name, filepath.Ext(name))
				baseSet[base] = struct{}{}
			}
		}
	} else {
		for _, item := range files {
			candidate := item
			if !strings.Contains(item, string(os.PathSeparator)) && !filepath.IsAbs(item) {
				candidate = filepath.Join(cr3Path, item)
			}

			if !strings.HasSuffix(strings.ToLower(candidate), ".cr3") {
				return nil, fmt.Errorf("file is not a CR3: %s", item)
			}
			if _, err := os.Stat(candidate); err != nil {
				return nil, fmt.Errorf("CR3 file not found: %s", candidate)
			}

			base := strings.TrimSuffix(filepath.Base(candidate), filepath.Ext(candidate))
			baseSet[base] = struct{}{}
		}
	}

	if len(baseSet) == 0 {
		return nil, fmt.Errorf("no CR3 files found to clear")
	}

	bases := make([]string, 0, len(baseSet))
	for base := range baseSet {
		bases = append(bases, base)
	}
	sort.Strings(bases)
	return bases, nil
}

func ensureDir(p string) error {
	st, err := os.Stat(p)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return fmt.Errorf("not a directory: %s", p)
	}
	return nil
}
