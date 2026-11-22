#!/usr/bin/env python3
"""
Lab 02 — Tamper-Evident SHA-256 Hash Chain & Signed Export Manifest
(RBI FREE-AI Rec 22 / EU AI Act Article 12).

Demonstrates why an audit trail must catch its own operator: editing row 2
invalidates the hash chain at seq 2, and the export manifest honestly reports
`chain_verified: False` rather than hiding the break.

Run:
    python3 labs/audit.py
    python3 labs/audit.py --test
"""

import hashlib
import json
import sys


def row_hash(prev: str, entry: dict) -> str:
    body = json.dumps({"prev": prev, **entry}, sort_keys=True).encode()
    return hashlib.sha256(body).hexdigest()


def append(log: list[dict], actor: str, action: str, target: str) -> None:
    prev = log[-1]["row_hash"] if log else "genesis"
    e = {"seq": len(log) + 1, "actor": actor, "action": action, "target": target}
    e["prev_hash"] = prev
    e["row_hash"] = row_hash(prev, e)
    log.append(e)


def verify(log: list[dict]) -> tuple[bool, int | None, str | None]:
    prev = "genesis"
    for e in log:
        if e["prev_hash"] != prev:
            return False, e["seq"], "prev_hash"
        chk = {k: v for k, v in e.items() if k != "row_hash"}
        if row_hash(prev, chk) != e["row_hash"]:
            return False, e["seq"], "row_hash"
        prev = e["row_hash"]
    return True, None, None


def main() -> int:
    log: list[dict] = []
    append(log, "system", "consent.grant", "user-1")
    append(log, "advisor-7", "portfolio.read", "user-1")
    append(log, "payment_orch", "refund.issue", "txn-9931")

    fresh_ok, fresh_seq, fresh_field = verify(log)
    print(
        "fresh chain          ",
        "OK" if fresh_ok else f"BROKEN at seq {fresh_seq} ({fresh_field})",
    )

    log[1]["actor"] = "advisor-99"  # someone edits history
    tampered_ok, tampered_seq, tampered_field = verify(log)
    print(
        "after one edited row ",
        "OK" if tampered_ok else f"BROKEN at seq {tampered_seq} ({tampered_field})",
    )

    payload = "\n".join(json.dumps(e, sort_keys=True) for e in log).encode()
    manifest = {
        "count": len(log),
        "chain_verified": tampered_ok,
        "sha256": hashlib.sha256(payload).hexdigest()[:24] + "…",
    }
    print()
    print("export manifest      ", manifest)
    print()
    print("The export does not hide the break. It reports it.")

    if "--test" in sys.argv:
        assert fresh_ok is True
        assert (tampered_ok, tampered_seq, tampered_field) == (False, 2, "row_hash")
        assert manifest["count"] == 3 and manifest["chain_verified"] is False
        print("\n[OK] All hash-chain tamper assertions passed.")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
