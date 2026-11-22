# The Five-Minute AI Vendor Procurement Checklist

> *"A control nobody can run is a claim. It is a fair test to put to any AI vendor, and a fairer one to put to yourself."*  
> Reference: [Slide 07](https://build-for-rbi-talk.vercel.app/#7) · [Slide 21](https://build-for-rbi-talk.vercel.app/#21)

Hand this 1-page checklist to any AI vendor (or internal platform team) before signing an SOW or moving to production in a regulated financial entity.

---

## Part 1 — The Five-Minute Runnable Bar

| # | Question to Ask the Vendor | Pass Criteria | Score |
| :-: | --- | --- | :-: |
| **1** | **Can I clone and run your governance controls locally on a laptop with zero network?** | Reproducible test suite (`go test ./...` / `pytest`) passes offline with no cloud credentials or external API keys | ☐ Pass ☐ Fail |
| **2** | **Is policy enforced at one pre-dispatch gate or scattered across agents?** | Single Composite Policy Gate evaluates every message before any agent executes | ☐ Pass ☐ Fail |
| **3** | **Can you prove per-message data residency without forking the deployment per country?** | Demo shows PII payload routed to disallowed region (e.g. `us`) refused before dispatch with audit reason | ☐ Pass ☐ Fail |
| **4** | **What happens to an in-flight request when a customer revokes consent?** | Consent checked against an append-only ledger; request fails closed and audit log records `consent` refusal | ☐ Pass ☐ Fail |
| **5** | **How do you catch silent approval drift after a refactor?** | Runtime check compares approved risk tier vs. current declared risk tier (`ErrRiskChanged`) | ☐ Pass ☐ Fail |
| **6** | **Can your own operator quietly edit the audit trail?** | SHA-256 hash-chained log + export manifest that explicitly reports `chain_verified: false` if any row is altered | ☐ Pass ☐ Fail |
| **7** | **Does your inventory read the live registry or a manually maintained spreadsheet?** | Live runtime registry exporter (`Rec 23` / `Annex IV`); shadow agents cannot execute | ☐ Pass ☐ Fail |
| **8** | **Does your red-team corpus run automatically on every PR?** | CI pipeline executes adversarial residency, RBAC, consent, and explainability tests on every commit | ☐ Pass ☐ Fail |

---

## Part 2 — Honest Scope & Contract Questions

| # | Scoping Check | Red Flag to Watch For | Score |
| :-: | --- | --- | :-: |
| **9** | **Did the proposal price 17 RBI FREE-AI RE controls or all 26?** | Pricing all 26 means you are paying for the regulator/SRO's own roadmap (`Recs 1, 3, 5, 7, 9–13`) | ☐ Pass ☐ Fail |
| **10** | **Does the SOW explicitly define Deployer vs. Provider under the EU AI Act?** | Fine-tuning on bank data or white-labeling under the bank's brand silently shifts provider liability to you | ☐ Pass ☐ Fail |
| **11** | **Can the vendor name at least two honest gaps in their current coverage?** | Vendors claiming "100% RBI + EU AI Act coverage" haven't tested against the June 2026 Draft MRM kill switch or 10-year retention rules | ☐ Pass ☐ Fail |
| **12** | **Can they emit both RBI Annexure VI and EU AI Act evidence from one control set?** | Duplicate compliance stacks double maintenance cost and diverge under audit | ☐ Pass ☐ Fail |

---

## Scoring Guide

- **11–12 Pass:** Production-grade engineering governance — proceed to technical deep-dive.
- **8–10 Pass:** Workable foundation — close missing controls as binding SOW milestones before go-live.
- **< 8 Pass:** Compliance by slide deck — require a runnable 5-minute proof before contracting.
