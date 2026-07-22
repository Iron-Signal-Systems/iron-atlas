SHELL := /usr/bin/env bash

.PHONY: format tidy test validate toolchain portable-validation evidence-check canonical-step2 gate phase0-gate phase1-step1-gate phase1-step2-gate database-test build clean

format:
	gofmt -w cmd internal modules integrations

tidy:
	go mod tidy

test:
	./test-framework/run_all.sh

validate:
	./tools/validation/validate_repository.sh

toolchain:
	python3 tools/validation/validate_toolchain.py

portable-validation:
	python3 tools/validation/validate_portable_acceptance.py
	./test-framework/portability/test_portable_validation.sh

evidence-check:
	python3 tools/validation/validate_committed_evidence.py

canonical-step2:
	./tools/validation/verify_canonical_clone.sh "$$(git rev-parse HEAD)" ./tools/validation/phase-gates/validate_phase1_step2.sh

gate: phase1-step2-gate

phase0-gate:
	./tools/validation/phase-gates/validate_phase0_acceptance.sh

phase1-step1-gate:
	./tools/validation/phase-gates/validate_phase1_step1_acceptance.sh

phase1-step2-gate:
	./tools/validation/phase-gates/validate_phase1_step2.sh

database-test:
	./test-framework/database/run_disposable_postgres.sh

build:
	mkdir -p build
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/atlas ./cmd/atlasd
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/atlasctl ./cmd/atlasctl

clean:
	rm -rf build dist
