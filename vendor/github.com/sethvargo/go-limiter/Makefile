GOFMT_FILES = $(shell go list -f '{{.Dir}}' ./...)

benchmarks:
	@(cd benchmarks/ && go test -bench=. -benchmem -benchtime=1s ./...)
.PHONY: benchmarks

fmtcheck:
	@command -v goimports > /dev/null 2>&1 || (cd tools/ && go install golang.org/x/tools/cmd/goimports@latest)
	@CHANGES="$$(goimports -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run goimports -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi
	@# Annoyingly, goimports does not support the simplify flag.
	@CHANGES="$$(gofmt -s -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run gofmt -s -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi
.PHONY: fmtcheck

spellcheck:
	@command -v misspell > /dev/null 2>&1 || (cd tools/ && go install github.com/client9/misspell/cmd/misspell@latest)
	@misspell -locale="US" -error -source="text" **/*
.PHONY: spellcheck

staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || (cd tools/ && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@staticcheck -checks="all" -tests $(GOFMT_FILES)
.PHONY: staticcheck

test:
	@go test \
		-count=1 \
		-shuffle=on \
		-short \
		-timeout=5m \
		-vet=all \
		./...
.PHONY: test

test-acc:
	@go test \
		-count=1 \
		-shuffle=on \
		-race \
		-timeout=10m \
		-vet=all \
		./...
.PHONY: test-acc
