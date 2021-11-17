GOMAXPROCS = 4

PROJECT    = "github.com/hashicorp/vault-plugin-secrets-gcpkms"
NAME       = $(shell go run version/cmd/main.go name)
VERSION    = $(shell go run version/cmd/main.go version)
COMMIT     = $(shell git rev-parse --short HEAD)

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

LDFLAGS = \
	-s \
	-w \
	-X ${PROJECT}/version.GitCommit=${COMMIT}

# XC_* are the platforms for cross-compiling. Customize these values to suit
# your needs.
XC_OS      = darwin linux windows
XC_ARCH    = 386 amd64
XC_EXCLUDE =

# default is the default make command
default: test

fmt:
	gofmt -w $(GOFMT_FILES)

# deps updates the project deps using golang/dep
deps:
	@dep ensure -v -update
.PHONY: deps

# dev builds and installs the plugin for local development
dev:
	@env \
		CGO_ENABLED=0 \
		go install \
			-ldflags="${LDFLAGS}" \
			./cmd/...
.PHONY: dev

# test runs the tests
test:
	@go test -timeout=60s -parallel=10 ./...
.PHONY: test

# xc compiles all the binaries using the local go installation
xc:
	@for OS in $(XC_OS); do \
		for ARCH in $(XC_ARCH); do \
			env \
				CGO_ENABLED=0 \
				GOOS=$${OS} \
				GOARCH=$${ARCH} \
				go build \
					-a \
					-o "pkg/$${OS}_$${ARCH}/${NAME}" \
					-ldflags "${LDFLAGS}"
					./cmd/... ; \
		done \
	done
.PHONY: xc

website: website-sync
	@(cd .vault-master/website && make website)
.PHONY: website

website-sync:
	@if [[ -d .vault-master ]]; then \
		cd .vault-master; \
		git reset --quiet --hard HEAD; \
		git clean --quiet -ffxd; \
		git pull --quiet origin master; \
	else \
		git clone --quiet --depth=1 https://github.com/hashicorp/vault .vault-master; \
	fi
	rsync --quiet -a website/ .vault-master/website
