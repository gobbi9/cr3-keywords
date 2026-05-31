---
name: shell-completions
description: Keep CLI shell completion/module generators and install paths in sync with command/flag changes.
---

# Shell Completions

Use this skill whenever CLI commands, subcommands, flags, or argument shapes change.

## Goal

Keep shell completion definitions accurate across:

- `internal/cli/shell/nushell.go`
- `internal/cli/shell/zsh.go`
- `internal/cli/shell/bash.go`
- `internal/cli/shell/shell.go` (target routing + install path conventions)
- `internal/cli/shell/completion.go` (script installation)

## What to update

1. Add/remove/update commands and flags in all supported shells.
2. Keep install target lists aligned (`nushell`, `zsh`, `bash`).
3. Keep generated help strings and argument expectations consistent with CLI parser behavior.
4. When install paths change, update `shell.InstallPath` and any related docs.
5. Run formatting/tests after edits:
   - `gofmt -w <changed-go-files>`
   - `go test ./...`

## Constraints

- Keep behavior equivalent across shells as closely as each shell allows.
- If a shell cannot express an exact CLI shape (for example optional flag values), document the limitation in generated output comments.
- Do not introduce new completion dependencies unless explicitly requested.
