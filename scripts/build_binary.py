#!/usr/bin/env python3
"""Build one release binary for a target GOOS/GOARCH."""

from __future__ import annotations

import argparse
import subprocess
from pathlib import Path

from shared.logging import logger
from shared.version import parse_version_arg


def main() -> int:
    parser = argparse.ArgumentParser(description="Build one release binary")
    parser.add_argument("--goos", required=True)
    parser.add_argument("--goarch", required=True)
    parser.add_argument("--asset-suffix", required=True)
    parser.add_argument("--ext", default="")
    parser.add_argument("--version", required=True, type=parse_version_arg)
    args = parser.parse_args()

    version = args.version

    output = Path("dist") / f"cr3_{version}_{args.asset_suffix}{args.ext}"

    command = [
        "make",
        "build",
        f"VERSION={version}",
        f"OUT={output.as_posix()}",
        f"GOOS={args.goos}",
        f"GOARCH={args.goarch}",
    ]
    subprocess.run(command, check=True)

    logger.info("Built %s", output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
