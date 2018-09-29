TOOL?=vault-plugin-auth-jwt
TEST?=$$(go list ./... | grep -v /vendor/)
EXTERNAL_TOOLS=\
	github.com/mitchellh/gox
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

# test runs all tests
test: generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout 10m

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

# deps updates all dependencies for this project.
deps:
	@echo "==> Updating deps for ${TOOL}"
	@dep ensure -update

.PHONY: bin default generate test bootstrap fmt deps
