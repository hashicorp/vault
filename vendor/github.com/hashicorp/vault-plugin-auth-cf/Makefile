TOOL?=vault-plugin-auth-cf
TEST?=$$(go list ./...)
TESTARGS  ?= '-test.v -count=1'
EXTERNAL_TOOLS=
BUILD_TAGS?=${TOOL}
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# bin generates the releaseable binaries for this plugin
bin: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

default: dev

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin.
dev: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# testshort runs the quick unit tests and vets the code
testshort: fmtcheck generate
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -short -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=4

# test runs the unit tests and vets the code
test: fmtcheck generate
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test ./... -v -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=4

testcompile: fmtcheck generate
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	go generate $(go list ./... | grep -v /vendor/)

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

proto:
	protoc *.proto --go_out=plugins=grpc:.

tools:
	go install ./...

.PHONY: bin default generate test vet bootstrap fmt fmtcheck
