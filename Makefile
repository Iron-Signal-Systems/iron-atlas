SHELL := /usr/bin/env bash

.PHONY: format test validate gate phase0-gate phase1-gate database-test build clean

format:
	gofmt -w cmd internal modules integrations

test:
	./test-framework/run_all.sh

validate:
	./tools/validation/validate_repository.sh

gate: phase1-gate

phase0-gate:
	./tools/validation/phase-gates/validate_phase0_acceptance.sh

phase1-gate:
	./tools/validation/phase-gates/validate_phase1_step1.sh

database-test:
	./test-framework/database/run_disposable_postgres.sh

build:
	mkdir -p build
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/iron-atlas ./cmd/atlasd
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/atlasctl ./cmd/atlasctl

clean:
	rm -rf build dist
