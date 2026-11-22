.PHONY: demo test serve validate-crosswalk

demo:
	@echo "=== 1. Composite Policy Gate Demo ==="
	python3 labs/gate_demo.py
	@echo ""
	@echo "=== 2. Lab 01: Approval Drift (approval.py) ==="
	python3 labs/approval.py
	@echo ""
	@echo "=== 3. Lab 02: Tamper-Evident Audit Chain (audit.py) ==="
	python3 labs/audit.py

test: validate-crosswalk
	python3 labs/gate_demo.py --test
	python3 labs/approval.py --test
	python3 labs/audit.py --test

validate-crosswalk:
	python3 -c 'import json, csv; d=json.load(open("crosswalk/rbi-eu-ai-act.json")); assert len(d["rbi_free_ai_recommendations"])==26; assert sum(1 for r in d["rbi_free_ai_recommendations"] if r["owner"]=="regulated_entity")==17; assert len(d["rbi_to_eu_crosswalk"])==7; assert len(d["rbi_draft_mrm_scorecard"])==6; rows=list(csv.DictReader(open("crosswalk/rbi-eu-ai-act.csv"))); assert len(rows)==7; print("[OK] Cross-walk JSON & CSV verified.")'

serve:
	python3 -m http.server 8000
