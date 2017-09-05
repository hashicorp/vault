TEST?=$$(go list ./... | grep -v /vendor/)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
EXTERNAL_TOOLS=\
	github.com/mitchellh/gox \
	github.com/kardianos/govendor
BUILD_TAGS?=vault
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: dev

# bin generates the releaseable binaries for Vault
bin: fmtcheck prep
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin, except for quickdev which
# is only put into /bin/
quickdev: prep
	@CGO_ENABLED=0 go build -i -tags='$(BUILD_TAGS)' -o bin/vault
dev: fmtcheck prep
	@CGO_ENABLED=0 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"
dev-dynamic: prep
	@CGO_ENABLED=1 BUILD_TAGS='$(BUILD_TAGS)' VAULT_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests and vets the code
test: fmtcheck prep
	CGO_ENABLED=0 VAULT_TOKEN= VAULT_ACC= go test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -timeout=20m -parallel=4

testcompile: fmtcheck prep
	@for pkg in $(TEST) ; do \
		go test -v -c -tags='$(BUILD_TAGS)' $$pkg -parallel=4 ; \
	done

# testacc runs acceptance tests
testacc: fmtcheck prep
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package"; \
		exit 1; \
	fi
	VAULT_ACC=1 go test -tags='$(BUILD_TAGS)' $(TEST) -v $(TESTARGS) -timeout 45m

# testrace runs the race checker
testrace: fmtcheck prep
	CGO_ENABLED=1 VAULT_TOKEN= VAULT_ACC= go test -tags='$(BUILD_TAGS)' -race $(TEST) $(TESTARGS) -timeout=45m -parallel=4

cover:
	./scripts/coverage.sh --html

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| grep -v '.*github.com/hashicorp/vault$$' \
		| xargs go tool vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# prep runs `go generate` to build the dynamically generated
# source files.
prep:
	go generate $(go list ./... | grep -v /vendor/)
	cp .hooks/* .git/hooks/

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done

proto:
	protoc -I helper/forwarding -I vault -I ../../.. vault/*.proto --go_out=plugins=grpc:vault
	protoc -I helper/forwarding -I vault -I ../../.. helper/forwarding/types.proto --go_out=plugins=grpc:helper/forwarding

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

mysql-database-plugin:
	@CGO_ENABLED=0 go build -o bin/mysql-database-plugin ./plugins/database/mysql/mysql-database-plugin

mysql-legacy-database-plugin:
	@CGO_ENABLED=0 go build -o bin/mysql-legacy-database-plugin ./plugins/database/mysql/mysql-legacy-database-plugin

cassandra-database-plugin:
	@CGO_ENABLED=0 go build -o bin/cassandra-database-plugin ./plugins/database/cassandra/cassandra-database-plugin

postgresql-database-plugin:
	@CGO_ENABLED=0 go build -o bin/postgresql-database-plugin ./plugins/database/postgresql/postgresql-database-plugin

mssql-database-plugin:
	@CGO_ENABLED=0 go build -o bin/mssql-database-plugin ./plugins/database/mssql/mssql-database-plugin

hana-database-plugin:
	@CGO_ENABLED=0 go build -o bin/hana-database-plugin ./plugins/database/hana/hana-database-plugin

mongodb-database-plugin:
	@CGO_ENABLED=0 go build -o bin/mongodb-database-plugin ./plugins/database/mongodb/mongodb-database-plugin

.PHONY: bin default prep test vet bootstrap fmt fmtcheck mysql-database-plugin mysql-legacy-database-plugin cassandra-database-plugin postgresql-database-plugin mssql-database-plugin hana-database-plugin mongodb-database-plugin
