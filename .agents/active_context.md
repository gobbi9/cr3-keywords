# Active Context

## Current Task

- Session request completed: shell integration UX was extended and stabilized across `nushell`, `zsh`, and `bash`.
- Implemented `cr3 completion <target>` and `cr3 install <target>` command flows with explicit required targets.
- Refactored completion generation/install code into `internal/cli/shell/*` and fixed Nushell-specific completion behavior.

## Open Issues

- No blocking issues identified.
- Known follow-up opportunity: Nushell completion parsing still uses space-token splitting; quoted paths containing spaces may need dedicated handling if this becomes a user scenario.

## Next Steps

1. If CLI flags/subcommands change again, update all shell generators (`nushell.go`, `zsh.go`, `bash.go`) in lockstep.
2. Keep `README.md` usage/install examples aligned with parser behavior and shell completions.
3. Continue applying end-session flow: `memento` -> `shell-completions` -> `ai-janitor` -> `tests` -> `docs`.

## Recently Changed

- Added shell completion command surface:
  - `cr3 completion nushell|zsh|bash`
  - `cr3 install nushell|zsh|bash`
- Removed implicit default for `cr3 install` (target is now required).
- Fixed Nushell extern duplication error by switching to subcommand externs (`"cr3 install"`, `"cr3 completion"`).
- Improved Nushell completers:
  - first positional completion lists directories (excluding dot-directories)
  - subsequent file completion is scoped to selected CR3 directory and filters `.cr3` case-insensitively.
