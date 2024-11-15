GOLANG_CI_LINT_VERSION := $(shell golangci-lint --version 2>/dev/null)

.PHONY: all
all: test lint

.PHONY: clean
clean: ## Clean testcache and delete build output
	go clean -testcache

.PHONY: test
test: ## Run the unit tests
	go test -v -race

.PHONY: generate
generate: ## Generate fakes
	go generate

.PHONY: lint-prepare
lint-prepare:
ifdef GOLANG_CI_LINT_VERSION
	@echo "Found golangci-lint $(GOLANG_CI_LINT_VERSION)"
else
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	@echo "[OK] golangci-lint installed"
endif

.PHONY: lint
lint: lint-prepare ## Run the golangci linter
	golangci-lint run

.PHONY: tidy
tidy: ## Remove unused dependencies
	go mod tidy

.PHONY: list
list: ## Print the current module's dependencies.
	go list -m all

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'