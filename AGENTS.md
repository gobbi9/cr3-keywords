# AGENTS.md

Project-specific instructions for coding agents working in `cr3-keywords`.

## Scope and priorities

- Keep changes focused and minimal.
- Preserve the core privacy-first design: local processing for `.CR3` metadata generation.
- Treat `README.md` as the user-facing source of truth for CLI behavior.
- Keep `.agents` memory artifacts accurate when workflow or architecture changes.

## Project snapshot

- Language: Go (`.mise.toml`).
- CLI binary: `cr3`.
- Pipeline: `CR3 -> JPG preview -> TXT (LLM output) -> XMP sidecar`.
- Primary model backend: LM Studio local API (`http://localhost:1234`).
- Optional runtime tools:
    - `lms` CLI (server start + model fallback)
    - `exiftool` (fallback or forced `--exif` extraction)

## Repository map

- Entry point: `cmd/cr3-keyword/main.go`
- CLI parsing/modes: `internal/cli/*`
- Shell completions/install: `internal/cli/shell/*`
- Pipeline orchestration: `internal/pipeline/runner.go`
- XMP generation: `internal/pipeline/3_xmp.go`
- Packaging scripts: `scripts/*.py`

## Required development conventions

- Use `mise` runtime/tooling from repo root.
- Prefer running project commands with `mise exec -- ...`.
- Keep Go code idiomatic and `gofmt`-compatible.
- Do not introduce cloud-only assumptions; local/offline-first behavior is intentional.
- Preserve strict model output contract (keywords line, blank line, caption) unless explicitly changing parser/format semantics.

## Behavior-change checklist

When changing CLI behavior, flags, or help text:

1. Update docs in `README.md`.
2. Keep shell completion generators aligned (`internal/cli/shell/*`).
3. Keep command semantics explicit (`cr3 install <target>`, `cr3 completion <target>`).
4. Preserve Nushell completion UX expectations:
    - curated examples first for `cr3 <TAB>`
    - first positional suggestions are directories (excluding dot-directories)
    - later positional suggestions are `.cr3` files under selected directory

When changing metadata/geotagging behavior:

1. Verify Lightroom field mapping remains consistent.
2. Keep `--gps` semantics strict when explicitly passed.
3. Do not weaken the guarantee that `.CR3` originals are never modified.

## Validation expectations

Run targeted checks for touched areas, then broader checks when appropriate.

Suggested commands (from repo root):

```bash
mise exec -- go test ./...
mise exec -- go build ./...
```

If shell completion code changes, also validate generated output paths/behavior for Nushell, Zsh, and Bash consistency.

## Release and packaging guardrails

- Release flow is documented in `MAINTAINERS.md`.
- Packaging artifacts are generated via `scripts/generate_packaging.py`.
- Keep Homebrew/Scoop manifest templates and generated outputs consistent with release assets.
- Do not change release credential assumptions in workflows without documenting maintainers impact.

## `.agents` workflow expectations

- Follow `.agents/README.md` for skills/memory coordination.
- Keep these files coherent after meaningful architectural/workflow changes:
    - `.agents/memory.md`
    - `.agents/active_context.md`
    - `.agents/decisions.md`
    - `.agents/failures.md`
- Log only real failed attempts in `.agents/failures.md`.
- Keep `README.md` updates scoped to sessions that changed non-`.agents` files.

## Non-goals

- Do not add unrelated refactors during focused fixes.
- Do not create new unit-test suites when only alignment of existing tests is needed.
- Do not break backward-compatible CLI usage unless explicitly requested and documented.
