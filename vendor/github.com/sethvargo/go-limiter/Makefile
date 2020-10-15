VETTERS = "asmdecl,assign,atomic,bools,buildtag,cgocall,composites,copylocks,errorsas,httpresponse,loopclosure,lostcancel,nilfunc,printf,shift,stdmethods,structtag,tests,unmarshal,unreachable,unsafeptr,unusedresult"
GOFMT_FILES = $(shell go list -f '{{.Dir}}' ./...)

benchmarks:
	@(cd benchmarks/ && go test -bench=. -benchmem -benchtime=1s ./...)
.PHONY: benchmarks

fmtcheck:
	@command -v goimports > /dev/null 2>&1 || (cd tools/ && go get golang.org/x/tools/cmd/goimports)
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
	@command -v misspell > /dev/null 2>&1 || (cd tools/ && go get github.com/client9/misspell/cmd/misspell)
	@misspell -locale="US" -error -source="text" **/*
.PHONY: spellcheck

staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || (cd tools/ && go get honnef.co/go/tools/cmd/staticcheck)
	@staticcheck -checks="all" -tests $(GOFMT_FILES)
.PHONY: staticcheck

test:
	@go test \
		-count=1 \
		-short \
		-timeout=5m \
		-vet="${VETTERS}" \
		./...
.PHONY: test

test-acc:
	@go test \
		-count=1 \
		-race \
		-timeout=10m \
		-vet="${VETTERS}" \
		./...
.PHONY: test-acc
