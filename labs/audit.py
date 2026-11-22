#!/usr/bin/env python3
"""Thin wrapper delegating to canonical code/python/audit.py."""
from pathlib import Path
import runpy

runpy.run_path(str(Path(__file__).resolve().parent.parent / "code" / "python" / "audit.py"), run_name="__main__")
