# Local Lightroom AI Keywording Pipeline (Go CLI)

## Overview

`cr3-keywords` is a privacy-first, local AI keywording and captioning pipeline for Canon `.cr3` photos.

The main goal of this project is to create keywords and captions from `.cr3` files using a local LLM with LM Studio to preserve privacy.

- Built around a Go CLI (`cr3`) for a fast, reproducible workflow.
- Pipeline is: **CR3 -> JPG -> TXT -> XMP** (keywords + caption + optional geotagging).
- Progress bar and colored step output are preserved.
- Logging uses Go structured logging (`log/slog`) instead of `println`.
- Intermediate directories (`jpgs`, `outputs`, `tmp`) are written under the macOS temp folder:
  - `$TEMPDIR/cr3-keywords/jpgs/<folder>`
  - `$TEMPDIR/cr3-keywords/outputs/<folder>`
  - `$TEMPDIR/cr3-keywords/tmp/<folder>`

---

## Requirements

- LM Studio running local API server on `http://localhost:1234`

Start LM Studio server:

```bash
lms start server
```

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

## Prompt file behavior

- Default prompt file is `~/.cr3-keywords/prompt.md`.
- You can override with `--prompt`.
- You can edit prompt in terminal before processing with `--edit-prompt`.
  - Uses `$EDITOR`, defaults to `nano`.
- When geotagging is available, the CLI prepends location context (`city` and `sublocation`) to the beginning of the prompt **in-memory only** before sending to the LLM.
- Geotagged XMP sidecars also include IPTC location fields used by Lightroom Classic: sublocation, city, state/province, country/region, and ISO country code.
  - The prompt file on disk (default or passed with `--prompt`) is never modified.

---

## Model selection

Model is optional.

If not provided, the CLI auto-detects available models in this order:
1. LM Studio HTTP API (`GET /v1/models`)
2. `lms ls --json`

Then it picks the best model using a simple vision-priority scoring heuristic.

You can force a specific model with `--model`.

---

## Usage

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

### Flags

```bash
cr3 --verbose <cr3_path> IMG_0150.CR3
cr3 --dry-run <cr3_path>
cr3 --edit-prompt <cr3_path>
cr3 --model qwen2.5-vl <cr3_path>
cr3 --prompt ~/.cr3-keywords/prompt.md <cr3_path>
cr3 --gps <cr3_path>
cr3 --gps ~/.cr3-keywords/track.gpx <cr3_path>
cr3 --exif <cr3_path>
cr3 --clear
cr3 --clear <cr3_path>
cr3 --clear <cr3_path> IMG_0150.CR3
cr3 --help
```

`--gps` enables GPX geotagging.

- You can pass it with no value (`--gps`) to force usage of the default track file: `~/.cr3-keywords/track.gpx`.
- You can pass a custom file (`--gps /path/to/file.gpx`) to use that track.
- If `--gps` is present (with or without path) and the resolved file does not exist, the CLI exits with an error.
- If `--gps` is not provided and the default file is missing, the CLI does **not** fail and simply skips geotagging.

`--exif` forces CR3 preview extraction via `exiftool` (fast path when available).

`--clear` deletes generated temporary files for matching CR3 files:

- Temp JPG: `$TEMPDIR/cr3-keywords/jpgs/<folder>/<base>.jpg`
- Temp TXT: `$TEMPDIR/cr3-keywords/outputs/<folder>/<base>.txt`
- Sidecar XMP: `<cr3_path>/<base>.xmp`

Matching behavior:

- `cr3 --clear <cr3_path> IMG_0150.CR3` clears only files for that CR3 base name.
- `cr3 --clear <cr3_path>` clears files for all `.CR3` files found in that folder.
- `cr3 --clear` (no path) uses the last successful run path from `~/.cr3-keywords/last-run-cr3-path.txt`.

### Show version

```bash
cr3 --version
```

Version output includes app version + short git commit hash + build timestamp.

### Clear generated temporary files and XMP sidecars

```bash
cr3 --clear
cr3 --clear <cr3_path>
cr3 --clear <cr3_path> IMG_0150.CR3
```

- Clears only files matching `.CR3` names from the provided folder.
- Deletes corresponding temp `.jpg`, temp `.txt`, and sidecar `.xmp` files.
- With no filenames, clears for all `.CR3` files in that folder.
- `cr3 --clear` (no path) uses the last successful run folder.

---

## Lightroom import

1. Select files in Lightroom (`Cmd + A` for all).
2. Right click -> **Metadata** -> **Read Metadata from File(s)**.
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
- [goenv](https://github.com/go-nv/goenv)
- Go `1.22.5` (see `.go-version`)
- LM Studio running local API server on `http://localhost:1234`
- Optional: `exiftool` (fallback when built-in CR3 preview extraction fails, or required when using `--exif`)

Install tooling:

```bash
brew install goenv
# optional (fallback / --exif mode):
# brew install exiftool

goenv install
goenv local
# only needed for IDEs
goenv global 1.22.5
```

Build using the Makefile (uses `CGO_ENABLED=0` by default):

```bash
make tidy
make build
```

Install globally (default path `/usr/local/bin`):

```bash
make install
# use sudo if needed for permissions
# sudo make install
```
