# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
REPO_DIR := $(shell basename $(CURDIR))

PLUGIN_NAME := $(shell command ls cmd/)

.PHONY: default
default: dev

.PHONY: dev
dev:
	@echo "Skip"

.PHONY: test
test: fmtcheck
	CGO_ENABLED=0 go test ./... $(TESTARGS) -timeout=20m

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: fmt
fmt:
	gofumpt -l -w .

mocks:
	go install github.com/vektra/mockery/v2@v2.34.2
	mockery --srcpkg github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service --name=ClientService
	mockery --srcpkg github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service --name=ClientService
	mockery --srcpkg github.com/hashicorp/hcp-sdk-go/clients/cloud-iam/stable/2019-12-10/client/iam_service --name=ClientService
