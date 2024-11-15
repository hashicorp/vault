PROJECT="github.com/hashicorp/vault-plugin-secrets-gcpkms"
PLUGIN_NAME=$(shell go run version/cmd/main.go name)
VERSION=$(shell go run version/cmd/main.go version)
COMMIT=$(shell git rev-parse --short HEAD)

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# default is the default make command
default: test

fmt:
	gofmt -w $(GOFMT_FILES)

# deps updates the project deps using golang/dep
deps:
	@dep ensure -v -update
.PHONY: deps

# dev builds the plugin for local development
dev:
	CGO_ENABLED=0 go build -o bin/$(PLUGIN_NAME) cmd/$(PLUGIN_NAME)/main.go
.PHONY: dev

# test runs the tests
test:
	@go test -timeout=240s ./... $(TESTARGS)
.PHONY: test
