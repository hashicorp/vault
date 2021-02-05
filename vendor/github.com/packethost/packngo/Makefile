IMG ?= golang:1.15

# enable go modules, disabled CGO

GOENV ?= GO111MODULE=on CGO_ENABLED=0
export GO111MODULE=on
export CGO_ENABLED=0

# we build in a docker image, unless we are set to BUILD=local
GO ?= docker run --rm -v $(PWD):/app -w /app $(IMG) env $(GOENV)
ifeq ($(BUILD),local)
GO = 
endif


build:
	$(GO) go build -i -v ./...

golangci-lint:
ifeq (, $(shell which golangci-lint))
	$(GO) go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.33.0
endif

golint:
ifeq (, $(shell which golint))
	$(GO) go get -u golang.org/x/lint/golint
endif

lint: golint golangci-lint
	$(GO) golangci-lint run ./...
	$(GO) go vet ./...
	$(GO) gofmt -d .

test:
	$(GO) test ./...

