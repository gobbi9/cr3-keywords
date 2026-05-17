# Active Context

## Current Task

- Session request completed: updated root `README.md` Development tooling docs to add `mise` as the first/recommended option and keep `goenv` as an alternative.
- Documented required `mise` setup commands:
  - `mise settings add idiomatic_version_file_enable_tools go`
  - `mise install`
- Finalizing end-session memory/decision records and refreshing `.ai/README.md`.

## Open Issues

- No blocking issues identified.
- `tests` skill: no behavior/code changes this session; test updates/runs were not required.
- `docs` skill: satisfied (session changed a non-`.ai` file and `README.md` was updated accordingly).

## Next Steps

1. If development tooling changes again, keep both `mise` and `goenv` instructions aligned with actual project version-file behavior.
2. Continue applying end-session maintenance flow (`memento` -> `ai-janitor` -> `tests` -> `docs`).
3. Keep `.ai` memory artifacts concise and synchronized with doc/config workflow changes.

## Recently Changed

- Updated `README.md` Development section:
  - Added `mise` link and positioned it as recommended.
  - Kept `goenv` as alternative.
  - Added mandatory `mise` setup commands for `.go-version` idiomatic tool handling.