#!/usr/bin/env python3
"""Thin wrapper delegating to canonical code/python/approval.py."""
from pathlib import Path
import runpy

runpy.run_path(str(Path(__file__).resolve().parent.parent / "code" / "python" / "approval.py"), run_name="__main__")
