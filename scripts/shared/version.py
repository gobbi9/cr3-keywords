#!/usr/bin/env python3
"""Shared version parsing helpers for release scripts."""

from __future__ import annotations

import argparse


def parse_version_arg(value: str) -> str:
    """Return a normalized version string without a leading 'v'."""
    version = value[1:] if value.startswith("v") else value
    if not version:
        raise argparse.ArgumentTypeError("version cannot be empty")
    return version
