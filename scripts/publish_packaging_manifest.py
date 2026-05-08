#!/usr/bin/env python3
"""Publish generated packaging manifest to an external GitHub repository."""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
import tempfile
from pathlib import Path

from shared.logging import logger


def run(
    command: list[str], *, cwd: Path | None = None, check: bool = True
) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=str(cwd) if cwd else None,
        check=check,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )


def main() -> int:
    token = os.environ.get("PACKAGING_PUSH_TOKEN", "")
    if not token:
        logger.error("PACKAGING_PUSH_TOKEN env var is required.")
        return 1

    parser = argparse.ArgumentParser(
        description="Publish one packaging manifest file to an external repo."
    )
    parser.add_argument(
        "--target-repo", required=True, help="GitHub repo in owner/name form"
    )
    parser.add_argument(
        "--source-file", required=True, help="Path of generated source file"
    )
    parser.add_argument(
        "--target-file", required=True, help="Destination path inside target repo"
    )
    parser.add_argument("--display-name", required=True, help="Friendly name for logs")
    parser.add_argument("--commit-message", required=True, help="Commit message")
    args = parser.parse_args()

    source_file = Path(args.source_file)
    if not source_file.is_file():
        logger.error("Source file not found: %s", source_file)
        return 1

    run(["git", "config", "--global", "user.name", "github-actions[bot]"])
    run(
        [
            "git",
            "config",
            "--global",
            "user.email",
            "41898282+github-actions[bot]@users.noreply.github.com",
        ]
    )

    logger.info("Publishing %s to %s", args.display_name, args.target_repo)

    with tempfile.TemporaryDirectory() as tmp:
        work_dir = Path(tmp)
        repo_dir = work_dir / "repo"

        clone_url = f"https://x-access-token:{token}@github.com/{args.target_repo}.git"
        run(["git", "clone", clone_url, str(repo_dir)])

        destination = repo_dir / args.target_file
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(source_file, destination)

        diff_result = run(
            ["git", "diff", "--quiet", "--", args.target_file],
            cwd=repo_dir,
            check=False,
        )
        if diff_result.returncode == 0:
            logger.warning("No %s changes to commit.", args.display_name)
            return 0

        run(["git", "add", args.target_file], cwd=repo_dir)
        run(["git", "commit", "-m", args.commit_message], cwd=repo_dir)
        push_result = run(["git", "push"], cwd=repo_dir)
        if push_result.stdout:
            logger.info("%s", push_result.stdout.strip())

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
