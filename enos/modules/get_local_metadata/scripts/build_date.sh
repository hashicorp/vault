#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -eu -o pipefail

pushd "$(git rev-parse --show-toplevel)" > /dev/null
make ci-get-date
popd > /dev/null
