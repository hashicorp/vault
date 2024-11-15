GO_CMD?=go
CGO_ENABLED?=0
TOOL?=vault-plugin-secrets-terraform
TEST?=$$($(GO_CMD) list ./... | grep -v /vendor/ | grep -v /integ)
EXTERNAL_TOOLS=
BUILD_TAGS?=${TOOL}
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# bin generates the releaseable binaries for this plugin
bin: generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

default: dev

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin, except for quickdev which
# is only put into /bin/
quickdev: generate
	@CGO_ENABLED=0 go build -i -tags='$(BUILD_TAGS)' -o bin/${TOOL}

dev: generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

testcompile: generate
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# test runs the unit tests and vets the code
test: generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -v -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=10m -parallel=4

# testacc runs acceptance tests
testacc:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	CGO_ENABLED=0 VAULT_ACC=1 VAULT_TOKEN= $(GO_CMD) test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout=10m

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	@go generate $(go list ./... | grep -v /vendor/)

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: bin default generate test bootstrap fmt deps
