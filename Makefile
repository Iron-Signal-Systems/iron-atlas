SHELL := /usr/bin/env bash

.PHONY: format test validate gate build clean

format:
	gofmt -w cmd internal modules integrations

test:
	go test -race ./...

validate:
	./tools/validation/validate_repository.sh

gate:
	./tools/validation/phase-gates/validate_phase0_step1.sh

build:
	mkdir -p build
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/iron-atlas ./cmd/atlasd
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o build/atlasctl ./cmd/atlasctl

clean:
	rm -rf build dist
