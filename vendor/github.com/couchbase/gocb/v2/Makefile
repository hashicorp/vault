devsetup:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/vektra/mockery/.../
	git submodule update --remote --init --recursive

test:
	go test ./
fasttest:
	go test -short ./

cover:
	go test -coverprofile=cover.out ./

lint:
	golangci-lint run -v

check: lint
	go test -short -cover -race ./

bench:
	go test -bench=. -run=none --disable-logger=true

updatetestcases:
	git submodule update --remote --init --recursive

updatemocks:
	mockery -name connectionManager -output . -testonly -inpkg
	mockery -name kvProvider -output . -testonly -inpkg
	mockery -name httpProvider -output . -testonly -inpkg
	mockery -name diagnosticsProvider -output . -testonly -inpkg
	mockery -name mgmtProvider -output . -testonly -inpkg
	mockery -name analyticsProvider -output . -testonly -inpkg
	mockery -name queryProvider -output . -testonly -inpkg
	mockery -name searchProvider -output . -testonly -inpkg
	mockery -name viewProvider -output . -testonly -inpkg
	mockery -name waitUntilReadyProvider -output . -testonly -inpkg
	# pendingOp is manually mocked

.PHONY: all test devsetup fasttest lint cover check bench updatetestcases updatemocks
