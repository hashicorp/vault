devsetup:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/vektra/mockery/.../

test:
	go test ./...
fasttest:
	go test -short ./...

cover:
	go test -coverprofile=cover.out ./...

lint:
	golangci-lint run -v

check: lint
	go test -cover -race ./...

updatemocks:
	mockery -name dispatcher -output . -testonly -inpkg
	mockery -name tracerManager -output . -testonly -inpkg
	mockery -name configManager -output . -testonly -inpkg

.PHONY: all test devsetup fasttest lint cover checkerrs checkfmt checkvet checkiea checkspell check updatemocks
