devsetup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	go install github.com/vektra/mockery/v2@v2.26.1

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
	mockery --name=connectionManager --output=. --testonly --inpackage
	mockery --name=kvProvider --output=. --testonly --inpackage
	mockery --name=httpProvider --output=. --testonly --inpackage
	mockery --name=diagnosticsProvider --output=. --testonly --inpackage
	mockery --name=mgmtProvider --output=. --testonly --inpackage
	mockery --name=analyticsProvider --output=. --testonly --inpackage
	mockery --name=queryProvider --output=. --testonly --inpackage
	mockery --name=searchProvider --output=. --testonly --inpackage
	mockery --name=viewProvider --output=. --testonly --inpackage
	mockery --name=waitUntilReadyProvider --output=. --testonly --inpackage
	mockery --name=kvCapabilityVerifier --output=. --testonly --inpackage
	mockery --name=kvProviderCoreProvider --output=. --testonly --inpackage
	mockery --name=queryProviderCoreProvider --output=. --testonly --inpackage
	mockery --name=searchProviderCoreProvider --output=. --testonly --inpackage
	mockery --name=searchCapabilityVerifier --output=. --testonly --inpackage
	mockery --name=viewProviderCoreProvider --output=. --testonly --inpackage
	mockery --name=analyticsProviderCoreProvider --output=. --testonly --inpackage
	mockery --name=diagnosticsProviderCoreProvider --output=. --testonly --inpackage
	# pendingOp is manually mocked

.PHONY: all test devsetup fasttest lint cover check bench updatetestcases updatemocks
