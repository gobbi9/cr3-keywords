# Project Memory

## Project Overview

- `cr3-keywords` is a privacy-first local photo metadata pipeline for Canon `.CR3` files.
- The project evolved from an initial zsh script into a Go CLI named `cr3`.
- Core workflow is `CR3 -> JPG preview -> LLM keywords/caption TXT -> XMP sidecar`.
- LM Studio local API (`http://localhost:1234`) is the primary inference backend.

## Architecture

- Entry point: `cmd/cr3-keyword/main.go`.
- CLI parsing and command-mode behavior: `internal/cli/*`.
- Pipeline orchestration: `internal/pipeline/runner.go`.
- Stage 1 preview extraction:
  - Pure-Go CR3 preview extraction is preferred.
  - `exiftool` is fallback and can be forced with `--exif`.
- Stage 2 keyword/caption generation:
  - LM Studio client in `internal/lm/client.go`.
  - Vision model auto-detection when `--model` is not provided.
  - Server readiness is checked once before batch generation.
- Stage 3 metadata writing:
  - XMP generation in `internal/pipeline/3_xmp.go`.
  - Lightroom-compatible location field mapping included.
- Temp workspace uses macOS temp dir via `os.TempDir()` under `cr3-keywords/{jpgs,outputs,tmp}/<folder>`.
- Packaging/release helper scripts moved from shell to Python (`scripts/*.py`).

## Key Constraints

- Designed for local-only processing (privacy and offline control).
- Pipeline output format for model response is strict: first line keywords, blank line, caption.
- Geotagging is optional but strict when explicitly requested (`--gps` must resolve to an existing file).
- Must preserve CR3 files; generated artifacts are sidecars and temp files only.

## Important APIs

- LM Studio endpoints:
  - `GET /api/v1/models` for model discovery.
  - Chat completions endpoint for caption/keyword generation.
- Optional LM Studio CLI fallback (`lms`) for server start/model lookup.
- Release/distribution automation targets GitHub Releases, Homebrew tap, and Scoop bucket.

## Conventions

- Go formatting conventions (`gofmt`/`gopls`) and documentation comments were standardized.
- CLI favors explicit flags over legacy positional argument modes.
- `--clear` is the cleanup interface (instead of subcommand style).
- README acts as source-of-truth for operational behavior and user docs.
- Development docs now prefer [`mise`](https://mise.jdx.dev) first, with `goenv` as an alternative.
- For `mise`, setup must include `mise settings add idiomatic_version_file_enable_tools go` before `mise install` so `.go-version` is honored idiomatically.
- End-session `.agents` workflow applies `memento`, `ai-janitor`, `tests`, and `docs` skills in that order.
- `docs` skill updates root `README.md` only when changes include at least one non-`.agents` file.
- `tests` skill updates/runs existing unit tests only; it must not create new test files.
- `failures.md` should contain only real failed attempts; if none occurred in a session, add no failure entry.
- Mermaid diagrams in Markdown should be GitHub-safe: use simple alphanumeric node IDs (e.g., `endSession` instead of reserved words like `end`), keep display text in brackets, and avoid Mermaid reserved keywords as node identifiers.

## Current Priorities

- Keep release/packaging automation reliable across GitHub Releases + Homebrew + Scoop.
- Maintain robust LM Studio client behavior and diagnostics.
- Keep `.agents` memory files synchronized with thread history and recent architectural decisions.
- Preserve incremental Zed-thread ingestion state in `.agents/skills/zed-threads/state/nu.cursor`.
- Keep end-session guardrails active for test maintenance and conditional README updates.

## Historical Thread Digest (seeded)

- Processed Zed thread records incrementally for this project from cursor `0` to `32` (exclusive end cursor now `32`).
- Recurring themes: CLI ergonomics, release automation, CR3 extraction performance, geotagging/XMP correctness, and documentation polish.
- Notable evolution arc: zsh script -> Go CLI hardening -> packaging/release system maturity -> `.agents` memory tooling.
