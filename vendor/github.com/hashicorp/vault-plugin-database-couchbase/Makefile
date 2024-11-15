TOOL?=vault-plugin-database-couchbase
TEST?=$$(go list ./... | grep -v /vendor/ | grep -v teamcity)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
BUILD_TAGS?=${TOOL}
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GO_TEST_CMD?=go test -v

# bin generates the releaseable binaries for this plugin
bin: fmtcheck
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

default: dev

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin.
dev: fmtcheck
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# dev-vault starts up `vault` from your $PATH, then builds the couchbase
# plugin, registers it with vault and enables it.
# A ./tmp dir is created for configs and binaries, and cleaned up on exit.
dev-vault: fmtcheck
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build_with_vault.sh'"

# test runs the unit tests and vets the code
test: fmtcheck
	CGO_ENABLED=0 VAULT_TOKEN= ${GO_TEST_CMD} -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=5m -parallel=4

testacc: fmtcheck
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC=1 ${GO_TEST_CMD} -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -count=1 -timeout=20m

testcompile: fmtcheck
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg ; \
	done

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: bin default dev dev-vault test testacc testcompile fmtcheck fmt
