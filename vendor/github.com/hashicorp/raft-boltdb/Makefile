DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

.PHONY: test deps

test:
	go test -timeout=30s ./...

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

