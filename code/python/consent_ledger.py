#!/usr/bin/env python3
"""
Lab 03 — Append-Only Consent Ledger (Slide 03: "Revoke on Tuesday, fail in-flight").

Demonstrates an append-only consent ledger where a customer grants PII consent
on Monday and revokes it on Tuesday. An in-flight PII request fails closed, and
the immutable ledger preserves the complete GRANT -> REVOKE history so auditors
can prove why the request failed.

Run:
    python3 code/python/consent_ledger.py
    python3 code/python/consent_ledger.py --test
"""

from dataclasses import dataclass
import sys


@dataclass(frozen=True)
class ConsentEvent:
    seq: int
    day: str
    subject: str
    category: str
    action: str  # "GRANT" | "REVOKE"
    reason: str


class AppendOnlyConsentLedger:
    def __init__(self) -> None:
        self._events: list[ConsentEvent] = []

    def grant(self, subject: str, category: str, day: str, reason: str) -> ConsentEvent:
        ev = ConsentEvent(len(self._events) + 1, day, subject, category, "GRANT", reason)
        self._events.append(ev)
        return ev

    def revoke(self, subject: str, category: str, day: str, reason: str) -> ConsentEvent:
        ev = ConsentEvent(len(self._events) + 1, day, subject, category, "REVOKE", reason)
        self._events.append(ev)
        return ev

    def check(self, subject: str, category: str) -> tuple[bool, str]:
        for ev in reversed(self._events):
            if ev.subject == subject and ev.category == category:
                if ev.action == "GRANT":
                    return True, f"live consent (seq {ev.seq}, {ev.day})"
                return False, f"DENY [consent] revoked at seq {ev.seq} ({ev.day}: {ev.reason})"
        return False, "DENY [consent] no consent record"

    def history(self, subject: str) -> list[ConsentEvent]:
        return [e for e in self._events if e.subject == subject]


def main() -> int:
    ledger = AppendOnlyConsentLedger()
    ledger.grant("user-1", "pii", "Monday 09:00", "Open Banking onboarding opt-in")

    mon_ok, mon_why = ledger.check("user-1", "pii")
    print(f"Monday  in-flight PII request : {'ALLOW' if mon_ok else 'DENY '} — {mon_why}")

    ledger.revoke("user-1", "pii", "Tuesday 14:30", "customer revoked on Tuesday")

    tue_ok, tue_why = ledger.check("user-1", "pii")
    print(f"Tuesday in-flight PII request : {'ALLOW' if tue_ok else 'DENY '} — {tue_why}")

    print("\nAppend-only history for user-1:")
    for ev in ledger.history("user-1"):
        print(f"  seq {ev.seq} | {ev.day:<14} | {ev.action:<6} | {ev.reason}")

    if "--test" in sys.argv:
        assert mon_ok is True
        assert tue_ok is False and "revoked at seq 2" in tue_why
        assert len(ledger.history("user-1")) == 2
        print("\n[OK] Append-only consent ledger assertions passed.")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
