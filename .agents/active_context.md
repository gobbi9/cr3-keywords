# Active Context

## Current Task

- Added Nushell-only top-of-list examples for `cr3 <TAB>` in `internal/cli/shell/nushell.go` using record completions (`value` + `description`).
- Updated those examples to five requested workflows:
  1. recommended `--gps` + `--exif`
  2. same flow with `--edit-prompt`
  3. full control run with `--gps --exif --edit-prompt --model --prompt`
  4. test mode with `--verbose --dry-run`
  5. cleanup flow with `--clear`
- Preserved existing Nushell completion behavior after the examples (directory-first, then `.cr3` suggestions from selected directory).
- Updated root `README.md` completion section to document the new Nushell example-first behavior.

## Open Issues

- No active blockers identified.

## Next Steps

1. If desired, tune example paths/model names in Nushell examples to user-specific defaults.
2. Keep shell completion docs aligned with future UX tweaks.
3. Continue applying end-session flow: `memento` -> `shell-completions` -> `ai-janitor` -> `tests` -> `docs`.

## Recently Changed

- `internal/cli/shell/nushell.go`: added curated example completions and merged with existing directory completions for `cr3 ` context.
- `README.md`: documented Nushell curated examples before standard positional suggestions.
- Validation: `mise exec -- go test ./...` passed.
