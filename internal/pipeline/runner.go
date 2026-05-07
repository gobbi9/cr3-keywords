package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"cr3-keywords/internal/cli"
	"cr3-keywords/internal/lm"

	"github.com/fatih/color"
)

type Runner struct {
	opts   cli.Options
	logger *slog.Logger
}

func NewRunner(opts cli.Options, logger *slog.Logger) *Runner {
	return &Runner{opts: opts, logger: logger}
}

func (r *Runner) Run(ctx context.Context) error {
	if err := ensureDir(r.opts.CR3Path); err != nil {
		return fmt.Errorf("CR3 path not found: %w", err)
	}
	if err := ensureFile(r.opts.PromptPath); err != nil {
		return fmt.Errorf("prompt file not found: %w", err)
	}

	lmClient := lm.NewClient()
	model := r.opts.Model
	if strings.TrimSpace(model) == "" {
		auto, err := lmClient.DetectBestModel(ctx)
		if err != nil {
			return err
		}
		model = auto
		r.logger.Debug("auto-detected model", "model", model)

		fmt.Print("Auto-detected model: ")
		color.New(color.FgHiYellow, color.Bold).Printf("%s\n", model)
	}

	folderName := folderNameFromPath(r.opts.CR3Path)
	jpgDir := filepath.Join(os.TempDir(), "cr3-keywords", "jpgs", folderName)
	outputDir := filepath.Join(os.TempDir(), "cr3-keywords", "outputs", folderName)
	tmpDir := filepath.Join(os.TempDir(), "cr3-keywords", "tmp", folderName)

	if !r.opts.DryRun {
		for _, d := range []string{jpgDir, outputDir, tmpDir} {
			if err := os.MkdirAll(d, 0o755); err != nil {
				return fmt.Errorf("create dir %s: %w", d, err)
			}
		}
	}

	cr3Files, err := r.resolveCR3Files()
	if err != nil {
		return err
	}
	if len(cr3Files) == 0 {
		return fmt.Errorf("no CR3 files found to process")
	}

	selectedJPGs := make([]string, 0, len(cr3Files))
	selectedTXTs := make([]string, 0, len(cr3Files))
	for _, cr3 := range cr3Files {
		base := strings.TrimSuffix(filepath.Base(cr3), filepath.Ext(cr3))
		selectedJPGs = append(selectedJPGs, filepath.Join(jpgDir, base+".jpg"))
		selectedTXTs = append(selectedTXTs, filepath.Join(outputDir, base+".txt"))
	}

	progress := NewProgressBar()
	stepColor := color.New(color.FgHiCyan, color.Bold)

	stepColor.Println("=== Step 1: CR3 → JPG ===")
	if err := cr3ToJPGs(ctx, r.logger, progress, jpgDir, cr3Files, r.opts.DryRun, r.opts.UseExif); err != nil {
		return err
	}

	stepColor.Println("=== Step 2: JPG → TXT ===")
	if err := batchCaption(ctx, r.logger, progress, lmClient, jpgDir, outputDir, model, r.opts.PromptPath, selectedJPGs, r.opts.DryRun); err != nil {
		return err
	}

	stepColor.Println("=== Step 3: TXT → XMP ===")
	if err := txtToXMP(r.logger, progress, outputDir, r.opts.CR3Path, selectedTXTs, r.opts.DryRun); err != nil {
		return err
	}

	color.New(color.FgHiGreen, color.Bold).Println("=== Done ===")
	return nil
}

func (r *Runner) resolveCR3Files() ([]string, error) {
	if len(r.opts.Files) == 0 {
		entries, err := os.ReadDir(r.opts.CR3Path)
		if err != nil {
			return nil, fmt.Errorf("read CR3 path: %w", err)
		}
		files := make([]string, 0)
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			lower := strings.ToLower(name)
			if strings.HasSuffix(lower, ".cr3") {
				files = append(files, filepath.Join(r.opts.CR3Path, name))
			}
		}
		return files, nil
	}

	files := make([]string, 0, len(r.opts.Files))
	for _, item := range r.opts.Files {
		candidate := item
		if !strings.Contains(item, string(os.PathSeparator)) && !filepath.IsAbs(item) {
			candidate = filepath.Join(r.opts.CR3Path, item)
		}
		if _, err := os.Stat(candidate); err == nil {
			files = append(files, candidate)
		} else {
			r.logger.Debug("skipping missing CR3", "file", candidate)
		}
	}
	return files, nil
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

func ensureFile(p string) error {
	st, err := os.Stat(p)
	if err != nil {
		return err
	}
	if st.IsDir() {
		return fmt.Errorf("expected file but got directory: %s", p)
	}
	return nil
}
