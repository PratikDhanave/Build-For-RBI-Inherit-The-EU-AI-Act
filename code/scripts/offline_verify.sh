#!/usr/bin/env bash
# Five-minute offline verification script (100% pure Go — no cloud account, no API key, no network)
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "=== 1. Go Vet, Race Detector & 20-Case OWASP ASI Red-Team Suite ==="
(cd "$ROOT_DIR/code" && go vet ./... && go test -race -count=1 -v ./...)

echo ""
echo "=== 2. Go End-to-End Offline CLI Pipeline (go run ./cmd/demo) ==="
(cd "$ROOT_DIR/code" && go run ./cmd/demo)

echo ""
echo "=== 3. Go Lab: Composite Policy Gate (labs/gate_demo/main.go) ==="
go run "$ROOT_DIR/labs/gate_demo/main.go" --test

echo ""
echo "=== 4. Go Lab 01: Approval Drift (labs/approval/main.go) ==="
go run "$ROOT_DIR/labs/approval/main.go" --test

echo ""
echo "=== 5. Go Lab 02: Tamper-Evident Hash Chain (labs/audit/main.go) ==="
go run "$ROOT_DIR/labs/audit/main.go" --test

echo ""
echo "=== 6. Go Lab 03: Append-Only Consent Ledger (labs/consent_ledger/main.go) ==="
go run "$ROOT_DIR/labs/consent_ledger/main.go" --test

echo ""
echo "=== 7. Go WebAssembly Evaluator Tests (labs/wasm) ==="
(cd "$ROOT_DIR/labs" && go vet ./... && go test -race -count=1 -v ./...)

echo ""
echo "=== 8. Go Cross-Walk, Policy Parity, Mermaid & Asset Verification ==="
(cd "$ROOT_DIR/code" && go run ./cmd/verify --root "$ROOT_DIR")

echo ""
echo "=== 9. Go Dual Evidence Pack Generator (go run ./cmd/evidencepack) ==="
(cd "$ROOT_DIR/code" && go run ./cmd/evidencepack --out /tmp/rbi-eu-evidence-pack)

echo ""
echo "[OK] Offline verification complete — 100% Go, zero network required."
