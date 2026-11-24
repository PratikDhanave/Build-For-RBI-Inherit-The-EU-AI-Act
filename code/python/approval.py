#!/usr/bin/env python3
"""
Lab 01 — Five Agents, One Bad Approval (RBI FREE-AI Rec 17).

Demonstrates silent governance drift: `var_calculator` was approved at medium
risk, then refactored to declare high risk. Nobody rejected or withdrew it —
the approval simply stopped describing the system it had approved.

Run:
    python3 labs/approval.py
    python3 labs/approval.py --test
"""

from datetime import datetime, timedelta
import sys

NOW = datetime(2026, 11, 28, 10, 0)
KNOWN_RISK_TIERS = {"low", "medium", "high"}

REGISTER = {
    "payment_orchestrator": dict(
        status="approved", risk="high", expires=NOW + timedelta(days=60)
    ),
    "kyc_orchestrator": dict(
        status="approved", risk="high", expires=NOW - timedelta(days=1)
    ),
    "aml_monitor": dict(status="rejected", risk="high"),
    "var_calculator": dict(status="approved", risk="medium"),  # signed off at medium
    "educator": None,  # never reviewed
}


def allow(agent_id: str, declares_now: str) -> tuple[bool, str]:
    risk = declares_now.strip().lower()
    if risk not in KNOWN_RISK_TIERS:
        return False, f"ErrUnknownRisk ({declares_now!r})"
    if risk != "high":
        return True, "low/medium — no approval required"
    rec = REGISTER.get(agent_id)
    if rec is None:
        return False, "ErrNoApproval"
    if rec["status"].strip().lower() != "approved":
        return False, f"ErrNotApproved ({rec['status']})"
    if rec.get("expires") and NOW > rec["expires"]:
        return False, "ErrExpired"
    if rec["risk"].strip().lower() != risk:
        return False, (
            f"ErrRiskChanged (approved as {rec['risk']}, now declares {risk})"
        )
    return True, "live"


TODAY = [
    ("payment_orchestrator", "high"),
    ("kyc_orchestrator", "high"),
    ("aml_monitor", "high"),
    ("var_calculator", "high"),
    ("educator", "low"),
]


def main() -> int:
    outcomes = {}
    for agent_id, risk in TODAY:
        ok, why = allow(agent_id, risk)
        outcomes[agent_id] = (ok, why)
        print(f"{'RUN ' if ok else 'STOP'}  {agent_id:<22} {why}")

    print()
    print("var_calculator is the one I would have missed.")

    if "--test" in sys.argv:
        assert outcomes["payment_orchestrator"][0] is True
        assert outcomes["kyc_orchestrator"] == (False, "ErrExpired")
        assert outcomes["aml_monitor"][0] is False and "ErrNotApproved" in outcomes["aml_monitor"][1]
        assert outcomes["var_calculator"][0] is False and "ErrRiskChanged" in outcomes["var_calculator"][1]
        assert outcomes["educator"][0] is True
        # Security regression checks: padded uppercase and unknown risk tiers fail closed
        assert allow("var_calculator", " HIGH ")[0] is False and "ErrRiskChanged" in allow("var_calculator", " HIGH ")[1]
        assert allow("educator", "critical")[0] is False and "ErrUnknownRisk" in allow("educator", "critical")[1]
        print("\n[OK] All approval-drift assertions passed.")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
