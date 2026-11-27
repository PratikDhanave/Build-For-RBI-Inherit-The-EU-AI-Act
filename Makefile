.PHONY: demo test validate-crosswalk wasm serve

demo:
	bash code/scripts/offline_verify.sh

test: validate-crosswalk
	(cd code && go vet ./... && go test -race -count=1 -v ./...)
	(cd code && go run ./cmd/demo >/dev/null)
	(cd labs && go vet ./... && go test -race -count=1 -v ./...)
	go run labs/gate_demo/main.go --test
	go run labs/approval/main.go --test
	go run labs/audit/main.go --test
	go run labs/consent_ledger/main.go --test

validate-crosswalk:
	(cd code && go run ./cmd/verify --root ..)

wasm:
	(cd labs && GOOS=js GOARCH=wasm go build -buildvcs=false -ldflags="-s -w" -o ../assets/labs.wasm ./wasm)
	@GOROOT="$$(go env GOROOT)"; \
	if [ -f "$$GOROOT/lib/wasm/wasm_exec.js" ]; then cp "$$GOROOT/lib/wasm/wasm_exec.js" assets/wasm_exec.js; \
	else cp "$$GOROOT/misc/wasm/wasm_exec.js" assets/wasm_exec.js; fi
	@echo "[OK] Built assets/labs.wasm and synced assets/wasm_exec.js"

serve:
	(cd code && go run ./cmd/serve --dir .. --addr :8000)
