# Decisions

### 2026-06-27

Decision:
- Keep man-page source in `docs/man/cr3.1` and stage a separate repo-local manpath (`.man/man1`) during `make build` for verification, rather than writing build output back into docs source paths.

Reason:
- Preserves clean separation between authored docs and build artifacts while enabling local `man` checks without system-wide install.

Consequences:
- `overlay.nu` provides explicit local-man helper behavior and `.man/` is ignored in git.

### 2026-06-27

Decision:
- Add `{{config_root}}/bin` to project `.mise.toml` PATH injection and remove explicit `cr3` alias from repo overlay.

Reason:
- Simplifies overlay maintenance and keeps binary discovery consistent for repo-local tooling.

Consequences:
- `cr3` resolves from repo `bin/` when using `mise` in the project.

### 2026-06-27

Decision:
- Track external tool commands in man-page `RELATED COMMANDS`; include only tools with real manpages in `SEE ALSO`.

Reason:
- Keeps CLI operational dependencies discoverable while preserving standard manpage semantics.

Consequences:
- `lms` and concrete `exiftool` command invocations are documented under `RELATED COMMANDS`; `SEE ALSO` remains strict (`exiftool(1)`).


### 2026-06-27

Decision:
- Consolidate shell-completion ownership into `docs` skill and remove the dedicated `shell-completions` skill.

Reason:
- Completion/help/man documentation and user-facing command docs must stay aligned, and a single docs-oriented skill reduces coordination overhead.

Consequences:
- End-session skill chain now omits `shell-completions`.
- `.agents` docs and workflow references must point to `docs/SKILL.md` for completion/man guidance.

### 2026-06-27

Decision:
- Externalize shell completion script bodies into template files (`*.tmpl`) rendered by embedded Go templates.

Reason:
- Improves readability and maintainability of shell script content without changing completion install/routing semantics.

Consequences:
- `internal/cli/shell/{bash,zsh,nushell}.go` now delegate rendering.
- New shared renderer in `internal/cli/shell/templates.go` and template assets are embedded in the binary.

### 2026-06-27

Decision:
- Introduce man-page distribution (`docs/man/cr3.1`) and include it in release + Homebrew packaging; treat Scoop as non-man-page platform and provide `--help` guidance there.

Reason:
- Align CLI UX with Unix expectations (`man cr3`) while keeping Windows packaging pragmatic.

Consequences:
- Release workflow now ships `cr3.1` asset.
- Packaging manifest generation includes `SHA256_MANPAGE`.
- Homebrew formula installs man page and links `exiftool`; Scoop suggests `exiftool` and documents `cr3 --help`.


### 2026-06-07

Decision:
- For Nushell only, prepend `cr3 <TAB>` completions with curated workflow examples that include descriptions, while preserving existing directory and `.cr3` completion behavior after those examples.

Reason:
- Improves discoverability of common command patterns directly in interactive completion without changing parser semantics or other shells.

Consequences:
- `internal/cli/shell/nushell.go` now includes an examples completer that returns structured completion records (`value`, `description`) merged ahead of directory suggestions in root command context.
- README completion docs include the Nushell example-first behavior.

### 2026-06-07

Decision:
- Replace `.go-version` with project-local `mise` configuration in `.mise.toml` (`[tools] go = "1.22.5"`).
- Simplify development docs to `mise install` without `idiomatic_version_file_enable_tools` setup.

Reason:
- The project now defines runtime versions directly in `mise` native config, so idiomatic-version-file bridging is unnecessary.
- Reduces setup friction and keeps runtime source-of-truth in one place.

Consequences:
- Contributors should rely on `.mise.toml` for Go version selection.
- `gopls` issues like `No active builds contain ... (go list)` should be debugged as editor PATH/runtime environment problems when module layout is otherwise valid.


### 2026-05-31

Decision:
- Standardize shell integration commands to require explicit targets: `cr3 install <nushell|zsh|bash>` and `cr3 completion <nushell|zsh|bash>`.

Reason:
- Removes parser ambiguity and keeps command semantics consistent across install/completion flows.

Consequences:
- Help/docs and all shell completion generators must stay aligned with explicit target requirements.
- `cr3 install` without a target now returns a user-facing parser error.

### 2026-05-31

Decision:
- Consolidate shell completion generation/install logic under `internal/cli/shell/*` and keep Nushell extern definitions subcommand-scoped (`"cr3 install"`, `"cr3 completion"`).

Reason:
- Separate shell-specific responsibilities cleanly and avoid Nushell duplicate-command parse errors caused by multiple `export extern "cr3"` blocks.

Consequences:
- Main command wiring now calls `shell.Script` / `shell.Install` directly.
- Completion maintenance has a single, explicit package boundary and skill ownership (`shell-completions`).

### 2026-05-31

Decision:
- Nushell completion behavior for the main `cr3` flow is context-aware:
  - first positional argument suggests directories (excluding dot-directories)
  - subsequent file suggestions are `.cr3` files from the selected first directory, case-insensitive.

Reason:
- Match expected interactive UX for photo-folder selection and avoid irrelevant file suggestions.

Consequences:
- Nushell completer functions now parse command context and apply scoped filtering logic.
- Future quoted-path support may require more robust context tokenization.

### 2026-05-17

Decision:
- In Development docs, prefer `mise` as the first toolchain manager option and keep `goenv` as an alternative.
- Document `mise` bootstrap requirements explicitly: `mise settings add idiomatic_version_file_enable_tools go` and `mise install`.

Reason:
- The project relies on `.go-version`; `mise` needs idiomatic version-file tool support enabled so Go version resolution works as expected.
- Keep setup guidance modern while preserving compatibility for contributors already using `goenv`.

Consequences:
- README onboarding for contributors is clearer and less error-prone for `mise` users.
- Development setup now has two documented, explicit paths (`mise` recommended, `goenv` alternative).

### 2026-05-16

Decision:
- Extend end-session workflow with `tests` and `docs` skills, applied after `ai-janitor`.

Reason:
- Make session close-out enforce two guardrails:
  - keep existing unit tests aligned and executed when code changes occur
  - update root `README.md` only when the session changed files outside `.agents/`

Consequences:
- End-session behavior now includes test/doc maintenance policy checks.
- Documentation updates are explicitly gated to non-`.agents` changes, reducing unnecessary README churn.

### 2026-05-16

Decision:
- Add `.agents` memory workflow with `theory-of-mind` and `ai-janitor` skills.

Reason:
- Preserve project context across sessions and prevent loss of architectural/history knowledge.

Consequences:
- Added recurring maintenance responsibility for `.agents/README.md` and long-term memory summaries.

### 2026-05-14 to 2026-05-15

Decision:
- Improve LM Studio client robustness (server check strategy, diagnostics/logging behavior).

Reason:
- Reduce repeated overhead and improve troubleshooting clarity for API/server issues.

Consequences:
- Better runtime efficiency and clearer debug information during failures.

### 2026-05-14

Decision:
- Add GPX-aware geotagging and prompt injection of location context.

Reason:
- Improve caption/keyword quality and write Lightroom-compatible location metadata.

Consequences:
- Added optional geospatial dependency path (`--gps`) and expanded XMP mapping responsibilities.

### 2026-05-14

Decision:
- Keep only `--version` flag (remove `version` subcommand style).

Reason:
- Align with modern CLI conventions and reduce parser branching.

Consequences:
- Simpler help/usage and reduced command-mode complexity.

### 2026-05-09

Decision:
- Remove positional argument mode for model/prompt and keep flag-driven interface.

Reason:
- Positional mode created ambiguity and increased parser complexity.

Consequences:
- CLI semantics became more explicit and docs became easier to keep accurate.

### 2026-05-07 to 2026-05-08

Decision:
- Build GitHub release automation and package distribution to Homebrew/Scoop.

Reason:
- Project was preparing for public/open-source consumption and needed reproducible distribution.

Consequences:
- Added workflow/pipeline complexity and credential/permissions requirements, but enabled repeatable cross-platform delivery.

### 2026-05-07

Decision:
- Migrate core pipeline from zsh script into a Go CLI (`cr3`).

Reason:
- Better maintainability, clearer structure, typed error handling, and easier distribution.

Consequences:
- New module layout (`cmd`, `internal/*`), explicit option parsing, and standardized logging/progress behavior.

### 2026-05-07

Decision:
- Add explicit timing/progress UX and introduce cleanup workflow via `--clear`.

Reason:
- Step 2 (LLM prompting) is the dominant latency source; users needed clear runtime feedback and cleanup controls.

Consequences:
- Better operator confidence during long runs and easier post-run artifact cleanup.

### 2026-05-07

Decision:
- Standardize binary naming (`cr3`) and implement version metadata in build/release flow.

Reason:
- Simpler CLI ergonomics and cleaner release story for end users.

Consequences:
- Improved install UX, clearer `--version` output, and reduced release confusion.

### 2026-05-07

Decision:
- Implement pure-Go CR3 preview extraction as default; keep `exiftool` fallback and optional forced mode `--exif`.

Reason:
- Reduce external dependency burden while retaining performance/reliability fallback.

Consequences:
- Better out-of-box behavior with optional fast path for users with `exiftool` installed.

### 2026-05-06

Decision:
- Start from script-first implementation, then progressively formalize into maintainable CLI architecture.

Reason:
- Needed quick functional progress for CR3 -> JPG -> TXT -> XMP flow before investing in stronger structure.

Consequences:
- Early velocity was high, followed by iterative refactors to reduce script complexity and edge-case fragility.
