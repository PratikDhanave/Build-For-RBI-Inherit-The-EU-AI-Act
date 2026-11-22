.PHONY: demo test serve validate-crosswalk

demo:
	bash code/scripts/offline_verify.sh

test: validate-crosswalk
	(cd code && go test -v ./...)
	(cd code && go run ./cmd/demo >/dev/null)
	python3 code/python/gate_demo.py --test
	python3 code/python/approval.py --test
	python3 code/python/audit.py --test
	python3 code/python/consent_ledger.py --test
	python3 labs/gate_demo.py --test
	python3 labs/approval.py --test
	python3 labs/audit.py --test

validate-crosswalk:
	python3 -c 'import json, csv; d=json.load(open("crosswalk/rbi-eu-ai-act.json")); assert len(d["rbi_free_ai_recommendations"])==26; assert sum(1 for r in d["rbi_free_ai_recommendations"] if r["owner"]=="regulated_entity")==17; assert len(d["rbi_to_eu_crosswalk"])==7; assert len(d["rbi_draft_mrm_scorecard"])==6; rows=list(csv.DictReader(open("crosswalk/rbi-eu-ai-act.csv"))); assert len(rows)==7; json.load(open("examples/annexure-vi-sample.json")); json.load(open("examples/inventory-export.json")); json.load(open("code/redteam/corpus.json")); print("[OK] Cross-walk, sample exports & red-team JSON verified.")'

serve:
	python3 -m http.server 8000
