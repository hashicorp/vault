#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
#

#
# Copyright 2019 Joyent, Inc.
#

TEST?=$$(go list ./... |grep -Ev 'vendor|examples|testutils')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

.PHONY: all
all:

.PHONY: tools
tools: ## Download and install all dev/code tools
	@echo "==> Installing dev tools"
	go get -u github.com/golang/dep/cmd/dep

.PHONY: build
build:
	@govvv build

.PHONY: install
install:
	@govvv install

.PHONY: test
test: ## Run unit tests
	@echo "==> Running unit test with coverage"
	@./scripts/go-test-with-coverage.sh

.PHONY: testacc
testacc: ## Run acceptance tests
	@echo "==> Running acceptance tests"
	TRITON_TEST=1 go test $(TEST) -v $(TESTARGS) -run -timeout 120m

.PHONY: check
check:
	scripts/gofmt-check.sh

.PHONY: help
help: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
