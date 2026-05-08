#!/usr/bin/env python3
"""Generate packaging manifests from downloaded release artifacts."""

from __future__ import annotations

import argparse
import hashlib
import sys
from pathlib import Path

from shared.version import parse_version_arg


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Generate Homebrew/Scoop manifests from dist artifacts."
    )
    parser.add_argument("--version", required=True, type=parse_version_arg)
    args = parser.parse_args()

    version = args.version

    root = Path.cwd()
    dist_dir = root / "dist"
    out_dir = root / "packaging" / "generated"

    homebrew_template = root / "packaging" / "homebrew" / "cr3-keywords.rb"
    scoop_template = root / "packaging" / "scoop" / "cr3-keywords.json"

    if not homebrew_template.is_file() or not scoop_template.is_file():
        print("Run this script from repository root (cr3-keywords).", file=sys.stderr)
        return 1

    artifacts = {
        "SHA256_DARWIN_ARM64": dist_dir / f"cr3_{version}_darwin-arm64",
        "SHA256_DARWIN_AMD64": dist_dir / f"cr3_{version}_darwin-amd64",
        "SHA256_LINUX_AMD64": dist_dir / f"cr3_{version}_linux-amd64",
        "SHA256_WINDOWS_AMD64": dist_dir / f"cr3_{version}_windows-amd64.exe",
    }

    missing = [str(path) for path in artifacts.values() if not path.is_file()]
    if missing:
        for artifact in missing:
            print(f"Missing artifact: {artifact}", file=sys.stderr)
        print("Download release assets into 'dist/' first.", file=sys.stderr)
        return 1

    replacements = {"{{VERSION}}": version}
    replacements.update(
        {"{{" + name + "}}": sha256_file(path) for name, path in artifacts.items()}
    )

    out_dir.mkdir(parents=True, exist_ok=True)

    homebrew_output = out_dir / "cr3-keywords.rb"
    scoop_output = out_dir / "cr3-keywords.json"

    homebrew_text = homebrew_template.read_text(encoding="utf-8")
    scoop_text = scoop_template.read_text(encoding="utf-8")

    for token, value in replacements.items():
        homebrew_text = homebrew_text.replace(token, value)
        scoop_text = scoop_text.replace(token, value)

    homebrew_output.write_text(homebrew_text, encoding="utf-8")
    scoop_output.write_text(scoop_text, encoding="utf-8")

    print("Generated:")
    print(f"  {homebrew_output}")
    print(f"  {scoop_output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
