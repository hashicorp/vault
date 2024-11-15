devsetup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/vektra/mockery/v2@v2.38.0

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

bench:
	go test -run ^$$ -bench . --disable-logger

updatemocks:
	mockery --name dispatcher --output . --testonly --inpackage
	mockery --name configManager --output . --testonly --inpackage
	mockery --name httpComponentInterface --output . --testonly --inpackage

.PHONY: all test devsetup fasttest lint cover checkerrs checkfmt checkvet checkiea checkspell check bench updatemocks
