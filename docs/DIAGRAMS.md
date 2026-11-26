# Visual Architecture Atlas — Regulated AI Governance Diagrams

[![Live Web Deck](https://img.shields.io/badge/Live_Deck-build--for--rbi--talk.vercel.app-4285F4?style=for-the-badge&logo=vercel&logoColor=FFFFFF)](https://build-for-rbi-talk.vercel.app)
[![Diagrams](https://img.shields.io/badge/Mermaid_Diagrams-25_Verified-34A853?style=for-the-badge)](#table-of-diagrams)
[![Palette](https://img.shields.io/badge/Style-Google_Cloud_Palette-F9AB00?style=for-the-badge)](#color-legend)

> Visual companion to [**Build for RBI, Inherit the EU AI Act**](../README.md). Every diagram uses the **Google Cloud color palette** and **Google Sans / Roboto / Roboto Mono** typography, verified against Mermaid v11.

### Color Legend
- **Blue (`#4285F4`)** — Policy Gate, Core Infrastructure, RBI Controls, Runtime Modules
- **Green (`#34A853`)** — Allowed Flows, Verified States, Regulated Entity Scope, Built Controls
- **Yellow/Amber (`#F9AB00`)** — Reviews, Regulator/SRO Scope, Drift Warnings, Pending Milestones
- **Red (`#EA4335`)** — Blocked/Denied Flows, Tamper Alerts, EU AI Act Obligations, Open Gaps

---

## Table of Diagrams

1. [End-to-End System Context Architecture](#1-end-to-end-system-context-architecture)
2. [RBI FREE-AI Complete 3-Layer Taxonomy (7 Sutras · 6 Pillars · 9/17 Split)](#2-rbi-free-ai-complete-3-layer-taxonomy)
3. [Agent Governance Approval Lifecycle State Machine](#3-agent-governance-approval-lifecycle-state-machine)
4. [Append-Only Consent Ledger & In-Flight Revocation Sequence](#4-append-only-consent-ledger--in-flight-revocation-sequence)
5. [Cryptographic Audit Chain, Ed25519 Manifest & TSA Anchoring](#5-cryptographic-audit-chain-ed25519-manifest--tsa-anchoring)
6. [Governance & Compliance Data Model (ER Diagram)](#6-governance--compliance-data-model-er-diagram)
7. [Go Governance Packages & Interfaces (Class Diagram)](#7-go-governance-packages--interfaces-class-diagram)
8. [First-Event AI Disclosure Session Flow (RBI Rec 18, 25 / EU Art 50)](#8-first-event-ai-disclosure-session-flow)
9. [CI/CD Pull Request Red-Team & Governance Gate](#9-cicd-pull-request-red-team--governance-gate)
10. [Incident Escalation: RBI Annexure VI ↔ EU AI Act Article 73](#10-incident-escalation-rbi-annexure-vi--eu-ai-act-article-73)
11. [RFC-002 Architecture Options: 10-Year Retention vs. Zero-Drift](#11-rfc-002-architecture-options-10-year-retention-vs-zero-drift)
12. [Dual Evidence Pack Assembly Pipeline (`cmd/evidencepack`)](#12-dual-evidence-pack-assembly-pipeline)

*(Plus 13 additional diagrams in [`README.md`](../README.md), [`docs/deep-dive/`](./deep-dive), [`docs/SOW-DEPLOYER-VS-PROVIDER.md`](./SOW-DEPLOYER-VS-PROVIDER.md), and [`RFC-001`](./open-problems/RFC-001-ai-kill-switch.md).)*

---

## 1. End-to-End System Context Architecture

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368", "clusterBkg": "#F8F9FA", "clusterBorder": "#DADCE0"}}}%%
flowchart TB
    subgraph CALLERS["Callers & Stakeholders"]
        U["Retail Customer"]
        ADV["Bank Advisor"]
        PB["Partner Bank API"]
        AUD["RBI / EU Auditor"]
    end

    subgraph CONTROL["Governance Control Plane"]
        YAML["Board Policy YAML<br/>(Rec 14 / Art 17)"]
        REG["Live Agent Registry<br/>(62 Agents · Zero Drift)"]
        APR["Approval Register<br/>(Rec 17 Drift Guard)"]
        CON["Append-Only Consent Ledger<br/>(Rec 15)"]
    end

    BUS["Composite Policy Gate on Message Bus<br/>10 Deterministic Pre-Dispatch Policies"]

    subgraph EXEC["Regional Execution Plane"]
        IN["India Pool (in)"]
        UAE["ADGM / DIFC Pool"]
        KSA["SAMA Pool"]
    end

    subgraph EVIDENCE["Tamper-Evident Evidence Plane"]
        LOG["SHA-256 Hash-Chained Log<br/>+ Ed25519 Signed Manifest"]
        TSA["External RFC 3161 TSA / WORM"]
        INC["Annexure VI & Art 73<br/>Incident Generator"]
        PACK["Dual Evidence Pack<br/>(RBI + EU AI Act)"]
    end

    U --> BUS
    ADV --> BUS
    PB --> BUS

    YAML --> BUS
    REG --> BUS
    APR --> BUS
    CON --> BUS

    BUS -->|"ALLOW"| IN
    BUS -->|"ALLOW"| UAE
    BUS -->|"ALLOW"| KSA

    BUS -->|"ALLOW / DENY"| LOG
    LOG --> TSA
    LOG --> INC
    LOG --> PACK
    REG --> PACK
    YAML --> PACK
    PACK --> AUD

    style BUS fill:#4285F4,stroke:#1A73E8,stroke-width:2px,color:#FFFFFF
    style IN fill:#E6F4EA,stroke:#34A853,color:#137333
    style UAE fill:#E6F4EA,stroke:#34A853,color:#137333
    style KSA fill:#E6F4EA,stroke:#34A853,color:#137333
    style LOG fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style TSA fill:#FEF7E0,stroke:#F9AB00,color:#B06000
    style INC fill:#FCE8E6,stroke:#EA4335,color:#A50E0E
    style PACK fill:#34A853,stroke:#137333,color:#FFFFFF
```

---

## 2. RBI FREE-AI Complete 3-Layer Taxonomy

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368", "clusterBkg": "#F8F9FA", "clusterBorder": "#DADCE0"}}}%%
flowchart TD
    ROOT["RBI FREE-AI Framework (13 Aug 2025)<br/>7 Sutras · 6 Pillars · 26 Recommendations"]

    subgraph SUTRAS["7 Guiding Sutras"]
        S1["1. Trust · 2. People First · 3. Innovation over Restraint<br/>4. Fairness · 5. Accountability · 6. Understandable · 7. Safety & Resilience"]
    end

    subgraph ENABLE["3 Enablement Pillars"]
        P1["Infrastructure"]
        P2["Policy"]
        P3["Capacity"]
    end

    subgraph MITIGATE["3 Risk Mitigation Pillars"]
        P4["Governance"]
        P5["Protection"]
        P6["Assurance"]
    end

    REG["Regulator & SRO Scope<br/>9 Recs: 1, 3, 5, 7, 9-13<br/>Exclude from RE SOW"]
    RE["Regulated Entity Scope<br/>17 Recs: 2, 4, 6, 8, 14-26<br/>17 / 17 Implemented"]

    ROOT --> SUTRAS
    SUTRAS --> ENABLE
    SUTRAS --> MITIGATE
    ENABLE --> REG
    ENABLE --> RE
    MITIGATE --> RE

    style ROOT fill:#4285F4,stroke:#1A73E8,stroke-width:2px,color:#FFFFFF
    style S1 fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style REG fill:#FEF7E0,stroke:#F9AB00,stroke-width:2px,color:#B06000
    style RE fill:#E6F4EA,stroke:#34A853,stroke-width:2px,color:#137333
```

---

## 3. Agent Governance Approval Lifecycle State Machine

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
stateDiagram-v2
    [*] --> Unreviewed: New high-risk agent registered
    Unreviewed --> UnderReview: Submit risk & red-team dossier
    Unreviewed --> Blocked_NoApproval: Dispatch attempt (ErrNoApproval)

    UnderReview --> Approved_Live: CRO signs off risk_class + expiry
    UnderReview --> Rejected: Review failed

    Rejected --> Blocked_Rejected: Dispatch attempt (ErrNotApproved)

    Approved_Live --> Blocked_Expired: now > expires_at (ErrExpired)
    Approved_Live --> Blocked_Drifted: Refactor changes risk_class (ErrRiskChanged)
    Approved_Live --> Killed_RFC001: Operator triggers kill switch (RFC-001)
    Approved_Live --> Decommissioned_RFC002: Agent retired (10-yr retention RFC-002)

    Blocked_Expired --> UnderReview: Re-validation submitted
    Blocked_Drifted --> UnderReview: New risk tier sign-off required
    Killed_RFC001 --> UnderReview: Dual-custody re-arm & audit
    Decommissioned_RFC002 --> [*]
```

---

## 4. Append-Only Consent Ledger & In-Flight Revocation Sequence

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "actorFontFamily": "Google Sans, Roboto, sans-serif", "noteFontFamily": "Google Sans, Roboto, sans-serif", "messageFontFamily": "Roboto Mono, ui-monospace, monospace", "messageFontSize": "13px", "actorBkg": "#E8F0FE", "actorBorder": "#4285F4", "actorTextColor": "#174EA6", "actorLineColor": "#DADCE0", "signalColor": "#5F6368", "signalTextColor": "#202124", "noteBkgColor": "#FEF7E0", "noteBorderColor": "#F9AB00", "noteTextColor": "#202124"}}}%%
sequenceDiagram
    autonumber
    actor Customer as Customer (user-1)
    participant Ledger as Append-Only Consent Ledger
    participant Gate as Composite Policy Gate
    participant Audit as SHA-256 Audit Log
    participant Agent as Portfolio Agent

    rect rgb(230, 244, 234)
    Note over Customer,Agent: Monday 09:00 - Consent Granted
    Customer->>Ledger: Grant("user-1", "pii", "onboarding opt-in")
    Ledger-->>Customer: seq=1 (GRANT recorded)
    Gate->>Ledger: IsLive("user-1", "pii")
    Ledger-->>Gate: true (seq=1 GRANT)
    Gate->>Agent: Dispatch PII request
    end

    rect rgb(252, 232, 230)
    Note over Customer,Agent: Tuesday 14:30 - Customer Revokes Consent Mid-Flight
    Customer->>Ledger: Revoke("user-1", "pii", "customer revoked on Tuesday")
    Ledger-->>Customer: seq=2 (REVOKE appended, seq=1 preserved)
    Gate->>Ledger: IsLive("user-1", "pii")
    Ledger-->>Gate: false (seq=2 REVOKE)
    Gate->>Audit: RecordDeny("user-1", "consent", "revoked on Tuesday")
    Gate-->>Customer: 403 Refused (agent never invoked, proof logged)
    end
```

---

## 5. Cryptographic Audit Chain, Ed25519 Manifest & TSA Anchoring

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
flowchart LR
    B1["seq 1 · genesis<br/>consent.grant<br/>row_hash: h1"] -->|"prev=h1"| B2["seq 2 · TAMPERED<br/>actor edited to #99<br/>hash mismatch!"]
    B2 -.->|"BROKEN LINK"| B3["seq 3 · INVALID<br/>refund.issue<br/>prev != recomputed h2"]
    B3 -.->|"BROKEN LINK"| B4["seq 4 · INVALID<br/>export.sign<br/>chain severed"]

    B4 --> VER["Verify() Walk<br/>ok=false, broken_at_seq=2<br/>field=row_hash"]
    VER --> MAN["Ed25519 SignedManifest<br/>chain_verified: false<br/>signature_valid: true"]
    MAN --> TSA["External RFC 3161 TSA<br/>or WORM Bucket Anchor"]

    style B1 fill:#E6F4EA,stroke:#34A853,stroke-width:2px,color:#137333
    style B2 fill:#FCE8E6,stroke:#EA4335,stroke-width:2.5px,color:#A50E0E
    style B3 fill:#FCE8E6,stroke:#EA4335,color:#A50E0E
    style B4 fill:#FCE8E6,stroke:#EA4335,color:#A50E0E
    style VER fill:#FEF7E0,stroke:#F9AB00,color:#B06000
    style MAN fill:#E8F0FE,stroke:#4285F4,stroke-width:2px,color:#174EA6
    style TSA fill:#4285F4,stroke:#1A73E8,color:#FFFFFF
```

---

## 6. Governance & Compliance Data Model (ER Diagram)

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "13px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
erDiagram
    POLICY_VERSION ||--o{ AUDIT_ENTRY : governs
    AGENT_REGISTRY ||--o| APPROVAL_RECORD : requires
    AGENT_REGISTRY ||--o{ AUDIT_ENTRY : target_of
    AGENT_REGISTRY ||--o{ INCIDENT_ANNEXURE_VI : triggers
    CONSENT_EVENT }o--|| AUDIT_ENTRY : referenced_by
    AUDIT_ENTRY }o--|| SIGNED_MANIFEST : summarized_in

    POLICY_VERSION {
        string version PK
        string approved_by
        date effective_date
        string hash_algorithm
    }
    AGENT_REGISTRY {
        string agent_id PK
        string domain
        string declared_risk
        string owner
        string vendor_model
        string disclosure_mode
    }
    APPROVAL_RECORD {
        string agent_id PK
        string status
        string approved_risk_class
        timestamp expires_at
    }
    CONSENT_EVENT {
        int seq PK
        string subject_id
        string category
        string action
        string reason
        timestamp created_at
    }
    AUDIT_ENTRY {
        int seq PK
        string actor
        string action
        string target
        string outcome
        string prev_hash
        string row_hash
    }
    SIGNED_MANIFEST {
        string payload_sha256 PK
        int count
        bool chain_verified
        int broken_at_seq
        string ed25519_sig
        string tsa_anchor
    }
    INCIDENT_ANNEXURE_VI {
        string incident_id PK
        string agent_id FK
        string severity
        string control_triggered
        bool audit_chain_verified
    }
```

---

## 7. Go Governance Packages & Interfaces (Class Diagram)

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "13px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
classDiagram
    class CompositeGate {
        +Policies []Policy
        +Audit AuditLogger
        +Allow(ctx, msg) error
    }
    class Policy {
        <<interface>>
        +Name() string
        +Check(ctx, msg) error
    }
    class AuditLogger {
        <<interface>>
        +RecordAllow(msg) error
        +RecordDeny(msg, policy, reason) error
    }
    class RegistryPolicy {
        +IsRegistered func(string) bool
        +Check(ctx, msg) error
    }
    class ApprovalDriftPolicy {
        +Register map~string,*ApprovalRecord~
        +Now func() time.Time
        +Check(ctx, msg) error
    }
    class ResidencyPolicy {
        +AllowedPIIRegions map~string,bool~
        +DenyRegions map~string,bool~
        +Check(ctx, msg) error
    }
    class ExplainabilityPolicy {
        +Check(ctx, msg) error
    }
    class RBACPolicy {
        +HighRiskRoles map~string,bool~
        +Check(ctx, msg) error
    }
    class ConsentPolicy {
        +LedgerCheck func
        +Check(ctx, msg) error
    }
    class AuditLog {
        -entries []Entry
        +Append(...) Entry
        +Snapshot() []Entry
        +Verify() (bool, int, string)
        +ExportSignedManifest(priv, anchor) (SignedManifest, error)
    }
    class ConsentLedger {
        +Grant(subject, cat, reason, at) (Event, error)
        +Revoke(subject, cat, reason, at) (Event, error)
        +IsLive(subject, cat) (bool, Event, bool)
    }
    class DisclosureTracker {
        +MaxSessions int
        +Wrap(sessionID, payload) Envelope
    }
    class LiveRegistry {
        +Register(desc)
        +Has(agentID) bool
        +Export(now) InventoryExport
    }

    CompositeGate --> Policy : evaluates in order
    CompositeGate --> AuditLogger : logs allow/deny
    Policy <|.. RegistryPolicy
    Policy <|.. ApprovalDriftPolicy
    Policy <|.. ResidencyPolicy
    Policy <|.. ExplainabilityPolicy
    Policy <|.. RBACPolicy
    Policy <|.. ConsentPolicy
    AuditLogger <|.. AuditLog
    RegistryPolicy ..> LiveRegistry : checks active registration
    ConsentPolicy ..> ConsentLedger : queries live state
```

---

## 8. First-Event AI Disclosure Session Flow

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "actorFontFamily": "Google Sans, Roboto, sans-serif", "noteFontFamily": "Google Sans, Roboto, sans-serif", "messageFontFamily": "Roboto Mono, ui-monospace, monospace", "messageFontSize": "13px", "actorBkg": "#E8F0FE", "actorBorder": "#4285F4", "actorTextColor": "#174EA6", "actorLineColor": "#DADCE0", "signalColor": "#5F6368", "signalTextColor": "#202124", "noteBkgColor": "#E6F4EA", "noteBorderColor": "#34A853", "noteTextColor": "#137333"}}}%%
sequenceDiagram
    autonumber
    actor User as Customer (sess-101)
    participant Tracker as disclosure.Tracker (Art 50)
    participant Agent as Educator Agent

    rect rgb(230, 244, 234)
    Note over User,Agent: Turn 1 - First Event of Session (Mandatory Disclosure)
    User->>Agent: "What is my savings balance?"
    Agent->>Tracker: Wrap("sess-101", "Balance is INR 1,42,000")
    Tracker-->>User: Envelope(ai_disclosed=true, notice="AI Disclosure RBI Rec 18,25 / EU Art 50")
    end

    rect rgb(232, 240, 254)
    Note over User,Agent: Turn 2 - Same Session (No Duplicate Banner)
    User->>Agent: "Show last 3 transactions"
    Agent->>Tracker: Wrap("sess-101", "Last 3 transactions...")
    Tracker-->>User: Envelope(ai_disclosed=false, payload="Last 3 transactions...")
    end
```

---

## 9. CI/CD Pull Request Red-Team & Governance Gate

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
flowchart LR
    PR["Pull Request Opened"] --> VET["go vet & go test<br/>All 9 Go Packages"]
    VET --> RT["20-Case Red-Team Corpus<br/>RT-01..RT-20 (Rec 20 / Art 15)"]
    RT --> PY["4 Python Labs<br/>Gate · Drift · Audit · Consent"]
    PY --> XW["Cross-Walk & YAML Check<br/>JSON · CSV · policy.yaml"]
    XW --> MMD["Mermaid & Web Deck Lint<br/>index.html + Diagrams"]
    MMD --> PASS["CI Green Badge<br/>Auto-Deploy to Vercel"]

    RT -.->|"Any bypass fails PR"| FAIL["Block Merge"]

    style PR fill:#F8F9FA,stroke:#5F6368,color:#202124
    style VET fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style RT fill:#FEF7E0,stroke:#F9AB00,stroke-width:2px,color:#B06000
    style PY fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style XW fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style MMD fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style PASS fill:#34A853,stroke:#137333,color:#FFFFFF
    style FAIL fill:#FCE8E6,stroke:#EA4335,color:#A50E0E
```

---

## 10. Incident Escalation: RBI Annexure VI ↔ EU AI Act Article 73

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "actorFontFamily": "Google Sans, Roboto, sans-serif", "noteFontFamily": "Google Sans, Roboto, sans-serif", "messageFontFamily": "Roboto Mono, ui-monospace, monospace", "messageFontSize": "13px", "actorBkg": "#E8F0FE", "actorBorder": "#4285F4", "actorTextColor": "#174EA6", "actorLineColor": "#DADCE0", "signalColor": "#5F6368", "signalTextColor": "#202124", "noteBkgColor": "#FCE8E6", "noteBorderColor": "#EA4335", "noteTextColor": "#A50E0E"}}}%%
sequenceDiagram
    autonumber
    participant Gate as Policy / Drift Gate
    participant Audit as SHA-256 + Ed25519 Log
    participant Inc as incident.NewAnnexureVI
    participant RBI as RBI Supervisory Portal
    participant EU as EU Market Surveillance

    rect rgb(252, 232, 230)
    Note over Gate,Inc: Control Failure or Drift Detected Pre-Dispatch
    Gate->>Audit: RecordDeny("var_calculator", "ErrRiskChanged")
    Audit->>Inc: ExportSignedManifest(priv, tsa_anchor)
    Inc->>Inc: Build structured Annexure VI JSON + Markdown
    end

    rect rgb(232, 240, 254)
    Note over Inc,EU: One Incident Artefact, Two Regulatory Filings
    Inc-->>RBI: Submit RBI FREE-AI Rec 22 Annexure VI Form
    Inc-->>EU: Submit EU AI Act Article 73 Serious Incident Notice
    end
```

---

## 11. RFC-002 Architecture Options: 10-Year Retention vs. Zero-Drift

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368", "clusterBkg": "#F8F9FA", "clusterBorder": "#DADCE0"}}}%%
flowchart TB
    CONFLICT["Architectural Conflict (RFC-002)<br/>Live Inventory Cannot Drift, and Cannot Remember 10 Years"]

    subgraph OPT_A["Option A: Event Log + Live Projection"]
        A1["Append-Only Lifecycle Log<br/>REGISTER / UPDATE / DECOMMISSION"] --> A2["Live View (ACTIVE only)"]
        A1 --> A3["10-Year Historical Query"]
    end

    subgraph OPT_B["Option B: Signed WORM Snapshots"]
        B1["Live Registry (Zero Drift)"] --> B2["Daily Ed25519 Export"]
        B2 --> B3["10-Yr S3/GCS Object Lock"]
    end

    subgraph OPT_C["Option C: Registry Tombstones"]
        C1["Single Registry Table"] --> C2["ACTIVE Rows (Dispatchable)"]
        C1 --> C3["DECOMMISSIONED Tombstones<br/>(Must Never Dispatch)"]
    end

    CONFLICT --> OPT_A
    CONFLICT --> OPT_B
    CONFLICT --> OPT_C

    style CONFLICT fill:#FCE8E6,stroke:#EA4335,stroke-width:2px,color:#A50E0E
    style A1 fill:#E6F4EA,stroke:#34A853,color:#137333
    style B1 fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style C3 fill:#FEF7E0,stroke:#F9AB00,color:#B06000
```

---

## 12. Dual Evidence Pack Assembly Pipeline

```mermaid
%%{init: {"theme": "base", "themeVariables": {"fontFamily": "Google Sans, Roboto, Helvetica Neue, Arial, sans-serif", "fontSize": "14px", "primaryColor": "#E8F0FE", "primaryTextColor": "#202124", "primaryBorderColor": "#4285F4", "lineColor": "#5F6368"}}}%%
flowchart LR
    M1["code/config<br/>policy.yaml"] --> CLI["go run ./cmd/evidencepack<br/>Offline Assembler (< 1s)"]
    M2["code/registry<br/>Live 62 Agents"] --> CLI
    M3["code/gate<br/>Drift Report"] --> CLI
    M4["code/disclosure<br/>Art 50 Checks"] --> CLI
    M5["code/audit<br/>SHA-256 + Ed25519"] --> CLI
    M6["code/incident<br/>Annexure VI"] --> CLI

    CLI --> OUT["rbi-eu-evidence-pack/<br/>INDEX.md + 6 Folders"]
    OUT --> PACK1["RBI FREE-AI Evidence<br/>Recs 14, 16, 17, 18, 22, 23, 25"]
    OUT --> PACK2["EU AI Act Evidence<br/>Arts 9, 12, 14, 17, 50, 73 & Annex IV"]

    style CLI fill:#4285F4,stroke:#1A73E8,stroke-width:2px,color:#FFFFFF
    style OUT fill:#E6F4EA,stroke:#34A853,stroke-width:2px,color:#137333
    style PACK1 fill:#E8F0FE,stroke:#4285F4,color:#174EA6
    style PACK2 fill:#FCE8E6,stroke:#EA4335,color:#A50E0E
```
