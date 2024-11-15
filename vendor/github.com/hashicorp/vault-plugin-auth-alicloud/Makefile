TOOL?=vault-plugin-auth-alicloud
TEST?=$$(go list ./... | grep -v /vendor/ | grep -v teamcity)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
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
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -short -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=20m -parallel=4

# test runs the unit tests and vets the code
test: fmtcheck generate
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -v -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=20m -parallel=4

# testacc runs the acceptance tests and vets the code
testacc: fmtcheck generate
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC=1 go test -v -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=20m -parallel=4

testcompile: fmtcheck generate
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	go generate $(go list ./... | grep -v /vendor/)


fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

proto:
	protoc *.proto --go_out=plugins=grpc:.

.PHONY: bin default generate test vet fmt fmtcheck
