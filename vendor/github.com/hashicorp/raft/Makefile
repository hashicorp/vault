DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
TEST_RESULTS_DIR?=/tmp/test-results

test:
	go test $(TESTARGS) -timeout=60s -race .
	go test $(TESTARGS) -timeout=60s -tags batchtest -race .

integ: test
	INTEG_TESTS=yes go test $(TESTARGS) -timeout=25s -run=Integ .
	INTEG_TESTS=yes go test $(TESTARGS) -timeout=25s -tags batchtest -run=Integ .

ci.test-norace:
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s -tags batchtest

ci.test:
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s -race .
	gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-test.xml -- -timeout=60s -race -tags batchtest .

ci.integ: ci.test
	INTEG_TESTS=yes gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-integ.xml -- -timeout=25s -run=Integ .
	INTEG_TESTS=yes gotestsum --format=short-verbose --junitfile $(TEST_RESULTS_DIR)/gotestsum-report-integ.xml -- -timeout=25s -run=Integ -tags batchtest .

fuzz:
	go test $(TESTARGS) -timeout=500s ./fuzzy
	go test $(TESTARGS) -timeout=500s -tags batchtest ./fuzzy

deps:
	go get -t -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

cov:
	INTEG_TESTS=yes gocov test github.com/hashicorp/raft | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

.PHONY: test cov integ deps
