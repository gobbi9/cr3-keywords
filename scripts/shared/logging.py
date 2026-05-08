#!/usr/bin/env python3
"""Shared logging configuration for release scripts."""

from __future__ import annotations

import logging
import sys

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    stream=sys.stdout,
    force=True,
)

logger = logging.getLogger(__name__)
