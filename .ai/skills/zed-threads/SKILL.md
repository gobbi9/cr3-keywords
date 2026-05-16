---
name: zed-threads
description: Read and decode Zed thread records from the local SQLite database one by one, maintaining an automatic cursor so historical threads can be processed safely without manual tracking.
---

# Zed Threads

You ingest Zed thread history from the local `threads.db` and return decoded thread data incrementally.

## Source

Use this Nushell pipeline to read one thread at a time:

```/dev/null/zed-threads.nu#L1-12
openn `~/Library/Application Support/Zed/threads/threads.db`
| get threads
| where folder_paths =~ 'cr3-keywords'
| update data {|row|
    $row.data
    | into binary
    | ^zstd -d
    | decode utf-8
    | from json
  }
| reject folder_paths folder_paths_order parent_id id data_type
| skip 0
| first 1
```

## Shell handling

1. If the current shell is `zsh` or `bash`, start `nu` first and run the pipeline inside Nushell.
2. If already in Nushell, run the pipeline directly.

## Cursor/state management (no manual tracking)

Persist progress in a cursor file so the user never has to track indices manually.

- Cursor file: `.ai/skills/zed-threads/state/nu.cursor`
- Value stored: integer `i` (0-based)

Workflow:

1. Ensure `.ai/skills/zed-threads/state` exists.
2. Read `i` from `.ai/skills/zed-threads/state/nu.cursor`.
   - If missing/invalid, default to `0`.
3. Execute the pipeline with:
   - `skip i`
   - `first (i + 1)`
4. Process only the newly included tail thread (the current thread at index `i`) to keep context small.
5. Increment cursor to `i + 1` and write it back to `.ai/skills/zed-threads/state/nu.cursor`.
6. Stop when no new thread is returned.

## Requirements

1. Always process threads sequentially, one at a time.
2. Never load all threads into context at once.
3. Keep extracted summaries compact and actionable for memory seeding.
4. Be resilient to decode/command errors: report the failure and continue with the next thread when possible.
5. Do not ask the user to manually provide the next index.
6. The cursor file is the source of truth for progress.