.PHONY: all build compile fmt vet lint test test-cover-html clean ci run run-dev gen-wire-deps help
export GO111MODULE=on

APP=eta-service
APP_VERSION:=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
APP_COMMIT:=$(shell git rev-parse HEAD 2>/dev/null || echo unknown)
APP_EXECUTABLE=./out/$(APP)
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor/)

all: clean build test

help:
	@echo "Available targets:"
	@echo "  build          - compile, fmt, vet, lint"
	@echo "  compile        - build the binary into out/$(APP)"
	@echo "  run            - build and run the service (./out/$(APP) start [CONFIG])"
	@echo "  run-dev        - go run the service without building a binary (CONFIG=test to use config/test.yaml)"
	@echo "  gen-wire-deps  - regenerate cmd/$(APP)/wire_gen.go via wire"
	@echo "  fmt            - go fmt all packages"
	@echo "  vet            - go vet all packages"
	@echo "  lint           - golangci-lint (skipped with a warning if not installed)"
	@echo "  test           - go test with coverage (out/cover.out)"
	@echo "  test-cover-html- generate an HTML coverage report (out/coverage.html)"
	@echo "  clean          - remove build artifacts"
	@echo "  ci             - clean, build, test"

# Dependency injection (regenerates cmd/eta-service/wire_gen.go from di.go)
gen-wire-deps:
	cd cmd/$(APP) && wire

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE) -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" ./cmd/$(APP)

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

# Uses golangci-lint if available; skips (with a warning) rather than failing
# the build when it isn't installed. Install: https://golangci-lint.run/usage/install/
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping lint (see https://golangci-lint.run/usage/install/)"; \
	fi

build: compile fmt vet lint

clean:
	rm -rf out/

# -vet=off avoids re-running go vet during tests since `make vet` already
# covers it as a separate target; keeps `go test` failures scoped to actual
# test failures rather than vet findings.
test:
	mkdir -p out/
	go test -vet=off -race -covermode=atomic -coverprofile=out/cover.out $(ALL_PACKAGES)
	@go tool cover -func=out/cover.out | tail -1

test-cover-html: test
	go tool cover -html=out/cover.out -o out/coverage.html
	@echo "Coverage report: out/coverage.html"

ci: clean build test

run: compile
	$(APP_EXECUTABLE) start $(CONFIG)

run-dev:
	go run ./cmd/$(APP) start $(CONFIG)
