TOOL?=vault-plugin-secrets-azure
TEST?=$$(go list ./... | grep -v /vendor/)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS=
BUILD_TAGS?=${TOOL}
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
PLUGIN_NAME := $(shell command ls cmd/)
PLUGIN_DIR ?= $$GOPATH/vault-plugins
PLUGIN_PATH ?= local-secrets-azure

# Acceptance test variables
WITH_DEV_PLUGIN?=1
AZURE_TENANT_ID?=
SKIP_TEARDOWN?=
TESTS_OUT_FILE?=
TESTS_FILTER?='.*'

# bin generates the releaseable binaries for this plugin
bin: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

default: dev

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin, except for quickdev which
# is only put into /bin/
quickdev: generate
	@CGO_ENABLED=0 go build -i -tags='$(BUILD_TAGS)' -o bin/vault-plugin-secrets-azure
dev: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-dynamic: generate
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-acceptance: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD= XC_OSARCH=linux/amd64 sh -c "'$(CURDIR)/scripts/build.sh'"

testcompile: fmtcheck generate
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# test runs all unit tests
test: fmtcheck generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC= go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout 10m

testacc: fmtcheck generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout 45m

# test-acceptance runs all acceptance tests
test-acceptance: $(if $(WITH_DEV_PLUGIN), dev-acceptance)
	 WITH_DEV_PLUGIN=$(WITH_DEV_PLUGIN) bats -f $(TESTS_FILTER) $(CURDIR)/tests/acceptance/basic.bats

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

setup-env:
	cd bootstrap/terraform && terraform init && terraform apply -auto-approve

teardown-env:
	cd bootstrap/terraform && terraform init && terraform destroy -auto-approve

configure: dev
	@./bootstrap/configure.sh \
	$(PLUGIN_DIR) \
	$(PLUGIN_NAME) \
	$(PLUGIN_PATH)

.PHONY: bin default generate test vet bootstrap fmt fmtcheck setup-env teardown-env configure
