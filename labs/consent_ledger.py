#!/usr/bin/env python3
"""Thin wrapper delegating to canonical code/python/consent_ledger.py."""
from pathlib import Path
import runpy

runpy.run_path(str(Path(__file__).resolve().parent.parent / "code" / "python" / "consent_ledger.py"), run_name="__main__")
