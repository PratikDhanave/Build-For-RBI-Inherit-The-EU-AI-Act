#!/usr/bin/env python3
"""
Composite Policy Gate Demo — One Gate, Every Message, Before Dispatch.

Demonstrates how a single pre-dispatch gate enforces residency, explainability,
RBAC, and append-only consent policies across 62 registered agents without
duplicating compliance logic inside individual agent code.

Run:
    python3 labs/gate_demo.py
    python3 labs/gate_demo.py --test
"""

from dataclasses import dataclass
import sys


@dataclass(frozen=True)
class Message:
    label: str
    actor_role: str
    target_risk: str
    contains_pii: bool = False
    target_region: str = "in"
    is_advice: bool = False
    has_rationale: bool = True
    consent_live: bool = True


ALLOWED_PII_REGIONS = {"in", "adgm", "difc", "sama"}


def evaluate_gate(msg: Message) -> tuple[bool, str]:
    """Evaluate ordered policies before any agent sees the message."""
    # Policy 1: Per-message data residency
    if msg.contains_pii and msg.target_region not in ALLOWED_PII_REGIONS:
        return False, f"[residency] pii bound for region {msg.target_region}"

    # Policy 2: Explainability requirement for advice/recommendations
    if msg.is_advice and not msg.has_rationale:
        return False, "[explainability] recommendation with no rationale"

    # Policy 3: Role-based access control for high-risk agents
    if msg.target_risk == "high" and msg.actor_role not in {"advisor", "operator"}:
        return False, "[rbac] role may not address a high-risk agent"

    # Policy 4: Append-only live consent verification
    if msg.contains_pii and not msg.consent_live:
        return False, "[consent] no live consent for this category"

    return True, "dispatched to agent"


MESSAGES = [
    Message("ordinary question", actor_role="user", target_risk="low"),
    Message(
        "PAN to a US model",
        actor_role="advisor",
        target_risk="medium",
        contains_pii=True,
        target_region="us",
    ),
    Message(
        "advice with no reason",
        actor_role="advisor",
        target_risk="medium",
        is_advice=True,
        has_rationale=False,
    ),
    Message(
        "user hits a high agent",
        actor_role="user",
        target_risk="high",
    ),
    Message(
        "PII, consent revoked",
        actor_role="advisor",
        target_risk="medium",
        contains_pii=True,
        target_region="in",
        consent_live=False,
    ),
    Message("balance lookup", actor_role="user", target_risk="low"),
]


def main() -> int:
    results = []
    for msg in MESSAGES:
        allowed, reason = evaluate_gate(msg)
        results.append((msg.label, allowed, reason))
        verdict = "ALLOW" if allowed else "DENY "
        detail = "" if allowed else f" {reason}"
        print(f"{msg.label:<26} {verdict}{detail}")

    print()
    print("Nothing above is agent code. The gate runs before any agent sees it.")

    if "--test" in sys.argv:
        expected = {
            "ordinary question": True,
            "PAN to a US model": False,
            "advice with no reason": False,
            "user hits a high agent": False,
            "PII, consent revoked": False,
            "balance lookup": True,
        }
        for label, allowed, _ in results:
            assert allowed == expected[label], f"Unexpected verdict for {label}: {allowed}"
        print("\n[OK] All composite policy gate assertions passed.")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
