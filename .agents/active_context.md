# Active Context

## Current Task

- Completed migration from `.go-version` to `mise` tool definition via `.mise.toml` (`go = "1.22.5"`).
- Updated root `README.md` development setup to reference `.mise.toml` and removed the now-unneeded `mise settings add idiomatic_version_file_enable_tools go` step.
- Investigated Zed `gopls` warning (`No active builds contain ... (go list)`) and verified root cause was editor runtime environment (`go` not found on PATH for Zed process), not module/package layout.

## Open Issues

- No repository code blockers.
- Environment note: GUI-launched Zed may still need PATH initialization aligned with `mise` (or fallback Go on PATH) so `gopls` can run `go list` reliably.

## Next Steps

1. If `gopls` warnings recur, launch Zed from a shell initialized with `mise` or ensure GUI app PATH includes Go.
2. Keep runtime docs aligned with actual project runtime source of truth (`.mise.toml`).
3. Continue applying end-session flow: `memento` -> `shell-completions` -> `ai-janitor` -> `tests` -> `docs`.

## Recently Changed

- Added `.mise.toml` with Go tool version pin.
- Removed `.go-version`.
- Updated development docs in `README.md` to match the new `mise` workflow.
- Confirmed `mise exec -- go list ./...` includes `internal/cli` while plain `go list ./...` fails when `go` is missing from PATH.
