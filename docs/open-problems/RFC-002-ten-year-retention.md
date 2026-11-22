# RFC-002: Ten-Year Decommissioned Model Retention vs. Zero-Drift Live Inventory

- **Status:** Open Architectural Conflict (`GAP 02`)
- **Regulatory Trigger:** RBI Draft Model Risk Management Guidance (24 June 2026) — 10-Year Inventory Retention (including decommissioned models)
- **Related Slides:** [Slide 18](https://build-for-rbi-talk.vercel.app/#18), [Slide 19](https://build-for-rbi-talk.vercel.app/#19)

---

## 1. The Architectural Conflict

Today, the model & agent inventory is generated **directly from the live runtime registry**:
- **Guarantee:** Because the inventory *is* the live registry, it **cannot drift** from what is actually running in production. There are no shadow agents outside the register.
- **Conflict:** The RBI Draft MRM Guidance requires keeping the full model inventory — **including decommissioned models** — for **10 years**.

> **A live inventory cannot drift, and cannot remember.**

This is not a backlog ticket; it is a genuine tension between two competing invariants:
1. **Zero-drift truth:** Inventory reflects *only* live, executable code.
2. **Ten-year historical memory:** Inventory preserves *retired* models, weights/prompt versions, approvals, risk tiers, and decommission dates for a decade.

## 2. Candidate Designs Under Evaluation

| Approach | How It Works | Strengths | Trade-Offs |
| --- | --- | --- | --- |
| **A. Append-Only Lifecycle Event Log + Live Projection** | Every registration, risk change, and decommission emits a hash-chained event; live registry is a materialized view (`status=ACTIVE`), 10-year archive queries full history | Preserves zero-drift for active set while keeping immutable 10-year lineage | Requires strict reconciliation so no code can boot without emitting a registration event |
| **B. Signed Periodic Snapshots (WORM Archive)** | Daily signed export of the live registry written to WORM object storage (10-year retention lock) | Simple, zero coupling to hot path, verifiable hash chain | Misses ephemeral models deployed and retired between snapshot windows |
| **C. Tombstone Records Inside Registry** | Decommissioned agents remain in registry with `state=DECOMMISSIONED` and immutable metadata | Single source of truth for auditors | Bloats runtime registry over 10 years; must strictly ensure tombstoned entries can never be dispatched |

## 3. Desired Invariants for Any Solution

1. **No shadow execution:** An unregistered model must still be impossible to invoke.
2. **Tamper-evident history:** Decommissioned entries must be cryptographically bound (consistent with Lab 02's SHA-256 hash chain).
3. **Offline verifiability:** An auditor must still be able to run a 5-minute offline check across both live and decommissioned records.

> If you have solved this cleanly in a regulated bank or financial platform, [open a proposal issue](../../issues/new?template=open-problem-proposal.md) referencing `RFC-002`.
