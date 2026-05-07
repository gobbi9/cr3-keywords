# Local Lightroom AI Keywording Pipeline

## Overview

`cr3-keyword.sh` converts Canon `.CR3` files to JPEG previews, sends them to your local LM Studio model for keyword/caption generation, and writes `.xmp` metadata files for Lightroom.

---

## Requirements

Install dependencies:

```/dev/null/install.sh#L1-1
brew install imagemagick exiftool jq
```

Start LM Studio server:

```/dev/null/lms.sh#L1-1
lms start server
```

---

## Setup

From the project directory:

```/dev/null/setup.sh#L1-2
chmod +x ./cr3-keyword.sh
ls -l ./cr3-keyword.sh
```

Create or edit your prompt file (default: `prompt.md`).

---

## Usage

### Default model + default prompt

```/dev/null/usage.sh#L1-1
./cr3-keyword.sh <cr3_path>
```

### Process selected files only

```/dev/null/usage.sh#L3-3
./cr3-keyword.sh <cr3_path> IMG_0150.CR3 IMG_0151.CR3
```

### Custom model + custom prompt

```/dev/null/usage.sh#L5-5
./cr3-keyword.sh <model> <prompt_file> <cr3_path> IMG_0150.CR3 IMG_0151.CR3
```

### Flags

```/dev/null/usage.sh#L7-10
./cr3-keyword.sh --verbose <cr3_path> IMG_0150.CR3
./cr3-keyword.sh --dry-run <cr3_path> IMG_0150.CR3
./cr3-keyword.sh --verbose --dry-run <cr3_path> IMG_0150.CR3
./cr3-keyword.sh --help
```

---

## Output locations

For input path:

`~/Pictures/raw/braunschweig-20260503`

The script uses:

- JPEGs: `./jpgs/braunschweig-20260503`
- TXT outputs: `./outputs/braunschweig-20260503`
- TMP files: `./tmp/braunschweig-20260503`
- XMP files: `~/Pictures/raw/braunschweig-20260503`

---

## Runtime behavior

- Steps run sequentially (no parallel processing)
- One headline is printed at start of each step
- Progress bar is shown during each step
- Per-file logs are shown only with `--verbose`
- `--dry-run` performs no file writes and sends no API requests

---

## Lightroom import

1. Select files (or press `Cmd + A`)
2. Right click → **Metadata** → **Read Metadata from File(s)**
3. Confirm import

You can delete `.xmp` files afterwards if desired.
