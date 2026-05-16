# Active Context

## Current Task

- Finalize end-session workflow updates in `.ai`.
- Ensure new skills are present and wired in `prompts/end-session.md`:
  - `tests/SKILL.md`
  - `docs/SKILL.md`
- Keep memory artifacts synchronized with this session's outcomes.

## Open Issues

- No blocking issues identified.
- `docs` skill must only update root `README.md` when any changed file is outside `.ai/`.
- `tests` skill must only update/run existing unit tests and never create new test files.

## Next Steps

1. On code-changing sessions, apply `tests` skill to adjust existing unit tests and run the test suite.
2. On sessions with non-`.ai` file changes, apply `docs` skill to update root `README.md`.
3. Continue maintaining `.ai/README.md` whenever `.ai` prompts/skills/state files change.

## Recently Changed

- Added `.ai/skills/tests/SKILL.md`.
- Added `.ai/skills/docs/SKILL.md`.
- Updated `.ai/prompts/end-session.md` to apply both skills after `ai-janitor`.
- Refreshed `.ai/README.md` to include the new skills and end-session dependency links.