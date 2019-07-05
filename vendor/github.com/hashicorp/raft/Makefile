DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
TEST_RESULTS_DIR?=/tmp/test-results

test:
	go test -timeout=60s -race .

integ: test
	INTEG_TESTS=yes go test -timeout=25s -run=Integ .

ci.test-norace:
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s

ci.test:
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s -race .

ci.integ: ci.test
	INTEG_TESTS=yes gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-integ.xml -- -timeout=25s -run=Integ .

fuzz:
	go test -timeout=300s ./fuzzy

deps:
	go get -t -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

cov:
	INTEG_TESTS=yes gocov test github.com/hashicorp/raft | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

.PHONY: test cov integ deps
