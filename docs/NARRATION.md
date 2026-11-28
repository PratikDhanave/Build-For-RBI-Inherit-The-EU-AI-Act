# Speaker Guide — Build for RBI, Inherit the EU AI Act (38 Slides)

Matches the live deck at **https://build-for-rbi-talk.vercel.app/** (`35` talk slides + `3` Appendix slides `36–38`).
**Recommended slot:** ~35 minutes (or ~20 minutes by skipping the 5 `[SKIP IF TIGHT]` slides) · **5 live Go WebAssembly demos** (`S18`, `S20`, `S23`, `S25`, `S26`).

> **Pre-stage checklist (1 minute before going live):**
> 1. Open the deck 60 seconds early and confirm Slide 18 shows `Go WebAssembly ready · real compiled Go` (`assets/labs.wasm` is 4.2 MB).
> 2. If you click inside an editable code box on stage, press **`Esc`** before `r` or `←/→` (or press `Cmd+Enter` inside the editor to run).
> 3. **Running tight on time?** Safely skip **`S12`**, **`S13`**, **`S15`**, **`S31`**, and **`S33`** (each deepens a neighbouring slide) — and treat **`S20`** and **`S26`** as quick 15-second one-click `r` hits.
> 4. **Facts to say accurately:** `approval.go` (`S23`) prints **2 RUN, 3 STOP** (`payment_orchestrator` & `educator` run; `kyc`, `aml`, `var_calculator` stop). Consent revocation is **Tuesday 2:30 PM → 2:31 PM**. Never claim Wi-Fi-off or a "9-second audit".

---

## Opening & Mental Models (Slides 01–15)

* **01 · Cover** — "What India's banking AI rules ask of your codebase — and how much of Europe's AI law comes with it. Five live Go checks run right inside these slides."
* **02 · Before we begin** — Ground rules: ask questions, plain English first, real Go code, honest gaps named openly.
* **03 · Huge shoutout** — Pause for the FDExHyderabad Operations Team. *"A talk is the visible part of a great deal of good work, and this one stands on theirs."*
* **04 · The honest score, up front** — `17 / 26` RBI codebase controls, `7` EU AI Act controls inherited, `4 built` / `2 open gaps` against RBI Draft MRM (24 June 2026).
* **05 · The words we will use** — Only 2 terms (`Policy gate`, `Audit chain`) + 6 plain jobs. All regulatory numbers wait until the end.
* **06 · Three knocks on our door** — Meet **Bodhyanga**: (1) PAN stays in India per message, (2) **Revoke at 2:30 PM, fail at 2:31 PM**, (3) signed proof in 3 seconds (*"Encrypted at rest"* is not an answer).
* **07 · When it goes wrong in public** — Four real incidents (`$1 car`, `Air Canada tribunal`, `Samsung code leak`, `Knight Capital $440M in 45 min`). *"Three of these were not clever attacks — systems moving fast with nobody in the way."*
* **08 · The handler everyone writes first** — **Read the 8 lines of Go slowly, then walk down the six flaws one at a time.** Eight lines of code, six missing controls — and we are about to build all six around this exact handler.
* **09 · The one idea** — Pause 2 seconds: *"AI compliance is not a legal checklist. It is ordinary engineering around your AI."*
* **10 · Nine belong to the referee — seventeen are yours** — RBI FREE-AI (13 Aug 2025): `13` direct RE Risk Mitigation recs (`14–26`) + `4` RE technical-readiness tracks (`2, 4, 6, 8`) = `17` codebase-scope controls; `9` belong to regulators/SROs & non-codebase org training (`Rec 10`).
* **11 · Think of it as an airport** — Clean mental model mapping all 6 controls; honest gap: airport gates are steel, ours are software so we fail closed and test every PR.
* **12 · You already rely on all six (`[SKIP IF TIGHT]`)** — Everyday proof (`UPI autopay`, `scooter vs. 20-tonne truck`, `medicine seal`, `fire drill`, `one passport`).
* **13 · What you rent vs. what you own (`[SKIP IF TIGHT]`)** — Cloud gives WAF/KMS/WORM; your codebase owns the banking judgement.
* **14 · The whole thing, on one slide** — Request top-left (`Gate → Consent → Drift → 62 Agents → Disclosure/Registry → SHA-256 Audit → Dual Evidence`), wrapped by `policy.yaml`.
* **15 · The six, and where they live (`[SKIP IF TIGHT]`)** — Plain name → RBI/EU name → Go package (`code/gate/`, `consent/`, `audit/`, `redteam/`, `registry/`).

---

## Phase 1 · The Doorway (Slides 16–20)

* **16 · Phase 1: The Doorway** — Putting checks inside 62 agents broke; moving **six plain rules** into one pre-dispatch gate fixed all 62 at once.
* **17 · 1 · Check every message before any agent sees it** — `Without it` vs `With it`. Pure Go gate runs in nanoseconds (under a millisecond), fails closed, never guesses like an LLM judge.
* **18 · 1 · LIVE DEMO: Six messages, one gate** — Press **`r`**. Point at the verdict: `2 allowed, 4 blocked before any agent woke up.` **Then click the third button — `Try made-up risk "critical"` — to show a sneaky bypass trick failing closed.**
* **19 · 2 · Revoke on Tuesday, fail in-flight** — Never overwrite `consent = false`; append `GRANT` (`seq 1`) and `REVOKE` (`seq 2`) so history survives.
* **20 · 2 · LIVE DEMO: Monday grant, Tuesday revoke** — Press **`r`**: Monday `ALLOW`, Tuesday 14:31 `DENY [consent] revoked at seq 2`, both rows preserved.

---

## Phase 2 · The Proof (Slides 21–27)

* **21 · Phase 2: The Proof** — Two quiet bugs that slip past every checklist: approval drift and insider log edits.
* **22 · 3 · When an old permit stops matching the code** — `var_calculator` approved at `medium`, refactored to `high` — still `approved`, still unexpired. Valid scooter licence driving a 20-tonne truck.
* **23 · 3 · LIVE DEMO: Five agents, one trap** — Press **`r`**. **2 RUN, 3 STOP** (`payment_orchestrator` & `educator` run; `kyc` expired, `aml` rejected, **`var_calculator` stopped by `ErrRiskChanged`**). Click **Re-approve as high** to show live Go re-evaluation.
* **24 · 4 · A lock that cannot catch its own team is not a lock** — Encryption hides from outsiders, not DB admins. **Write the record first, then act** (if the audit write fails, the gate refuses to dispatch). Three layers: SHA-256 chain + caller-pinned Ed25519 + external WORM anchor.
* **25 · 4 · LIVE DEMO: Break one row** — Press **`r`**. Fresh chain `OK` → insider edits Row 2 → `BROKEN at seq 2 (row_hash)`, manifest honestly sets `chain_verified: false`. *"The export does not hide the break. It reports it."*
* **26 · 5 · LIVE DEMO: 20 attack tests on every PR** — Press **`r`**: scores all 20 OWASP ASI cases live (`20/20 clean sweep — 13 attacks blocked, 7 safe paths allowed`). Honest leak: tests gate rules, not LLM jailbreaks.
* **27 · A control nobody can run is just a claim** — Honest leak: browser WASM runs real compiled Go in memory; production adds WORM + HSM. **The documentation is tested too** (`make test` checks crosswalk, YAML parity, 11 assets, 25 Mermaid diagrams).

---

## Phase 3 · Europe, Honest Gaps & Monday (Slides 28–35)

* **28 · Phase 3: Europe & Gaps** — What you inherit free for Europe, and what RBI's June 2026 Draft MRM exposed as still missing.
* **29 · 6 · Seven controls, two reports** — Pause 5 seconds so room photographs the table. Highlight Row 1 (`Rec 22 ↔ Art. 12`) and Row 6 (`Art. 50` already binding law since **2 Aug 2026**).
* **30 · Four built — two honest gaps** — RBI Draft MRM (24 June 2026): 4 built, 2 open RFCs (`RFC-001` operator AI kill switch referencing **CrowdStrike Falcon Guardian**; `RFC-002` 10-year decommissioned retention without polluting live inventory).
* **31 · Two open problems worth solving right (`[SKIP IF TIGHT]`)** — Architecture sketch of `RFC-001` (3-step stop) and `RFC-002` (live catalog + cold WORM tombstones).
* **32 · User vs. maker — decide before signing** — EU AI Act Article 25: rebranding or substantially modifying a high-risk AI flips your bank from **Deployer** to **Provider**. Lock it into the SOW Day 1.
* **33 · Check yourself (`[SKIP IF TIGHT]`)** — Click-to-reveal recap of the 3 hardest audit questions.
* **34 · Three things to take to work Monday** — (01) Build the 17, not 26. (02) One boring gate outside your agents + tamper-evident chain. (03) Name your gaps and SOW boundary Day 1.
* **35 · Thank you** — Point to both QR codes (live deck + open-source Go repo). Stop cleanly for Q&A.

---

## Appendix — Q&A Reference (`36↵`–`38↵`)

* **36 (`36↵`) · The words — the engineering** — Composite Policy Gate, Fail closed, Append-only ledger, SHA-256 chain, Pinned-key Ed25519, WORM, Approval drift.
* **37 (`37↵`) · The words — the regulations** — RBI FREE-AI, RBI Draft MRM, EU Digital Omnibus `2026/1744`, Art. 50, Art. 25, Annexure VI, DPDP '25 / MeitY '26.
* **38 (`38↵`) · One rulebook, five jurisdictions** — `examples/policy.yaml` loaded across **stateless horizontally scaled gate replicas behind your load balancer** (zero single point of failure, zero country forks).
