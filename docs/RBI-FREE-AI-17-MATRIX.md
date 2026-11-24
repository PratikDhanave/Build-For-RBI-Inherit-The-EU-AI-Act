# RBI FREE-AI — Complete 17-Recommendation Traceability Matrix

RBI's [**FREE-AI Framework** (13 August 2025 — Official RBI PDF)](https://rbidocs.rbi.org.in/rdocs/PublicationReport/Pdfs/FREEAIR130820250A24FF2D4578453F824C72ED9F5D5851.PDF) contains **26 recommendations**:
- **9 belong to the Regulator & SRO:** `Recs 1, 3, 5, 7, 9, 10, 11, 12, 13`
- **17 belong to the Regulated Entity (your engineering scope):** `Recs 2, 4, 6, 8, 14–26`

This matrix traces all **17 Regulated Entity recommendations** to concrete architectural controls, runnable artefacts in this repository, EU AI Act convergence, and RBI Draft MRM alignment.

| Rec | Regulated Entity Scope | Control Pattern Implemented | Artefact / Verification in Repo | EU AI Act Convergence |
| :-: | --- | --- | --- | --- |
| **2** | Foundational AI adoption & governance | Board-approved governance baseline + single pre-dispatch bus | [`code/config/policy.yaml`](../code/config/policy.yaml) | Art 17 |
| **4** | Data governance, lineage & residency | Per-message residency policy (`in`, `adgm`, `difc`, `sama`) | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) (`ResidencyPolicy`) | Art 10 |
| **6** | Internal testing & sandbox readiness | Zero-network offline deterministic verification suite | [`code/scripts/offline_verify.sh`](../code/scripts/offline_verify.sh) | Art 15 |
| **8** | AI literacy & operational readiness | Runnable interactive labs (`approval.py`, `audit.py`) + RFCs | [`labs/`](../labs) · [`docs/open-problems/`](./open-problems) | Art 4 |
| **14** | Board-approved AI policy | Versioned declarative policy-as-config loaded at runtime | [`code/config/policy.yaml`](../code/config/policy.yaml) · [`loader.go`](../code/config/loader.go) | **Article 17** |
| **15** | Privacy & consent management | Append-only consent ledger; in-flight failure on revocation | [`code/consent/ledger.go`](../code/consent/ledger.go) · [`ConsentPolicy`](../code/gate/composite_gate.go) | Art 10 |
| **16** | Risk classification & tiering | Risk class checked pre-dispatch with RBAC role guard | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) (`RBACPolicy`) | **Articles 9, 14** |
| **17** | Approval lifecycle & drift prevention | Runtime risk-drift guard (`ErrRiskChanged`, `ErrExpired`) | [`code/gate/approval_gate.go`](../code/gate/approval_gate.go) · [`code/python/approval.py`](../code/python/approval.py) | Articles 9, 14 |
| **18** | Consumer transparency & disclosure | Mandatory AI disclosure emitted on first event | [`code/disclosure/disclosure.go`](../code/disclosure/disclosure.go) (`Tracker`) | **Article 50** *(law Aug 2026)* |
| **19** | Explainability & rationale | Rationale required on all advice/credit/loan recommendations | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) (`ExplainabilityPolicy`) | Art 13, 86 |
| **20** | Adversarial testing & red-teaming | 20-case OWASP ASI-tagged red-team corpus executed automatically on every PR | [`code/redteam/corpus.json`](../code/redteam/corpus.json) · [`redteam_test.go`](../code/redteam/redteam_test.go) | **Article 15** |
| **21** | Third-party / vendor accountability | Vendor model metadata tracked in live registry + 5-min test | [`docs/VENDOR-CHECKLIST.md`](./VENDOR-CHECKLIST.md) · [`examples/inventory-export.json`](../examples/inventory-export.json) | Art 25 |
| **22** | Auditability & incident reporting | SHA-256 hash chain + pinned-key Ed25519 signed manifest + RBI Annexure VI | [`code/audit/chain.go`](../code/audit/chain.go) · [`code/incident/annexure_vi.go`](../code/incident/annexure_vi.go) | **Articles 12, 73** |
| **23** | Full AI/model inventory | Live runtime registry exporter + pre-dispatch `RegistryPolicy` (see `RFC-002` for 10-yr gap) | [`code/registry/registry.go`](../code/registry/registry.go) · [`examples/inventory-export.json`](../examples/inventory-export.json) | **Annex IV** |
| **24** | Continuous monitoring & fallbacks | Circuit breakers, caps & fallbacks (see `RFC-001` for kill switch gap) | [`docs/open-problems/RFC-001-ai-kill-switch.md`](./open-problems/RFC-001-ai-kill-switch.md) | Art 72 |
| **25** | Customer grievance & first-event notice | First-event disclosure + auditable refusal traces for disputes | [`code/audit/chain.go`](../code/audit/chain.go) | **Article 50** |
| **26** | Independent assurance & evidence packs | Dual RBI + EU AI Act evidence pack generator + machine-readable cross-walk | [`code/cmd/evidencepack/main.go`](../code/cmd/evidencepack/main.go) · [`crosswalk/rbi-eu-ai-act.json`](../crosswalk/rbi-eu-ai-act.json) · `make test` | Art 11, 17 |

---

### Machine-Readable Format
Download [`crosswalk/rbi-eu-ai-act.json`](../crosswalk/rbi-eu-ai-act.json) or [`crosswalk/rbi-eu-ai-act.csv`](../crosswalk/rbi-eu-ai-act.csv) for direct import into Jira, Confluence, or GRC tooling.
