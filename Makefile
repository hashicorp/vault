# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

TEST?=$$($(GO_CMD) list ./... | grep -v /vendor/ | grep -v /integ)
TEST_TIMEOUT?=45m
EXTENDED_TEST_TIMEOUT=60m
INTEG_TEST_TIMEOUT=120m
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS_CI=\
	golang.org/x/tools/cmd/goimports \
	github.com/golangci/revgrep/cmd/revgrep
EXTERNAL_TOOLS=\
	github.com/client9/misspell/cmd/misspell
GOFMT_FILES?=$$(find . -name '*.go' | grep -v pb.go | grep -v vendor)
SED?=$(shell command -v gsed || command -v sed)


GO_VERSION_MIN=$$(cat $(CURDIR)/.go-version)
PROTOC_VERSION_MIN=3.21.12
GO_CMD?=go
CGO_ENABLED?=0
ifneq ($(FDB_ENABLED), )
	CGO_ENABLED=1
	BUILD_TAGS+=foundationdb
endif

default: dev

# bin generates the releasable binaries for Vault
bin: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS) ui' sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-ui: assetcheck prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS) ui' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-dynamic: prep
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

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
docker-dev: prep
	docker build --build-arg VERSION=$(GO_VERSION_MIN) --build-arg BUILD_TAGS="$(BUILD_TAGS)" -f scripts/docker/Dockerfile -t vault:dev .

docker-dev-ui: prep
	docker build --build-arg VERSION=$(GO_VERSION_MIN) --build-arg BUILD_TAGS="$(BUILD_TAGS)" -f scripts/docker/Dockerfile.ui -t vault:dev-ui .

# test runs the unit tests and vets the code
test: prep
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	VAULT_DEV_ROOT_TOKEN_ID= \
	VAULT_ACC= \
	$(GO_CMD) test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=$(TEST_TIMEOUT) -parallel=20

testcompile: prep
	@for pkg in $(TEST) ; do \
		$(GO_CMD) test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# testacc runs acceptance tests
testacc: prep
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 $(GO_CMD) test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout=$(EXTENDED_TEST_TIMEOUT)

# testrace runs the race checker
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

# godoctests builds the custom analyzer to check for godocs for tests
godoctests:
	@$(GO_CMD) build -o ./tools/godoctests/.bin/godoctests ./tools/godoctests

# vet-godoctests runs godoctests on the test functions. All output gets piped to revgrep
# which will only return an error if a new function is missing a godoc
vet-godoctests: godoctests
	@$(GO_CMD) vet -vettool=./tools/godoctests/.bin/godoctests $(TEST) 2>&1 | revgrep

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder
lint:
	@$(GO_CMD) list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs golangci-lint run; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi
# for ci jobs, runs lint against the changed packages in the commit
ci-lint:
	@golangci-lint run --deadline 10m --new-from-rev=HEAD~

# prep runs `go generate` to build the dynamically generated
# source files.
prep: fmtcheck
	@sh -c "'$(CURDIR)/scripts/goversioncheck.sh' '$(GO_VERSION_MIN)'"
	@$(GO_CMD) generate $($(GO_CMD) list ./... | grep -v /vendor/)
	@if [ -d .git/hooks ]; then cp .hooks/* .git/hooks/; fi

# bootstrap the build by downloading additional tools needed to build
ci-bootstrap:
	@for tool in  $(EXTERNAL_TOOLS_CI) ; do \
		echo "Installing/Updating $$tool" ; \
		GO111MODULE=off $(GO_CMD) get -u $$tool; \
	done

# bootstrap the build by downloading additional tools that may be used by devs
bootstrap: ci-bootstrap
	go generate -tags tools tools/tools.go

# Note: if you have plugins in GOPATH you can update all of them via something like:
# for i in $(ls | grep vault-plugin-); do cd $i; git remote update; git reset --hard origin/master; dep ensure -update; git add .; git commit; git push; cd ..; done
update-plugins:
	grep vault-plugin- go.mod | cut -d ' ' -f 1 | while read -r P; do echo "Updating $P..."; go get -v "$P"; done

static-assets-dir:
	@mkdir -p ./http/web_ui

install-ui-dependencies:
	@echo "--> Installing JavaScript assets"
	@cd ui && yarn --ignore-optional

test-ember: install-ui-dependencies
	@echo "--> Running ember tests"
	@cd ui && yarn run test:oss

test-ember-enos: install-ui-dependencies
	@echo "--> Running ember tests with a real backend"
	@cd ui && yarn run test:enos

check-vault-in-path:
	@VAULT_BIN=$$(command -v vault) || { echo "vault command not found"; exit 1; }; \
		[ -x "$$VAULT_BIN" ] || { echo "$$VAULT_BIN not executable"; exit 1; }; \
		printf "Using Vault at %s:\n\$$ vault version\n%s\n" "$$VAULT_BIN" "$$(vault version)"

ember-dist: install-ui-dependencies
	@cd ui && npm rebuild node-sass
	@echo "--> Building Ember application"
	@cd ui && yarn run build
	@rm -rf ui/if-you-need-to-delete-this-open-an-issue-async-disk-cache

ember-dist-dev: install-ui-dependencies
	@cd ui && npm rebuild node-sass
	@echo "--> Building Ember application"
	@cd ui && yarn run build:dev

static-dist: ember-dist
static-dist-dev: ember-dist-dev

proto: bootstrap
	@sh -c "'$(CURDIR)/scripts/protocversioncheck.sh' '$(PROTOC_VERSION_MIN)'"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative vault/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative vault/activity/activity_log.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helper/storagepacker/types.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helper/forwarding/types.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sdk/logical/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative physical/raft/types.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helper/identity/mfa/types.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helper/identity/types.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sdk/database/dbplugin/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sdk/database/dbplugin/v5/proto/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sdk/plugin/pb/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative vault/tokens/token.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sdk/helper/pluginutil/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative vault/hcp_link/proto/*/*.proto

	# No additional sed expressions should be added to this list. Going forward
	# we should just use the variable names choosen by protobuf. These are left
	# here for backwards compatability, namely for SDK compilation.
	$(SED) -i -e 's/Id/ID/' vault/request_forwarding_service.pb.go
	$(SED) -i -e 's/Idp/IDP/' -e 's/Url/URL/' -e 's/Id/ID/' -e 's/IDentity/Identity/' -e 's/EntityId/EntityID/' -e 's/Api/API/' -e 's/Qr/QR/' -e 's/Totp/TOTP/' -e 's/Mfa/MFA/' -e 's/Pingid/PingID/' -e 's/namespaceId/namespaceID/' -e 's/Ttl/TTL/' -e 's/BoundCidrs/BoundCIDRs/' helper/identity/types.pb.go helper/identity/mfa/types.pb.go helper/storagepacker/types.pb.go sdk/plugin/pb/backend.pb.go sdk/logical/identity.pb.go vault/activity/activity_log.pb.go

	# This will inject the sentinel struct tags as decorated in the proto files.
	protoc-go-inject-tag -input=./helper/identity/types.pb.go
	protoc-go-inject-tag -input=./helper/identity/mfa/types.pb.go

fmtcheck:
	@true
#@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gofumpt -w

semgrep:
	semgrep --include '*.go' --exclude 'vendor' -a -f tools/semgrep .

semgrep-ci:
	semgrep --error --include '*.go' --exclude 'vendor' -f tools/semgrep/ci .

assetcheck:
	@echo "==> Checking compiled UI assets..."
	@sh -c "'$(CURDIR)/scripts/assetcheck.sh'"

spellcheck:
	@echo "==> Spell checking website..."
	@misspell -error -source=text website/source

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

.PHONY: ci-config
ci-config:
	@$(MAKE) -C .circleci ci-config
.PHONY: ci-verify
ci-verify:
	@$(MAKE) -C .circleci ci-verify

.PHONY: bin default prep test vet bootstrap ci-bootstrap fmt fmtcheck mysql-database-plugin mysql-legacy-database-plugin cassandra-database-plugin influxdb-database-plugin postgresql-database-plugin mssql-database-plugin hana-database-plugin mongodb-database-plugin ember-dist ember-dist-dev static-dist static-dist-dev assetcheck check-vault-in-path packages build build-ci semgrep semgrep-ci

.NOTPARALLEL: ember-dist ember-dist-dev

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

.PHONY: ci-filter-matrix
ci-filter-matrix:
	@$(CURDIR)/scripts/ci-helper.sh matrix-filter-file

.PHONY: ci-get-artifact-basename
ci-get-artifact-basename:
	@$(CURDIR)/scripts/ci-helper.sh artifact-basename

.PHONY: ci-get-date
ci-get-date:
	@$(CURDIR)/scripts/ci-helper.sh date

.PHONY: ci-get-matrix-group-id
ci-get-matrix-group-id:
	@$(CURDIR)/scripts/ci-helper.sh matrix-group-id

.PHONY: ci-get-revision
ci-get-revision:
	@$(CURDIR)/scripts/ci-helper.sh revision

.PHONY: ci-get-version
ci-get-version:
	@$(CURDIR)/scripts/ci-helper.sh version

.PHONY: ci-get-version-base
ci-get-version-base:
	@$(CURDIR)/scripts/ci-helper.sh version-base

.PHONY: ci-get-version-major
ci-get-version-major:
	@$(CURDIR)/scripts/ci-helper.sh version-major

.PHONY: ci-get-version-meta
ci-get-version-meta:
	@$(CURDIR)/scripts/ci-helper.sh version-meta

.PHONY: ci-get-version-minor
ci-get-version-minor:
	@$(CURDIR)/scripts/ci-helper.sh version-minor

.PHONY: ci-get-version-package
ci-get-version-package:
	@$(CURDIR)/scripts/ci-helper.sh version-package

.PHONY: ci-get-version-patch
ci-get-version-patch:
	@$(CURDIR)/scripts/ci-helper.sh version-patch

.PHONY: ci-get-version-pre
ci-get-version-pre:
	@$(CURDIR)/scripts/ci-helper.sh version-pre

.PHONY: ci-prepare-legal
ci-prepare-legal:
	@$(CURDIR)/scripts/ci-helper.sh prepare-legal
