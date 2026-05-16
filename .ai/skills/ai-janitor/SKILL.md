---
name: ai-janitor
description: Keep .ai/README.md accurate and up to date whenever files in .ai change, including prompts, skills, and memory artifacts.
---

# AI Janitor

You maintain `.ai/README.md` as the canonical index of the `.ai` workspace.

## Responsibility

Whenever anything in `.ai` changes (new file, deletion, rename, or content update), refresh `.ai/README.md`.

## What to maintain in `.ai/README.md`

1. A short description for every prompt in `.ai/prompts`.
2. A short description for every skill in `.ai/skills`.
3. A short explanation for every other `.md` file in `.ai`.
4. A Mermaid dependency diagram showing prompt -> skill -> `.ai` file relationships.

## Diagram rules

1. Use simple file names in node labels.
2. Use relative path fragments only when needed to avoid ambiguity (for example, `memento/SKILL.md`).
3. Keep the graph readable and focused on real dependencies.
4. Include state files when they are part of a skill workflow (for example, `nu.cursor`).

## Operating guidelines

1. Keep descriptions concise and practical.
2. Avoid stale entries and duplicates.
3. Keep ordering stable for easy diffs.
4. Do not document files outside `.ai` in `.ai/README.md`.
5. Treat `.ai/README.md` as required maintenance after `.ai` updates.