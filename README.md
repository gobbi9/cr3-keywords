# Local Lightroom AI Keywording Pipeline (Go CLI)

## Overview

`cr3-keywords` is a privacy-first, local AI keywording and captioning pipeline for Canon `.cr3` photos.

This project generates keywords and captions from `.cr3` files using a local LLM with [LM Studio](https://lmstudio.ai) to preserve privacy.

- Built around a Go CLI (`cr3`) for fast, reproducible workflows.
- Pipeline is: **CR3 -> JPG -> TXT -> XMP** (keywords + caption + optional geotagging).
- The CLI shows a progress bar and colorized step output. It also shows a timer for the JPG -> TXT step, since this is the most time-consuming part (i.e., prompting the LLM).
- Intermediate directories (`jpgs`, `outputs`, `tmp`) are written under the macOS temp folder:
  - `$TEMPDIR/cr3-keywords/jpgs/<folder>`
  - `$TEMPDIR/cr3-keywords/outputs/<folder>`
  - `$TEMPDIR/cr3-keywords/tmp/<folder>`

---

## Requirements

- LM Studio local API server on `http://localhost:1234`
- Optional but recommended: LM Studio `lms` CLI available in `PATH`

You can start LM Studio server manually:

```bash
lms server start
```

`cr3` first checks whether `http://localhost:1234` responds with HTTP `200 OK`.
If not, it tries `lms server start` automatically.
If `lms` is not available and the server is not already running, startup fails.

`cr3` extracts embedded CR3 JPEG previews using a built-in pure-Go path.
`exiftool` is optional and used as a fallback if built-in extraction fails for a file.
You can force `exiftool` mode with `--exif` (often faster, requires `exiftool` in `PATH`).

## Install via package managers

### Homebrew (macOS/Linux)

```bash
brew tap gobbi9/tap https://github.com/gobbi9/tap
brew install cr3-keywords
```

### Scoop (Windows)

```bash
scoop bucket add gobbi9 https://github.com/gobbi9/scoop-bucket
scoop install gobbi9/cr3-keywords
```

## Install prebuilt binaries (GitHub Releases)

If you don't want to build from source, download a release binary from:

- https://github.com/gobbi9/cr3-keywords/releases/latest

Assets are published as:

- `cr3_<version>_darwin-arm64` (macOS Apple Silicon)
- `cr3_<version>_darwin-amd64` (macOS Intel)
- `cr3_<version>_linux-amd64` (Linux x86_64)
- `cr3_<version>_windows-amd64.exe` (Windows x86_64)

Quick install examples:

```bash
# macOS / Linux
# Pick one ASSET value:
#   darwin-arm64   (macOS Apple Silicon)
#   darwin-amd64   (macOS Intel)
#   linux-amd64    (Linux x86_64)
VERSION=0.2.0
ASSET=darwin-arm64
curl -L -o cr3 "https://github.com/gobbi9/cr3-keywords/releases/download/v${VERSION}/cr3_${VERSION}_${ASSET}"
chmod +x cr3
sudo mv cr3 /usr/local/bin/cr3
```

```bash
# Windows (PowerShell)
$version = "0.2.0"
Invoke-WebRequest -Uri "https://github.com/gobbi9/cr3-keywords/releases/download/v$version/cr3_${version}_windows-amd64.exe" -OutFile "cr3.exe"
# Move cr3.exe to a folder in PATH, for example:
# Move-Item .\cr3.exe "$env:USERPROFILE\bin\cr3.exe"
```

---

## Usage

### Flags

```bash
cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path>
cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path> IMG_0150.CR3
cr3 [--verbose] [--dry-run] [--model MODEL] [--prompt ~/.cr3-keywords/prompt.md] [--edit-prompt] [--gps [~/.cr3-keywords/track.gpx]] [--exif] <cr3_path> IMG_0150.CR3 IMG_0151.CR3
cr3 --clear
cr3 --clear <cr3_path>
cr3 --clear <cr3_path> IMG_0150.CR3
cr3 install nushell|zsh|bash
cr3 completion nushell|zsh|bash
cr3 --version

Commands:
install TARGET    Install shell integration
completion TARGET Print shell completion/module script to stdout
                  TARGET: nushell|zsh|bash

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
```

`--gps` enables GPX geotagging.

- You can pass it with no value (`--gps`) to force usage of the default track file: `~/.cr3-keywords/track.gpx`.
- You can pass a custom file (`--gps /path/to/file.gpx`) to use that track.
- If `--gps` is present (with or without path) and the resolved file does not exist, the CLI exits with an error.
- If `--gps` is not provided and the default file is missing, the CLI does **not** fail and simply skips geotagging.

`--exif` forces CR3 preview extraction via `exiftool` (fast path when available).

`--clear` deletes generated temporary files for matching CR3 files. More information below.

Version output includes the app version and build timestamp.

### Default auto model + default prompt (+ optional default GPX)

```bash
cr3 <cr3_path>
```

By default, the CLI looks for GPX track file at:

- `~/.cr3-keywords/track.gpx`

If that file exists, geotagging data is added to generated `.xmp` files (including Lightroom-compatible IPTC location fields) and location context is prepended to the LLM prompt in-memory.
If it does not exist, processing continues without geotagging data.

### Process selected files only

```bash
cr3 <cr3_path> IMG_0150.CR3 IMG_0151.CR3
```

### Custom model + prompt + GPX (flags)

```bash
cr3 --model qwen2.5-vl --prompt ~/.cr3-keywords/prompt.md --gps <cr3_path> IMG_0150.CR3
cr3 --model qwen2.5-vl --prompt ~/.cr3-keywords/prompt.md --gps ~/.cr3-keywords/track.gpx <cr3_path> IMG_0150.CR3
```


### Shell completion scripts

Print script/module to stdout:

```bash
cr3 completion nushell
cr3 completion zsh
cr3 completion bash
```

Install to per-user defaults:

```bash
cr3 install nushell
cr3 install zsh
cr3 install bash
```

Install locations:

- Nushell
  - macOS: `~/Library/Application Support/nushell/vendor/autoload/cr3.nu`
  - Linux: `~/.config/nushell/vendor/autoload/cr3.nu` (or `$XDG_CONFIG_HOME/nushell/vendor/autoload/cr3.nu`)
  - Windows: `%APPDATA%\nushell\vendor\autoload\cr3.nu`
- Zsh: `$ZDOTDIR/completions/_cr3` (fallback: `~/.zsh/completions/_cr3`)
- Bash:
  - Linux/macOS: `$XDG_DATA_HOME/bash-completion/completions/cr3` (fallback: `~/.local/share/bash-completion/completions/cr3`)
  - Windows: `%LOCALAPPDATA%\bash-completion\completions\cr3`

After install, restart your shell (or source the generated file manually).

Nushell-specific behavior for `cr3` completion:

- first positional argument suggests directories only (dot-directories are excluded)
- subsequent positional arguments suggest `.cr3` files from the selected first directory (case-insensitive extension matching)

### Clear generated temporary files and XMP sidecars

```bash
cr3 --clear
cr3 --clear <cr3_path>
cr3 --clear <cr3_path> IMG_0150.CR3
```

#### Matching behavior

- `cr3 --clear <cr3_path> IMG_0150.CR3` clears only files for that CR3 base name.
- `cr3 --clear <cr3_path>` clears temp files for all `.CR3` files found in that folder.
- `cr3 --clear` (no path) uses the last successful run path from `~/.cr3-keywords/last-run-cr3-path.txt`.

Deletes

- Temp JPG: `$TEMPDIR/cr3-keywords/jpgs/<folder>/<base>.jpg`
- Temp TXT: `$TEMPDIR/cr3-keywords/outputs/<folder>/<base>.txt`
- Sidecar XMP: `<cr3_path>/<base>.xmp`

> **It does NOT delete or modify any .CR3 files.**

---

## Prompt file behavior

- Default prompt file is `~/.cr3-keywords/prompt.md`.
- You can override with `--prompt`.
- You can edit the prompt in a terminal before processing with `--edit-prompt`; it uses `$EDITOR`.
- When geotagging is available, the CLI prepends location context (`city` and `sublocation`) to the beginning of the prompt **in-memory only** before sending to the LLM.
- Geotagged XMP sidecars also include IPTC location fields used by Lightroom Classic: sublocation, city, state/province, country/region, and ISO country code.
  - The prompt file on disk (default or passed with `--prompt`) is never modified.

If no prompt is found, the user is asked to edit the default prompt:

```text
[Location data from gpx track is added dynamically here to improve the model's understanding of the photo context]

Context:
This photo is part of a series of photos.

Task:
1. Describe the image
2. Generate 10–20 simple keywords (comma-separated)
3. Write a short caption (1 sentence), using the image description from (task 1).

Avoid generic terms like "image" or "photo".

Output should be result of task 2, empty line, result of task 3.
Make sure output does not contain the full description,
only comma separated keywords in the first line, an empty line and the caption.
```

> It is necessary to keep the output format as described above to ensure the CLI can parse the results correctly. But you are free to edit the prompt to suit your needs, for instance more or less keywords, longer captions.

In practice, this CLI uses a single prompt for all images. If you can invest time writing a custom prompt for each image, you could likely create keywords and captions manually as well.

It is currently not possible for the LLM to share context across multiple images or process multiple images in a single request. As a trade-off, you may see many similar keywords repeated across images.

---

## Model selection

Model is optional.

If not provided, the CLI auto-detects available models in this order:
1. LM Studio native API (`GET /api/v1/models`) using `key` + `capabilities.vision`
2. `lms ls --json` using `modelKey` + `vision`

Only vision-capable models are considered. Then it picks the best model using a simple vision-priority scoring heuristic.

You can force a specific model with `--model`.

---

## Adobe Lightroom Classic import

It is recommended to do this before editing your CR3 files in Lightroom. This should only override keywords, captions, and geotagging metadata (when applicable). However, Lightroom metadata handling can be inconsistent, so you might lose data if your CR3 files have already been edited in Lightroom.

1. Select files in Lightroom (`Cmd + A` for all).
2. Right-click -> **Metadata** -> **Read Metadata from File(s)**.
3. Confirm import.

XMP files are written next to your CR3 files. After the import, you can delete generated sidecars/temp files with `--clear`.

Location field mapping in generated XMP:

- Lightroom **Sublocation** -> `Iptc4xmpCore:Location`
- Lightroom **City** -> `photoshop:City`
- Lightroom **State / Province** -> `photoshop:State`
- Lightroom **Country / Region** -> `photoshop:Country`
- Lightroom **ISO Country Code** -> `Iptc4xmpCore:CountryCode`

---

## Development

- macOS/Linux
- [mise](https://mise.jdx.dev) (recommended) or [goenv](https://github.com/go-nv/goenv)
- Go `1.22.5` (see `.go-version`)
- [LM Studio](https://lmstudio.ai) local API server on `http://localhost:1234`
- Optional but recommended: LM Studio `lms` CLI in `PATH` (used for auto-start and CLI model fallback)
- Optional: `exiftool` (fallback when built-in CR3 preview extraction fails, or required when using `--exif`, faster than the built-in CR3 preview extraction)

Install tooling:

```bash
# Option A (recommended): mise
brew install mise
brew install exiftool # optional (fallback / --exif mode)
mise settings add idiomatic_version_file_enable_tools go
mise install

# Option B: goenv
brew install goenv
brew install exiftool # optional (fallback / --exif mode)
goenv install
goenv local
```

Build using the Makefile (uses `CGO_ENABLED=0` by default):

```bash
make tidy
make build
```

Install globally (default path `/usr/local/bin`):

```bash
sudo make install
```

If you wish to contribute, see [MAINTAINERS.md](MAINTAINERS.md).

## AI Guidelines

- For skills and project memory, see: [.agents/README.md](.agents/README.md)
