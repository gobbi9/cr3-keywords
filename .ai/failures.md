# Failures and Dead Ends

### 2026-05-16

Attempted:
- Memory seeding without persistent thread-ingestion state.

Why it failed:
- Historical processing became easy to duplicate or skip without reliable progress tracking.

What was learned:
- Cursor-based incremental ingestion (`.ai/skills/zed-threads/state/nu.cursor`) is required for safe, repeatable memory seeding.

### 2026-05-14 to 2026-05-15

Attempted:
- Re-checking LM Studio server readiness too frequently (per request) and limited raw-response visibility during debugging.

Why it failed:
- Added avoidable overhead and slowed diagnosis of malformed/edge responses.

What was learned:
- Perform readiness checks once per batch and keep targeted debug logging available for client/API interactions.

### 2026-05-14

Attempted:
- Initial Lightroom location mapping in XMP with incomplete/incorrect field mapping.

Why it failed:
- Some metadata fields did not align with Lightroom expectations (notably sublocation mapping behavior).

What was learned:
- Validate metadata mappings against target app semantics, not only schema appearance.

### 2026-05-09

Attempted:
- Supporting mixed positional + flag argument styles in CLI parsing.

Why it failed:
- Increased ambiguity and parser complexity; harder to document and reason about edge cases.

What was learned:
- Prefer explicit flags for optional values and keep argument grammar narrow.

### 2026-05-07 to 2026-05-08

Attempted:
- Fast-path release/publishing automation without fully accounting for GitHub expression and permission constraints.

Why it failed:
- Workflow conditionals and secrets handling caused YAML/runtime validation issues; branch/ruleset expectations were initially conflated with workflow behavior.

What was learned:
- Validate expression contexts early, isolate CI auth assumptions, and explicitly test tag-triggered paths.

### 2026-05-07

Attempted:
- Early Go progress output with incorrect format-string argument alignment.

Why it failed:
- ANSI/color reset argument did not match format placeholders, producing malformed output (`%!d(...)`-style artifacts).

What was learned:
- Keep progress rendering format strings minimal and test with real terminal output after each change.
