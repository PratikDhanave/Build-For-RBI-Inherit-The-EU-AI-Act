#!/usr/bin/env python3
"""Thin wrapper delegating to canonical code/python/gate_demo.py."""
from pathlib import Path
import runpy

runpy.run_path(str(Path(__file__).resolve().parent.parent / "code" / "python" / "gate_demo.py"), run_name="__main__")
