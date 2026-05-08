package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options contains parsed CLI flags, command mode, and positional arguments.
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
	// Clear runs the clear command to remove last-run XMP files.
	Clear bool
	// UseExif forces exiftool-based CR3 preview extraction.
	UseExif bool

	// Model is the selected model name. Empty means auto-detect.
	Model string
	// PromptPath is the path to the prompt markdown file.
	PromptPath string
	// CR3Path is the CR3 source directory path.
	CR3Path string
	// Files is an optional list of specific CR3 files to process.
	Files []string
}

const defaultModel = ""

// Parse converts CLI arguments into Options and validates supported command
// forms and flag combinations.
func Parse(args []string) (Options, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Options{}, fmt.Errorf("resolve home directory: %w", err)
	}

	defaultPromptPath := filepath.Join(home, ".cr3-keywords", "prompt.md")
	opts := Options{
		Model:      defaultModel,
		PromptPath: defaultPromptPath,
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

	if len(rest) == 1 && rest[0] == "version" {
		opts.Version = true
		return opts, nil
	}

	if len(rest) < 1 {
		return Options{}, errors.New("missing required arguments\n\n" + Usage())
	}

	if len(rest) == 1 && rest[0] == "clear" {
		opts.Clear = true
		return opts, nil
	}

	first := expandPath(rest[0], home)
	if isDir(first) {
		opts.CR3Path = first
		if len(rest) > 1 {
			opts.Files = rest[1:]
		}
		return opts, nil
	}

	if opts.Model != "" {
		return Options{}, errors.New("do not mix --model with positional <model> argument")
	}
	if opts.PromptPath != defaultPromptPath {
		return Options{}, errors.New("do not mix --prompt with positional <prompt_file> argument")
	}
	if len(rest) < 3 {
		return Options{}, errors.New("invalid arguments\n\n" + Usage())
	}

	opts.Model = rest[0]
	opts.PromptPath = expandPath(rest[1], home)
	opts.CR3Path = expandPath(rest[2], home)
	if len(rest) > 3 {
		opts.Files = rest[3:]
	}

	return opts, nil
}

// Usage returns the CLI help text.
func Usage() string {
	return `Usage:
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--exif] <cr3_path>
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--exif] <cr3_path> IMG_0150.CR3
  cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--exif] <cr3_path> IMG_0150.CR3 IMG_0151.CR3
  cr3 [--verbose] [--dry-run] <model> <prompt_file> <cr3_path> IMG_0150.CR3 IMG_0151.CR3
  cr3 clear
  cr3 version
  cr3 --version

Flags:
  -v, --verbose       Print detailed per-file logs
  -n, --dry-run       Simulate actions without writing files or sending requests
  -m, --model         Optional model name (if omitted, auto-detected)
  -p, --prompt        Prompt file path (default: ~/.cr3-keywords/prompt.md)
  -e, --edit-prompt   Open prompt file in terminal editor before running
      --exif          Force exiftool for CR3 preview extraction (faster, requires exiftool)
  -h, --help          Show this help
      --version       Show version and exit

Commands:
  clear               Delete all .xmp files from the CR3 path of the last successful run
  version             Show version and exit
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

func isDir(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	return st.IsDir()
}
