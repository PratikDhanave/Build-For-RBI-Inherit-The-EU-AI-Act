# Contributing

Contributions, regulatory cross-walk refinements, adversarial red-team test cases, and architecture proposals for the two open problems are warmly welcome!

## High-Impact Areas to Contribute
1. **Open Problem Proposals:**
   - [`RFC-001` — Mandatory Operator-Triggered AI Kill Switch](./docs/open-problems/RFC-001-ai-kill-switch.md)
   - [`RFC-002` — 10-Year Decommissioned Model Retention vs. Zero-Drift Live Inventory](./docs/open-problems/RFC-002-ten-year-retention.md)
   - Use the [Open Problem Proposal Issue Template](../../issues/new?template=open-problem-proposal.md).
2. **Red-Team Corpus Expansion:**
   - Add new adversarial edge cases to [`code/redteam/corpus.json`](./code/redteam/corpus.json) and verify `go test ./...` passes.
3. **Additional Regulatory Instruments:**
   - Extend [`crosswalk/rbi-eu-ai-act.json`](./crosswalk/rbi-eu-ai-act.json) with mappings to MAS FEAT, ADGM, DIFC, SAMA, or NIST AI RMF.

## Verification Before Opening a PR

All controls must remain **offline-runnable with zero external dependencies**:

```bash
make test
```

GitHub Actions (`.github/workflows/verify.yml`) runs the exact same verification suite on every pull request.
