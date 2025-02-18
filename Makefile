# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

MAIN_PACKAGES=$$($(GO_CMD) list ./... | grep -v vendor/ )
SDK_PACKAGES=$$(cd $(CURDIR)/sdk && $(GO_CMD) list ./... | grep -v vendor/ )
API_PACKAGES=$$(cd $(CURDIR)/api && $(GO_CMD) list ./... | grep -v vendor/ )
ALL_PACKAGES=$(MAIN_PACKAGES) $(SDK_PACKAGES) $(API_PACKAGES)
TEST=$$(echo $(ALL_PACKAGES) | grep -v integ/ )
TEST_TIMEOUT?=45m
EXTENDED_TEST_TIMEOUT=60m
INTEG_TEST_TIMEOUT=120m
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
GOFMT_FILES?=$$(find . -name '*.go' | grep -v pb.go | grep -v vendor)
SED?=$(shell command -v gsed || command -v sed)

GO_VERSION_MIN=$$(cat $(CURDIR)/.go-version)
GO_CMD?=go
CGO_ENABLED?=0
ifneq ($(FDB_ENABLED), )
	CGO_ENABLED=1
	BUILD_TAGS+=foundationdb
endif

# Set BUILD_MINIMAL to a non-empty value to build a minimal version of Vault with only core features.
BUILD_MINIMAL ?=
ifneq ($(strip $(BUILD_MINIMAL)),)
	BUILD_TAGS+=minimal
endif

default: dev

# bin generates the releasable binaries for Vault
bin: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS) ui' sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: BUILD_TAGS+=testonly
dev: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-ui: BUILD_TAGS+=testonly
dev-ui: assetcheck prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS) ui' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-dynamic: BUILD_TAGS+=testonly
dev-dynamic: prep
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# quickdev creates binaries for testing Vault locally like dev, but skips
# the prep step.
quickdev: BUILD_TAGS+=testonly
quickdev:
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# *-mem variants will enable memory profiling which will write snapshots of heap usage
# to $TMP/vaultprof every 5 minutes. These can be analyzed using `$ go tool pprof <profile_file>`.
# Note that any build can have profiling added via: `$ BUILD_TAGS=memprofiler make ...`
dev-mem: BUILD_TAGS+=memprofiler
dev-mem: dev
dev-ui-mem: BUILD_TAGS+=memprofiler
dev-ui-mem: assetcheck dev-ui
dev-dynamic-mem: BUILD_TAGS+=memprofiler
dev-dynamic-mem: dev-dynamic

# Creates a Docker image by adding the compiled linux/amd64 binary found in ./bin.
# The resulting image is tagged "vault:dev".
docker-dev: BUILD_TAGS+=testonly
docker-dev: prep
	docker build --build-arg VERSION=$(GO_VERSION_MIN) --build-arg BUILD_TAGS="$(BUILD_TAGS)" -f scripts/docker/Dockerfile -t vault:dev .

docker-dev-ui: BUILD_TAGS+=testonly
docker-dev-ui: prep
	docker build --build-arg VERSION=$(GO_VERSION_MIN) --build-arg BUILD_TAGS="$(BUILD_TAGS)" -f scripts/docker/Dockerfile.ui -t vault:dev-ui .

# test runs the unit tests and vets the code
test: BUILD_TAGS+=testonly
test: prep
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	$(GO_CMD) test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=$(TEST_TIMEOUT) -parallel=20

testcompile: BUILD_TAGS+=testonly
testcompile: prep
	@for pkg in $(TEST) ; do \
		$(GO_CMD) test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# testacc runs acceptance tests
testacc: BUILD_TAGS+=testonly
testacc: prep
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 $(GO_CMD) test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT)

# testrace runs the race checker
testrace: BUILD_TAGS+=testonly
testrace: prep
	@CGO_ENABLED=1 \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	$(GO_CMD) test -tags='$(BUILD_TAGS)' -race $(TEST) $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT) -parallel=20

cover:
	./scripts/coverage.sh --html

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@$(GO_CMD) list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| grep -v '.*github.com/hashicorp/vault$$' \
		| xargs $(GO_CMD) vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# deprecations runs staticcheck tool to look for deprecations. Checks entire code to see if it
# has deprecated function, variable, constant or field
deprecations: bootstrap prep
	@BUILD_TAGS='$(BUILD_TAGS)' ./scripts/deprecations-checker.sh ""

# ci-deprecations runs staticcheck tool to look for deprecations. All output gets piped to revgrep
# which will only return an error if changes that is not on main has deprecated function, variable, constant or field
ci-deprecations: prep check-tools-external
	@BUILD_TAGS='$(BUILD_TAGS)' ./scripts/deprecations-checker.sh main

# vet-codechecker runs our custom linters on the test functions. All output gets
# piped to revgrep which will only return an error if new piece of code violates
# the check
vet-codechecker: check-tools-internal
	@echo "==> Running go vet with ./tools/codechecker..."
	@$(GO_CMD) vet -vettool=$$(which codechecker) -tags=$(BUILD_TAGS) ./... 2>&1 | revgrep

# vet-codechecker runs our custom linters on the test functions. All output gets
# piped to revgrep which will only return an error if new piece of code that is
# not on main violates the check
ci-vet-codechecker: tools-internal check-tools-external
	@echo "==> Running go vet with ./tools/codechecker..."
	@$(GO_CMD) vet -vettool=$$(which codechecker) -tags=$(BUILD_TAGS) ./... 2>&1 | revgrep origin/main

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder
lint: check-tools-external
	@$(GO_CMD) list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs golangci-lint run; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# for ci jobs, runs lint against the changed packages in the commit
ci-lint: check-tools-external
	@golangci-lint run --deadline 10m --new-from-rev=HEAD~

# Lint protobuf files
protolint: prep check-tools-external
	@echo "==> Linting protobufs..."
	@buf lint

# prep runs `go generate` to build the dynamically generated
# source files.
#
# n.b.: prep used to depend on fmtcheck, but since fmtcheck is
# now run as a pre-commit hook (and there's little value in
# making every build run the formatter), we've removed that
# dependency.
prep: check-go-version clean
	@echo "==> Running go generate..."
	@GOARCH= GOOS= $(GO_CMD) generate $(MAIN_PACKAGES)
	@GOARCH= GOOS= cd api && $(GO_CMD) generate $(API_PACKAGES)
	@GOARCH= GOOS= cd sdk && $(GO_CMD) generate $(SDK_PACKAGES)

# Git doesn't allow us to store shared hooks in .git. Instead, we make sure they're up-to-date
# whenever a make target is invoked.
.PHONY: hooks
hooks:
	@if [ -d .git/hooks ]; then cp .hooks/* .git/hooks/; fi

-include hooks # Make sure they're always up-to-date

# bootstrap the build by generating any necessary code and downloading additional tools that may
# be used by devs.
bootstrap: tools prep

# Note: if you have plugins in GOPATH you can update all of them via something like:
# for i in $(ls | grep vault-plugin-); do cd $i; git remote update; git reset --hard origin/master; dep ensure -update; git add .; git commit; git push; cd ..; done
update-plugins:
	grep vault-plugin- go.mod | cut -d ' ' -f 1 | while read -r P; do echo "Updating $P..."; go get -v "$P"; done

static-assets-dir:
	@mkdir -p ./http/web_ui

install-ui-dependencies:
	@echo "==> Installing JavaScript assets"
	@cd ui && yarn

test-ember: install-ui-dependencies
	@echo "==> Running ember tests"
	@cd ui && yarn run test:oss

test-ember-enos: install-ui-dependencies
	@echo "==> Running ember tests with a real backend"
	@cd ui && yarn run test:enos

ember-dist: install-ui-dependencies
	@cd ui && npm rebuild node-sass
	@echo "==> Building Ember application"
	@cd ui && yarn run build
	@rm -rf ui/if-you-need-to-delete-this-open-an-issue-async-disk-cache

ember-dist-dev: install-ui-dependencies
	@cd ui && npm rebuild node-sass
	@echo "==> Building Ember application"
	@cd ui && yarn run build:dev

static-dist: ember-dist
static-dist-dev: ember-dist-dev

proto: check-tools-external
	@echo "==> Generating Go code from protobufs..."
	buf generate

	# No additional sed expressions should be added to this list. Going forward
	# we should just use the variable names choosen by protobuf. These are left
	# here for backwards compatibility, namely for SDK compilation.
	$(SED) -i -e 's/Id/ID/' -e 's/SPDX-License-IDentifier/SPDX-License-Identifier/' vault/request_forwarding_service.pb.go
	$(SED) -i -e 's/Idp/IDP/' -e 's/Url/URL/' -e 's/Id/ID/' -e 's/IDentity/Identity/' -e 's/EntityId/EntityID/' -e 's/Api/API/' -e 's/Qr/QR/' -e 's/Totp/TOTP/' -e 's/Mfa/MFA/' -e 's/Pingid/PingID/' -e 's/namespaceId/namespaceID/' -e 's/Ttl/TTL/' -e 's/BoundCidrs/BoundCIDRs/' -e 's/SPDX-License-IDentifier/SPDX-License-Identifier/' helper/identity/types.pb.go helper/identity/mfa/types.pb.go helper/storagepacker/types.pb.go sdk/plugin/pb/backend.pb.go sdk/logical/identity.pb.go vault/activity/activity_log.pb.go

	# This will inject the sentinel struct tags as decorated in the proto files.
	protoc-go-inject-tag -input=./helper/identity/types.pb.go
	protoc-go-inject-tag -input=./helper/identity/mfa/types.pb.go

importfmt: check-tools-external
	find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gosimports -w

fmt: importfmt
	find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gofumpt -w

fmtcheck: check-go-fmt

.PHONY: go-mod-download
go-mod-download:
	@$(CURDIR)/scripts/go-helper.sh mod-download

.PHONY: go-mod-tidy
go-mod-tidy:
	@$(CURDIR)/scripts/go-helper.sh mod-tidy

protofmt:
	buf format -w

semgrep:
	semgrep --include '*.go' --exclude 'vendor' -a -f tools/semgrep .

assetcheck:
	@echo "==> Checking compiled UI assets..."
	@sh -c "'$(CURDIR)/scripts/assetcheck.sh'"

spellcheck:
	@echo "==> Spell checking website..."
	@misspell -error -source=text website/source

.PHONY check-go-fmt:
check-go-fmt:
	@$(CURDIR)/scripts/go-helper.sh check-fmt

.PHONY check-go-version:
check-go-version:
	@$(CURDIR)/scripts/go-helper.sh check-version $(GO_VERSION_MIN)

.PHONY: check-proto-fmt
check-proto-fmt:
	buf format -d --error-format github-actions --exit-code

.PHONY: check-proto-delta
check-proto-delta: prep
	@echo "==> Checking for a delta in proto generated Go files..."
	@echo "==> Deleting all *.pg.go files..."
	find . -type f -name '*.pb.go' -delete -print0
	@$(MAKE) -f $(THIS_FILE) proto
	@if ! git diff --exit-code; then echo "Go protobuf bindings need to be regenerated. Run 'make proto' to fix them." && exit 1; fi

.PHONY:check-sempgrep
check-sempgrep: check-tools-external
	@echo "==> Checking semgrep..."
	@semgrep --error --include '*.go' --exclude 'vendor' -f tools/semgrep/ci .

.PHONY: check-tools
check-tools:
	@$(CURDIR)/tools/tools.sh check

.PHONY: check-tools-external
check-tools-external:
	@$(CURDIR)/tools/tools.sh check-external

.PHONY: check-tools-internal
check-tools-internal:
	@$(CURDIR)/tools/tools.sh check-internal

.PHONY: check-tools-pipeline
check-tools-pipeline:
	@$(CURDIR)/tools/tools.sh check-pipeline

check-vault-in-path:
	@VAULT_BIN=$$(command -v vault) || { echo "vault command not found"; exit 1; }; \
		[ -x "$$VAULT_BIN" ] || { echo "$$VAULT_BIN not executable"; exit 1; }; \
		printf "Using Vault at %s:\n\$$ vault version\n%s\n" "$$VAULT_BIN" "$$(vault version)"

.PHONY: tools
tools:
	@$(CURDIR)/tools/tools.sh install

.PHONY: tools-external
tools-external:
	@$(CURDIR)/tools/tools.sh install-external

.PHONY: tools-internal
tools-internal:
	@$(CURDIR)/tools/tools.sh install-internal

.PHONY: tools-pipeline
tools-pipeline:
	@$(CURDIR)/tools/tools.sh install-pipeline

mysql-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/mysql-database-plugin ./plugins/database/mysql/mysql-database-plugin

mysql-legacy-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/mysql-legacy-database-plugin ./plugins/database/mysql/mysql-legacy-database-plugin

cassandra-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/cassandra-database-plugin ./plugins/database/cassandra/cassandra-database-plugin

influxdb-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/influxdb-database-plugin ./plugins/database/influxdb/influxdb-database-plugin

postgresql-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/postgresql-database-plugin ./plugins/database/postgresql/postgresql-database-plugin

mssql-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/mssql-database-plugin ./plugins/database/mssql/mssql-database-plugin

hana-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/hana-database-plugin ./plugins/database/hana/hana-database-plugin

mongodb-database-plugin:
	@CGO_ENABLED=0 $(GO_CMD) build -o bin/mongodb-database-plugin ./plugins/database/mongodb/mongodb-database-plugin

# These ci targets are used for used for building and testing in Github Actions
# workflows and for Enos scenarios.
.PHONY: ci-build
ci-build:
	@$(CURDIR)/scripts/ci-helper.sh build

.PHONY: ci-build-ui
ci-build-ui:
	@$(CURDIR)/scripts/ci-helper.sh build-ui

.PHONY: ci-bundle
ci-bundle:
	@$(CURDIR)/scripts/ci-helper.sh bundle

.PHONY: ci-get-artifact-basename
ci-get-artifact-basename:
	@$(CURDIR)/scripts/ci-helper.sh artifact-basename

.PHONY: ci-get-date
ci-get-date:
	@$(CURDIR)/scripts/ci-helper.sh date

.PHONY: ci-get-revision
ci-get-revision:
	@$(CURDIR)/scripts/ci-helper.sh revision

.PHONY: ci-get-version-package
ci-get-version-package:
	@$(CURDIR)/scripts/ci-helper.sh version-package

.PHONY: ci-prepare-ent-legal
ci-prepare-ent-legal:
	@$(CURDIR)/scripts/ci-helper.sh prepare-ent-legal

.PHONY: ci-prepare-ce-legal
ci-prepare-ce-legal:
	@$(CURDIR)/scripts/ci-helper.sh prepare-ce-legal

.PHONY: ci-copywriteheaders
ci-copywriteheaders:
	copywrite headers --plan
	# Special case for MPL headers in /api, /sdk, and /shamir
	cd api && $(CURDIR)/scripts/copywrite-exceptions.sh
	cd sdk && $(CURDIR)/scripts/copywrite-exceptions.sh
	cd shamir && $(CURDIR)/scripts/copywrite-exceptions.sh

.PHONY: all bin default prep test vet bootstrap fmt fmtcheck mysql-database-plugin mysql-legacy-database-plugin cassandra-database-plugin influxdb-database-plugin postgresql-database-plugin mssql-database-plugin hana-database-plugin mongodb-database-plugin ember-dist ember-dist-dev static-dist static-dist-dev assetcheck check-vault-in-path packages build build-ci semgrep semgrep-ci vet-codechecker ci-vet-codechecker clean dev

.NOTPARALLEL: ember-dist ember-dist-dev

.PHONY: all-packages
all-packages:
	@echo $(ALL_PACKAGES) | tr ' ' '\n'

.PHONY: clean
clean:
	@echo "==> Cleaning..."
