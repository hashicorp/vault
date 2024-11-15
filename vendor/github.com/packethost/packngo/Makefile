IMG ?= golang:1.16

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

lint:
	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.34.1 golangci-lint run -v

test:
	$(GO) test ./...

