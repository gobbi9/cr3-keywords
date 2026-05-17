---
name: docs
description: Update the root README.md only when session changes include at least one file outside .ai.
---

# Docs

You maintain the root project `README.md`.

## Responsibility

Update `README.md` if and only if the session changed at least one file outside `.ai`.

## Decision rule

1. Inspect changed files for the session.
2. If every changed file is under `.ai/`, do not modify `README.md`.
3. If any changed file is outside `.ai/`, update `README.md` to reflect those changes.

## Operating guidelines

1. Keep `README.md` accurate, concise, and user-focused.
2. Document externally visible behavior, setup, usage, and workflow impacts.
3. Avoid documenting internal-only churn that does not affect users.
4. Keep section ordering and style consistent with existing README conventions.
5. This skill governs update conditions for `README.md`; it does **not** change dependency relationships in `.ai/README.md` diagrams.
6. In `.ai/README.md` dependency graphs, keep the conceptual edge `docs/SKILL.md -> ../README.md` stable.