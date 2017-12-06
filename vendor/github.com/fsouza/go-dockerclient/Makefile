.PHONY: \
	all \
	lint \
	vet \
	fmt \
	fmtcheck \
	pretest \
	test \
	integration \
	clean

all: test

lint:
	@ go get -v github.com/golang/lint/golint
	[ -z "$$(golint . | grep -v 'type name will be used as docker.DockerInfo' | grep -v 'context.Context should be the first' | tee /dev/stderr)" ]

vet:
	go vet $$(go list ./... | grep -v vendor)

fmt:
	gofmt -s -w $$(go list ./... | grep -v vendor)

fmtcheck:
	[ -z "$$(gofmt -s -d $$(go list ./... | grep -v vendor) | tee /dev/stderr)" ]

testdeps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

pretest: testdeps lint vet fmtcheck

gotest:
	go test -race $$(go list ./... | grep -v vendor)

test: pretest gotest

integration:
	go test -tags docker_integration -run TestIntegration -v

clean:
	go clean ./...
