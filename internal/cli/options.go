package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options contains parsed CLI flags, command mode, and optional file filters.
type Options struct {
	// Verbose enables debug-level logging.
	Verbose bool
	// DryRun simulates pipeline actions without writing files or calling the model.
	DryRun bool
	// Help requests usage output.
	Help bool
	// Version requests version output.
	Version bool
	// EditPrompt opens the prompt file in an editor before execution.
	EditPrompt bool
	// Clear removes generated .jpg/.txt and sidecar .xmp files for matching CR3 files.
	Clear bool
	// UseExif forces exiftool-based CR3 preview extraction.
	UseExif bool
	// GPSPath is the path to the GPX track file.
	GPSPath string
	// GPSProvided tracks whether --gps was explicitly provided by user.
	GPSProvided bool

	// Model is the selected model name. Empty means auto-detect.
	Model string
	// PromptPath is the path to the prompt markdown file.
	PromptPath string
	// CR3Path is the CR3 source directory path.
	CR3Path string
	// Files is an optional list of specific CR3 files to process.
	Files []string
}

// Parse converts CLI arguments into Options and validates supported command
// forms and flag combinations.
func Parse(args []string) (Options, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Options{}, fmt.Errorf("resolve home directory: %w", err)
	}

	defaultPromptPath := filepath.Join(home, ".cr3-keywords", "prompt.md")
	defaultGPSPath := filepath.Join(home, ".cr3-keywords", "track.gpx")
	opts := Options{
		PromptPath: defaultPromptPath,
		GPSPath:    defaultGPSPath,
	}

	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-v", "--verbose":
			opts.Verbose = true
		case "-n", "--dry-run":
			opts.DryRun = true
		case "-h", "--help":
			opts.Help = true
		case "--version":
			opts.Version = true
		case "--clear":
			opts.Clear = true
		case "-e", "--edit-prompt":
			opts.EditPrompt = true
		case "--exif":
			opts.UseExif = true
		case "-m", "--model":
			i++
			if i >= len(args) {
				return Options{}, errors.New("missing value for --model")
			}
			opts.Model = args[i]
		case "-p", "--prompt":
			i++
			if i >= len(args) {
				return Options{}, errors.New("missing value for --prompt")
			}
			opts.PromptPath = expandPath(args[i], home)
		case "--gps":
			opts.GPSProvided = true
			if i+1 < len(args) {
				next := args[i+1]
				if !strings.HasPrefix(next, "-") {
					i++
					opts.GPSPath = expandPath(next, home)
				}
			}
		case "--":
			rest = append(rest, args[i+1:]...)
			i = len(args)
		default:
			if strings.HasPrefix(a, "-") {
				return Options{}, fmt.Errorf("unknown flag: %s", a)
			}
			rest = append(rest, a)
		}
	}

	if opts.Help || opts.Version {
		return opts, nil
	}

	if opts.Clear {
		if len(rest) > 0 {
			opts.CR3Path = expandPath(rest[0], home)
		}
		if len(rest) > 1 {
			opts.Files = rest[1:]
		}
		return opts, nil
	}

	if len(rest) < 1 {
		return Options{}, errors.New("missing required arguments\n\n" + Usage())
	}

	opts.CR3Path = expandPath(rest[0], home)
	if len(rest) > 1 {
		opts.Files = rest[1:]
	}

	return opts, nil
}

// Usage returns the CLI help text.
func Usage() string {
	return `Usage:
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path>
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path> IMG_0150.CR3
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path> IMG_0150.CR3 IMG_0151.CR3
  cr3 --clear
  cr3 --clear <cr3_path>
  cr3 --clear <cr3_path> IMG_0150.CR3
  cr3 --version

Flags:
  -v, --verbose       Print detailed per-file logs
  -n, --dry-run       Simulate actions without writing files or sending requests
  -m, --model         Optional model name (if omitted, auto-detected)
  -p, --prompt        Prompt file path (default: ~/.cr3-keywords/prompt.md)
  -e, --edit-prompt   Open prompt file in terminal editor before running
      --gps [PATH]    Enable GPX geotagging. Optional path (default: ~/.cr3-keywords/track.gpx)
      --exif          Force exiftool for CR3 preview extraction (faster, requires exiftool)
      --clear         Delete generated .jpg/.txt and sidecar .xmp for matching CR3 files
  -h, --help          Show this help
      --version       Show version and exit

`
}

func expandPath(p, home string) string {
	if p == "~" {
		return home
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
}
