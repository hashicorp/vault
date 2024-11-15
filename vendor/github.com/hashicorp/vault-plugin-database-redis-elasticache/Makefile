PLUGIN_NAME := $(shell command ls cmd/)
PLUGIN_DIR := $(GOPATH)/vault-plugins

.PHONY: default
default: dev

.PHONY: dev
dev:
	CGO_ENABLED=0 go build -o bin/$(PLUGIN_NAME) cmd/$(PLUGIN_NAME)/main.go

.PHONY: test
test:
	CGO_ENABLED=0 go test -v ./... $(TESTARGS) -timeout=20m

.PHONY: testacc
testacc:
	ACC_TEST_ENABLED=1 CGO_ENABLED=0 go test ./... $(TESTARGS) -timeout=20m

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: fmt
fmt:
	gofumpt -l -w . && cd bootstrap/terraform && terraform fmt

.PHONY: setup-env
setup-env:
	cd bootstrap/terraform && terraform init && terraform apply -auto-approve

.PHONY: teardown-env
teardown-env:
	cd bootstrap/terraform && terraform init && terraform destroy -auto-approve

.PHONY: configure
configure: dev
	@./bootstrap/configure.sh \
	$(PLUGIN_DIR) \
	$(PLUGIN_NAME) \
	$(TEST_ELASTICACHE_URL) \
	$(TEST_ELASTICACHE_REGION) \
	$(TEST_ELASTICACHE_ACCESS_KEY_ID) \
	$(TEST_ELASTICACHE_SECRET_ACCESS_KEY)
