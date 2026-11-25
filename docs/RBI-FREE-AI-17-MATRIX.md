# RBI FREE-AI — Complete 17-Recommendation Traceability Matrix

RBI's [**FREE-AI Framework** (13 August 2025 — Official RBI PDF)](https://rbidocs.rbi.org.in/rdocs/PublicationReport/Pdfs/FREEAIR130820250A24FF2D4578453F824C72ED9F5D5851.PDF) contains **26 recommendations**:
- **9 belong to the Regulator & SRO:** `Recs 1, 3, 5, 7, 9, 10, 11, 12, 13`
- **17 belong to the Regulated Entity (your engineering scope):** `Recs 2, 4, 6, 8, 14–26`

This matrix uses the **exact official RBI FREE-AI recommendation titles** (aligned 1:1 with `genie/docs/free-ai-mapping.md`) and traces all **17 Regulated Entity recommendations** to runnable artefacts in this companion repo and EU AI Act convergence.

| Rec | Official RBI FREE-AI Title (RE Scope) | Control Pattern Implemented | Artefact / Verification in Repo | EU AI Act Convergence |
| :-: | --- | --- | --- | --- |
| **2** | **AI Innovation Sandbox** | Zero-network offline runnable pipeline (`make demo` / `go run ./cmd/demo`) | [`code/scripts/offline_verify.sh`](../code/scripts/offline_verify.sh) · [`code/cmd/demo/main.go`](../code/cmd/demo/main.go) | Art 57 |
| **4** | **Indigenous AI Models** | Per-message region residency policy (`in`, `adgm`, `difc`, `sama`) + on-prem default | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) (`ResidencyPolicy`) | Art 10 |
| **6** | **Adaptive Policies** | Declarative YAML policy config + fail-closed loader so risk edits config, not code | [`code/config/policy.yaml`](../code/config/policy.yaml) · [`loader.go`](../code/config/loader.go) | Art 17 |
| **8** | **Graded Liability** | Severity grading & automatic Annexure VI incident capture on policy deny | [`code/incident/annexure_vi.go`](../code/incident/annexure_vi.go) · [`docs/VENDOR-CHECKLIST.md`](./VENDOR-CHECKLIST.md) | Art 25, 73 |
| **14** | **Board-Approved AI Policy (Annexure V)** | Versioned, board-approved `policy.yaml` loaded and validated at runtime | [`code/config/policy.yaml`](../code/config/policy.yaml) · [`loader.go`](../code/config/loader.go) | **Article 17** |
| **15** | **Data Lifecycle Governance** | Per-message residency + append-only consent ledger (in-flight revocation proof) | [`code/consent/ledger.go`](../code/consent/ledger.go) · [`ConsentPolicy`](../code/gate/composite_gate.go) | Art 10 |
| **16** | **AI System Governance + Autonomous Systems** | Pre-dispatch risk tiering, RBAC guard, caps & circuit breakers | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) (`RBACPolicy`) | **Articles 9, 14** |
| **17** | **Product Approval** | Approval register + runtime risk-drift guard (`ErrRiskChanged`, `ErrExpired`) | [`code/gate/approval_gate.go`](../code/gate/approval_gate.go) · [`code/python/approval.py`](../code/python/approval.py) | Articles 9, 14 |
| **18** | **Consumer Protection** | Mandatory first-event AI disclosure + non-empty rationale on advice | [`code/disclosure/disclosure.go`](../code/disclosure/disclosure.go) · [`ExplainabilityPolicy`](../code/gate/composite_gate.go) | **Article 50** *(law Aug 2026)* |
| **19** | **Cybersecurity** | Pre-dispatch RBAC, registry enforcement, input normalization & fail-closed audit | [`code/gate/composite_gate.go`](../code/gate/composite_gate.go) · [`docs/deep-dive/01-composite-policy-gate-architecture.md`](./deep-dive/01-composite-policy-gate-architecture.md) | Art 15 |
| **20** | **Red Teaming** | 20-case OWASP ASI-tagged adversarial corpus executed on every PR | [`code/redteam/corpus.json`](../code/redteam/corpus.json) · [`redteam_test.go`](../code/redteam/redteam_test.go) | **Article 15** |
| **21** | **BCP for AI** | Deterministic fallbacks & automated failover drills *(operator kill switch open in `RFC-001`)* | [`docs/open-problems/RFC-001-ai-kill-switch.md`](./open-problems/RFC-001-ai-kill-switch.md) | Art 15, 72 |
| **22** | **AI Incident Reporting (Annexure VI)** | Structured Annexure VI incident report + tamper-evident SHA-256 audit log | [`code/incident/annexure_vi.go`](../code/incident/annexure_vi.go) · [`code/audit/chain.go`](../code/audit/chain.go) | **Articles 12, 73** |
| **23** | **AI Inventory** | Zero-drift inventory from live runtime registry + pre-dispatch `RegistryPolicy` *(10-yr gap in `RFC-002`)* | [`code/registry/registry.go`](../code/registry/registry.go) · [`examples/inventory-export.json`](../examples/inventory-export.json) | **Annex IV** |
| **24** | **AI Audit Framework** | SHA-256 hash chain + pinned-key Ed25519 signed manifest + offline `shasum` parity | [`code/audit/chain.go`](../code/audit/chain.go) · [`code/python/audit.py`](../code/python/audit.py) | **Article 12** |
| **25** | **Disclosures** | First-event AI disclosure tracker (`Art. 50`) & auditable refusal traces | [`code/disclosure/disclosure.go`](../code/disclosure/disclosure.go) | **Article 50** |
| **26** | **AI Toolkit** | Dual RBI + EU AI Act evidence pack exporter + machine-readable cross-walk | [`code/cmd/evidencepack/main.go`](../code/cmd/evidencepack/main.go) · [`crosswalk/rbi-eu-ai-act.json`](../crosswalk/rbi-eu-ai-act.json) · `make test` | Art 11, 17 |

---

### Machine-Readable Format
Download [`crosswalk/rbi-eu-ai-act.json`](../crosswalk/rbi-eu-ai-act.json) or [`crosswalk/rbi-eu-ai-act.csv`](../crosswalk/rbi-eu-ai-act.csv) for direct import into Jira, Confluence, or GRC tooling.
