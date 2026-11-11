shell := /bin/sh

GO ?= go
GOBIN := $(shell $(GO) env GOPATH)/bin
GOIMPORTS ?= $(GOBIN)/goimports
GOLANGCI ?= $(GOBIN)/golangci-lint

LOCAL_PREFIX := raft-based-kv

.PHONY: all fmt fmt-check vet lint test test-race build check install-tools hooks

all: check

fmt:
	$(GOIMPORTS) -local $(LOCAL_PREFIX) -w $(GO_FILES)
	gofmt -w $(GOPATH)

fmt-check:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi
	@out=$$($(GOIMPORTS) -local $(LOCAL_PREFIX) -l .); \
	if [ -n "$$out" ]; then echo "goimports needed:"; echo "$$out"; exit 1; fi

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI) run ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

build:
	$(GO) build ./...

check: fmt-check vet lint test build

install-tools:
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

hooks:
	git config core.hooksPath .githooks
	@echo "Git hooks enabled from .githooks/"
