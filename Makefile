# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))

export RELEASE_GPG_KEY_FINGERPRINT := 91A6 E7F8 5D05 C656 30BE  F189 5185 2D87 348F FC4C

TEST?=$$($(GO_CMD) list ./... | grep -v /vendor/ | grep -v /integ)
TEST_TIMEOUT?=45m
EXTENDED_TEST_TIMEOUT=60m
INTEG_TEST_TIMEOUT=120m
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS_CI=\
	github.com/elazarl/go-bindata-assetfs/... \
	github.com/hashicorp/go-bindata/... \
	github.com/mitchellh/gox \
	golang.org/x/tools/cmd/goimports 
EXTERNAL_TOOLS=\
	github.com/client9/misspell/cmd/misspell
GOFMT_FILES?=$$(find . -name '*.go' | grep -v pb.go | grep -v vendor)


GO_VERSION_MIN=1.14.7
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
	docker build -f scripts/docker/Dockerfile -t vault:dev .

docker-dev-ui: prep
	docker build -f scripts/docker/Dockerfile.ui -t vault:dev-ui .

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
	@# Remove old (now broken) husky git hooks.
	@[ ! -d .git/hooks ] || grep -l '^# husky$$' .git/hooks/* | xargs rm -f
	@if [ -d .git/hooks ]; then cp .hooks/* .git/hooks/; fi

PACKAGES_LOCK_DIR := $(shell find . -mindepth 1 -maxdepth 1 \
	-type d -name 'packages*.lock')

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
	grep vault-plugin- vendor/vendor.json | cut -d '"' -f 4 | xargs govendor fetch

static-assets-dir:
	@mkdir -p ./pkg/web_ui

static-assets: static-assets-dir
	@echo "--> Generating static assets"
	@go-bindata-assetfs -o bindata_assetfs.go -pkg http -prefix pkg -modtime 1480000000 -tags ui ./pkg/web_ui/...
	@mv bindata_assetfs.go http
	@$(MAKE) -f $(THIS_FILE) fmt

test-ember:
	@echo "--> Installing JavaScript assets"
	@cd ui && yarn --ignore-optional
	@echo "--> Running ember tests"
	@cd ui && yarn run test:oss

ember-ci-test: # Deprecated, to be removed soon.
	@echo "ember-ci-test is deprecated in favour of test-ui-browserstack"
	@exit 1

check-vault-in-path:
	@VAULT_BIN=$$(command -v vault) || { echo "vault command not found"; exit 1; }; \
		[ -x "$$VAULT_BIN" ] || { echo "$$VAULT_BIN not executable"; exit 1; }; \
		printf "Using Vault at %s:\n\$$ vault version\n%s\n" "$$VAULT_BIN" "$$(vault version)"

check-browserstack-creds:
	@[ -n "$$BROWSERSTACK_ACCESS_KEY" ] || { echo "BROWSERSTACK_ACCESS_KEY not set"; exit 1; }
	@[ -n "$$BROWSERSTACK_USERNAME" ] || { echo "BROWSERSTACK_USERNAME not set"; exit 1; }

test-ui-browserstack: check-vault-in-path check-browserstack-creds
	@echo "--> Installing JavaScript assets"
	@cd ui && yarn --ignore-optional
	@echo "--> Running ember tests in Browserstack"
	@cd ui && yarn run test:browserstack

ember-dist:
	@echo "--> Installing JavaScript assets"
	@cd ui && yarn --ignore-optional
	@cd ui && npm rebuild node-sass
	@echo "--> Building Ember application"
	@cd ui && yarn run build
	@rm -rf ui/if-you-need-to-delete-this-open-an-issue-async-disk-cache

ember-dist-dev:
	@echo "--> Installing JavaScript assets"
	@cd ui && yarn --ignore-optional
	@cd ui && npm rebuild node-sass
	@echo "--> Building Ember application"
	@cd ui && yarn run build-dev

static-dist: ember-dist static-assets
static-dist-dev: ember-dist-dev static-assets

proto:
	protoc vault/*.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc helper/storagepacker/types.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc helper/forwarding/types.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc sdk/logical/*.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc physical/raft/types.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc helper/identity/mfa/types.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc helper/identity/types.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc sdk/database/dbplugin/*.proto --go_out=plugins=grpc,paths=source_relative:.
	protoc sdk/plugin/pb/*.proto --go_out=plugins=grpc,paths=source_relative:.
	sed -i -e 's/Id/ID/' vault/request_forwarding_service.pb.go
	sed -i -e 's/Idp/IDP/' -e 's/Url/URL/' -e 's/Id/ID/' -e 's/IDentity/Identity/' -e 's/EntityId/EntityID/' -e 's/Api/API/' -e 's/Qr/QR/' -e 's/Totp/TOTP/' -e 's/Mfa/MFA/' -e 's/Pingid/PingID/' -e 's/protobuf:"/sentinel:"" protobuf:"/' -e 's/namespaceId/namespaceID/' -e 's/Ttl/TTL/' -e 's/BoundCidrs/BoundCIDRs/' helper/identity/types.pb.go helper/identity/mfa/types.pb.go helper/storagepacker/types.pb.go sdk/plugin/pb/backend.pb.go sdk/logical/identity.pb.go 

fmtcheck:
	@true
#@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	goimports -w $(GOFMT_FILES)

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

# GPG_KEY_VARS sets FPR to the fingerprint with no spaces and GIT_GPG_KEY_ID to a git compatible gpg key id.
define GPG_KEY_VARS
	FPR=$$(echo "$(RELEASE_GPG_KEY_FINGERPRINT)" | sed 's/ //g') && GIT_GPG_KEY_ID="0x$$(printf $$FPR | tail -c 16)"
endef

define FAIL_IF_GPG_KEY_MISSING
	@$(GPG_KEY_VARS) && echo "Checking for key '$$FPR' (git key id: $$GIT_GPG_KEY_ID)"; \
	gpg --list-secret-keys "$$FPR" >/dev/null 2>&1 || { \
		echo "ERROR: Secret GPG key missing: $$FPR"; \
		exit 1; \
	}
endef

define FAIL_IF_UNCOMMITTED_CHANGES
	@{ git diff --exit-code && git diff --cached --exit-code; } >/dev/null 2>&1 || { \
		echo "ERROR: Uncommitted changes detected."; \
		exit 1; \
	}
endef

stage-commit:
	$(FAIL_IF_GPG_KEY_MISSING)
	$(FAIL_IF_UNCOMMITTED_CHANGES)
	@[ -n "$(STAGE_VERSION)" ] || { echo "You must set STAGE_VERSION to the version in semver-like format."; exit 1; }
	set -x; $(GPG_KEY_VARS) && git commit --allow-empty --gpg-sign=$$GIT_GPG_KEY_ID -m 'release: stage v$(STAGE_VERSION)' 

publish-commit:
	$(FAIL_IF_GPG_KEY_MISSING)
	$(FAIL_IF_UNCOMMITTED_CHANGES)
	@[ -n "$(PUBLISH_VERSION)" ] || { echo "You must set PUBLISH_VERSION to the version in semver-like format."; exit 1; }
	set -x; $(GPG_KEY_VARS) && git commit --allow-empty --gpg-sign=$$GIT_GPG_KEY_ID -m 'release: publish v$(PUBLISH_VERSION)'

# WRITE_GENERATED_FILE_HEADER overwrites the file specified, replacing its contents with
# the header warning people not to attempt to edit or merge the file. You should call this
# before writing the generated contents to the file.
# Args: 1: File to write; 2: Command to generate it; 3: Source files to edit/merge instead.
define WRITE_GENERATED_FILE_HEADER
	echo "### ***" > $(1); \
	echo "### WARNING: DO NOT manually EDIT or MERGE this file, it is generated by '$(2)'." >> $(1); \
	echo "### INSTEAD: Edit or merge the source in $(3) then run '$(2)'." >> $(1); \
	echo "### ***" >> $(1)
endef

## begin packagespec integration ##

# The plan is to generate this packagespec integration section automatically.
# By keeping it in a contiguous block for now, it will be easier to
# auto-generate when we get to it.

SPEC_FILE_PATTERN := packages*.yml
# SPEC is the human-managed description of which packages we are able to build.
SPEC := $(shell find . -mindepth 1 -maxdepth 1 -name '$(SPEC_FILE_PATTERN)')
ifneq ($(words $(SPEC)),1)
$(error Found $(words $(SPEC)) $(SPEC_FILE_PATTERN) files, need exactly 1: $(SPEC))
endif
SPEC_FILENAME := $(notdir $(SPEC))
SPEC_MODIFIER := $(SPEC_FILENAME:packages%.yml=%)
# LOCKDIR contains the lockfile and layer files.
LOCKDIR  := packages$(SPEC_MODIFIER).lock
LOCKFILE := $(LOCKDIR)/pkgs.yml

export PACKAGE_SPEC_ID LAYER_SPEC_ID PRODUCT_REVISION PRODUCT_VERSION

# PACKAGESPEC_TARGETS are convenience aliases for targets defined in $(LOCKDIR)/Makefile
PACKAGESPEC_TARGETS := \
	build build-all build-ci \
	aliases meta package \
	package-meta stage-config stage \
	watch-ci publish-config publish list-staged-builds

$(PACKAGESPEC_TARGETS): 
$(PACKAGESPEC_TARGETS):
	@PRODUCT_REPO_ROOT="$(shell git rev-parse --show-toplevel)" $(MAKE) -C $(LOCKDIR) $@

# packages regenerates $(LOCKDIR) from $(SPEC) using packagespec. This is only for
# internal HashiCorp use, as it has dependencies not available to OSS contributors.
packages:
	@command -v packagespec > /dev/null 2>&1 || { \
		echo "Please install packagespec."; \
		echo "Note: packagespec is only available to HashiCorp employees at present."; \
		exit 1; \
	}
	@packagespec lock -specfile $(SPEC) -lockdir $(LOCKDIR)
	@$(MAKE) ci-config

.PHONY: $(PACKAGESPEC_TARGETS) packages
## end packagespec integration ##

CI_WORKFLOW_TPL     := .circleci/config/@build-release.yml.tpl
CI_WORKFLOW         := .circleci/config/@build-release.yml

.PHONY: ci-update-release-packages $(CI_WORKFLOW)
ci-update-release-packages: $(CI_WORKFLOW)
	@echo $^

$(CI_WORKFLOW): $(LOCKFILE) $(CI_WORKFLOW_TPL)
	@\
		echo "==> Updating $@ to match $<"; \
		$(call WRITE_GENERATED_FILE_HEADER,$@,make $@,$^); \
		cat $< | gomplate -f $(CI_WORKFLOW_TPL) -d 'package-list=stdin://?type=application/yaml' >> $@

.PHONY: ci-config
ci-config: ci-update-release-packages
	@$(MAKE) -C .circleci ci-config
.PHONY: ci-verify
ci-verify:
	@$(MAKE) -C .circleci ci-verify

.PHONY: bin default prep test vet bootstrap ci-bootstrap fmt fmtcheck mysql-database-plugin mysql-legacy-database-plugin cassandra-database-plugin influxdb-database-plugin postgresql-database-plugin mssql-database-plugin hana-database-plugin mongodb-database-plugin static-assets ember-dist ember-dist-dev static-dist static-dist-dev assetcheck check-vault-in-path check-browserstack-creds test-ui-browserstack stage-commit publish-commit

.NOTPARALLEL: ember-dist ember-dist-dev static-assets
