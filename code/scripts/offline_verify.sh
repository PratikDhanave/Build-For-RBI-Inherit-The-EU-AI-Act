#!/usr/bin/env bash
# Five-minute offline verification script (no cloud account, no API key, no network)
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "=== 1. Go Composite Policy Gate & Approval Drift Tests ==="
(cd "$ROOT_DIR/code" && go test -v ./...)

echo ""
echo "=== 2. Python Composite Policy Gate Demo ==="
python3 "$ROOT_DIR/code/python/gate_demo.py" --test

echo ""
echo "=== 3. Python Lab 01: Approval Drift (approval.py) ==="
python3 "$ROOT_DIR/code/python/approval.py" --test

echo ""
echo "=== 4. Python Lab 02: Tamper-Evident Hash Chain (audit.py) ==="
python3 "$ROOT_DIR/code/python/audit.py" --test

echo ""
echo "[OK] Offline verification complete — zero network required."
