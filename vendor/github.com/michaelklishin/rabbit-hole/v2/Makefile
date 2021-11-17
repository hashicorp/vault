.DEFAULT_GOAL := test

.PHONY: all
all: test

.PHONY: test
test:
	go test -race -v
.PHONY: tests
tests: test

COVER_FILE := coverage
.PHONY: cover
cover:
	go test -v -test.coverprofile="$(COVER_FILE).prof"
	sed -i.bak 's|_'$(GOPATH)'|.|g' $(COVER_FILE).prof
	go tool cover -html=$(COVER_FILE).prof -o $(COVER_FILE).html
	rm $(COVER_FILE).prof*

.PHONY: ginkgo
ginkgo:
	command -v ginkgo || go install github.com/onsi/ginkgo/ginkgo
	ginkgo -race -v

.PHONY: docker
docker:
	docker run --rm \
	  --interactive --tty --entrypoint /bin/bash \
	  --volume $(CURDIR):/usr/src/app --workdir /usr/src/app \
	  golang:1.12
