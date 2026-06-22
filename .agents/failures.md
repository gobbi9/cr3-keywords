# Failures and Dead Ends

### 2026-06-07

Attempted:
- Investigated Zed `gopls` warning (`No active builds contain ... (go list)`) as a potential module/package inclusion issue.

Why it failed:
- Module/package configuration was valid; the actual failure was environment-level (`go` missing from Zed process PATH), so `go list` could not run.

What was learned:
- Validate tool availability in the editor runtime first (`go list ./...` vs `mise exec -- go list ./...`) before changing module/workspace files.


### 2026-05-31

Attempted:
- Defined Nushell install/completion signatures as multiple `export extern "cr3" [...]` blocks.

Why it failed:
- Nushell parser rejects duplicate command definitions within a block (`nu::parser::duplicate_command_def`).

What was learned:
- Define subcommand externs explicitly (`export extern "cr3 install"`, `export extern "cr3 completion"`) instead of redefining the root command.

### 2026-05-31

Attempted:
- Parsed Nushell completion context using `split words`.

Why it failed:
- Path tokens containing separators were fragmented during completion-context parsing, so directory/file suggestions were resolved against wrong paths.

What was learned:
- Use row-based token splitting for completion context (`split row " "` + filtering) and treat trailing `/` tokens explicitly for directory completion.

### 2026-05-16

Attempted:
- Memory seeding without persistent thread-ingestion state.

Why it failed:
- Historical processing became easy to duplicate or skip without reliable progress tracking.

What was learned:
- Durable progress tracking is required for safe, repeatable historical-memory seeding.

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
