GOTEST_PKGS=$(shell go list ./... | grep -v examples)

BENCHTIME ?= 2s
BENCHTESTS ?= .

BENCHFULL?=0
ifeq (${BENCHFULL},1)
BENCHFULL_ARG=-bench-full -timeout 60m
else
BENCHFULL_ARG=
endif

TEST_VERBOSE?=0
ifeq (${TEST_VERBOSE},1)
TEST_VERBOSE_ARG=-v
else
TEST_VERBOSE_ARG=
endif

TEST_RESULTS?="/tmp/test-results"

generate:
	@echo "Regenerating Parser"
	@go generate ./

test:
	@go test $(TEST_VERBOSE_ARG) $(GOTEST_PKGS)

test-ci:
	@gotestsum --junitfile $(TEST_RESULTS)/gotestsum-report.xml -- $(GOTEST_PKGS)

bench:
	@go test $(TEST_VERBOSE_ARG) -run DONTRUNTESTS -bench $(BENCHTESTS) $(BENCHFULL_ARG) -benchtime=$(BENCHTIME) $(GOTEST_PKGS)

coverage:
	@go test -coverprofile /tmp/coverage.out $(GOTEST_PKGS)
	@go tool cover -html /tmp/coverage.out

fmt:
	@gofmt -w -s

examples: simple expr-parse expr-eval filter

simple:
	@go build ./examples/simple

expr-parse:
	@go build ./examples/expr-parse

expr-eval:
	@go build ./examples/expr-eval

filter:
	@go build ./examples/filter

deps:
	@go get github.com/mna/pigeon@master
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/tools/cmd/cover
	@go mod tidy

.PHONY: generate test coverage fmt deps bench examples expr-parse expr-eval filter

