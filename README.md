# Local Lightroom AI Keywording Pipeline (Go CLI)

## Overview

This project now includes a Go CLI (`cr3`) that replaces the original shell pipeline while preserving behavior and flags.

- `cr3-keyword.sh` is still kept in this repository as reference.
- Pipeline remains: **CR3 -> JPG -> TXT -> XMP**.
- Progress bar and colored step output are preserved.
- Logging uses Go structured logging (`log/slog`) instead of `println`.
- Intermediate directories (`jpgs`, `outputs`, `tmp`) are written under macOS temp folder:
  - `/tmp/cr3-keywords/jpgs/<folder>`
  - `/tmp/cr3-keywords/outputs/<folder>`
  - `/tmp/cr3-keywords/tmp/<folder>`

---

## Requirements

- macOS
- [goenv](https://github.com/go-nv/goenv)
- Go `1.22.5` (see `.go-version`)
- `exiftool` (used to extract CR3 preview image)
- LM Studio running local API server on `http://localhost:1234`

Install tooling:

```bash
brew install goenv exiftool

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

Start LM Studio server:

```bash
lms start server
```

---

## Prompt file behavior

- Default prompt file is `~/.cr3-keywords/prompt.md`.
- You can override with `--prompt`.
- You can edit prompt in terminal before processing with `--edit-prompt`.
  - Uses `$EDITOR`, defaults to `nano`.

---

## Model selection

Model is optional.

If not provided, the CLI auto-detects available models in this order:
1. LM Studio HTTP API (`GET /v1/models`)
2. `lms ls --json`

Then it picks the best model using a simple vision-priority scoring heuristic.

You can force a specific model with `--model` or positional `<model>`.

---

## Usage

### Default auto model + default prompt

```bash
cr3 <cr3_path>
```

### Process selected files only

```bash
cr3 <cr3_path> IMG_0150.CR3 IMG_0151.CR3
```

### Custom model + prompt (flags)

```bash
cr3 --model qwen2.5-vl --prompt ~/.cr3-keywords/prompt.md <cr3_path> IMG_0150.CR3
```

### Legacy positional model + prompt mode

```bash
cr3 <model> <prompt_file> <cr3_path> IMG_0150.CR3 IMG_0151.CR3
```

### Flags

```bash
cr3 --verbose <cr3_path> IMG_0150.CR3
cr3 --dry-run <cr3_path>
cr3 --edit-prompt <cr3_path>
cr3 --model qwen2.5-vl <cr3_path>
cr3 --prompt ~/.cr3-keywords/prompt.md <cr3_path>
cr3 --help
```

### Show version

```bash
cr3 --version
cr3 version
```

Version output includes app version + short git commit hash + build timestamp.

### Clear generated XMP files from last run

```bash
cr3 clear
```

- Uses the CR3 folder from the last successful run.
- Finds all `.xmp` files in that folder.
- Asks for confirmation before deleting.
- Prints the number of deleted files.

---

## Notes on external tools/libraries

- Replaced ImageMagick with Go image processing (`github.com/disintegration/imaging`).
- Replaced `jq` usage with native Go JSON handling.
- Replaced most shell logic with native Go implementations.
- `exiftool` is still used specifically for CR3 preview extraction (practical fallback for CR3 support).

---

## Lightroom import

1. Select files in Lightroom (`Cmd + A` for all).
2. Right click -> **Metadata** -> **Read Metadata from File(s)**.
3. Confirm import.

XMP files are written next to your CR3 files. After the import, you should delete the XMP files, using the `clear` command.
