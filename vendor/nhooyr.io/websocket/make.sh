#!/bin/sh
set -eu
cd -- "$(dirname "$0")"

echo "=== fmt.sh"
./ci/fmt.sh
echo "=== lint.sh"
./ci/lint.sh
echo "=== test.sh"
./ci/test.sh "$@"
echo "=== bench.sh"
./ci/bench.sh
