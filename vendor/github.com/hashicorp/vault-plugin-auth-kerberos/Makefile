TOOL?=vault-plugin-auth-kerberos
TEST?=$$(go list ./... | grep -v /vendor/)
EXTERNAL_TOOLS=
BUILD_TAGS?=${TOOL}
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# bin generates the releaseable binaries for this plugin
bin: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

default: dev

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin, except for quickdev which
# is only put into /bin/
quickdev: generate
	@CGO_ENABLED=0 go build -i -tags='$(BUILD_TAGS)' -o bin/vault-plugin-auth-kerberos
dev: fmtcheck generate
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
	go build -o bin/login-kerb ./cmd/login-kerb
dev-linux-only: fmtcheck generate
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEVENV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
	@mkdir -p pkg/linux_amd64; cp bin/$(TOOL) pkg/linux_amd64/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/login-kerb ./cmd/login-kerb
dev-dynamic: generate
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests and vets the code
test: fmtcheck generate
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -v -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=4

test-acceptance: dev-linux-only
	 bats $(CURDIR)/test/acceptance/server-enterprise-basic-tests.bats

testcompile: fmtcheck generate
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# testacc runs acceptance tests
testacc: fmtcheck generate
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout 45m

# generate runs `go generate` to build the dynamically generated
# source files.
generate:
	go generate ./...

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

# builds for linux/amd64, then spins up a few containers for testing
dev-env: dev-linux-only
	./scripts/dev_env.sh

integration: dev-linux-only
	./scripts/integration_env.sh

.PHONY: bin default generate test test-acceptance vet bootstrap fmt fmtcheck
