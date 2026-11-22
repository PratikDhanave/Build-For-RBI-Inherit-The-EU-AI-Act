# Deep Dive 02: All 26 RBI FREE-AI Recommendations — The 7 Sutras, 6 Pillars & Exact 9 / 17 Split

> **Companion to [Slide 04](https://build-for-rbi-talk.vercel.app/#4) & [Slide 05](https://build-for-rbi-talk.vercel.app/#5)**  
> **Primary Source:** [RBI Report of the Committee to Develop a Framework for Responsible and Ethical Enablement of AI (FREE-AI), 13 August 2025 (Official PDF)](https://rbidocs.rbi.org.in/rdocs/PublicationReport/Pdfs/FREEAIR130820250A24FF2D4578453F824C72ED9F5D5851.PDF)

---

## 1. The Structural Map of RBI FREE-AI

RBI's FREE-AI committee report is organized into three layers:

### The 7 Sutras (Guiding Principles)
1. **Trust is the Foundation** — Trust cannot be retrofitted after an AI failure in financial services.
2. **People First** — AI augments human judgment; meaningful human oversight remains essential for high-impact decisions.
3. **Innovation over Restraint** — Enable responsible adoption rather than blanket bans.
4. **Fairness and Equity** — Prevent discriminatory outcomes across credit, pricing, and access.
5. **Accountability** — Regulated Entities remain accountable for every model they deploy, including third-party vendor models.
6. **Understandable by Design** — Decisions affecting customers must carry clear, intelligible explanations.
7. **Safety, Resilience and Sustainability** — Systems must withstand adversarial inputs, drift, and operational shocks.

### The 6 Pillars
- **Innovation Enablement (3 Pillars):** Infrastructure · Policy · Capacity
- **Risk Mitigation (3 Pillars):** Governance · Protection · Assurance

---

## 2. Complete 26-Recommendation Ownership Table (9 vs. 17)

When scoping a regulated AI engagement in India, **never quote for all 26 recommendations**. Nine belong to RBI, government bodies, or Self-Regulatory Organizations (SROs). Seventeen are the Regulated Entity's engineering and governance scope.

| Rec # | Owner | Include in RE SOW? | Summary Scope | Engineering Control / Notes |
| :-: | --- | :-: | --- | --- |
| **1** | Regulator / SRO | ❌ No | Shared AI compute & public digital infrastructure | Sector-level enablement |
| **2** | **Regulated Entity** | ✅ **Yes** | Foundational AI adoption & internal governance | Board policy baseline ([`policy.yaml`](../../code/config/policy.yaml)) |
| **3** | Regulator / SRO | ❌ No | Indigenous financial datasets & DPI integration | Sector-level dataset initiatives |
| **4** | **Regulated Entity** | ✅ **Yes** | Data governance, lineage & residency | Per-message `ResidencyPolicy` ([`composite_gate.go`](../../code/gate/composite_gate.go)) |
| **5** | Regulator / SRO | ❌ No | Regulatory AI Innovation Sandbox | RBI-operated sandbox |
| **6** | **Regulated Entity** | ✅ **Yes** | Internal testing, validation & pre-prod sandbox | Offline deterministic verification (`make test`) |
| **7** | Regulator / SRO | ❌ No | AI Centre of Excellence & ecosystem capacity | Sector-level institution |
| **8** | **Regulated Entity** | ✅ **Yes** | Staff AI literacy & operational readiness | Interactive labs (`code/python/`) + RFC playbooks |
| **9** | Regulator / SRO | ❌ No | Inter-regulatory coordination (RBI/SEBI/IRDAI/MeitY) | Inter-agency coordination |
| **10** | Regulator / SRO | ❌ No | Sectoral code of conduct | SRO deliverable |
| **11** | Regulator / SRO | ❌ No | Supervisory AI toolkit (SupTech) | Regulator supervisory tooling |
| **12** | Regulator / SRO | ❌ No | Sector-wide AI incident intelligence sharing | Centralised sectoral sharing |
| **13** | Regulator / SRO | ❌ No | Periodic regulatory framework review | RBI policy review cycle |
| **14** | **Regulated Entity** | ✅ **Yes** | Board-approved AI policy | Versioned [`policy.yaml`](../../code/config/policy.yaml) + [`loader.go`](../../code/config/loader.go) |
| **15** | **Regulated Entity** | ✅ **Yes** | Privacy, per-message residency & consent | Append-only [`consent/ledger.go`](../../code/consent/ledger.go) |
| **16** | **Regulated Entity** | ✅ **Yes** | Risk classification & tiering | Pre-dispatch `RBACPolicy` & risk checks |
| **17** | **Regulated Entity** | ✅ **Yes** | Approval lifecycle & drift detection | Runtime [`approval_gate.go`](../../code/gate/approval_gate.go) (`ErrRiskChanged`) |
| **18** | **Regulated Entity** | ✅ **Yes** | Consumer transparency & AI disclosure | First-event [`disclosure/disclosure.go`](../../code/disclosure/disclosure.go) |
| **19** | **Regulated Entity** | ✅ **Yes** | Explainability & mandatory rationale | `ExplainabilityPolicy` in Composite Gate |
| **20** | **Regulated Entity** | ✅ **Yes** | Adversarial red-teaming & robustness | 16-case PR red-team suite ([`code/redteam/`](../../code/redteam)) |
| **21** | **Regulated Entity** | ✅ **Yes** | Vendor / third-party model accountability | Live vendor metadata + [5-Min Vendor Checklist](../VENDOR-CHECKLIST.md) |
| **22** | **Regulated Entity** | ✅ **Yes** | Audit trails & Annexure VI incident reporting | [`audit/chain.go`](../../code/audit/chain.go) + [`incident/annexure_vi.go`](../../code/incident/annexure_vi.go) |
| **23** | **Regulated Entity** | ✅ **Yes** | Full AI/model inventory | Zero-drift [`registry/registry.go`](../../code/registry/registry.go) *(see `RFC-002`)* |
| **24** | **Regulated Entity** | ✅ **Yes** | Continuous monitoring, caps & fallbacks | Circuit breakers & fallbacks *(see `RFC-001`)* |
| **25** | **Regulated Entity** | ✅ **Yes** | Customer grievance redressal & disclosure | First-event disclosure + auditable refusal traces |
| **26** | **Regulated Entity** | ✅ **Yes** | Independent assurance & verifiable evidence | Machine-readable [`crosswalk/`](../../crosswalk) + `offline_verify.sh` |

---

## 3. Advisory vs. Binding: How FREE-AI Relates to RBI Draft MRM Guidance

FREE-AI (**13 August 2025**) is a committee report — extraordinarily clear architectural advice, but not a binding master direction. The binding instrument is the **RBI Draft Model Risk Management Guidance** (**24 June 2026**, comments closed **24 July 2026**, final circular pending). Building for FREE-AI's 17 RE recommendations early automatically satisfied **4 of the 6** core Draft MRM requirements and exposed the exact **2 gaps** (`RFC-001` Kill Switch, `RFC-002` 10-Year Retention) before the circular lands.
