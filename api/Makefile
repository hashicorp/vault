# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

include ../Makefile.defs

TEST=$$(go list ./... | grep -v /vendor/ | grep -v /integ)
GOFMT_FILES=$$(find . -name '*.go' | grep -v pb.go)

# test runs the unit tests and vets the code
test: prep
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	go test tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=$(TEST_TIMEOUT) -parallel=20

testcompile: prep
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# testacc runs acceptance tests
testacc: prep
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT)

# testrace runs the race checker
testrace: prep
	@CGO_ENABLED=1 \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	go test -tags='$(BUILD_TAGS)' -race $(TEST) $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT) -parallel=20

cover:
	../scripts/coverage.sh --html

fmt:
	goimports -w $(GOFMT_FILES)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go list -f '{{.Dir}}' ./... \
		| grep -v '.*github.com/hashicorp/vault/api$$' \
		| xargs go vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder
lint:
	@go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs golangci-lint run; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# prep runs `go generate` to build the dynamically generated
# source files.
prep: fmtcheck
	@sh -c "'$(CURDIR)/../scripts/goversioncheck.sh' '$(GO_VERSION_MIN)'"
	@go generate $(go list ./...)

fmtcheck:
	@true
#@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: default prep test testcompile testacc testrace cover fmt fmtcheck vet lint
