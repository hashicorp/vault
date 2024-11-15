.PHONY: vet fmt lint test mocks envvars generate

# Make target to generate resource scaffolding for specified RESOURCE
generate: check-resource
	@cd ./scripts/generate_resource; \
	go mod tidy; \
	go run . $(RESOURCE) ;

vet:
	go vet

fmt:
	gofmt -s -l -w .

fmtcheck:
	./scripts/gofmtcheck.sh

lint:
	golangci-lint run .

test:
	go test ./... $(TESTARGS) -timeout=30m

# Make target to generate mocks for specified FILENAME
mocks: check-filename
	@echo "mockgen -source=$(FILENAME) -destination=mocks/$(subst .go,_mocks.go,$(FILENAME)) -package=mocks" >> generate_mocks.sh
	./generate_mocks.sh

envvars:
	./scripts/setup-test-envvars.sh

check-filename:
ifndef FILENAME
	$(error Missing FILENAME param. Example usage: FILENAME=example_resource.go make mocks)
endif

check-resource:
ifndef RESOURCE
	$(error Missing RESOURCE param. Example usage: RESOURCE=foo_bar make generate)
endif

