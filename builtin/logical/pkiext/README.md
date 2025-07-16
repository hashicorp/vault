# What is `pkiext`?

`pkiext` exists to split the Docker tests into a separate package from the
main PKI tests. Because the Docker tests execute in a smaller runner with
fewer resources, and we were hitting timeouts waiting for the entire PKI
test suite to run, we need to split the larger non-Docker PKI tests from
the smaller Docker tests, to ensure the former can execute.

This package should lack any non-test related targets.
