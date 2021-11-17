## Setup
SHELL := /bin/bash
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

setup:
	@which golint &> /dev/null  || go get -u golang.org/x/lint/golint
	@which make2help &> /dev/null || go get github.com/Songmu/make2help/cmd/make2help
	@which staticcheck &> /dev/null || go get honnef.co/go/tools/cmd/staticcheck

## Install dependencies
deps: setup
	go mod tidy
	go mod vendor

## Show help
help:
	@make2help $(MAKEFILE_LIST)

# Format source codes (internally used)
cfmt: setup
	@gofmt -l -w $(SRC)

# Lint (internally used)
clint: deps
	@echo "Running staticcheck" && staticcheck
	@echo "Running go vet and lint"
	@for pkg in $$(go list ./... | grep -v /vendor/); do \
		echo "Verifying $$pkg"; \
		go vet $$pkg; \
		golint -set_exit_status $$pkg || exit $$?; \
	done

# Install (internally used)
cinstall:
	@export GOBIN=$$GOPATH/bin; \
	go install -tags=sfdebug $(CMD_TARGET).go

# Run (internally used)
crun: install
	$(CMD_TARGET)

.PHONY: setup help cfmt clint cinstall crun
