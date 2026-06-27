---
name: docs
description: Maintain user-facing docs, shell completion docs, and man-page docs when non-.agents files change.
---

# Docs

You maintain user-facing documentation.

## Responsibility

Update docs if and only if the session changed at least one file outside `.agents`.

## Decision rule

1. Inspect changed files for the session.
2. If every changed file is under `.agents/`, do not modify project docs.
3. If any changed file is outside `.agents/`, update docs to reflect those changes.

## Scope

When updates are required, keep these sources aligned:

- Root `README.md` (primary user-facing source of truth)
- Shell completion docs in `README.md` for `cr3 install` / `cr3 completion`
- Man-page workflow docs in `README.md` and/or `MAINTAINERS.md`
- CLI help text in `internal/cli/options.go` (`Usage()`)
- Man page source in `docs/man/cr3.1`

## External tool documentation policy

When external tools are invoked by project code (for example `lms`, `exiftool`):

1. Track exact commands used in code and keep them documented in `docs/man/cr3.1` under `RELATED COMMANDS`.
2. If a tool has a real man page, include it in `SEE ALSO` with section notation (for example `exiftool(1)`).
3. If a tool does not have a man page, do **not** fake section notation; document it only in `RELATED COMMANDS`.
4. Re-check these entries whenever command arguments change in code.

## Anti-drift policy (`--help` vs `man cr3`)

The `docs` skill is responsible for keeping `cr3 --help` and `man cr3` semantically aligned.

Required when CLI behavior/help changes:

1. Update `internal/cli/options.go` (`Usage()`).
2. Update `docs/man/cr3.1` in the same session.
3. Ensure `README.md` usage/help examples remain consistent.
4. Ensure release/install docs still reflect how man page is shipped (`docs/man/cr3.1` -> release asset -> package/install paths).

If exact wording differs due to format constraints (plain text vs roff), keep command/flag semantics identical.

Do not leave one updated without the other.

## Operating guidelines

1. Keep documentation accurate, concise, and user-focused.
2. Document externally visible behavior, setup, usage, and workflow impacts.
3. Avoid documenting internal-only churn that does not affect users.
4. Keep section ordering and style consistent with existing repo conventions.
5. Include shell completion behavior updates whenever CLI command/flag shapes change.
6. Include man-page maintenance/update instructions whenever help/man behavior or packaging/man assets change.
7. Explicitly check for `--help`/man drift whenever either file changes.
8. This skill governs doc update conditions; it does **not** change dependency relationships in `.agents/README.md` diagrams.
9. In `.agents/README.md` dependency graphs, keep conceptual edge `docs/SKILL.md -> ../README.md` stable.
