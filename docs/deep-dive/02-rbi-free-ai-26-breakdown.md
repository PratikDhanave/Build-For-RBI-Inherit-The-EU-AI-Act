# Deep Dive 02: All 26 RBI FREE-AI Recommendations — The 7 Sutras, 6 Pillars & Exact 9 / 17 Split

> **Companion to [Slide 06](https://build-for-rbi-talk.vercel.app/#6) & [Slide 07](https://build-for-rbi-talk.vercel.app/#7) (RBI FREE-AI · 9 / 17 Split)**  
> **Primary Source:** [RBI Report of the Committee to Develop a Framework for Responsible and Ethical Enablement of AI (FREE-AI), 13 August 2025 (Official PDF)](https://rbidocs.rbi.org.in/rdocs/PublicationReport/Pdfs/FREEAIR130820250A24FF2D4578453F824C72ED9F5D5851.PDF)

---

## 1. The Structural Map of RBI FREE-AI

RBI's FREE-AI committee report (chaired by **Dr. Pushpak Bhattacharyya, IIT Bombay**) is organized into three layers:

### The 7 Sutras (Guiding Principles)
1. **Trust is the Foundation** — Trust cannot be retrofitted after an AI failure in financial services.
2. **People First** — AI augments human judgment; meaningful human oversight remains essential for high-impact decisions.
3. **Innovation over Restraint** — Enable responsible adoption rather than blanket bans.
4. **Fairness and Equity** — Prevent discriminatory outcomes across credit, pricing, and access.
5. **Accountability** — Regulated Entities remain accountable for every model they deploy, including third-party vendor models.
6. **Understandable by Design** — Decisions affecting customers must carry clear, intelligible explanations.
7. **Safety, Resilience and Sustainability** — Systems must withstand adversarial inputs, drift, and operational shocks.

### The 6 Pillars
- **Innovation Enablement (3 Pillars · Recs 1–13):** Infrastructure · Policy · Capacity
- **Risk Mitigation (3 Pillars · Recs 14–26):** Governance · Protection · Assurance

---

## 2. Complete 26-Recommendation Ownership Table (9 vs. 17)

When scoping a regulated AI engineering engagement in India, **never quote for all 26 recommendations**. Nine belong to RBI, government bodies, SROs, or non-codebase ecosystem/training tracks. Seventeen represent the Regulated Entity's engineering & compliance scope (aligned 1:1 with [`RBI-FREE-AI-17-MATRIX.md`](../RBI-FREE-AI-17-MATRIX.md) and [`crosswalk/rbi-eu-ai-act.json`](../../crosswalk/rbi-eu-ai-act.json)).

| Rec # | Owner (Engineering SOW) | Include in RE SOW? | Official RBI FREE-AI Title | Engineering Control / Notes |
| :-: | --- | :-: | --- | --- |
| **1** | Regulator / Ecosystem | ❌ No | **Financial Sector Data Infrastructure** | Sector DPI integrated with *AI Kosh* under IndiaAI Mission |
| **2** | **Regulated Entity** | ✅ **Yes** | **AI Innovation Sandbox** | Offline runnable sandbox pipeline (`make demo` / [`cmd/demo`](../../code/cmd/demo/main.go)) |
| **3** | Regulator / Ecosystem | ❌ No | **Incentives and Funding Support** | RBI/sector shared compute & data funding incentives |
| **4** | **Regulated Entity** | ✅ **Yes** | **Indigenous Financial Sector-Specific AI Models** | Per-message `ResidencyPolicy` (`in`/`adgm`/`difc`/`sama`) + on-prem default ([`composite_gate.go`](../../code/gate/composite_gate.go)) |
| **5** | Regulator / Ecosystem | ❌ No | **Integrating AI with DPI** | Sector-level framework for AI + Digital Public Infrastructure |
| **6** | **Regulated Entity** | ✅ **Yes** | **Adaptive and Enabling Policies** | Declarative YAML policy config + fail-closed loader ([`policy.yaml`](../../code/config/policy.yaml), [`loader.go`](../../code/config/loader.go)) |
| **7** | Regulator / Ecosystem | ❌ No | **Enabling AI-Based Affirmative Action** | Regulators lowering compliance friction for financial inclusion |
| **8** | **Regulated Entity** | ✅ **Yes** | **AI Liability Framework (Graded Liability)** | Severity grading & auto-incident capture to qualify for RBI's tolerant first-offense stance ([`annexure_vi.go`](../../code/incident/annexure_vi.go)) |
| **9** | Regulator / SRO | ❌ No | **AI Institutional Framework** | Permanent multi-stakeholder Standing AI Committee under RBI |
| **10** | Organizational / HR *(Non-Codebase)* | ❌ No | **Capacity Building within REs** | Board, C-suite & workforce upskilling programs *(RE-owned organizational/HR track, not a software control)* |
| **11** | Regulator / SRO | ❌ No | **Capacity Building for Regulators and Supervisors** | Supervisory AI training & dedicated sector AI institute |
| **12** | SRO / Industry | ❌ No | **Framework for Sharing Best Practices** | SRO / IBA sector-wide use-case & lessons sharing |
| **13** | Regulator / Industry | ❌ No | **Recognise and Reward Responsible AI Innovation** | Sectoral recognition & awards for ethical financial AI |
| **14** | **Regulated Entity** | ✅ **Yes** | **Board-Approved AI Policy (Annexure V)** | Versioned [`policy.yaml`](../../code/config/policy.yaml) + [`loader.go`](../../code/config/loader.go) |
| **15** | **Regulated Entity** | ✅ **Yes** | **Data Lifecycle Governance** | Per-message residency + append-only [`consent/ledger.go`](../../code/consent/ledger.go) *(aligned with DPDP Rules Consent Managers, live 13 Nov 2026)* |
| **16** | **Regulated Entity** | ✅ **Yes** | **AI System Governance + Autonomous Systems** | Pre-dispatch risk tiering, `RBACPolicy`, caps & circuit breakers |
| **17** | **Regulated Entity** | ✅ **Yes** | **Product Approval Process** | Runtime [`approval_gate.go`](../../code/gate/approval_gate.go) (`ErrRiskChanged`, `ErrExpired`) |
| **18** | **Regulated Entity** | ✅ **Yes** | **Consumer Protection** | First-event [`disclosure/disclosure.go`](../../code/disclosure/disclosure.go) + `ExplainabilityPolicy` + human escalation readiness |
| **19** | **Regulated Entity** | ✅ **Yes** | **Cybersecurity Measures (AI-Specific)** | Pre-dispatch RBAC, registry enforcement, input normalization & fail-closed audit ([`composite_gate.go`](../../code/gate/composite_gate.go)) |
| **20** | **Regulated Entity** | ✅ **Yes** | **Red-Teaming** | 20-case OWASP ASI-tagged PR red-team suite ([`code/redteam/`](../../code/redteam)) |
| **21** | **Regulated Entity** | ✅ **Yes** | **Business Continuity Plan (BCP) for AI** | Deterministic fallbacks & failover drills *(operator kill switch tracked in [`RFC-001`](../open-problems/RFC-001-ai-kill-switch.md))* |
| **22** | **Regulated Entity** | ✅ **Yes** | **AI Incident Reporting (Annexure VI)** | Structured non-punitive good-faith reporting ([`incident/annexure_vi.go`](../../code/incident/annexure_vi.go)) + SHA-256 [`audit/chain.go`](../../code/audit/chain.go) |
| **23** | **Regulated Entity** | ✅ **Yes** | **AI Inventory** | Zero-drift [`registry/registry.go`](../../code/registry/registry.go) *(updated at least half-yearly per RBI PDF; 10-yr gap in [`RFC-002`](../open-problems/RFC-002-ten-year-retention.md))* |
| **24** | **Regulated Entity** | ✅ **Yes** | **AI Audit Framework** | SHA-256 hash chain + caller-pinned Ed25519 signed manifest + offline `shasum` verification ([`audit/chain.go`](../../code/audit/chain.go)) |
| **25** | **Regulated Entity** | ✅ **Yes** | **Disclosures by REs** | First-event disclosure ([`disclosure/disclosure.go`](../../code/disclosure/disclosure.go)) + auditable refusal traces |
| **26** | **Regulated Entity** | ✅ **Yes** | **AI Compliance Toolkit** | Dual RBI + EU AI Act evidence pack exporter ([`cmd/evidencepack`](../../code/cmd/evidencepack/main.go)) + machine-readable [`crosswalk/`](../../crosswalk) *(SRO/industry-curated toolkit adopted by REs)* |

> **Scoping Nuance for Q&A:** All **13 Risk Mitigation recommendations (`Recs 14–26`)** are direct RE operational obligations (with SRO toolkit support on Rec 26). Under Innovation Enablement (`Recs 1–13`), **Recs 2, 4, 6, 8** are regulator-enabled tracks where REs implement internal technical readiness (offline sandbox, residency/local models, adaptive YAML policies, graded incident capture), whereas **Rec 10** (*Capacity Building within REs*) is an RE-owned HR/Board training program outside software engineering scope.

---

## 3. Advisory vs. Binding: How FREE-AI Relates to RBI Draft MRM Guidance

FREE-AI (**13 August 2025**) is a committee report — extraordinarily clear architectural advice, but not a binding master direction. The binding instrument is the **RBI Draft Model Risk Management Guidance** (**24 June 2026**, comments closed **24 July 2026**, final circular pending). Building for FREE-AI's 17 RE recommendations early automatically satisfied **4 of the 6** core Draft MRM requirements and exposed the exact **2 gaps** (`RFC-001` Kill Switch, `RFC-002` 10-Year Retention) before the circular lands.
