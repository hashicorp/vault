TEST?=$$(go list ./... | grep -v /vendor/ | grep -v /integ)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
EXTERNAL_TOOLS=

default: dev

# bin generates the releaseable binaries for vault-plugin-database-elasticsearch
bin: fmtcheck generate
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' XC_ARCH="amd64" XC_OS="linux" XC_OSARCH="linux/amd64" sh -c "'$(CURDIR)/scripts/build.sh'"

dev: fmtcheck
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests and vets the code
test: fmtcheck generate
	CGO_ENABLED=1 go test -v -short -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=1

testacc: fmtcheck generate
	VAULT_ACC=1 VAULT_ADDR=http://localhost:8200 VAULT_TOKEN=root CGO_ENABLED=1 go test -v -race -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=1

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w -s $(GOFMT_FILES)

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done

.PHONY: bin default generate test testacc fmt fmtcheck dev bootstrap
