# Active Context

## Current Task

- Refactored shell completion generators to external template files for readability:
  - `internal/cli/shell/nushell.tmpl.nu`
  - `internal/cli/shell/zsh.tmpl.zsh`
  - `internal/cli/shell/bash.tmpl.sh`
  - shared renderer in `internal/cli/shell/templates.go`
- Updated CLI `--help` output to a man-style sectioned format (`NAME`, `SYNOPSIS`, `OPTIONS`, etc.).
- Added first-class man page source at `docs/man/cr3.1`.
- Added man-page packaging flow:
  - release workflow copies `docs/man/cr3.1` into `dist/cr3.1`
  - packaging generation now computes `SHA256_MANPAGE`
  - Homebrew formula installs `cr3.1` and links `exiftool`
  - Scoop manifest suggests `extras/exiftool` and notes to use `cr3 --help` on Windows
- Moved shell-completion maintenance responsibility into `.agents/skills/docs/SKILL.md` and removed `.agents/skills/shell-completions`.
- Nushell module keeps a conservative `help cr3` fallback behavior; local-man verification remains explicit via overlay helper command.

## Open Issues

- Local runtime validation of freshly built `cr3` binary is blocked in this agent environment by macOS loader error (`missing LC_UUID load command`).
- Nushell command-definition behavior (`extern`) can still shadow native external help ergonomics, so `^cr3 --help`/explicit man checks remain the reliable validation path.

## Next Steps

1. Cut a release to validate `dist/cr3.1` asset and packaging generation in CI.
2. Decide whether to keep manual dual-maintenance (`Usage()` + `docs/man/cr3.1`) or adopt generation to reduce drift risk.
3. Verify Homebrew/Scoop downstream rendering after next packaging publish.

## Recently Changed

- `internal/cli/options.go`
- `internal/cli/shell/*` (new template-based rendering)
- `docs/man/cr3.1`
- `packaging/homebrew/cr3-keywords.rb`
- `packaging/scoop/cr3-keywords.json`
- `scripts/generate_packaging.py`
- `.github/workflows/release.yml`
- `README.md`
- `MAINTAINERS.md`
- `.agents` docs/skills metadata files

## Validation

- `mise exec -- gofmt -w ...` (changed Go files)
- `mise exec -- go test ./...` passed
- `mise exec -- go build ./...` passed
- Nushell module source validation via:
  - sourcing generated module in Nu runtime
  - checking behavior around extern/help shadowing and explicit fallback paths
