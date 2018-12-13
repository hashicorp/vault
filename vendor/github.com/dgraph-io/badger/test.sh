#!/bin/bash

set -e

# Ensure that we can compile the binary.
pushd badger
go build -v .
popd

# Run the memory intensive tests first.
go test -v --manual=true -run='TestBigKeyValuePairs$'
go test -v --manual=true -run='TestPushValueLogLimit'

# Run the special Truncate test.
rm -R p || true
go test -v --manual=true -run='TestTruncateVlogNoClose$' .
truncate --size=4096 p/000000.vlog
go test -v --manual=true -run='TestTruncateVlogNoClose2$' .
go test -v --manual=true -run='TestTruncateVlogNoClose3$' .
rm -R p

# Then the normal tests.
go test -v --vlog_mmap=true -race ./...
go test -v --vlog_mmap=false -race ./...
