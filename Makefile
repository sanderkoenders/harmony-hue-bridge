.PHONY: help install-hooks deps init fmt lint

GOBIN := $(shell go env GOPATH)/bin

help:
	@echo "Targets:"
	@echo "  init           Install hooks and bootstrap Go deps"
	@echo "  install-hooks  Configure repo-local git hooks"
	@echo "  deps           Download Go modules and install tooling"
	@echo "  fmt            Format all Go files"
	@echo "  lint           Run golangci-lint"

install-hooks:
	@mkdir -p .githooks
	@chmod +x .githooks/pre-commit
	@git config --local core.hooksPath .githooks
	@echo "Git hooks installed at .githooks"

deps:
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Go dependencies and tooling installed"

init: install-hooks deps
	@echo "Project initialized successfully"

fmt:
	gofmt -w ./...

lint:
	$(GOBIN)/golangci-lint run --fast
